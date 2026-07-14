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

  public static class Elasticsearch {
    private String url, index, username, password;

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
  }

  public static class Milvus {
    private String url, collection, username, password;

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
  }

  public static class Embedding {
    private String baseUrl, apiKey, model;
    private int dimensions;

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
  }
}
