package com.team10.mes.ai.task;

import static org.junit.jupiter.api.Assertions.*;

import com.team10.mes.ai.service.AiService.Identity;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import org.junit.jupiter.api.Test;

class ChatTaskTest {
  @Test
  void cancelInterruptsExecutingThread() throws Exception {
    CountDownLatch started = new CountDownLatch(1);
    AtomicBoolean interrupted = new AtomicBoolean();
    ChatTask task =
        new ChatTask(
            1,
            new Identity(1, "worker"),
            false,
            ignored -> {
              started.countDown();
              try {
                Thread.sleep(30000);
              } catch (InterruptedException e) {
                interrupted.set(true);
                Thread.currentThread().interrupt();
              }
            },
            (t, s) -> {});
    Thread thread = Thread.ofVirtual().start(task);
    assertTrue(started.await(1, TimeUnit.SECONDS));
    task.cancel();
    thread.join(1000);
    assertTrue(interrupted.get());
    assertTrue(task.runtime().isCancelled());
    assertFalse(thread.isAlive());
  }
}
