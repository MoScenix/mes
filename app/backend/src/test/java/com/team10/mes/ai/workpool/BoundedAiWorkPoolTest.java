package com.team10.mes.ai.workpool;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import java.time.Duration;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;
import org.junit.jupiter.api.Test;

class BoundedAiWorkPoolTest {
  @Test
  void drainsAcceptedTasksBeforeStopping() throws Exception {
    AtomicInteger completed = new AtomicInteger();
    BoundedAiWorkPool pool = new BoundedAiWorkPool(config(1, 2, 4));

    for (int index = 0; index < 4; index++) {
      pool.submit(completed::incrementAndGet);
    }

    assertThat(pool.stop(Duration.ofSeconds(1))).isTrue();
    assertThat(completed).hasValue(4);
    assertThatThrownBy(() -> pool.submit(() -> {})).isInstanceOf(RejectedExecutionException.class);
  }

  @Test
  void rejectsTimedSubmissionWhenBoundedQueueStaysFull() throws Exception {
    CountDownLatch release = new CountDownLatch(1);
    BoundedAiWorkPool pool = new BoundedAiWorkPool(config(1, 1, 1));
    pool.submit(() -> await(release));
    awaitStats(pool, stats -> stats.running() == 1);
    pool.submit(() -> await(release));

    assertThatThrownBy(() -> pool.submit(() -> {}, Duration.ofMillis(30)))
        .isInstanceOf(RejectedExecutionException.class)
        .hasMessageContaining("queue is full");

    release.countDown();
    assertThat(pool.stop(Duration.ofSeconds(1))).isTrue();
  }

  @Test
  void scalesUpAndRetiresIdleWorkers() throws Exception {
    CountDownLatch release = new CountDownLatch(1);
    WorkPoolConfig config = new WorkPoolConfig(0, 3, 8, 1, 0, Duration.ofMillis(20));
    BoundedAiWorkPool pool = new BoundedAiWorkPool(config);
    for (int index = 0; index < 3; index++) {
      pool.submit(() -> await(release));
    }

    awaitStats(pool, stats -> stats.workers() == 3);
    release.countDown();
    awaitStats(pool, stats -> stats.workers() == 0 && stats.running() == 0);
    assertThat(pool.stop(Duration.ofSeconds(1))).isTrue();
  }

  @Test
  void extraWorkerWaitsForConfiguredIdleTimeout() throws Exception {
    BoundedAiWorkPool pool =
        new BoundedAiWorkPool(new WorkPoolConfig(0, 1, 2, 1, 0, Duration.ofMillis(350)));
    pool.submit(() -> {});
    awaitStats(pool, stats -> stats.workers() == 1 && stats.running() == 0);
    Thread.sleep(150);
    assertThat(pool.stats().workers()).isEqualTo(1);
    awaitStats(pool, stats -> stats.workers() == 0);
    assertThat(pool.stop(Duration.ofSeconds(1))).isTrue();
  }

  @Test
  void reportsTaskFailureWithoutLosingWorker() throws Exception {
    AtomicReference<Throwable> failure = new AtomicReference<>();
    AtomicInteger completed = new AtomicInteger();
    BoundedAiWorkPool pool =
        new BoundedAiWorkPool(config(1, 1, 2), (task, error) -> failure.set(error));
    pool.submit(
        () -> {
          throw new IllegalStateException("failed");
        });
    pool.submit(completed::incrementAndGet);

    assertThat(pool.stop(Duration.ofSeconds(1))).isTrue();
    assertThat(failure.get()).isInstanceOf(IllegalStateException.class);
    assertThat(completed).hasValue(1);
  }

  private static WorkPoolConfig config(int min, int max, int queue) {
    return new WorkPoolConfig(min, max, queue, 1, 0, Duration.ofMillis(20));
  }

  private static void await(CountDownLatch latch) {
    try {
      latch.await();
    } catch (InterruptedException exception) {
      Thread.currentThread().interrupt();
    }
  }

  private static void awaitStats(
      BoundedAiWorkPool pool, java.util.function.Predicate<WorkPoolStats> condition)
      throws InterruptedException {
    long deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(1);
    while (System.nanoTime() < deadline) {
      if (condition.test(pool.stats())) {
        return;
      }
      Thread.sleep(5);
    }
    throw new AssertionError("condition not met; stats=" + pool.stats());
  }
}
