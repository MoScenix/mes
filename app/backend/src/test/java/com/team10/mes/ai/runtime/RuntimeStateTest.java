package com.team10.mes.ai.runtime;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.IntStream;
import org.junit.jupiter.api.Test;

class RuntimeStateTest {
  @Test
  void bufferSupportsConcurrentWritesAndReplacement() {
    ConcurrentTextBuffer buffer = new ConcurrentTextBuffer();
    IntStream.range(0, 100).parallel().forEach(ignored -> buffer.append("x"));
    assertThat(buffer.value()).hasSize(100);

    buffer.set("answer");
    assertThat(buffer.value()).isEqualTo("answer");
    buffer.clear();
    assertThat(buffer.value()).isEmpty();
  }

  @Test
  void cancelAndStopHaveDistinctStateSemantics() {
    AtomicInteger cancelled = new AtomicInteger();
    RuntimeState state = new RuntimeState(cancelled::incrementAndGet);
    state.stop();
    assertThat(state.isCancelled()).isFalse();

    state.cancel();
    assertThat(state.isCancelled()).isTrue();
    assertThat(cancelled).hasValue(2);
  }

  @Test
  void runtimeScopesRestoreNestedStateAndCursor() {
    RuntimeState outer = new RuntimeState(null);
    RuntimeState inner = new RuntimeState(null);

    try (AiRuntimeContext.Scope ignored = AiRuntimeContext.open(outer)) {
      AiRuntimeContext.setControlCursor("1-0");
      try (AiRuntimeContext.Scope nested = AiRuntimeContext.open(inner)) {
        AiRuntimeContext.setControlCursor("2-0");
        AiRuntimeContext.cancel();
        assertThat(AiRuntimeContext.controlCursor()).isEqualTo("2-0");
      }
      assertThat(AiRuntimeContext.controlCursor()).isEqualTo("1-0");
      assertThat(AiRuntimeContext.isCancelled()).isFalse();
    }

    assertThat(AiRuntimeContext.current()).isEmpty();
  }
}
