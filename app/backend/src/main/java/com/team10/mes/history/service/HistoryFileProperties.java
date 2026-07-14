package com.team10.mes.history.service;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "mes.file")
public class HistoryFileProperties {
  private long bigThresholdBytes;
  private long chunkMinSize;
  private long chunkMaxSize;

  public long getBigThresholdBytes() {
    return bigThresholdBytes;
  }

  public void setBigThresholdBytes(long bigThresholdBytes) {
    this.bigThresholdBytes = bigThresholdBytes;
  }

  public long getChunkMinSize() {
    return chunkMinSize;
  }

  public void setChunkMinSize(long chunkMinSize) {
    this.chunkMinSize = chunkMinSize;
  }

  public long getChunkMaxSize() {
    return chunkMaxSize;
  }

  public void setChunkMaxSize(long chunkMaxSize) {
    this.chunkMaxSize = chunkMaxSize;
  }
}
