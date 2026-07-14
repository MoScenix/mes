package com.team10.mes.ai.node.designer;

import com.team10.mes.ai.service.AiService.PendingInterrupt;
import com.team10.mes.ai.service.AiService.Question;
import java.util.List;
import java.util.Optional;
import java.util.function.Consumer;

/** Spring AI replacement for the Eino designer agent turn. */
public final class DesignerNode {
  @FunctionalInterface
  public interface ModelRunner {
    String stream(String resumeContext, Consumer<String> chunks) throws Exception;
  }

  @FunctionalInterface
  public interface QuestionDetector {
    Optional<List<Question>> detect(String output);
  }

  public record Result(String output, List<Question> questions, PendingInterrupt interrupt) {
    public boolean interrupted() {
      return interrupt != null;
    }
  }

  public Result run(
      ModelRunner model,
      QuestionDetector detector,
      Consumer<String> chunks,
      java.util.function.Function<List<Question>, PendingInterrupt> interruptFactory)
      throws Exception {
    String output = model.stream("", chunks);
    Optional<List<Question>> questions = detector.detect(output);
    return questions
        .map(value -> new Result(output, value, interruptFactory.apply(value)))
        .orElseGet(() -> new Result(output, List.of(), null));
  }
}
