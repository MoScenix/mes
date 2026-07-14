package com.team10.mes.user.dal;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.time.Duration;
import java.util.Map;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Component;

@Component
public class UserCache {
  private static final String PREFIX = "team10:user:";

  private final StringRedisTemplate redis;
  private final ObjectMapper json;
  private final Duration ttl;

  public UserCache(
      StringRedisTemplate redis,
      ObjectMapper json,
      @Value("${mes.user.cache-ttl:1h}") Duration ttl) {
    this.redis = redis;
    this.json = json;
    this.ttl = ttl;
  }

  public Map<String, Object> get(long id) {
    try {
      String value = redis.opsForValue().get(key(id));
      return value == null || value.isBlank()
          ? null
          : json.readValue(value, new TypeReference<>() {});
    } catch (Exception ignored) {
      return null;
    }
  }

  public void put(long id, Map<String, Object> user) {
    try {
      redis.opsForValue().set(key(id), json.writeValueAsString(user), ttl);
    } catch (Exception ignored) {
      // User reads remain available from MySQL when Redis is temporarily unavailable.
    }
  }

  public void evict(long id) {
    try {
      redis.delete(key(id));
    } catch (RuntimeException ignored) {
      // A failed eviction must not roll back an already completed SQL delete.
    }
  }

  private String key(long id) {
    return PREFIX + id;
  }
}
