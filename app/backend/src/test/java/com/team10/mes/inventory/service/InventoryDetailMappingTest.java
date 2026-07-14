package com.team10.mes.inventory.service;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

import com.team10.mes.inventory.dal.InventoryMapper;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.time.LocalDateTime;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;

class InventoryDetailMappingTest {
  private final InventoryMapper mapper = mock(InventoryMapper.class);
  private final InventoryService service = new InventoryService(mapper);

  @Test
  void itemUnitDetailUsesOriginalCamelCaseContract() {
    LocalDateTime created = LocalDateTime.of(2026, 7, 1, 10, 0);
    LocalDateTime updated = LocalDateTime.of(2026, 7, 2, 11, 30);
    when(mapper.unit(9))
        .thenReturn(
            row(
                "id",
                9L,
                "item_id",
                3L,
                "stock_status",
                1,
                "quality_status",
                2,
                "engineering_order_id",
                7L,
                "created_at",
                created,
                "updated_at",
                updated));

    Map<String, Object> unit = data(service.unit(9));

    assertEquals(3L, unit.get("itemId"));
    assertEquals(1, unit.get("stockStatus"));
    assertEquals(2, unit.get("qualityStatus"));
    assertEquals(7L, unit.get("engineeringOrderId"));
    assertEquals(created, unit.get("createTime"));
    assertEquals(updated, unit.get("updateTime"));
    assertFalse(unit.containsKey("item_id"));
  }

  @Test
  void flowDetailMapsStatusTimesAndNestedRows() {
    when(mapper.flow(5))
        .thenReturn(
            row(
                "id",
                5L,
                "from_user_id",
                10L,
                "to_user_id",
                11L,
                "flow_type",
                1,
                "flow_status",
                3,
                "approved_at",
                LocalDateTime.of(2026, 7, 3, 9, 0),
                "updated_at",
                LocalDateTime.of(2026, 7, 3, 9, 1)));
    when(mapper.flowItems(5))
        .thenReturn(
            List.of(
                row(
                    "inventory_flow_id",
                    5L,
                    "item_id",
                    3L,
                    "apply_quantity",
                    4L,
                    "finished_quantity",
                    2L)));
    when(mapper.flowUnits(5)).thenReturn(List.of(row("id", 9L, "item_id", 3L, "stock_status", 1)));

    Map<String, Object> flow = data(service.flow(5));

    assertEquals(1, flow.get("flowType"));
    assertEquals(3, flow.get("flowStatus"));
    assertEquals(10L, flow.get("fromUserId"));
    assertTrue(flow.containsKey("approvedAt"));
    Map<String, Object> detail = first(flow, "items");
    assertEquals(5L, detail.get("inventoryFlowId"));
    assertEquals(4L, detail.get("applyQuantity"));
    assertEquals(2L, detail.get("finishedQuantity"));
    assertEquals(3L, first(flow, "itemUnits").get("itemId"));
  }

  @Test
  void traceListPassesItemUnitIdToMapperAndMapsFlowRows() {
    when(mapper.flows(isNull(), eq(false), isNull(), isNull(), eq(9L), isNull(), eq(0L), eq(50L)))
        .thenReturn(List.of(row("id", 5L, "flow_type", 2, "flow_status", 3)));

    Map<String, Object> page =
        data(service.flows(new LinkedHashMap<>(Map.of("itemUnitId", 9L, "pageSize", 50))));

    Map<String, Object> flow = first(page, "records");
    assertEquals(2, flow.get("flowType"));
    assertEquals(3, flow.get("flowStatus"));
    verify(mapper).flows(null, false, null, null, 9L, null, 0L, 50L);
  }

  @Test
  void globalFlowScopeHidesAllDrafts() {
    Map<String, Object> scoped =
        service.scopeFlow(new LinkedHashMap<>(Map.of("scope", 2, "pageSize", 30)), 12L, false);

    assertFalse(scoped.containsKey("userId"));
    assertEquals(0L, scoped.get("draftOwnerUserId"));

    when(mapper.flows(isNull(), eq(false), isNull(), isNull(), isNull(), eq(0L), eq(0L), eq(30L)))
        .thenReturn(List.of(row("id", 8L, "from_user_id", 13L, "flow_type", 2, "flow_status", 2)));

    Map<String, Object> page = data(service.flows(scoped));
    assertEquals(1, ((List<?>) page.get("records")).size());
    verify(mapper).flows(null, false, null, null, null, 0L, 0L, 30L);
  }

  @Test
  void mapperXmlFiltersTraceByBindingTable() throws Exception {
    try (InputStream stream =
        getClass().getResourceAsStream("/mapper/inventory/InventoryMapper.xml")) {
      assertNotNull(stream);
      String xml = new String(stream.readAllBytes(), StandardCharsets.UTF_8);
      assertTrue(xml.contains("fx.item_unit_id=#{itemUnitId}"));
      assertTrue(xml.contains("fx.inventory_flow_id=f.id"));
      assertTrue(xml.contains("f.flow_status!=1"));
      assertTrue(xml.contains("draftOwnerUserId &gt; 0"));
      assertTrue(xml.contains("OR f.from_user_id=#{draftOwnerUserId}"));
    }
  }

  @SuppressWarnings("unchecked")
  private static Map<String, Object> data(Map<String, Object> response) {
    return (Map<String, Object>) response.get("data");
  }

  @SuppressWarnings("unchecked")
  private static Map<String, Object> first(Map<String, Object> source, String key) {
    return ((List<Map<String, Object>>) source.get(key)).getFirst();
  }

  private static Map<String, Object> row(Object... pairs) {
    Map<String, Object> row = new LinkedHashMap<>();
    for (int i = 0; i < pairs.length; i += 2) row.put(pairs[i].toString(), pairs[i + 1]);
    return row;
  }
}
