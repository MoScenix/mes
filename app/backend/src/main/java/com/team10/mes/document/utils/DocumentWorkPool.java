package com.team10.mes.document.utils;

import com.team10.mes.ai.workpool.BoundedAiWorkPool;
import com.team10.mes.ai.workpool.WorkPoolConfig;
import jakarta.annotation.PreDestroy;
import java.time.Duration;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import org.springframework.stereotype.Component;

@Component
public final class DocumentWorkPool implements AutoCloseable {
  private final BoundedAiWorkPool delegate;

  public DocumentWorkPool(DocumentProperties properties) {
    DocumentProperties.Index index = properties.getIndex();
    this.delegate =
        new BoundedAiWorkPool(
            new WorkPoolConfig(
                positive(index.getMinWorkers(), 8),
                positive(index.getMaxWorkers(), 32),
                positive(index.getQueueSize(), 256),
                positive(index.getScaleUpThreshold(), 4),
                Math.max(0, index.getScaleDownThreshold()),
                Duration.ofSeconds(positive(index.getIdleTimeoutSeconds(), 30))));
  }

  public CompletableFuture<Void> submit(Runnable task) {
    Objects.requireNonNull(task, "task");
    CompletableFuture<Void> future = new CompletableFuture<>();
    try {
      delegate.submit(
          () -> {
            try {
              task.run();
              future.complete(null);
            } catch (Throwable error) {
              future.completeExceptionally(error);
            }
          });
    } catch (Throwable error) {
      if (error instanceof InterruptedException) Thread.currentThread().interrupt();
      future.completeExceptionally(error);
    }
    return future;
  }

  @PreDestroy
  @Override
  public void close() {
    delegate.close();
  }

  private static int positive(int value, int fallback) {
    return value > 0 ? value : fallback;
  }

  private static long positive(long value, long fallback) {
    return value > 0 ? value : fallback;
  }
}
