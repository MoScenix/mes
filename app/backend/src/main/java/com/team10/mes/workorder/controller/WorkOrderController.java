package com.team10.mes.workorder.controller;

import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import com.team10.mes.workorder.service.WorkOrderService;
import com.team10.mes.workorder.service.WorkOrderService.*;
import jakarta.servlet.http.HttpSession;
import java.util.Locale;
import java.util.Objects;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/mes/work-order")
public class WorkOrderController {
  private final WorkOrderService service;
  private final SessionIdentity identity;

  public WorkOrderController(WorkOrderService service, SessionIdentity identity) {
    this.service = service;
    this.identity = identity;
  }

  public record IdRequest(Long id) {}

  public record Response<T>(int code, String message, T data) {
    static <T> Response<T> ok(T data) {
      return new Response<>(0, "success", data);
    }
  }

  @PostMapping("/draft/create")
  public Response<Long> create(@RequestBody DraftRequest req, HttpSession session) {
    return Response.ok(service.create(user(session), req));
  }

  @PostMapping("/draft/update")
  public Response<Boolean> update(@RequestBody UpdateRequest req, HttpSession session) {
    service.update(user(session), admin(session), req);
    return Response.ok(true);
  }

  @PostMapping("/draft/delete")
  public Response<Boolean> delete(@RequestBody IdRequest req, HttpSession session) {
    service.delete(req.id(), user(session), admin(session));
    return Response.ok(true);
  }

  @PostMapping("/submit")
  public Response<Boolean> submit(@RequestBody IdRequest req, HttpSession session) {
    service.submit(req.id(), user(session), admin(session));
    return Response.ok(true);
  }

  @PostMapping("/read")
  public Response<Boolean> read(@RequestBody IdRequest req, HttpSession session) {
    service.read(req.id(), user(session), admin(session));
    return Response.ok(true);
  }

  @GetMapping("/get")
  public Response<View> get(@RequestParam long id, HttpSession session) {
    return Response.ok(service.get(id, user(session), admin(session)));
  }

  @PostMapping("/list")
  public Response<Page> list(@RequestBody ListRequest req, HttpSession session) {
    return Response.ok(service.list(req, user(session), admin(session)));
  }

  @ExceptionHandler(IllegalArgumentException.class)
  public Response<Void> badRequest(IllegalArgumentException e) {
    return new Response<>(1, e.getMessage(), null);
  }

  private long user(HttpSession session) {
    Long id = identity.userId(session);
    if (id == null || id <= 0) throw new UnauthorizedException();
    return id;
  }

  private boolean admin(HttpSession session) {
    String role = Objects.toString(identity.role(session), "").trim().toLowerCase(Locale.ROOT);
    return role.equals("admin") || role.equals("administrator") || role.equals("管理员");
  }
}
