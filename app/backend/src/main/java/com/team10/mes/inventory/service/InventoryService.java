package com.team10.mes.inventory.service;

import com.team10.mes.inventory.dal.InventoryMapper;
import java.util.*;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class InventoryService {
  private final InventoryMapper dal;

  public InventoryService(InventoryMapper dal) {
    this.dal = dal;
  }

  public void requireProcessOwner(long id, long uid, boolean admin, boolean update) {
    Map<String, Object> r = require(dal.process(id), "process");
    int status = (int) num(r, "status", 0);
    if (!update && (status == 2 || status == 3)) return;
    if (!admin && num(r, "owner_user_id", 0) != uid)
      throw new IllegalStateException("forbidden: no permission");
  }

  public void requireOrderLeader(long id, long uid, boolean admin, boolean update) {
    Map<String, Object> r = require(dal.order(id), "engineering order");
    if (!update && num(r, "status", 0) != 1) return;
    if (!admin && num(r, "leader_user_id", 0) != uid)
      throw new IllegalStateException("forbidden: no permission");
  }

  public void requireFlowAccess(long id, long uid, boolean admin, boolean update) {
    Map<String, Object> r = require(dal.flow(id), "inventory flow");
    long from = num(r, "from_user_id", 0), to = num(r, "to_user_id", 0);
    if (admin) return;
    if (update || num(r, "flow_status", 0) == 1) {
      if (from != uid) throw new IllegalStateException("forbidden: no permission");
    } else if (from != uid && to != uid)
      throw new IllegalStateException("forbidden: no permission");
  }

  public Map<String, Object> scopeProcess(Map<String, Object> source, long uid, boolean admin) {
    Map<String, Object> q = new HashMap<>(source);
    int scope = (int) num(q, "scope", 0), status = (int) num(q, "status", 0);
    if (scope == 1) q.put("ownerUserId", uid);
    else if (scope == 2) q.remove("ownerUserId");
    else if (!admin) {
      Long owner = optionalLong(q, "ownerUserId");
      if (owner != null && owner != uid)
        throw new IllegalStateException("forbidden: no permission");
      if (status == 1) q.put("ownerUserId", uid);
      if (status == 0 && owner == null) q.put("status", 2);
    }
    return q;
  }

  public Map<String, Object> scopeOrder(Map<String, Object> source, long uid, boolean admin) {
    Map<String, Object> q = new HashMap<>(source);
    int scope = (int) num(q, "scope", 0), status = (int) num(q, "status", 0);
    if (scope == 1) q.put("leaderUserId", uid);
    else if (scope == 2 || scope == 4) {
      q.remove("leaderUserId");
      q.put("draftOwnerLeaderId", uid);
    } else if (!admin) {
      Long leader = optionalLong(q, "leaderUserId");
      if (leader != null && leader != uid)
        throw new IllegalStateException("forbidden: no permission");
      if (status == 1) q.put("leaderUserId", uid);
      if (status == 0 && leader == null) q.put("status", 2);
    }
    return q;
  }

  public Map<String, Object> scopeFlow(Map<String, Object> source, long uid, boolean admin) {
    Map<String, Object> q = new HashMap<>(source);
    int scope = (int) num(q, "scope", 0), status = (int) num(q, "flowStatus", 0);
    if (admin && scope == 0) scope = 2;
    if (scope == 2) {
      q.remove("userId");
      q.put("draftOwnerUserId", admin ? uid : 0L);
    } else if (scope == 3) {
      q.remove("userId");
      q.put("draftOwnerUserId", 0L);
      if (status == 0) q.put("flowStatus", 2);
    } else {
      if (Boolean.TRUE.equals(q.get("isTo")) && status == 1 && !admin)
        throw new IllegalStateException("forbidden: no permission");
      Long requested = optionalLong(q, "userId");
      if (!admin && requested != null && requested != uid)
        throw new IllegalStateException("forbidden: no permission");
      long owner = requested != null && admin ? requested : uid;
      q.put("userId", owner);
      q.put("draftOwnerUserId", owner);
    }
    return q;
  }

  private static long offset(Map<String, Object> q) {
    long p = num(q, "pageNum", 1), s = size(q);
    return (Math.max(1, p) - 1) * s;
  }

  private static long size(Map<String, Object> q) {
    return Math.min(100, Math.max(1, num(q, "pageSize", 10)));
  }

  private static long num(Map<String, Object> m, String k, long d) {
    Object v = m.get(k);
    return v == null ? d : Long.parseLong(v.toString());
  }

  private static Long optionalLong(Map<String, Object> m, String k) {
    long v = num(m, k, 0);
    return v > 0 ? v : null;
  }

  private static Integer optionalInt(Map<String, Object> m, String k) {
    int v = (int) num(m, k, 0);
    return v > 0 ? v : null;
  }

  private static void required(Map<String, Object> row, String... keys) {
    for (String k : keys)
      if (row.get(k) == null || row.get(k).toString().isBlank())
        throw new IllegalArgumentException(k + " is required");
  }

  private static Map<String, Object> result(Object data) {
    return Map.of("code", 0, "message", "success", "data", data);
  }

  private static Map<String, Object> resultView(Object data) {
    return result(toView(data));
  }

  private static Object toView(Object value) {
    if (value instanceof Map<?, ?> source) {
      Map<String, Object> target = new LinkedHashMap<>();
      for (var entry : source.entrySet()) {
        String key = entry.getKey().toString();
        if (key.equals("deleted_at")) continue;
        target.put(apiKey(key), toView(entry.getValue()));
      }
      return target;
    }
    if (value instanceof List<?> list) return list.stream().map(InventoryService::toView).toList();
    return value;
  }

  private static String apiKey(String key) {
    if (key.equals("created_at")) return "createTime";
    if (key.equals("updated_at")) return "updateTime";
    StringBuilder camel = new StringBuilder(key.length());
    boolean upper = false;
    for (int i = 0; i < key.length(); i++) {
      char c = key.charAt(i);
      if (c == '_') {
        upper = true;
      } else {
        camel.append(upper ? Character.toUpperCase(c) : c);
        upper = false;
      }
    }
    return camel.toString();
  }

  private static void changed(int n, String message) {
    if (n == 0) throw new IllegalStateException(message);
  }

  public Map<String, Object> addItem(Map<String, Object> r) {
    required(r, "name", "unit");
    r.putIfAbsent("description", "");
    dal.insertItem(r);
    return result(Map.of("id", r.get("id")));
  }

  public Map<String, Object> updateItem(Map<String, Object> r) {
    required(r, "id", "name", "unit");
    r.putIfAbsent("description", "");
    changed(dal.updateItem(r), "item not found");
    return result(true);
  }

  public Map<String, Object> item(long id) {
    return resultView(
        Optional.ofNullable(dal.item(id))
            .orElseThrow(() -> new NoSuchElementException("item not found")));
  }

  public Map<String, Object> items(Map<String, Object> q) {
    var rows = dal.items((String) q.get("namePrefix"), offset(q), size(q));
    return resultView(page(rows, q));
  }

  @Transactional
  public Map<String, Object> createProcess(Map<String, Object> r) {
    required(r, "ownerUserId", "itemId", "name");
    r.putIfAbsent("description", "");
    dal.insertProcess(r);
    replaceProcessItems(num(r, "id", 0), list(r, "items"));
    return result(Map.of("id", r.get("id")));
  }

  @Transactional
  public Map<String, Object> updateProcess(Map<String, Object> r) {
    changed(dal.updateProcess(r), "only draft process can be updated");
    replaceProcessItems(num(r, "id", 0), list(r, "items"));
    return result(true);
  }

  public Map<String, Object> deleteProcess(long id) {
    changed(dal.softDelete("processes", "status", id), "only draft process can be deleted");
    return result(true);
  }

  public Map<String, Object> submitProcess(long id) {
    if (dal.processItems(id).isEmpty())
      throw new IllegalStateException("process items are required");
    changed(dal.transition("processes", "status", id, 1, 2), "only draft process can be submitted");
    return result(true);
  }

  public Map<String, Object> process(long id) {
    var r = require(dal.process(id), "process");
    r.put("items", dal.processItems(id));
    return resultView(r);
  }

  public Map<String, Object> processes(Map<String, Object> q) {
    return resultView(
        page(
            dal.processes(
                optionalLong(q, "ownerUserId"),
                optionalLong(q, "itemId"),
                optionalInt(q, "status"),
                (String) q.get("namePrefix"),
                (String) q.get("createdDate"),
                offset(q),
                size(q)),
            q));
  }

  private void replaceProcessItems(long id, List<Map<String, Object>> xs) {
    dal.deleteProcessItems(id);
    for (var x : xs) {
      if (num(x, "quantity", 0) <= 0)
        throw new IllegalArgumentException("quantity must be positive");
      Map<String, Object> detail = new HashMap<>(x);
      detail.put("processId", id);
      dal.insertProcessItem(detail);
    }
  }

  @Transactional
  public Map<String, Object> addUnit(Map<String, Object> r) {
    required(r, "itemId", "qualityStatus");
    long itemId = num(r, "itemId", 0);
    long orderId = num(r, "engineeringOrderId", 0);
    require(dal.item(itemId), "item");
    int qualityStatus = (int) num(r, "qualityStatus", 0);
    if (!validQualityStatus(qualityStatus))
      throw new IllegalArgumentException("invalid item unit quality status");
    r.putIfAbsent("stockStatus", 3);
    int stockStatus = (int) num(r, "stockStatus", 3);
    if (!validStockStatus(stockStatus))
      throw new IllegalArgumentException("invalid item unit stock status");
    if (orderId > 0) {
      var order = require(dal.order(orderId), "engineering order");
      if (num(order, "status", 0) != 2)
        throw new IllegalStateException("only submitted engineering order can be bound");
      if (num(order, "item_id", 0) != itemId)
        throw new IllegalStateException("item unit item must match engineering order item");
      if (num(order, "produced_quantity", 0) + 1 > num(order, "expected_quantity", 0))
        throw new IllegalStateException(
            "engineering order produced quantity exceeds expected quantity");
      r.put("engineeringOrderId", orderId);
    } else {
      r.putIfAbsent("engineeringOrderId", null);
    }
    r.putIfAbsent("description", "");
    dal.insertUnit(r);
    changed(dal.addUnitItemCounts(itemId, stockStatus, qualityStatus), "item not found");
    if (orderId > 0)
      changed(dal.addUnitOrderCounts(orderId, qualityStatus), "engineering order not found");
    return result(Map.of("id", r.get("id")));
  }

  @Transactional
  public Map<String, Object> updateUnit(Map<String, Object> r) {
    Map<String, Object> current = require(dal.unit(num(r, "id", 0)), "item unit");
    int stockStatus = (int) num(r, "stockStatus", 0);
    int qualityStatus = (int) num(r, "qualityStatus", 0);
    if (!validStockStatus(stockStatus) || !validQualityStatus(qualityStatus))
      throw new IllegalArgumentException("invalid item unit status");
    changed(dal.updateUnit(r), "item unit not found");
    long itemId = num(current, "item_id", 0);
    int oldStockStatus = (int) num(current, "stock_status", 0);
    int oldQualityStatus = (int) num(current, "quality_status", 0);
    changed(
        dal.changeUnitItemCounts(
            itemId, oldStockStatus, stockStatus, oldQualityStatus, qualityStatus),
        "item not found");
    long orderId = num(current, "engineering_order_id", 0);
    if (orderId > 0 && oldQualityStatus != qualityStatus)
      changed(
          dal.changeUnitOrderCounts(orderId, oldQualityStatus, qualityStatus),
          "engineering order not found");
    return result(true);
  }

  public Map<String, Object> unit(long id) {
    return resultView(require(dal.unit(id), "item unit"));
  }

  public Map<String, Object> units(Map<String, Object> q) {
    return resultView(
        page(
            dal.units(
                optionalLong(q, "itemId"),
                (String) q.get("itemNamePrefix"),
                optionalInt(q, "stockStatus"),
                optionalInt(q, "qualityStatus"),
                optionalLong(q, "engineeringOrderId"),
                optionalLong(q, "inventoryFlowId"),
                (String) q.get("createdDate"),
                offset(q),
                size(q)),
            q));
  }

  @Transactional
  public Map<String, Object> createFlow(Map<String, Object> r) {
    validateFlow(r);
    dal.insertFlow(r);
    replaceFlowDetails(num(r, "id", 0), r);
    return result(Map.of("id", r.get("id")));
  }

  @Transactional
  public Map<String, Object> updateFlow(Map<String, Object> r) {
    validateFlow(r);
    changed(dal.updateFlow(r), "only draft flow can be updated");
    replaceFlowDetails(num(r, "id", 0), r);
    return result(true);
  }

  public Map<String, Object> deleteFlow(long id) {
    changed(dal.softDelete("inventory_flows", "flow_status", id), "only draft flow can be deleted");
    return result(true);
  }

  public Map<String, Object> submitFlow(long id) {
    if (dal.flowItems(id).isEmpty()) throw new IllegalStateException("flow items are required");
    changed(
        dal.transition("inventory_flows", "flow_status", id, 1, 2),
        "only draft flow can be submitted");
    return result(true);
  }

  @Transactional
  public Map<String, Object> auditFlow(Map<String, Object> r) {
    long id = num(r, "id", 0);
    boolean approved = Boolean.TRUE.equals(r.get("approved"));
    var flow = require(dal.flow(id), "inventory flow");
    if (num(flow, "flow_status", 0) != 2)
      throw new IllegalStateException("only submitted flow can be audited");
    if (approved && num(flow, "flow_type", 0) == 2) {
      for (Map<String, Object> item : dal.flowItems(id)) {
        long itemId = num(item, "item_id", 0);
        long quantity = num(item, "apply_quantity", 0);
        if (itemId <= 0 || quantity <= 0) throw new IllegalStateException("invalid flow item");
        changed(dal.reserveItem(itemId, quantity), "insufficient available item quantity");
      }
    }
    int next = approved ? 3 : 4;
    changed(dal.auditFlow(id, num(r, "approvedBy", 0), next), "only submitted flow can be audited");
    return result(true);
  }

  @Transactional
  public Map<String, Object> completeFlow(long id, List<Long> ids) {
    var f = require(dal.flow(id), "inventory flow");
    if (num(f, "flow_status", 0) != 3)
      throw new IllegalStateException("only approved flow can be completed");
    if (ids == null || ids.isEmpty())
      throw new IllegalArgumentException("item unit ids are required");
    int stock = num(f, "flow_type", 0) == 1 ? 1 : 3;
    Map<Long, Map<String, Object>> completedUnits = new HashMap<>();
    for (Map<String, Object> unit : dal.flowUnits(id)) completedUnits.put(num(unit, "id", 0), unit);
    Map<Long, Map<String, Object>> detailByItem = new HashMap<>();
    for (Map<String, Object> item : dal.flowItems(id)) {
      long itemId = num(item, "item_id", 0);
      if (itemId <= 0) continue;
      detailByItem.put(itemId, item);
    }
    if (detailByItem.isEmpty()) throw new IllegalStateException("flow items are required");
    Set<Long> selected = new HashSet<>();
    Map<Long, Long> finishedInRequest = new HashMap<>();
    for (long uid : ids) {
      if (!selected.add(uid)) throw new IllegalStateException("duplicated item unit");
      if (completedUnits.containsKey(uid))
        throw new IllegalStateException("item unit has already been completed in this flow");
      Map<String, Object> unit = require(dal.unit(uid), "item unit");
      long itemId = num(unit, "item_id", 0);
      Map<String, Object> detail = detailByItem.get(itemId);
      if (detail == null)
        throw new IllegalStateException("item unit does not belong to current flow item");
      long nextFinished =
          num(detail, "finished_quantity", 0) + finishedInRequest.getOrDefault(itemId, 0L) + 1;
      if (nextFinished > num(detail, "apply_quantity", 0))
        throw new IllegalStateException("flow item quantity exceeded");
      if (num(f, "flow_type", 0) == 1) {
        if (num(unit, "stock_status", 0) != 3)
          throw new IllegalStateException("only out-stock units can be received");
        if (num(unit, "quality_status", 0) != 2)
          throw new IllegalStateException("only qualified units can be received");
      } else {
        if (num(unit, "stock_status", 0) != 1)
          throw new IllegalStateException("only in-stock units can be released");
        if (num(unit, "quality_status", 0) != 2)
          throw new IllegalStateException("only qualified units can be released");
      }
      dal.bindFlowUnit(id, uid);
      dal.setUnitStock(uid, stock);
      changed(
          dal.completeItemFlow(
              itemId, (int) num(f, "flow_type", 0), num(unit, "quality_status", 0) == 2),
          "item not found");
      finishedInRequest.put(itemId, finishedInRequest.getOrDefault(itemId, 0L) + 1);
      completedUnits.put(uid, unit);
    }
    dal.finishFlowItems(id);
    validateFlowProgress(id);
    return result(true);
  }

  public Map<String, Object> flow(long id) {
    var r = require(dal.flow(id), "inventory flow");
    r.put("items", dal.flowItems(id));
    r.put("itemUnits", dal.flowUnits(id));
    return resultView(r);
  }

  public Map<String, Object> flows(Map<String, Object> q) {
    if (Boolean.TRUE.equals(q.get("onlyDraft"))) q.put("flowStatus", 1);
    List<Map<String, Object>> rows =
        dal.flows(
            optionalLong(q, "userId"),
            Boolean.TRUE.equals(q.get("isTo")),
            optionalInt(q, "flowStatus"),
            optionalInt(q, "businessType"),
            Objects.toString(q.getOrDefault("keyword", q.get("namePrefix")), null),
            (String) q.get("createdDate"),
            optionalLong(q, "itemUnitId"),
            q.containsKey("draftOwnerUserId") ? num(q, "draftOwnerUserId", 0) : null,
            offset(q),
            size(q));
    for (Map<String, Object> row : rows) {
      long flowId = num(row, "id", 0);
      if (flowId > 0) row.put("items", dal.flowItems(flowId));
    }
    return resultView(page(rows, q));
  }

  private void validateFlow(Map<String, Object> r) {
    required(r, "fromUserId", "toUserId", "flowType", "businessType", "name");
    int businessType = (int) num(r, "businessType", 0);
    if (businessType < 1 || businessType > 3)
      throw new IllegalArgumentException("invalid inventory flow business type");
    int expectedFlowType = businessType == 2 ? 2 : 1;
    if (num(r, "flowType", 0) != expectedFlowType)
      throw new IllegalArgumentException("inventory flow direction does not match business type");
    if (num(r, "fromUserId", 0) == num(r, "toUserId", 0))
      throw new IllegalArgumentException("from and to user must differ");
    r.putIfAbsent("description", "");
  }

  private void replaceFlowDetails(long id, Map<String, Object> r) {
    dal.deleteFlowItems(id);
    dal.deleteFlowUnits(id);
    for (var x : list(r, "items")) {
      Map<String, Object> detail = new HashMap<>(x);
      detail.put("flowId", id);
      dal.insertFlowItem(detail);
    }
  }

  public Map<String, Object> createOrder(Map<String, Object> r) {
    validateOrder(r);
    dal.insertOrder(r);
    return result(Map.of("id", r.get("id")));
  }

  public Map<String, Object> updateOrder(Map<String, Object> r) {
    validateOrder(r);
    changed(dal.updateOrder(r), "only draft order can be updated");
    return result(true);
  }

  public Map<String, Object> deleteOrder(long id) {
    changed(dal.softDelete("engineering_orders", "status", id), "only draft order can be deleted");
    return result(true);
  }

  public Map<String, Object> submitOrder(long id) {
    changed(
        dal.transition("engineering_orders", "status", id, 1, 2),
        "only draft order can be submitted");
    return result(true);
  }

  public Map<String, Object> order(long id) {
    var r = require(dal.order(id), "engineering order");
    r.put("itemUnits", dal.units(null, null, null, null, id, null, null, 0, 100));
    return resultView(r);
  }

  public Map<String, Object> orders(Map<String, Object> q) {
    return resultView(
        page(
            dal.orders(
                optionalLong(q, "leaderUserId"),
                optionalLong(q, "itemId"),
                optionalLong(q, "processId"),
                optionalInt(q, "status"),
                optionalInt(q, "progressStatus"),
                optionalLong(q, "draftOwnerLeaderId"),
                Objects.toString(q.getOrDefault("keyword", q.get("namePrefix")), null),
                (String) q.get("createdDate"),
                offset(q),
                size(q)),
            q));
  }

  private void validateOrder(Map<String, Object> r) {
    required(r, "leaderUserId", "processId", "itemId", "name", "expectedQuantity");
    if (num(r, "expectedQuantity", 0) <= 0)
      throw new IllegalArgumentException("expected quantity must be positive");
    r.putIfAbsent("description", "");
  }

  private static Map<String, Object> require(Map<String, Object> r, String what) {
    if (r == null) throw new NoSuchElementException(what + " not found");
    return new LinkedHashMap<>(r);
  }

  @SuppressWarnings("unchecked")
  private static List<Map<String, Object>> list(Map<String, Object> r, String k) {
    return (List<Map<String, Object>>) r.getOrDefault(k, List.of());
  }

  @SuppressWarnings("unchecked")
  private static List<Long> listLong(Map<String, Object> r, String k) {
    return ((List<Object>) r.getOrDefault(k, List.of()))
        .stream().map(x -> Long.parseLong(x.toString())).toList();
  }

  private void validateFlowProgress(long id) {
    for (Map<String, Object> item : dal.flowItems(id)) {
      if (num(item, "finished_quantity", 0) > num(item, "apply_quantity", 0))
        throw new IllegalStateException("flow item quantity exceeded");
    }
  }

  private static boolean validStockStatus(int status) {
    return status == 1 || status == 2 || status == 3;
  }

  private static boolean validQualityStatus(int status) {
    return status == 1 || status == 2 || status == 3;
  }

  private static Map<String, Object> page(List<Map<String, Object>> rows, Map<String, Object> q) {
    return Map.of(
        "records",
        rows,
        "pageNumber",
        num(q, "pageNum", 1),
        "pageSize",
        size(q),
        "totalRow",
        rows.size(),
        "hasMore",
        rows.size() == size(q));
  }
}
