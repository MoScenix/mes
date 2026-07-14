package com.team10.mes.ai.control;

import com.team10.mes.ai.service.AiService.AiEvent;
import com.team10.mes.ai.state.RedisAiStore;
import java.util.concurrent.atomic.AtomicBoolean;

/** Redis control-stream watcher equivalent to node/control/watch.go. */
public final class AiControlWatcher implements Runnable, AutoCloseable {
  public interface Handler {
    void onPush(AiEvent event);

    void onCancel(AiEvent event);

    void onAnswer(AiEvent event);
  }

  @FunctionalInterface
  public interface ControlReader {
    java.util.List<AiEvent> read(long historyId, String cursor, long blockMs, int count);
  }

  private final ControlReader reader;
  private final long historyId;
  private final Handler handler;
  private final long blockMs;
  private final int count;
  private final AtomicBoolean running = new AtomicBoolean(true);
  private String cursor;

  public AiControlWatcher(
      RedisAiStore store, long historyId, String cursor, long blockMs, int count, Handler handler) {
    this(store::controls, historyId, cursor, blockMs, count, handler);
  }

  AiControlWatcher(
      ControlReader reader,
      long historyId,
      String cursor,
      long blockMs,
      int count,
      Handler handler) {
    this.reader = reader;
    this.historyId = historyId;
    this.cursor = cursor == null || cursor.isBlank() ? "$" : cursor;
    this.blockMs = Math.max(1, blockMs);
    this.count = Math.max(1, count);
    this.handler = handler;
  }

  @Override
  public void run() {
    while (running.get() && !Thread.currentThread().isInterrupted()) {
      try {
        for (AiEvent event : reader.read(historyId, cursor, blockMs, count)) {
          cursor = event.id();
          switch (event.type()) {
            case "push" -> handler.onPush(event);
            case "cancel" -> handler.onCancel(event);
            case "answer" -> handler.onAnswer(event);
            default -> {}
          }
        }
      } catch (RuntimeException error) {
        if (!running.get() || Thread.currentThread().isInterrupted()) return;
        throw error;
      }
    }
  }

  public String cursor() {
    return cursor;
  }

  @Override
  public void close() {
    running.set(false);
  }
}
