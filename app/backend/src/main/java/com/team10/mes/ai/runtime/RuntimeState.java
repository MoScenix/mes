package com.team10.mes.ai.runtime;

import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicReference;

public final class RuntimeState {
  private final ConcurrentTextBuffer buffer = new ConcurrentTextBuffer();
  private final AtomicReference<String> controlCursor = new AtomicReference<>("0");
  private final AtomicBoolean cancelled = new AtomicBoolean();
  private final Runnable cancelAction;

  public RuntimeState(Runnable cancelAction) {
    this.cancelAction = cancelAction == null ? () -> {} : cancelAction;
  }

  public ConcurrentTextBuffer buffer() {
    return buffer;
  }

  public void cancel() {
    cancelled.set(true);
    cancelAction.run();
  }

  public void stop() {
    cancelAction.run();
  }

  public boolean isCancelled() {
    return cancelled.get();
  }

  public String controlCursor() {
    return controlCursor.get();
  }

  public void setControlCursor(String cursor) {
    if (cursor != null && !cursor.isEmpty()) {
      controlCursor.set(cursor);
    }
  }
}
