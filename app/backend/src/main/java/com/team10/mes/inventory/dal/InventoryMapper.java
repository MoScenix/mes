package com.team10.mes.inventory.dal;

import java.util.List;
import java.util.Map;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

@Mapper
public interface InventoryMapper {
  int insertItem(Map<String, Object> row);

  int updateItem(Map<String, Object> row);

  Map<String, Object> item(long id);

  List<Map<String, Object>> items(
      @Param("name") String name, @Param("offset") long offset, @Param("limit") long limit);

  int insertProcess(Map<String, Object> row);

  int updateProcess(Map<String, Object> row);

  int insertProcessItem(Map<String, Object> row);

  int deleteProcessItems(long id);

  List<Map<String, Object>> processItems(long id);

  Map<String, Object> process(long id);

  List<Map<String, Object>> processes(
      @Param("ownerUserId") Long ownerUserId,
      @Param("itemId") Long itemId,
      @Param("status") Integer status,
      @Param("offset") long offset,
      @Param("limit") long limit);

  int insertUnit(Map<String, Object> row);

  int updateUnit(Map<String, Object> row);

  Map<String, Object> unit(long id);

  List<Map<String, Object>> units(
      @Param("itemId") Long itemId,
      @Param("itemName") String itemName,
      @Param("stockStatus") Integer stockStatus,
      @Param("qualityStatus") Integer qualityStatus,
      @Param("orderId") Long orderId,
      @Param("flowId") Long flowId,
      @Param("offset") long offset,
      @Param("limit") long limit);

  int insertFlow(Map<String, Object> row);

  int updateFlow(Map<String, Object> row);

  int insertFlowItem(Map<String, Object> row);

  int deleteFlowItems(long id);

  int deleteFlowUnits(long id);

  int bindFlowUnit(@Param("flowId") long flowId, @Param("unitId") long unitId);

  int countFlowUnit(@Param("flowId") long flowId, @Param("unitId") long unitId);

  List<Map<String, Object>> flowItems(long id);

  List<Map<String, Object>> flowUnits(long id);

  Map<String, Object> flow(long id);

  List<Map<String, Object>> flows(
      @Param("userId") Long userId,
      @Param("isTo") boolean isTo,
      @Param("status") Integer status,
      @Param("name") String name,
      @Param("itemUnitId") Long itemUnitId,
      @Param("draftOwnerUserId") Long draftOwnerUserId,
      @Param("offset") long offset,
      @Param("limit") long limit);

  int insertOrder(Map<String, Object> row);

  int updateOrder(Map<String, Object> row);

  Map<String, Object> order(long id);

  List<Map<String, Object>> orders(
      @Param("leaderId") Long leaderId,
      @Param("itemId") Long itemId,
      @Param("processId") Long processId,
      @Param("status") Integer status,
      @Param("offset") long offset,
      @Param("limit") long limit);

  int transition(
      @Param("table") String table,
      @Param("statusColumn") String statusColumn,
      @Param("id") long id,
      @Param("expected") int expected,
      @Param("next") int next);

  int softDelete(
      @Param("table") String table,
      @Param("statusColumn") String statusColumn,
      @Param("id") long id);

  int auditFlow(@Param("id") long id, @Param("by") long by, @Param("status") int status);

  int setUnitStock(@Param("id") long id, @Param("status") int status);

  int finishFlowItems(long id);

  int addUnitItemCounts(
      @Param("id") long id,
      @Param("stockStatus") int stockStatus,
      @Param("qualityStatus") int qualityStatus);

  int changeUnitItemCounts(
      @Param("id") long id,
      @Param("oldStockStatus") int oldStockStatus,
      @Param("newStockStatus") int newStockStatus,
      @Param("oldQualityStatus") int oldQualityStatus,
      @Param("newQualityStatus") int newQualityStatus);

  int addUnitOrderCounts(@Param("id") long id, @Param("qualityStatus") int qualityStatus);

  int changeUnitOrderCounts(
      @Param("id") long id,
      @Param("oldQualityStatus") int oldQualityStatus,
      @Param("newQualityStatus") int newQualityStatus);

  int reserveItem(@Param("id") long id, @Param("quantity") long quantity);

  int completeItemFlow(
      @Param("id") long id,
      @Param("flowType") int flowType,
      @Param("qualified") boolean qualified);
}
