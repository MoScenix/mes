package com.team10.mes.ai.workpool;

import java.time.Duration;

public record WorkPoolConfig(
    int minWorkers,
    int maxWorkers,
    int queueSize,
    int scaleUpThreshold,
    int scaleDownThreshold,
    Duration idleTimeout) {

  public WorkPoolConfig {
    minWorkers = Math.max(0, minWorkers);
    maxWorkers = Math.max(1, maxWorkers);
    minWorkers = Math.min(minWorkers, maxWorkers);
    queueSize = queueSize > 0 ? queueSize : maxWorkers;
    scaleUpThreshold =
        scaleUpThreshold > 0 && scaleUpThreshold <= queueSize ? scaleUpThreshold : queueSize;
    scaleDownThreshold = Math.max(0, scaleDownThreshold);
    scaleDownThreshold = Math.min(scaleDownThreshold, scaleUpThreshold - 1);
    idleTimeout =
        idleTimeout == null || idleTimeout.isNegative() || idleTimeout.isZero()
            ? Duration.ofSeconds(30)
            : idleTimeout;
  }
}
