package com.team10.mes.ai.graph;

import static org.junit.jupiter.api.Assertions.*;

import com.team10.mes.ai.node.coder.CoderNode;
import com.team10.mes.ai.node.designer.DesignerNode;
import com.team10.mes.ai.service.AiService.*;
import java.util.*;
import java.util.concurrent.atomic.AtomicReference;
import org.junit.jupiter.api.Test;

class AiGraphTest {
  private final AiGraph graph = new AiGraph(new DesignerNode(), new CoderNode());

  @Test
  void persistsInterruptAndResumesToCoder() throws Exception {
    AtomicReference<QuestionCheckpoint> saved = new AtomicReference<>();
    AtomicReference<String> committed = new AtomicReference<>();
    var interrupted =
        graph.run(
            (context, chunks) -> "ask",
            text -> Optional.of(List.of(new Question("Which?", List.of("A")))),
            chunk -> {},
            qs -> new PendingInterrupt("q1", "Assistant", "Which?", "{}"),
            p -> new QuestionCheckpoint("cp", List.of(p), "ask", "1-0", 1),
            saved::set,
            committed::set);
    assertEquals(AiGraph.Status.INTERRUPTED, interrupted.status());
    assertNotNull(saved.get());
    assertNull(committed.get());
    AtomicReference<String> resumeContext = new AtomicReference<>();
    var done =
        graph.resume(
            saved.get(),
            "A",
            "partial",
            (context, chunks) -> {
              resumeContext.set(context);
              return "final";
            },
            text -> Optional.empty(),
            chunk -> {},
            qs -> null,
            p -> null,
            saved::set,
            committed::set);
    assertEquals(AiGraph.Status.DONE, done.status());
    assertEquals("final", committed.get());
    assertTrue(resumeContext.get().contains("partial"));
    assertTrue(resumeContext.get().contains("A"));
  }
}
