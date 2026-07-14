package com.team10.mes.ai.workpool;

import jakarta.annotation.PreDestroy;
import java.time.Duration;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public final class AiWorkPool {
  private final BoundedAiWorkPool delegate;

  public AiWorkPool(
      @Value("${mes.ai.workpool.min-workers:1}") int min,
      @Value("${mes.ai.workpool.max-workers:4}") int max,
      @Value("${mes.ai.workpool.queue-size:128}") int queue,
      @Value("${mes.ai.workpool.scale-up-threshold:4}") int up,
      @Value("${mes.ai.workpool.scale-down-threshold:0}") int down,
      @Value("${mes.ai.workpool.idle-timeout-seconds:30}") long idle) {
    delegate =
        new BoundedAiWorkPool(
            new WorkPoolConfig(min, max, queue, up, down, Duration.ofSeconds(idle)));
  }

  public void submit(Runnable task) {
    try {
      delegate.submit(task);
    } catch (InterruptedException e) {
      Thread.currentThread().interrupt();
      throw new IllegalStateException("AI task enqueue interrupted", e);
    }
  }

  public WorkPoolStats stats() {
    return delegate.stats();
  }

  @PreDestroy
  public void close() {
    delegate.close();
  }
}
