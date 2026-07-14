package com.team10.mes.user.controller;

import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import com.team10.mes.user.service.UserService;
import jakarta.servlet.http.HttpSession;
import java.util.Map;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/user")
public class UserController {
  private final UserService service;
  private final SessionIdentity identity;

  public UserController(UserService service, SessionIdentity identity) {
    this.service = service;
    this.identity = identity;
  }

  @PostMapping("/register")
  public Map<String, Object> register(@RequestBody Map<String, Object> r) {
    return service.register(r);
  }

  @PostMapping("/login")
  public Map<String, Object> login(@RequestBody Map<String, Object> r, HttpSession session) {
    Map<String, Object> user = service.login(r);
    identity.login(session, user);
    return user;
  }

  @PostMapping("/add")
  public Map<String, Object> add(@RequestBody Map<String, Object> r) {
    service.add(r);
    return Map.of("success", true);
  }

  @PostMapping("/update")
  public Map<String, Object> update(@RequestBody Map<String, Object> r) {
    service.update(r);
    return Map.of("success", true);
  }

  @PostMapping("/delete")
  public Map<String, Object> delete(@RequestBody Map<String, Object> r) {
    service.delete(((Number) r.get("id")).longValue());
    return Map.of("success", true);
  }

  @GetMapping("/get")
  public Map<String, Object> get(@RequestParam Long id) {
    return service.get(id);
  }

  @GetMapping("/get/vo")
  public Map<String, Object> getVo(@RequestParam Long id) {
    return service.get(id);
  }

  @GetMapping("/get/login")
  public Map<String, Object> getLogin(HttpSession session) {
    Map<String, Object> user = identity.user(session);
    if (user == null) throw new UnauthorizedException();
    return user;
  }

  @PostMapping("/logout")
  public Map<String, Object> logout(HttpSession session) {
    session.invalidate();
    return Map.of("success", true);
  }

  @PostMapping("/list/page/vo")
  public Map<String, Object> page(@RequestBody Map<String, Object> r) {
    return service.page(r);
  }
}
