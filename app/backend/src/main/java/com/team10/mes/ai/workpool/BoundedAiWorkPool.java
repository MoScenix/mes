package com.team10.mes.ai.workpool;

import java.time.Duration;
import java.util.HashSet;
import java.util.Objects;
import java.util.Set;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.BiConsumer;

public final class BoundedAiWorkPool implements AutoCloseable {
  private static final Duration WORKER_POLL_INTERVAL = Duration.ofMillis(100);

  private final WorkPoolConfig config;
  private final ArrayBlockingQueue<Runnable> tasks;
  private final BiConsumer<Runnable, Throwable> errorHandler;
  private final Set<Thread> workers = new HashSet<>();
  private final AtomicInteger running = new AtomicInteger();
  private int nextWorkerId;
  private volatile boolean closed;
  private volatile boolean shutdownNow;

  public BoundedAiWorkPool(WorkPoolConfig config) {
    this(config, (task, error) -> {});
  }

  public BoundedAiWorkPool(WorkPoolConfig config, BiConsumer<Runnable, Throwable> errorHandler) {
    this.config = Objects.requireNonNull(config, "config");
    this.errorHandler = Objects.requireNonNull(errorHandler, "errorHandler");
    this.tasks = new ArrayBlockingQueue<>(config.queueSize());
    synchronized (workers) {
      resizeLocked(config.minWorkers());
    }
  }

  public void submit(Runnable task) throws InterruptedException {
    submit(task, null);
  }

  public void submit(Runnable task, Duration timeout) throws InterruptedException {
    Objects.requireNonNull(task, "task");
    if (closed) {
      throw new RejectedExecutionException("AI work pool is closed");
    }

    boolean accepted = timeout == null ? putUnlessClosed(task) : offerUnlessClosed(task, timeout);
    if (!accepted) {
      throw new RejectedExecutionException("AI work pool queue is full");
    }
    scaleAfterSubmit();
  }

  public WorkPoolStats stats() {
    synchronized (workers) {
      return new WorkPoolStats(workers.size(), running.get(), tasks.size(), closed);
    }
  }

  public boolean stop(Duration timeout) throws InterruptedException {
    Objects.requireNonNull(timeout, "timeout");
    closed = true;
    long deadline = System.nanoTime() + timeout.toNanos();
    synchronized (workers) {
      workers.notifyAll();
      while (!workers.isEmpty()) {
        long remaining = deadline - System.nanoTime();
        if (remaining <= 0) {
          return false;
        }
        TimeUnit.NANOSECONDS.timedWait(workers, remaining);
      }
    }
    return true;
  }

  public void shutdownNow() {
    closed = true;
    shutdownNow = true;
    tasks.clear();
    synchronized (workers) {
      for (Thread worker : workers) {
        worker.interrupt();
      }
    }
  }

  @Override
  public void close() {
    try {
      if (!stop(Duration.ofSeconds(30))) {
        shutdownNow();
      }
    } catch (InterruptedException exception) {
      shutdownNow();
      Thread.currentThread().interrupt();
    }
  }

  private boolean putUnlessClosed(Runnable task) throws InterruptedException {
    while (!closed) {
      if (tasks.offer(task, WORKER_POLL_INTERVAL.toMillis(), TimeUnit.MILLISECONDS)) {
        return true;
      }
    }
    throw new RejectedExecutionException("AI work pool is closed");
  }

  private boolean offerUnlessClosed(Runnable task, Duration timeout) throws InterruptedException {
    if (timeout.isNegative()) {
      throw new IllegalArgumentException("timeout must not be negative");
    }
    long deadline = System.nanoTime() + timeout.toNanos();
    do {
      if (closed) {
        throw new RejectedExecutionException("AI work pool is closed");
      }
      long remaining = deadline - System.nanoTime();
      if (remaining <= 0) {
        return false;
      }
      long waitNanos = Math.min(remaining, WORKER_POLL_INTERVAL.toNanos());
      if (tasks.offer(task, waitNanos, TimeUnit.NANOSECONDS)) {
        return true;
      }
    } while (true);
  }

  private void scaleAfterSubmit() {
    synchronized (workers) {
      if (!closed
          && tasks.size() >= config.scaleUpThreshold()
          && workers.size() < config.maxWorkers()) {
        startWorkerLocked();
      }
    }
  }

  private void resizeLocked(int target) {
    int boundedTarget = Math.max(config.minWorkers(), Math.min(target, config.maxWorkers()));
    while (workers.size() < boundedTarget) {
      startWorkerLocked();
    }
  }

  private void startWorkerLocked() {
    Thread worker = new Thread(this::workerLoop, "ai-work-pool-" + nextWorkerId++);
    worker.setDaemon(true);
    workers.add(worker);
    worker.start();
  }

  private void workerLoop() {
    Thread current = Thread.currentThread();
    long idleSince = System.nanoTime();
    try {
      while (!shutdownNow) {
        if (closed && tasks.isEmpty()) {
          return;
        }
        Runnable task = tasks.poll(pollDuration().toNanos(), TimeUnit.NANOSECONDS);
        if (task != null) {
          runTask(task);
          idleSince = System.nanoTime();
          continue;
        }
        synchronized (workers) {
          if (!closed
              && workers.size() > config.minWorkers()
              && tasks.size() <= config.scaleDownThreshold()
              && System.nanoTime() - idleSince >= config.idleTimeout().toNanos()) {
            return;
          }
        }
      }
    } catch (InterruptedException exception) {
      if (!shutdownNow) {
        Thread.currentThread().interrupt();
      }
    } finally {
      synchronized (workers) {
        workers.remove(current);
        workers.notifyAll();
      }
    }
  }

  private Duration pollDuration() {
    return WORKER_POLL_INTERVAL;
  }

  private void runTask(Runnable task) {
    running.incrementAndGet();
    try {
      task.run();
    } catch (Throwable error) {
      try {
        errorHandler.accept(task, error);
      } catch (Throwable ignored) {
        // A broken observer must not terminate a worker.
      }
    } finally {
      running.decrementAndGet();
    }
  }
}
