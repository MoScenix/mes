package com.team10.mes.ai.state;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.team10.mes.ai.service.AiService.AiEvent;
import com.team10.mes.ai.service.AiService.AiState;
import java.time.Duration;
import java.util.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.connection.stream.*;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Component;

@Component
public class RedisAiStore {
  private final StringRedisTemplate redis;
  private final ObjectMapper json;
  private final Duration terminalTtl;
  private final long maxBlockMs;
  private final int maxReadCount;

  public RedisAiStore(
      StringRedisTemplate redis,
      ObjectMapper json,
      @Value("${mes.ai.redis.terminal-ttl-seconds:10}") long ttl,
      @Value("${mes.ai.redis.max-block-ms:30000}") long maxBlockMs,
      @Value("${mes.ai.redis.max-read-count:100}") int maxReadCount) {
    this.redis = redis;
    this.json = json;
    this.terminalTtl = Duration.ofSeconds(Math.max(1, ttl));
    this.maxBlockMs = Math.max(0, maxBlockMs);
    this.maxReadCount = Math.max(1, maxReadCount);
  }

  public AiState state(long appId) {
    String raw = redis.opsForValue().get(key(appId, "state"));
    if (raw == null) return null;
    try {
      return json.readValue(raw, AiState.class);
    } catch (Exception e) {
      throw new IllegalStateException("invalid AI state", e);
    }
  }

  public void saveState(long appId, AiState state) {
    try {
      redis.opsForValue().set(key(appId, "state"), json.writeValueAsString(state));
    } catch (Exception e) {
      throw new IllegalStateException(e);
    }
  }

  public void saveCheckpoint(long appId, Object checkpoint) {
    try {
      redis.opsForValue().set(key(appId, "checkpoint"), json.writeValueAsString(checkpoint));
    } catch (Exception e) {
      throw new IllegalStateException(e);
    }
  }

  public <T> T checkpoint(long appId, Class<T> type) {
    String raw = redis.opsForValue().get(key(appId, "checkpoint"));
    if (raw == null) return null;
    try {
      return json.readValue(raw, type);
    } catch (Exception e) {
      throw new IllegalStateException("invalid AI checkpoint", e);
    }
  }

  public void deleteCheckpoint(long appId) {
    redis.delete(key(appId, "checkpoint"));
  }

  public String addEvent(long appId, AiEvent event) {
    return add(key(appId, "stream"), event);
  }

  public String addControl(long appId, AiEvent event) {
    return add(key(appId, "control"), event);
  }

  private String add(String key, AiEvent event) {
    try {
      Map<String, String> fields = new LinkedHashMap<>();
      fields.put("data", json.writeValueAsString(event));
      RecordId id = redis.opsForStream().add(StreamRecords.string(fields).withStreamKey(key));
      return id == null ? "" : id.getValue();
    } catch (Exception e) {
      throw new IllegalStateException(e);
    }
  }

  public List<AiEvent> events(long appId, String lastId, long blockMs, int count) {
    return read(key(appId, "stream"), lastId, blockMs, count);
  }

  public List<AiEvent> controls(long appId, String lastId, long blockMs, int count) {
    return read(key(appId, "control"), lastId, blockMs, count);
  }

  private List<AiEvent> read(String key, String lastId, long blockMs, int count) {
    StreamReadOptions options =
        StreamReadOptions.empty().count(Math.max(1, Math.min(maxReadCount, count)));
    long wait = Math.max(0, Math.min(maxBlockMs, blockMs));
    if (wait > 0) options = options.block(Duration.ofMillis(wait));
    ReadOffset offset =
        "$".equals(lastId) ? ReadOffset.latest() : ReadOffset.from(normalize(lastId));
    List<MapRecord<String, Object, Object>> records =
        redis.opsForStream().read(options, StreamOffset.create(key, offset));
    if (records == null) return List.of();
    List<AiEvent> out = new ArrayList<>();
    for (var record : records) {
      Object raw = record.getValue().get("data");
      if (raw == null) continue;
      try {
        AiEvent e = json.readValue(raw.toString(), AiEvent.class);
        out.add(
            new AiEvent(
                record.getId().getValue(),
                e.projectId(),
                e.type(),
                e.agent(),
                e.content(),
                e.targetId(),
                e.name(),
                e.status(),
                e.payloadJson(),
                e.createdAt(),
                e.questions()));
      } catch (Exception ignored) {
      }
    }
    return out;
  }

  public void resetEvents(long appId) {
    redis.delete(key(appId, "stream"));
  }

  public void expireTerminal(long appId) {
    for (String suffix : List.of("state", "stream", "control", "checkpoint"))
      redis.expire(key(appId, suffix), terminalTtl);
  }

  private static String normalize(String id) {
    return id == null || id.isBlank() || "0".equals(id) ? "0-0" : id;
  }

  private static String key(long appId, String suffix) {
    return "project:" + appId + ":" + suffix;
  }
}
