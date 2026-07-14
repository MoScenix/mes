package com.team10.mes.ai.workpool;

public record WorkPoolStats(int workers, int running, int queued, boolean closed) {}
