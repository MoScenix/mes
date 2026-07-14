package com.team10.mes.inventory.controller;

import com.team10.mes.inventory.service.InventoryService;
import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import jakarta.servlet.http.HttpSession;
import java.util.*;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/mes")
public class InventoryController {
  private final InventoryService service;
  private final SessionIdentity identity;

  public InventoryController(InventoryService service, SessionIdentity identity) {
    this.service = service;
    this.identity = identity;
  }

  @PostMapping("/item/add")
  public Map<String, Object> addItem(@RequestBody Map<String, Object> r, HttpSession s) {
    user(s);
    return service.addItem(r);
  }

  @PostMapping("/item/update")
  public Map<String, Object> updateItem(@RequestBody Map<String, Object> r, HttpSession s) {
    user(s);
    return service.updateItem(r);
  }

  @GetMapping("/item/get")
  public Map<String, Object> item(@RequestParam long id, HttpSession s) {
    user(s);
    return service.item(id);
  }

  @PostMapping("/item/list")
  public Map<String, Object> items(@RequestBody Map<String, Object> q, HttpSession s) {
    user(s);
    return service.items(q);
  }

  @GetMapping("/item/search")
  public Map<String, Object> search(
      @RequestParam(required = false) String namePrefix, HttpSession s) {
    user(s);
    return service.items(new HashMap<>(Map.of("namePrefix", namePrefix == null ? "" : namePrefix)));
  }

  @PostMapping("/process/draft/create")
  public Map<String, Object> createProcess(@RequestBody Map<String, Object> r, HttpSession s) {
    r.put("ownerUserId", user(s));
    return service.createProcess(r);
  }

  @PostMapping("/process/draft/update")
  public Map<String, Object> updateProcess(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireProcessOwner(num(r), user(s), admin(s), true);
    r.put("ownerUserId", user(s));
    return service.updateProcess(r);
  }

  @PostMapping("/process/draft/delete")
  public Map<String, Object> deleteProcess(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireProcessOwner(num(r), user(s), admin(s), true);
    return service.deleteProcess(num(r));
  }

  @PostMapping("/process/submit")
  public Map<String, Object> submitProcess(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireProcessOwner(num(r), user(s), admin(s), true);
    return service.submitProcess(num(r));
  }

  @GetMapping("/process/get")
  public Map<String, Object> process(@RequestParam long id, HttpSession s) {
    service.requireProcessOwner(id, user(s), admin(s), false);
    return service.process(id);
  }

  @PostMapping("/process/list")
  public Map<String, Object> processes(@RequestBody Map<String, Object> q, HttpSession s) {
    return service.processes(service.scopeProcess(q, user(s), admin(s)));
  }

  @PostMapping("/item-unit/add")
  public Map<String, Object> addUnit(@RequestBody Map<String, Object> r, HttpSession s) {
    user(s);
    return service.addUnit(r);
  }

  @PostMapping("/item-unit/status/update")
  public Map<String, Object> updateUnit(@RequestBody Map<String, Object> r, HttpSession s) {
    user(s);
    return service.updateUnit(r);
  }

  @GetMapping("/item-unit/get")
  public Map<String, Object> unit(@RequestParam long id, HttpSession s) {
    user(s);
    return service.unit(id);
  }

  @PostMapping("/item-unit/list")
  public Map<String, Object> units(@RequestBody Map<String, Object> q, HttpSession s) {
    user(s);
    return service.units(q);
  }

  @PostMapping("/inventory-flow/draft/create")
  public Map<String, Object> createFlow(@RequestBody Map<String, Object> r, HttpSession s) {
    r.put("fromUserId", user(s));
    return service.createFlow(r);
  }

  @PostMapping("/inventory-flow/draft/update")
  public Map<String, Object> updateFlow(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireFlowAccess(num(r), user(s), admin(s), true);
    r.put("fromUserId", user(s));
    return service.updateFlow(r);
  }

  @PostMapping("/inventory-flow/draft/delete")
  public Map<String, Object> deleteFlow(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireFlowAccess(num(r), user(s), admin(s), true);
    return service.deleteFlow(num(r));
  }

  @PostMapping("/inventory-flow/submit")
  public Map<String, Object> submitFlow(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireFlowAccess(num(r), user(s), admin(s), true);
    return service.submitFlow(num(r));
  }

  @PostMapping("/inventory-flow/audit")
  public Map<String, Object> auditFlow(@RequestBody Map<String, Object> r, HttpSession s) {
    long uid = user(s);
    if (!auditor(s)) throw new IllegalStateException("forbidden: no permission");
    r.put("approvedBy", uid);
    return service.auditFlow(r);
  }

  @PostMapping("/inventory-flow/complete")
  public Map<String, Object> completeFlow(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireFlowAccess(num(r), user(s), admin(s), false);
    return service.completeFlow(num(r), listLong(r, "itemUnitIds"));
  }

  @GetMapping("/inventory-flow/get")
  public Map<String, Object> flow(@RequestParam long id, HttpSession s) {
    service.requireFlowAccess(id, user(s), admin(s), false);
    return service.flow(id);
  }

  @PostMapping("/inventory-flow/list")
  public Map<String, Object> flows(@RequestBody Map<String, Object> q, HttpSession s) {
    return service.flows(service.scopeFlow(q, user(s), admin(s)));
  }

  @PostMapping("/engineering-order/draft/create")
  public Map<String, Object> createOrder(@RequestBody Map<String, Object> r, HttpSession s) {
    r.put("leaderUserId", user(s));
    return service.createOrder(r);
  }

  @PostMapping("/engineering-order/draft/update")
  public Map<String, Object> updateOrder(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireOrderLeader(num(r), user(s), admin(s), true);
    r.put("leaderUserId", user(s));
    return service.updateOrder(r);
  }

  @PostMapping("/engineering-order/draft/delete")
  public Map<String, Object> deleteOrder(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireOrderLeader(num(r), user(s), admin(s), true);
    return service.deleteOrder(num(r));
  }

  @PostMapping("/engineering-order/submit")
  public Map<String, Object> submitOrder(@RequestBody Map<String, Object> r, HttpSession s) {
    service.requireOrderLeader(num(r), user(s), admin(s), true);
    return service.submitOrder(num(r));
  }

  @GetMapping("/engineering-order/get")
  public Map<String, Object> order(@RequestParam long id, HttpSession s) {
    service.requireOrderLeader(id, user(s), admin(s), false);
    return service.order(id);
  }

  @PostMapping("/engineering-order/list")
  public Map<String, Object> orders(@RequestBody Map<String, Object> q, HttpSession s) {
    return service.orders(service.scopeOrder(q, user(s), admin(s)));
  }

  private long user(HttpSession s) {
    Long id = identity.userId(s);
    if (id == null || id <= 0) throw new UnauthorizedException();
    return id;
  }

  private boolean admin(HttpSession s) {
    String r = role(s);
    return r.equals("admin") || r.equals("administrator") || r.equals("管理员");
  }

  private boolean auditor(HttpSession s) {
    String r = role(s);
    return admin(s) || r.equals("warehouse") || r.equals("warehouse_admin") || r.equals("仓库管理员");
  }

  private String role(HttpSession s) {
    return Objects.toString(identity.role(s), "").trim().toLowerCase(Locale.ROOT);
  }

  private static long num(Map<String, Object> r) {
    return Long.parseLong(r.get("id").toString());
  }

  @SuppressWarnings("unchecked")
  private static List<Long> listLong(Map<String, Object> r, String k) {
    return ((List<Object>) r.getOrDefault(k, List.of()))
        .stream().map(x -> Long.parseLong(x.toString())).toList();
  }
}
