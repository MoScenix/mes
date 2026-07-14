package com.team10.mes.ai.service;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.team10.mes.document.service.DocumentService;
import com.team10.mes.inventory.service.InventoryService;
import com.team10.mes.user.service.UserService;
import com.team10.mes.workorder.service.WorkOrderService;
import java.util.*;
import org.springframework.ai.tool.annotation.Tool;
import org.springframework.ai.tool.annotation.ToolParam;

/** Per-request Spring AI tool object bound to the authenticated operator. */
public final class MesAiTools {
  private final WorkOrderService workOrders;
  private final InventoryService inventory;
  private final UserService users;
  private final DocumentService documents;
  private final long historyId;
  private final long userId;
  private final boolean admin;

  public MesAiTools(
      WorkOrderService workOrders,
      InventoryService inventory,
      UserService users,
      DocumentService documents,
      long historyId,
      long userId,
      boolean admin) {
    this.workOrders = workOrders;
    this.inventory = inventory;
    this.users = users;
    this.documents = documents;
    this.historyId = historyId;
    this.userId = userId;
    this.admin = admin;
  }

  @Tool(
      name = "list_work_orders",
      description = "Get latest time-ordered work orders. Default limit is 30.")
  public WorkOrderService.Page listWorkOrders(ListWorkOrdersInput in) {
    long target = admin && in.userId() > 0 ? in.userId() : userId;
    return workOrders.list(
        new WorkOrderService.ListRequest(
            1L,
            limit(in.limit()),
            target,
            in.isTo(),
            in.unread(),
            null,
            null,
            null,
            null,
            in.namePrefix(),
            workOrderStatus(in.status()),
            null),
        userId,
        admin);
  }

  @Tool(name = "mark_work_order_read", description = "Mark a work order as read.")
  public Map<String, Object> markWorkOrderRead(IdInput in) {
    workOrders.read(in.id(), userId, admin);
    return Map.of("success", true);
  }

  @Tool(name = "create_work_order_draft", description = "Create a work order draft.")
  public Map<String, Object> createWorkOrderDraft(WorkOrderDraftInput in) {
    return Map.of(
        "id",
        workOrders.create(
            userId, new WorkOrderService.DraftRequest(in.toUserId(), in.description(), in.name())));
  }

  @Tool(name = "update_work_order_draft", description = "Update an existing work order draft.")
  public Map<String, Object> updateWorkOrderDraft(UpdateWorkOrderInput in) {
    workOrders.update(
        userId,
        admin,
        new WorkOrderService.UpdateRequest(in.id(), in.toUserId(), in.description(), in.name()));
    return Map.of("success", true);
  }

  @Tool(
      name = "search_users",
      description =
          "Search assignable users by id, name, account, or role; returns no credential fields.")
  public Map<String, Object> searchUsers(SearchUsersInput in) {
    if (in.id() > 0) return Map.of("users", List.of(users.get(in.id())), "total", 1);
    Map<String, Object> q = new HashMap<>();
    q.put("pageNum", 1);
    q.put("pageSize", limit10(in.pageSize()));
    q.put("userName", in.name());
    q.put("account", in.account());
    return users.page(q);
  }

  @Tool(
      name = "search_history_file",
      description =
          "Search an uploaded large file in the current history by file_id and query. Use this whenever the user asks about an uploaded file's contents; returns matched parent text. top_k should be 3 to 5; default is 3.")
  public Map<String, Object> searchHistoryFile(SearchHistoryFileInput in) {
    return documents.searchWithParents(historyId, in.fileId(), in.query(), in.topK());
  }

  @Tool(
      name = "create_engineering_order_draft",
      description = "Create an engineering order draft for planned production.")
  public Map<String, Object> createEngineeringOrderDraft(EngineeringOrderInput in) {
    Map<String, Object> m = orderMap(in);
    return inventory.createOrder(m);
  }

  @Tool(name = "update_engineering_order_draft", description = "Update an engineering order draft.")
  public Map<String, Object> updateEngineeringOrderDraft(UpdateEngineeringOrderInput in) {
    Map<String, Object> m = orderMap(in.value());
    m.put("id", in.id());
    return inventory.updateOrder(m);
  }

  @Tool(name = "list_engineering_orders", description = "List engineering orders.")
  public Map<String, Object> listEngineeringOrders(ListOrdersInput in) {
    Map<String, Object> q = page(in.pageNum(), in.pageSize());
    q.put("namePrefix", in.namePrefix());
    q.put("leaderUserId", admin && in.leaderUserId() > 0 ? in.leaderUserId() : userId);
    if (in.itemId() > 0) q.put("itemId", in.itemId());
    return inventory.orders(q);
  }

  @Tool(name = "get_engineering_order", description = "Get engineering order details.")
  public Map<String, Object> getEngineeringOrder(IdInput in) {
    return inventory.order(in.id());
  }

  @Tool(
      name = "create_inventory_flow_draft",
      description =
          "Create an inbound or outbound inventory flow draft. The backend always uses the authenticated current user as from_user; do not ask for from_user.")
  public Map<String, Object> createInventoryFlowDraft(InventoryFlowInput in) {
    Map<String, Object> m = new HashMap<>();
    m.put("name", in.name());
    m.put("fromUserId", userId);
    m.put("toUserId", in.toUserId());
    m.put("flowType", flowType(in.flowType()));
    m.put("description", in.description());
    m.put(
        "items",
        in.items().stream()
            .map(
                x ->
                    Map.<String, Object>of(
                        "itemId", x.itemId(), "applyQuantity", x.applyQuantity()))
            .toList());
    m.put("itemUnitIds", in.itemUnitIds());
    return inventory.createFlow(m);
  }

  @Tool(
      name = "list_inventory_flows",
      description = "Get latest time-ordered inventory flows by scope and status.")
  public Map<String, Object> listInventoryFlows(ListFlowsInput in) {
    Map<String, Object> q = page(1, limit(in.limit()));
    q.put("userId", admin && in.userId() > 0 ? in.userId() : userId);
    q.put("namePrefix", in.namePrefix());
    Integer s = flowStatus(in.flowStatus());
    if (s != null) q.put("flowStatus", s);
    return inventory.flows(q);
  }

  @Tool(name = "get_inventory_flow", description = "Get inventory flow details.")
  public Map<String, Object> getInventoryFlow(IdInput in) {
    return inventory.flow(in.id());
  }

  @Tool(
      name = "search_items",
      description =
          "Search item definitions. availableCount is usable stock; totalCount is all concrete units.")
  public Map<String, Object> searchItems(SearchItemsInput in) {
    Map<String, Object> q = page(in.pageNum(), in.pageSize());
    q.put("namePrefix", in.namePrefix());
    return inventory.items(q);
  }

  @Tool(name = "get_item", description = "Get item definition details.")
  public Map<String, Object> getItem(IdInput in) {
    return inventory.item(in.id());
  }

  @Tool(name = "list_item_units", description = "List concrete item units in inventory.")
  public Map<String, Object> listItemUnits(ListUnitsInput in) {
    Map<String, Object> q = page(in.pageNum(), in.pageSize());
    if (in.itemId() > 0) q.put("itemId", in.itemId());
    Integer s = stockStatus(in.stockStatus());
    if (s != null) q.put("stockStatus", s);
    Integer quality = qualityStatus(in.qualityStatus());
    if (quality != null) q.put("qualityStatus", quality);
    return inventory.units(q);
  }

  @Tool(
      name = "list_pending_inventory_flows",
      description = "List submitted inventory flows pending warehouse processing.")
  public Map<String, Object> listPendingInventoryFlows(ListPendingInput in) {
    return listInventoryFlows(
        new ListFlowsInput(in.limit(), in.namePrefix(), "to_me", "submitted", in.userId()));
  }

  @Tool(
      name = "inventory_check",
      description = "Read-only inventory check for item stock and concrete units.")
  public Map<String, Object> inventoryCheck(InventoryCheckInput in) {
    return Map.of(
        "items",
        searchItems(new SearchItemsInput(in.namePrefix(), in.pageNum(), in.pageSize())),
        "itemUnits",
        listItemUnits(
            new ListUnitsInput(
                in.itemId(), in.stockStatus(), in.qualityStatus(), in.pageNum(), in.pageSize())));
  }

  private Map<String, Object> orderMap(EngineeringOrderInput in) {
    Map<String, Object> m = new HashMap<>();
    m.put("name", in.name());
    m.put("leaderUserId", admin && in.leaderUserId() > 0 ? in.leaderUserId() : userId);
    m.put("processId", in.processId());
    m.put("itemId", in.itemId());
    m.put("expectedQuantity", in.expectedQuantity());
    m.put("qualifiedQuantity", in.qualifiedQuantity());
    m.put("description", in.description());
    return m;
  }

  private static Map<String, Object> page(long p, long s) {
    Map<String, Object> q = new HashMap<>();
    q.put("pageNum", p <= 0 ? 1 : p);
    q.put("pageSize", limit(s));
    return q;
  }

  private static long limit(long n) {
    return n <= 0 ? 30 : Math.min(n, 100);
  }

  private static long limit10(long n) {
    return n <= 0 ? 10 : Math.min(n, 100);
  }

  private static Integer workOrderStatus(String s) {
    return "draft".equalsIgnoreCase(s) ? 1 : "submitted".equalsIgnoreCase(s) ? 2 : null;
  }

  private static Integer flowStatus(String s) {
    return switch (Objects.toString(s, "").toLowerCase()) {
      case "draft" -> 1;
      case "submitted" -> 2;
      case "approved" -> 3;
      case "rejected" -> 4;
      default -> null;
    };
  }

  private static int flowType(String s) {
    return "out".equalsIgnoreCase(s) ? 2 : 1;
  }

  private static Integer stockStatus(String s) {
    return switch (Objects.toString(s, "").toLowerCase()) {
      case "in_stock" -> 1;
      case "reserved" -> 2;
      case "out_stock" -> 3;
      default -> null;
    };
  }

  private static Integer qualityStatus(String s) {
    return switch (Objects.toString(s, "").toLowerCase()) {
      case "pending" -> 1;
      case "qualified" -> 2;
      case "unqualified" -> 3;
      default -> null;
    };
  }

  public record IdInput(long id) {}

  public record ListWorkOrdersInput(
      long limit,
      @JsonProperty("name_prefix") String namePrefix,
      @JsonProperty("user_id") long userId,
      @JsonProperty("is_to") boolean isTo,
      boolean unread,
      String status) {}

  public record WorkOrderDraftInput(
      String name, @JsonProperty("to_user_id") long toUserId, String description) {}

  public record UpdateWorkOrderInput(
      long id, String name, @JsonProperty("to_user_id") long toUserId, String description) {}

  public record SearchUsersInput(
      long id,
      String name,
      String account,
      String role,
      @JsonProperty("page_size") long pageSize) {}

  public record SearchHistoryFileInput(
      @JsonProperty("file_id") long fileId,
      String query,
      @JsonProperty("top_k")
          @ToolParam(
              required = false,
              description =
                  "Number of parent chunks to return. Recommended range is 3 to 5; default is 3 and values above 5 are capped.")
          long topK) {}

  public record EngineeringOrderInput(
      String name,
      @JsonProperty("leader_user_id") long leaderUserId,
      @JsonProperty("process_id") long processId,
      @JsonProperty("item_id") long itemId,
      @JsonProperty("expected_quantity") long expectedQuantity,
      @JsonProperty("qualified_quantity") long qualifiedQuantity,
      String description) {}

  public record UpdateEngineeringOrderInput(
      long id,
      String name,
      @JsonProperty("leader_user_id") long leaderUserId,
      @JsonProperty("process_id") long processId,
      @JsonProperty("item_id") long itemId,
      @JsonProperty("expected_quantity") long expectedQuantity,
      @JsonProperty("qualified_quantity") long qualifiedQuantity,
      String description) {
    EngineeringOrderInput value() {
      return new EngineeringOrderInput(
          name, leaderUserId, processId, itemId, expectedQuantity, qualifiedQuantity, description);
    }
  }

  public record ListOrdersInput(
      @JsonProperty("page_num") long pageNum,
      @JsonProperty("page_size") long pageSize,
      @JsonProperty("name_prefix") String namePrefix,
      @JsonProperty("leader_user_id") long leaderUserId,
      @JsonProperty("item_id") long itemId) {}

  public record FlowItem(
      @JsonProperty("item_id") long itemId, @JsonProperty("apply_quantity") long applyQuantity) {}

  public record InventoryFlowInput(
      String name,
      @JsonProperty("to_user_id") long toUserId,
      @JsonProperty("flow_type") String flowType,
      String description,
      List<FlowItem> items,
      @JsonProperty("item_unit_ids") List<Long> itemUnitIds) {
    public InventoryFlowInput {
      items = items == null ? List.of() : items;
      itemUnitIds = itemUnitIds == null ? List.of() : itemUnitIds;
    }
  }

  public record ListFlowsInput(
      long limit,
      @JsonProperty("name_prefix") String namePrefix,
      String scope,
      @JsonProperty("flow_status") String flowStatus,
      @JsonProperty("user_id") long userId) {}

  public record SearchItemsInput(
      @JsonProperty("name_prefix") String namePrefix,
      @JsonProperty("page_num") long pageNum,
      @JsonProperty("page_size") long pageSize) {}

  public record ListUnitsInput(
      @JsonProperty("item_id") long itemId,
      @JsonProperty("stock_status") String stockStatus,
      @JsonProperty("quality_status") String qualityStatus,
      @JsonProperty("page_num") long pageNum,
      @JsonProperty("page_size") long pageSize) {}

  public record ListPendingInput(
      long limit,
      @JsonProperty("name_prefix") String namePrefix,
      @JsonProperty("user_id") long userId) {}

  public record InventoryCheckInput(
      @JsonProperty("name_prefix") String namePrefix,
      @JsonProperty("item_id") long itemId,
      @JsonProperty("stock_status") String stockStatus,
      @JsonProperty("quality_status") String qualityStatus,
      @JsonProperty("page_num") long pageNum,
      @JsonProperty("page_size") long pageSize) {}
}
