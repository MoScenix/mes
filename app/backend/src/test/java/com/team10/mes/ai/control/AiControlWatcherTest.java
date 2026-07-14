package com.team10.mes.ai.control;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;

import com.team10.mes.ai.service.AiService.AiEvent;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;
import org.junit.jupiter.api.Test;

class AiControlWatcherTest {
  @Test
  void dispatchesPushAndAdvancesCursor() {
    AiEvent push =
        new AiEvent("1-0", "7", "push", "Assistant", "next", "", "", "", "", 1, List.of());
    AtomicInteger calls = new AtomicInteger();
    AiControlWatcher[] holder = new AiControlWatcher[1];
    holder[0] =
        new AiControlWatcher(
            (app, cursor, block, count) -> List.of(push),
            7,
            "0-0",
            100,
            10,
            new AiControlWatcher.Handler() {
              public void onPush(AiEvent e) {
                calls.incrementAndGet();
                holder[0].close();
              }

              public void onCancel(AiEvent e) {}

              public void onAnswer(AiEvent e) {}
            });
    holder[0].run();
    assertEquals(1, calls.get());
    assertEquals("1-0", holder[0].cursor());
  }

  @Test
  void exitsQuietlyWhenClosedBlockedReadIsInterrupted() {
    AiControlWatcher[] holder = new AiControlWatcher[1];
    holder[0] =
        new AiControlWatcher(
            (app, cursor, block, count) -> {
              holder[0].close();
              throw new RuntimeException("Redis command interrupted");
            },
            7,
            "0-0",
            100,
            10,
            new AiControlWatcher.Handler() {
              public void onPush(AiEvent e) {}

              public void onCancel(AiEvent e) {}

              public void onAnswer(AiEvent e) {}
            });
    assertDoesNotThrow(() -> holder[0].run());
  }
}
