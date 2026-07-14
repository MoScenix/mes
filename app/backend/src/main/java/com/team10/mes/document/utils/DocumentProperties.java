package com.team10.mes.document.utils;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "mes.document")
public class DocumentProperties {
  private String root;
  private final Elasticsearch elasticsearch = new Elasticsearch();
  private final Milvus milvus = new Milvus();
  private final Embedding embedding = new Embedding();
  private final Index index = new Index();

  public String getRoot() {
    return root;
  }

  public void setRoot(String root) {
    this.root = root;
  }

  public Elasticsearch getElasticsearch() {
    return elasticsearch;
  }

  public Milvus getMilvus() {
    return milvus;
  }

  public Embedding getEmbedding() {
    return embedding;
  }

  public Index getIndex() {
    return index;
  }

  public static class Elasticsearch {
    private String url, index, username, password;
    private int bulkBatchSize;

    public String getUrl() {
      return url;
    }

    public void setUrl(String v) {
      url = v;
    }

    public String getIndex() {
      return index;
    }

    public void setIndex(String v) {
      index = v;
    }

    public String getUsername() {
      return username;
    }

    public void setUsername(String v) {
      username = v;
    }

    public String getPassword() {
      return password;
    }

    public void setPassword(String v) {
      password = v;
    }

    public int getBulkBatchSize() {
      return bulkBatchSize;
    }

    public void setBulkBatchSize(int v) {
      bulkBatchSize = v;
    }
  }

  public static class Milvus {
    private String url, collection, username, password;
    private int insertBatchSize;

    public String getUrl() {
      return url;
    }

    public void setUrl(String v) {
      url = v;
    }

    public String getCollection() {
      return collection;
    }

    public void setCollection(String v) {
      collection = v;
    }

    public String getUsername() {
      return username;
    }

    public void setUsername(String v) {
      username = v;
    }

    public String getPassword() {
      return password;
    }

    public void setPassword(String v) {
      password = v;
    }

    public int getInsertBatchSize() {
      return insertBatchSize;
    }

    public void setInsertBatchSize(int v) {
      insertBatchSize = v;
    }
  }

  public static class Embedding {
    private String baseUrl, apiKey, model;
    private int dimensions, batchSize;

    public String getBaseUrl() {
      return baseUrl;
    }

    public void setBaseUrl(String v) {
      baseUrl = v;
    }

    public String getApiKey() {
      return apiKey;
    }

    public void setApiKey(String v) {
      apiKey = v;
    }

    public String getModel() {
      return model;
    }

    public void setModel(String v) {
      model = v;
    }

    public int getDimensions() {
      return dimensions;
    }

    public void setDimensions(int v) {
      dimensions = v;
    }

    public int getBatchSize() {
      return batchSize;
    }

    public void setBatchSize(int v) {
      batchSize = v;
    }
  }

  public static class Index {
    private int taskChunkSize;
    private int minWorkers;
    private int maxWorkers;
    private int queueSize;
    private int scaleUpThreshold;
    private int scaleDownThreshold;
    private long idleTimeoutSeconds;

    public int getTaskChunkSize() {
      return taskChunkSize;
    }

    public void setTaskChunkSize(int v) {
      taskChunkSize = v;
    }

    public int getMinWorkers() {
      return minWorkers;
    }

    public void setMinWorkers(int v) {
      minWorkers = v;
    }

    public int getMaxWorkers() {
      return maxWorkers;
    }

    public void setMaxWorkers(int v) {
      maxWorkers = v;
    }

    public int getQueueSize() {
      return queueSize;
    }

    public void setQueueSize(int v) {
      queueSize = v;
    }

    public int getScaleUpThreshold() {
      return scaleUpThreshold;
    }

    public void setScaleUpThreshold(int v) {
      scaleUpThreshold = v;
    }

    public int getScaleDownThreshold() {
      return scaleDownThreshold;
    }

    public void setScaleDownThreshold(int v) {
      scaleDownThreshold = v;
    }

    public long getIdleTimeoutSeconds() {
      return idleTimeoutSeconds;
    }

    public void setIdleTimeoutSeconds(long v) {
      idleTimeoutSeconds = v;
    }
  }
}
