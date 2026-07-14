package com.team10.mes.ai.runtime;

import java.util.concurrent.locks.ReentrantLock;

public final class ConcurrentTextBuffer {
  private final ReentrantLock lock = new ReentrantLock();
  private final StringBuilder value = new StringBuilder();

  public void append(String text) {
    if (text == null || text.isEmpty()) {
      return;
    }
    lock.lock();
    try {
      value.append(text);
    } finally {
      lock.unlock();
    }
  }

  public String value() {
    lock.lock();
    try {
      return value.toString();
    } finally {
      lock.unlock();
    }
  }

  public void clear() {
    set("");
  }

  public void set(String text) {
    lock.lock();
    try {
      value.setLength(0);
      if (text != null) {
        value.append(text);
      }
    } finally {
      lock.unlock();
    }
  }
}
