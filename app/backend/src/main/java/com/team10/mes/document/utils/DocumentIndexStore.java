package com.team10.mes.document.utils;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.net.URI;
import java.net.http.*;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicReference;
import org.springframework.stereotype.Component;

@Component
public class DocumentIndexStore {
  public record Child(
      long historyId, long fileId, long chunkId, List<Long> parentIds, String content) {}

  private final HttpClient http =
      HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(10)).build();
  private final ObjectMapper json;
  private final DocumentWorkPool workPool;
  private final String es, index, milvus, collection, embeddingUrl, embeddingKey, embeddingModel;
  private final int dimension, embeddingBatchSize, esBulkBatchSize, milvusInsertBatchSize;
  private final int taskChunkSize;
  private final Map<String, String> esHeaders, milvusHeaders;

  public DocumentIndexStore(ObjectMapper json, DocumentProperties p, DocumentWorkPool workPool) {
    this.json = json;
    this.workPool = workPool;
    this.es = trim(p.getElasticsearch().getUrl());
    this.index = p.getElasticsearch().getIndex();
    this.milvus = trim(p.getMilvus().getUrl());
    this.collection = p.getMilvus().getCollection();
    this.embeddingUrl = trim(p.getEmbedding().getBaseUrl());
    this.embeddingKey = p.getEmbedding().getApiKey();
    this.embeddingModel = p.getEmbedding().getModel();
    this.dimension = p.getEmbedding().getDimensions();
    this.embeddingBatchSize = positive(p.getEmbedding().getBatchSize(), 25);
    this.esBulkBatchSize = positive(p.getElasticsearch().getBulkBatchSize(), 200);
    this.milvusInsertBatchSize = positive(p.getMilvus().getInsertBatchSize(), 200);
    this.taskChunkSize = positive(p.getIndex().getTaskChunkSize(), 200);
    this.esHeaders = basic(p.getElasticsearch().getUsername(), p.getElasticsearch().getPassword());
    this.milvusHeaders = token(p.getMilvus().getUsername(), p.getMilvus().getPassword());
  }

  public void index(List<Child> children) {
    if (children.isEmpty()) return;
    ensureEs();
    ensureMilvus();
    AtomicBoolean cancelled = new AtomicBoolean(false);
    AtomicReference<Throwable> firstError = new AtomicReference<>();
    List<CompletableFuture<Void>> tasks = new ArrayList<>();
    for (int start = 0; start < children.size(); start += taskChunkSize) {
      int batchNumber = tasks.size() + 1;
      List<Child> batch =
          List.copyOf(children.subList(start, Math.min(children.size(), start + taskChunkSize)));
      CompletableFuture<Void> task =
          workPool
              .submit(
                  () -> {
                    if (cancelled.get()) return;
                    try {
                      indexChildBatch(batch, cancelled);
                    } catch (Throwable error) {
                      cancelled.set(true);
                      throw new IllegalStateException(
                          "index child batch " + batchNumber + " failed: " + error.getMessage(),
                          error);
                    }
                  })
              .whenComplete(
                  (ignored, error) -> {
                    if (error != null) {
                      cancelled.set(true);
                      firstError.compareAndSet(null, unwrap(error));
                    }
                  });
      tasks.add(task);
    }
    waitAll(tasks, firstError, "document index");
  }

  private void indexChildBatch(List<Child> children, AtomicBoolean cancelled) {
    List<List<Double>> vectors = new ArrayList<>(children.size());
    for (int start = 0; start < children.size(); start += embeddingBatchSize) {
      if (cancelled.get()) return;
      List<Child> batch =
          children.subList(start, Math.min(children.size(), start + embeddingBatchSize));
      vectors.addAll(embed(batch.stream().map(Child::content).toList()));
    }
    if (vectors.size() != children.size())
      throw new IllegalStateException(
          "embedding response size mismatch: expected "
              + children.size()
              + ", got "
              + vectors.size());
    int writeBatchSize = Math.max(1, Math.min(esBulkBatchSize, milvusInsertBatchSize));
    for (int start = 0; start < children.size(); start += writeBatchSize) {
      if (cancelled.get()) return;
      int end = Math.min(children.size(), start + writeBatchSize);
      indexBatch(children.subList(start, end), vectors.subList(start, end));
    }
  }

  private void indexBatch(List<Child> children, List<List<Double>> vectors) {
    if (vectors.size() != children.size())
      throw new IllegalStateException(
          "embedding response size mismatch: expected "
              + children.size()
              + ", got "
              + vectors.size());
    StringBuilder bulk = new StringBuilder();
    List<Map<String, Object>> rows = new ArrayList<>();
    for (int i = 0; i < children.size(); i++) {
      Child c = children.get(i);
      String id = id(c);
      bulk.append(write(Map.of("index", Map.of("_index", index, "_id", id))))
          .append('\n')
          .append(write(doc(c)))
          .append('\n');
      rows.add(
          new LinkedHashMap<>(
              Map.of(
                  "id",
                  id,
                  "content",
                  c.content(),
                  "history_id",
                  c.historyId(),
                  "file_id",
                  c.fileId(),
                  "chunk_id",
                  c.chunkId(),
                  "parent_ids",
                  c.parentIds(),
                  "vector",
                  vectors.get(i))));
    }
    Map<String, Object> bulkResult =
        read(
            request(
                es + "/_bulk?refresh=true",
                "POST",
                bulk.toString(),
                "application/x-ndjson",
                esHeaders),
            "elasticsearch bulk");
    if (Boolean.TRUE.equals(bulkResult.get("errors")))
      throw new IllegalStateException(
          "bulk index es child chunks failed: "
              + truncate(String.valueOf(bulkResult.get("items"))));
    milvus("/v2/vectordb/entities/insert", Map.of("collectionName", collection, "data", rows));
  }

  public List<Long> search(long history, long file, String query, long topK) {
    int k = normalizeTopK(topK);
    ensureEs();
    ensureMilvus();
    AtomicReference<Throwable> firstError = new AtomicReference<>();
    CompletableFuture<List<List<Long>>> esTask =
        supplyAsync(() -> searchEsParents(history, file, query, k), firstError);
    CompletableFuture<List<List<Long>>> milvusTask =
        supplyAsync(() -> searchMilvusParents(history, file, query, k), firstError);
    waitAll(
        List.of(esTask.thenAccept(ignored -> {}), milvusTask.thenAccept(ignored -> {})),
        firstError,
        "document search");
    return DocumentText.fuse(List.of(esTask.join(), milvusTask.join()), k);
  }

  private List<List<Long>> searchEsParents(long history, long file, String query, int k) {
    Map<String, Object> esBody =
        Map.of(
            "size",
            k,
            "query",
            Map.of(
                "bool",
                Map.of(
                    "filter",
                    List.of(
                        Map.of("term", Map.of("historyId", history)),
                        Map.of("term", Map.of("fileId", file))),
                    "must",
                    List.of(Map.of("match", Map.of("content", query))))));
    Map<String, Object> esResp =
        read(
            request(
                es + "/" + index + "/_search", "POST", write(esBody), "application/json", Map.of()),
            "elasticsearch search");
    List<Map<String, Object>> hits = maps(map(esResp, "hits").get("hits"));
    return hits.stream().map(h -> longs(map(h, "_source").get("parentIds"))).toList();
  }

  private List<List<Long>> searchMilvusParents(long history, long file, String query, int k) {
    List<Double> vector = embed(List.of(query)).getFirst();
    Map<String, Object> mv =
        milvus(
            "/v2/vectordb/entities/search",
            Map.of(
                "collectionName",
                collection,
                "data",
                List.of(vector),
                "annsField",
                "vector",
                "filter",
                "history_id == " + history + " && file_id == " + file,
                "limit",
                k,
                "outputFields",
                List.of("parent_ids"),
                "consistencyLevel",
                "Strong"));
    return maps(mv.get("data")).stream().map(x -> longs(x.get("parent_ids"))).toList();
  }

  public void deleteHistory(long history) {
    ensureEs();
    ensureMilvus();
    request(
        es + "/" + index + "/_delete_by_query?conflicts=proceed&refresh=true",
        "POST",
        write(Map.of("query", Map.of("term", Map.of("historyId", history)))),
        "application/json",
        Map.of());
    milvus(
        "/v2/vectordb/entities/delete",
        Map.of("collectionName", collection, "filter", "history_id == " + history));
  }

  private void ensureEs() {
    HttpResponse<String> exists =
        request(es + "/" + index, "HEAD", null, "application/json", Map.of(), true);
    if (exists.statusCode() == 200) return;
    if (exists.statusCode() != 404)
      throw new IllegalStateException(
          "check es index failed: " + exists.statusCode() + " " + exists.body());
    request(
        es + "/" + index,
        "PUT",
        write(
            Map.of(
                "mappings",
                Map.of(
                    "properties",
                    Map.of(
                        "historyId",
                        Map.of("type", "long"),
                        "fileId",
                        Map.of("type", "long"),
                        "chunkId",
                        Map.of("type", "long"),
                        "parentIds",
                        Map.of("type", "long"),
                        "content",
                        Map.of("type", "text"))))),
        "application/json",
        Map.of());
  }

  private void ensureMilvus() {
    Map<String, Object> has =
        milvus("/v2/vectordb/collections/has", Map.of("collectionName", collection));
    if (Boolean.TRUE.equals(map(has, "data").get("has"))) return;
    List<Map<String, Object>> fields =
        List.of(
            field("id", "VarChar", true, Map.of("max_length", "128")),
            field("content", "VarChar", false, Map.of("max_length", "65535")),
            field("history_id", "Int64", false, Map.of()),
            field("file_id", "Int64", false, Map.of()),
            field("chunk_id", "Int64", false, Map.of()),
            arrayField("parent_ids", "Int64", Map.of("max_capacity", "4")),
            field("vector", "FloatVector", false, Map.of("dim", Integer.toString(dimension))));
    List<Map<String, Object>> indexes =
        List.of(
            Map.of(
                "fieldName",
                "vector",
                "indexName",
                "vector",
                "metricType",
                "COSINE",
                "params",
                Map.of("index_type", "HNSW", "M", 16, "efConstruction", 200)),
            Map.of(
                "fieldName",
                "history_id",
                "indexName",
                "history_id",
                "params",
                Map.of("index_type", "STL_SORT")),
            Map.of(
                "fieldName",
                "file_id",
                "indexName",
                "file_id",
                "params",
                Map.of("index_type", "STL_SORT")));
    milvus(
        "/v2/vectordb/collections/create",
        Map.of(
            "collectionName",
            collection,
            "schema",
            Map.of("autoId", false, "enableDynamicField", false, "fields", fields),
            "indexParams",
            indexes));
    milvus("/v2/vectordb/collections/load", Map.of("collectionName", collection));
  }

  private Map<String, Object> field(
      String name, String type, boolean primary, Map<String, Object> params) {
    Map<String, Object> m = new LinkedHashMap<>();
    m.put("fieldName", name);
    m.put("dataType", type);
    m.put("isPrimary", primary);
    m.put("elementTypeParams", params);
    return m;
  }

  private Map<String, Object> arrayField(
      String name, String elementType, Map<String, Object> params) {
    Map<String, Object> field = field(name, "Array", false, params);
    field.put("elementDataType", elementType);
    return field;
  }

  private List<List<Double>> embed(List<String> texts) {
    Map<String, String> headers =
        embeddingKey.isBlank() ? Map.of() : Map.of("Authorization", "Bearer " + embeddingKey);
    Map<String, Object> response =
        read(
            request(
                embeddingUrl + "/embeddings",
                "POST",
                write(Map.of("model", embeddingModel, "input", texts)),
                "application/json",
                headers),
            "embedding");
    if (response.containsKey("error"))
      throw new IllegalStateException(
          "embedding error: " + truncate(String.valueOf(response.get("error"))));
    List<Map<String, Object>> data = maps(response.get("data"));
    if (data.isEmpty() && !texts.isEmpty())
      throw new IllegalStateException(
          "embedding response missing data: " + truncate(write(response)));
    return data.stream()
        .sorted(Comparator.comparingInt(x -> Integer.parseInt(x.get("index").toString())))
        .map(
            x ->
                list(x.get("embedding")).stream()
                    .map(v -> Double.parseDouble(v.toString()))
                    .toList())
        .toList();
  }

  private Map<String, Object> milvus(String path, Object body) {
    Map<String, Object> r =
        read(
            request(milvus + path, "POST", write(body), "application/json", milvusHeaders),
            "milvus " + path);
    Object code = r.get("code");
    if (code != null && Integer.parseInt(code.toString()) != 0)
      throw new IllegalStateException("milvus " + path + " error: " + truncate(write(r)));
    return r;
  }

  private HttpResponse<String> request(
      String uri, String method, String body, String contentType, Map<String, String> headers) {
    return request(uri, method, body, contentType, headers, false);
  }

  private HttpResponse<String> request(
      String uri,
      String method,
      String body,
      String contentType,
      Map<String, String> headers,
      boolean allowError) {
    try {
      HttpRequest.Builder b =
          HttpRequest.newBuilder(URI.create(uri))
              .timeout(Duration.ofSeconds(60))
              .header("Content-Type", contentType);
      headers.forEach(b::header);
      b.method(
          method,
          body == null
              ? HttpRequest.BodyPublishers.noBody()
              : HttpRequest.BodyPublishers.ofString(body, StandardCharsets.UTF_8));
      HttpResponse<String> r = http.send(b.build(), HttpResponse.BodyHandlers.ofString());
      if (!allowError && r.statusCode() >= 300)
        throw new IllegalStateException(
            method
                + " "
                + sanitize(uri)
                + " returned "
                + r.statusCode()
                + ": "
                + truncate(r.body()));
      return r;
    } catch (IllegalStateException e) {
      throw e;
    } catch (Exception e) {
      throw new IllegalStateException(
          "document index request failed: "
              + method
              + " "
              + sanitize(uri)
              + ": "
              + e.getClass().getSimpleName()
              + ": "
              + Objects.toString(e.getMessage(), ""),
          e);
    }
  }

  private String write(Object value) {
    try {
      return json.writeValueAsString(value);
    } catch (Exception e) {
      throw new IllegalStateException(e);
    }
  }

  private Map<String, Object> read(HttpResponse<String> response) {
    return read(response, "document index");
  }

  private Map<String, Object> read(HttpResponse<String> response, String context) {
    try {
      return json.readValue(response.body(), new TypeReference<>() {});
    } catch (Exception e) {
      throw new IllegalStateException(
          context
              + " returned invalid JSON: status="
              + response.statusCode()
              + ", body="
              + truncate(response.body()),
          e);
    }
  }

  @SuppressWarnings("unchecked")
  private static Map<String, Object> map(Map<String, Object> m, String key) {
    Object v = m.get(key);
    return v instanceof Map<?, ?> ? (Map<String, Object>) v : Map.of();
  }

  @SuppressWarnings("unchecked")
  private static List<Map<String, Object>> maps(Object v) {
    return v instanceof List<?> l ? (List<Map<String, Object>>) (List<?>) l : List.of();
  }

  private static List<?> list(Object v) {
    return v instanceof List<?> l ? l : List.of();
  }

  private static List<Long> longs(Object v) {
    return list(v).stream().map(x -> Long.parseLong(x.toString())).filter(x -> x > 0).toList();
  }

  private Map<String, Object> doc(Child c) {
    return Map.of(
        "historyId",
        c.historyId(),
        "fileId",
        c.fileId(),
        "chunkId",
        c.chunkId(),
        "parentIds",
        c.parentIds(),
        "content",
        c.content());
  }

  private <T> CompletableFuture<T> supplyAsync(
      java.util.function.Supplier<T> supplier, AtomicReference<Throwable> firstError) {
    CompletableFuture<T> result = new CompletableFuture<>();
    workPool
        .submit(
            () -> {
              try {
                result.complete(supplier.get());
              } catch (Throwable error) {
                firstError.compareAndSet(null, error);
                result.completeExceptionally(error);
              }
            })
        .whenComplete(
            (ignored, error) -> {
              if (error != null) {
                Throwable cause = unwrap(error);
                firstError.compareAndSet(null, cause);
                result.completeExceptionally(cause);
              }
            });
    return result;
  }

  private static void waitAll(
      List<CompletableFuture<Void>> tasks, AtomicReference<Throwable> firstError, String context) {
    try {
      CompletableFuture.allOf(tasks.toArray(CompletableFuture[]::new)).join();
    } catch (CompletionException error) {
      Throwable cause = firstError.get() == null ? unwrap(error) : firstError.get();
      if (cause instanceof RuntimeException runtime) throw runtime;
      throw new IllegalStateException(context + " failed: " + cause.getMessage(), cause);
    }
    Throwable error = firstError.get();
    if (error != null) {
      if (error instanceof RuntimeException runtime) throw runtime;
      throw new IllegalStateException(context + " failed: " + error.getMessage(), error);
    }
  }

  private static String id(Child c) {
    return c.historyId() + ":" + c.fileId() + ":" + c.chunkId();
  }

  private static Throwable unwrap(Throwable error) {
    if (error instanceof CompletionException && error.getCause() != null) return error.getCause();
    return error;
  }

  private static int positive(int value, int fallback) {
    return value > 0 ? value : fallback;
  }

  private static int normalizeTopK(long topK) {
    if (topK <= 0) return 3;
    return (int) Math.max(3, Math.min(5, topK));
  }

  private static String truncate(String value) {
    if (value == null) return "";
    if (value.length() <= 2000) return value;
    return value.substring(0, 2000) + "... (+" + (value.length() - 2000) + " chars)";
  }

  private static String sanitize(String uri) {
    try {
      URI parsed = URI.create(uri);
      return new URI(
              parsed.getScheme(),
              null,
              parsed.getHost(),
              parsed.getPort(),
              parsed.getPath(),
              parsed.getQuery(),
              null)
          .toString();
    } catch (Exception ignored) {
      return uri;
    }
  }

  private static String trim(String s) {
    return s.replaceAll("/+$", "");
  }

  private static Map<String, String> basic(String user, String password) {
    if (user == null || user.isBlank()) return Map.of();
    String value =
        Base64.getEncoder()
            .encodeToString(
                (user + ":" + Objects.toString(password, "")).getBytes(StandardCharsets.UTF_8));
    return Map.of("Authorization", "Basic " + value);
  }

  private static Map<String, String> token(String user, String password) {
    return user == null || user.isBlank()
        ? Map.of()
        : Map.of("Authorization", "Bearer " + user + ":" + Objects.toString(password, ""));
  }
}
