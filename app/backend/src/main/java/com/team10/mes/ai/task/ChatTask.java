package com.team10.mes.ai.task;

import com.team10.mes.ai.runtime.AiRuntimeContext;
import com.team10.mes.ai.runtime.RuntimeState;
import com.team10.mes.ai.service.AiService.Identity;
import java.util.concurrent.atomic.AtomicBoolean;

/** Chat task lifecycle equivalent to task/chat.go. */
public final class ChatTask implements Runnable {
  @FunctionalInterface
  public interface Runner {
    void run(ChatTask task) throws Exception;
  }

  @FunctionalInterface
  public interface StateMarker {
    void mark(ChatTask task, String status);
  }

  private final long appId;
  private final Identity identity;
  private final boolean resume;
  private final Runner runner;
  private final StateMarker marker;
  private final AtomicBoolean initialized = new AtomicBoolean();
  private RuntimeState runtime;
  private volatile Thread executingThread;

  public ChatTask(
      long appId, Identity identity, boolean resume, Runner runner, StateMarker marker) {
    this.appId = appId;
    this.identity = identity;
    this.resume = resume;
    this.runner = runner;
    this.marker = marker;
  }

  public synchronized RuntimeState init() {
    if (initialized.compareAndSet(false, true))
      runtime =
          new RuntimeState(
              () -> {
                Thread thread = executingThread;
                if (thread != null) thread.interrupt();
              });
    return runtime;
  }

  public void enqueue() {
    init();
    marker.mark(this, "queued");
  }

  @Override
  public void run() {
    init();
    if (runtime.isCancelled()) return;
    executingThread = Thread.currentThread();
    marker.mark(this, "running");
    try (AiRuntimeContext.Scope ignored = AiRuntimeContext.open(runtime)) {
      runner.run(this);
    } catch (RuntimeException e) {
      throw e;
    } catch (Exception e) {
      throw new IllegalStateException(e);
    } finally {
      executingThread = null;
      runtime.stop();
    }
  }

  public void resume() {
    run();
  }

  public void cancel() {
    init().cancel();
  }

  public long appId() {
    return appId;
  }

  public Identity identity() {
    return identity;
  }

  public boolean needsResume() {
    return resume;
  }

  public RuntimeState runtime() {
    return init();
  }
}
