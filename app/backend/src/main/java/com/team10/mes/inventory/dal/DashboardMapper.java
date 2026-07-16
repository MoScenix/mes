package com.team10.mes.inventory.dal;

import java.util.List;
import java.util.Map;
import org.apache.ibatis.annotations.Mapper;

@Mapper
public interface DashboardMapper {
  Map<String, Object> productionSummary();

  Map<String, Object> planSummary();

  List<Map<String, Object>> dailyProduction();
}
