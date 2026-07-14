package com.team10.mes.ai.runtime;

import java.util.Optional;

public final class AiRuntimeContext {
  private static final ThreadLocal<RuntimeState> CURRENT = new ThreadLocal<>();

  private AiRuntimeContext() {}

  public static Scope open(RuntimeState state) {
    RuntimeState previous = CURRENT.get();
    CURRENT.set(state);
    return new Scope(previous);
  }

  public static Optional<RuntimeState> current() {
    return Optional.ofNullable(CURRENT.get());
  }

  public static boolean isCancelled() {
    return current().map(RuntimeState::isCancelled).orElse(false);
  }

  public static String controlCursor() {
    return current().map(RuntimeState::controlCursor).orElse("");
  }

  public static void setControlCursor(String cursor) {
    current().ifPresent(state -> state.setControlCursor(cursor));
  }

  public static void cancel() {
    current().ifPresent(RuntimeState::cancel);
  }

  public static final class Scope implements AutoCloseable {
    private final RuntimeState previous;
    private boolean closed;

    private Scope(RuntimeState previous) {
      this.previous = previous;
    }

    @Override
    public void close() {
      if (closed) {
        return;
      }
      closed = true;
      if (previous == null) {
        CURRENT.remove();
      } else {
        CURRENT.set(previous);
      }
    }
  }
}
