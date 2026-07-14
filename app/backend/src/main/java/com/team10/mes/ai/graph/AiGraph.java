package com.team10.mes.ai.graph;

import com.team10.mes.ai.node.coder.CoderNode;
import com.team10.mes.ai.node.designer.DesignerNode;
import com.team10.mes.ai.service.AiService.PendingInterrupt;
import com.team10.mes.ai.service.AiService.QuestionCheckpoint;
import java.util.function.Consumer;

/** Explicit designer -> interrupt/resume -> coder graph. */
public final class AiGraph {
  public enum Status {
    DONE,
    INTERRUPTED
  }

  @FunctionalInterface
  public interface CheckpointWriter {
    void save(QuestionCheckpoint checkpoint);
  }

  public record Result(Status status, String output, QuestionCheckpoint checkpoint) {}

  private final DesignerNode designer;
  private final CoderNode coder;

  public AiGraph(DesignerNode designer, CoderNode coder) {
    this.designer = designer;
    this.coder = coder;
  }

  public Result run(
      DesignerNode.ModelRunner model,
      DesignerNode.QuestionDetector detector,
      Consumer<String> chunks,
      java.util.function.Function<
              java.util.List<com.team10.mes.ai.service.AiService.Question>, PendingInterrupt>
          interruptFactory,
      java.util.function.Function<PendingInterrupt, QuestionCheckpoint> checkpointFactory,
      CheckpointWriter checkpoints,
      CoderNode.Committer committer)
      throws Exception {
    DesignerNode.Result designed = designer.run(model, detector, chunks, interruptFactory);
    if (designed.interrupted()) {
      QuestionCheckpoint checkpoint = checkpointFactory.apply(designed.interrupt());
      checkpoints.save(checkpoint);
      return new Result(Status.INTERRUPTED, designed.output(), checkpoint);
    }
    return new Result(Status.DONE, coder.run(designed.output(), committer), null);
  }

  public Result resume(
      QuestionCheckpoint checkpoint,
      String answer,
      String buffer,
      DesignerNode.ModelRunner model,
      DesignerNode.QuestionDetector detector,
      Consumer<String> chunks,
      java.util.function.Function<
              java.util.List<com.team10.mes.ai.service.AiService.Question>, PendingInterrupt>
          interruptFactory,
      java.util.function.Function<PendingInterrupt, QuestionCheckpoint> checkpointFactory,
      CheckpointWriter checkpoints,
      CoderNode.Committer committer)
      throws Exception {
    if (checkpoint == null
        || checkpoint.pendingInterrupts() == null
        || checkpoint.pendingInterrupts().isEmpty())
      throw new IllegalStateException("graph has no interrupted checkpoint");
    String context =
        "Previous partial output:\n"
            + (buffer == null ? "" : buffer)
            + "\nUser answer:\n"
            + (answer == null ? "" : answer);
    DesignerNode.ModelRunner resumed = (ignored, sink) -> model.stream(context, sink);
    return run(
        resumed, detector, chunks, interruptFactory, checkpointFactory, checkpoints, committer);
  }
}
