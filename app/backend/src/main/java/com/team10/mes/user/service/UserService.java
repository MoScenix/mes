package com.team10.mes.user.service;

import com.team10.mes.user.dal.UserCache;
import com.team10.mes.user.dal.UserMapper;
import java.util.*;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class UserService {
  private final UserMapper mapper;
  private final UserCache cache;
  private final BCryptPasswordEncoder encoder = new BCryptPasswordEncoder();

  public UserService(UserMapper mapper, UserCache cache) {
    this.mapper = mapper;
    this.cache = cache;
  }

  public Map<String, Object> register(Map<String, Object> r) {
    String a = text(r, "userAccount"), p = text(r, "userPassword");
    if (a.isBlank() || p.isBlank())
      throw new IllegalArgumentException("account and password are required");
    if (!p.equals(text(r, "checkPassword")))
      throw new IllegalArgumentException("password confirmation does not match");
    if (mapper.countByAccount(a) > 0)
      throw new IllegalStateException("user account already exists");
    Map<String, Object> u = new HashMap<>();
    u.put("name", UUID.randomUUID().toString());
    u.put("passwordHash", encoder.encode(p));
    u.put("account", a);
    u.put("avatar", "");
    u.put("profile", "");
    u.put("role", "worker");
    mapper.insert(u);
    cache.put(
        ((Number) u.get("id")).longValue(),
        view(mapper.findById(((Number) u.get("id")).longValue())));
    return Map.of("userId", u.get("id"), "userRole", "worker");
  }

  public Map<String, Object> login(Map<String, Object> r) {
    Map<String, Object> u = mapper.findByAccount(text(r, "userAccount"));
    if (u == null
        || !encoder.matches(text(r, "userPassword"), String.valueOf(u.get("passwordHash"))))
      throw new IllegalArgumentException("invalid credentials");
    return view(u);
  }

  public Map<String, Object> get(Long id) {
    Map<String, Object> cached = cache.get(id);
    if (cached != null) return cached;
    Map<String, Object> u = mapper.findById(id);
    if (u == null) throw new NoSuchElementException("user not found");
    Map<String, Object> result = view(u);
    cache.put(id, result);
    return result;
  }

  public void add(Map<String, Object> r) {
    String password = text(r, "userPassword");
    register(
        Map.of(
            "userAccount",
            r.get("userAccount"),
            "userPassword",
            password,
            "checkPassword",
            password));
  }

  public void update(Map<String, Object> r) {
    long id = ((Number) r.get("id")).longValue();
    if (mapper.update(r) == 0) throw new NoSuchElementException("user not found");
    Map<String, Object> updated = mapper.findById(id);
    if (updated == null) cache.evict(id);
    else cache.put(id, view(updated));
  }

  public void delete(Long id) {
    if (mapper.softDelete(id) == 0) throw new NoSuchElementException("user not found");
    cache.evict(id);
  }

  public Map<String, Object> page(Map<String, Object> r) {
    int n = ((Number) r.getOrDefault("pageNum", 1)).intValue(),
        s = ((Number) r.getOrDefault("pageSize", 10)).intValue();
    String name = text(r, "userName"), a = text(r, "account");
    return Map.of(
        "records",
        mapper.page(name, a, Math.max(0, n - 1) * s, s),
        "total",
        mapper.pageCount(name, a),
        "current",
        n,
        "size",
        s);
  }

  private String text(Map<String, Object> m, String k) {
    return Objects.toString(m.get(k), "");
  }

  private Map<String, Object> view(Map<String, Object> u) {
    Map<String, Object> v = new HashMap<>(u);
    v.remove("passwordHash");
    v.remove("password_hash");
    v.put("userName", v.remove("name"));
    if (!v.containsKey("userAccount")) v.put("userAccount", v.remove("user_account"));
    if (!v.containsKey("userRole")) v.put("userRole", v.remove("user_role"));
    return v;
  }
}
