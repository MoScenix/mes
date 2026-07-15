package com.team10.mes.inventory.service;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

import com.team10.mes.inventory.dal.InventoryMapper;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;

class InventoryMutationConsistencyTest {
  private final InventoryMapper mapper = mock(InventoryMapper.class);
  private final InventoryService service = new InventoryService(mapper);

  @Test
  void addingBoundUnitMaintainsItemAndEngineeringOrderCounts() {
    when(mapper.item(3)).thenReturn(row("id", 3L));
    when(mapper.order(7))
        .thenReturn(
            row(
                "id", 7L,
                "item_id", 3L,
                "status", 2,
                "produced_quantity", 1L,
                "expected_quantity", 2L));
    Map<String, Object> request =
        new LinkedHashMap<>(Map.of("itemId", 3L, "engineeringOrderId", 7L, "qualityStatus", 2));
    doAnswer(
            invocation -> {
              request.put("id", 9L);
              return 1;
            })
        .when(mapper)
        .insertUnit(request);
    when(mapper.addUnitItemCounts(3, 3, 2)).thenReturn(1);
    when(mapper.addUnitOrderCounts(7, 2)).thenReturn(1);

    assertEquals(0, service.addUnit(request).get("code"));

    verify(mapper).addUnitItemCounts(3L, 3, 2);
    verify(mapper).addUnitOrderCounts(7L, 2);
  }

  @Test
  void inspectingBoundUnitMaintainsEngineeringOrderQualityCounts() {
    when(mapper.unit(9))
        .thenReturn(
            row(
                "id", 9L,
                "item_id", 3L,
                "engineering_order_id", 7L,
                "stock_status", 3,
                "quality_status", 1));
    when(mapper.updateUnit(anyMap())).thenReturn(1);
    when(mapper.changeUnitItemCounts(3, 3, 3, 1, 2)).thenReturn(1);
    when(mapper.changeUnitOrderCounts(7, 1, 2)).thenReturn(1);

    Map<String, Object> request =
        new LinkedHashMap<>(Map.of("id", 9L, "stockStatus", 3, "qualityStatus", 2));
    assertEquals(0, service.updateUnit(request).get("code"));

    verify(mapper).changeUnitItemCounts(3L, 3, 3, 1, 2);
    verify(mapper).changeUnitOrderCounts(7L, 1, 2);
  }

  @Test
  void completingInboundFlowMaintainsProgressAndItemCountsOnly() {
    when(mapper.flow(5)).thenReturn(row("id", 5L, "flow_status", 3, "flow_type", 1));
    when(mapper.flowUnits(5)).thenReturn(List.of());
    when(mapper.flowItems(5))
        .thenReturn(
            List.of(row("item_id", 3L, "apply_quantity", 1L, "finished_quantity", 0L)),
            List.of(row("item_id", 3L, "apply_quantity", 1L, "finished_quantity", 1L)));
    when(mapper.unit(9))
        .thenReturn(
            row(
                "id", 9L,
                "item_id", 3L,
                "engineering_order_id", 7L,
                "stock_status", 3,
                "quality_status", 2));
    when(mapper.bindFlowUnit(5, 9)).thenReturn(1);
    when(mapper.setUnitStock(9, 1)).thenReturn(1);
    when(mapper.completeItemFlow(3, 1, true)).thenReturn(1);

    assertEquals(0, service.completeFlow(5, List.of(9L)).get("code"));

    verify(mapper).finishFlowItems(5L);
    verify(mapper).completeItemFlow(3L, 1, true);
    verify(mapper, never()).changeUnitOrderCounts(anyLong(), anyInt(), anyInt());
  }

  @Test
  void approvingOutboundFlowValidatesAndReservesAvailableQuantity() {
    when(mapper.flow(5)).thenReturn(row("id", 5L, "flow_status", 2, "flow_type", 2));
    when(mapper.auditFlow(5, 11, 3)).thenReturn(1);
    when(mapper.flowItems(5)).thenReturn(List.of(row("item_id", 3L, "apply_quantity", 4L)));
    when(mapper.reserveItem(3, 4)).thenReturn(1);

    Map<String, Object> request =
        new LinkedHashMap<>(Map.of("id", 5L, "approvedBy", 11L, "approved", true));
    assertEquals(0, service.auditFlow(request).get("code"));

    verify(mapper).reserveItem(3L, 4L);
    verify(mapper, never()).addUnitOrderCounts(anyLong(), anyInt());
  }

  @Test
  void inboundFlowRejectsUnitThatIsNotQualified() {
    when(mapper.flow(5)).thenReturn(row("id", 5L, "flow_status", 3, "flow_type", 1));
    when(mapper.flowUnits(5)).thenReturn(List.of());
    when(mapper.flowItems(5))
        .thenReturn(List.of(row("item_id", 3L, "apply_quantity", 1L, "finished_quantity", 0L)));
    when(mapper.unit(9))
        .thenReturn(row("id", 9L, "item_id", 3L, "stock_status", 3, "quality_status", 3));

    IllegalStateException error =
        assertThrows(IllegalStateException.class, () -> service.completeFlow(5, List.of(9L)));

    assertEquals("only qualified units can be received", error.getMessage());
    verify(mapper, never()).bindFlowUnit(anyLong(), anyLong());
    verify(mapper, never()).setUnitStock(anyLong(), anyInt());
  }

  private static Map<String, Object> row(Object... pairs) {
    Map<String, Object> row = new LinkedHashMap<>();
    for (int i = 0; i < pairs.length; i += 2) row.put(pairs[i].toString(), pairs[i + 1]);
    return row;
  }
}
