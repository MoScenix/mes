package com.team10.mes.app.controller;

import com.team10.mes.app.dal.AppMapper;
import com.team10.mes.app.service.AppFileService;
import com.team10.mes.app.service.AppService;
import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import jakarta.servlet.http.HttpSession;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

@RestController
@RequestMapping("/app")
public class AppController {
  private final AppService service;
  private final AppFileService files;
  private final SessionIdentity identity;

  public AppController(AppService service, AppFileService files, SessionIdentity identity) {
    this.service = service;
    this.files = files;
    this.identity = identity;
  }

  @PostMapping("/add")
  public long create(@RequestBody CreateRequest request, HttpSession session) {
    return service.create(userId(session), request.initPrompt()).id();
  }

  @GetMapping("/get/vo")
  public AppMapper.AppRow get(@RequestParam long id, HttpSession session) {
    return owned(id, session);
  }

  @PostMapping("/my/list/page/vo")
  public AppService.Page listMine(@RequestBody ListRequest request, HttpSession session) {
    return service.list(userId(session), request.appName(), request.pageNum(), request.pageSize());
  }

  @PostMapping("/good/list/page/vo")
  public AppService.Page listGood(@RequestBody ListRequest request) {
    return service.list(request.userId(), request.appName(), request.pageNum(), request.pageSize());
  }

  @PostMapping("/admin/list/page/vo")
  public AppService.Page listAll(@RequestBody ListRequest request, HttpSession session) {
    requireAdmin(session);
    return service.list(request.userId(), request.appName(), request.pageNum(), request.pageSize());
  }

  @PostMapping("/update")
  public Result update(@RequestBody UpdateRequest request, HttpSession session) {
    owned(request.id(), session);
    return new Result(service.rename(request.id(), request.appName()));
  }

  @PostMapping("/delete")
  public Result delete(@RequestBody IdRequest request, HttpSession session)
      throws java.io.IOException {
    owned(request.id(), session);
    boolean deleted = service.delete(request.id());
    if (deleted) files.deleteProjectFiles(request.id());
    return new Result(deleted);
  }

  @GetMapping("/admin/get/vo")
  public AppMapper.AppRow adminGet(@RequestParam long id, HttpSession session) {
    requireAdmin(session);
    return required(service.get(id));
  }

  @PostMapping("/admin/update")
  public Result adminUpdate(@RequestBody UpdateRequest request, HttpSession session) {
    requireAdmin(session);
    return new Result(service.rename(request.id(), request.appName()));
  }

  @PostMapping("/admin/delete")
  public Result adminDelete(@RequestBody IdRequest request, HttpSession session)
      throws java.io.IOException {
    requireAdmin(session);
    boolean deleted = service.delete(request.id());
    if (deleted) files.deleteProjectFiles(request.id());
    return new Result(deleted);
  }

  private static <T> T required(T value) {
    if (value == null) throw new ResponseStatusException(HttpStatus.NOT_FOUND);
    return value;
  }

  private long userId(HttpSession session) {
    Long id = identity.userId(session);
    if (id == null) throw new UnauthorizedException();
    return id;
  }

  private void requireAdmin(HttpSession session) {
    if (!"admin".equalsIgnoreCase(identity.role(session)))
      throw new IllegalStateException("admin required");
  }

  private AppMapper.AppRow owned(long appId, HttpSession session) {
    AppMapper.AppRow app = required(service.get(appId));
    if (!app.userId().equals(userId(session)) && !"admin".equalsIgnoreCase(identity.role(session)))
      throw new IllegalStateException("forbidden");
    return app;
  }

  public record CreateRequest(String initPrompt) {}

  public record ListRequest(Long userId, String appName, int pageNum, int pageSize) {}

  public record UpdateRequest(long id, String appName) {}

  public record IdRequest(long id) {}

  public record Result(boolean success) {}
}
