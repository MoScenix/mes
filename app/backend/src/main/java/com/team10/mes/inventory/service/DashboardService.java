package com.team10.mes.inventory.service;

import com.team10.mes.inventory.dal.DashboardMapper;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.*;
import org.springframework.stereotype.Service;

@Service
public class DashboardService {
  private final DashboardMapper mapper;

  public DashboardService(DashboardMapper mapper) {
    this.mapper = mapper;
  }

  public Map<String, Object> overview() {
    Map<String, Object> production = mapper.productionSummary();
    Map<String, Object> plans = mapper.planSummary();
    long expected = number(plans, "expected_quantity");
    long completedQuantity = number(plans, "completed_quantity");

    Map<String, Long> quantitiesByDate = new HashMap<>();
    for (Map<String, Object> row : mapper.dailyProduction()) {
      quantitiesByDate.put(text(row, "production_date"), number(row, "quantity"));
    }
    List<Map<String, Object>> trend = new ArrayList<>();
    DateTimeFormatter formatter = DateTimeFormatter.ISO_LOCAL_DATE;
    for (int offset = 6; offset >= 0; offset--) {
      String date = LocalDate.now().minusDays(offset).format(formatter);
      trend.add(Map.of("date", date, "quantity", quantitiesByDate.getOrDefault(date, 0L)));
    }

    Map<String, Object> data = new LinkedHashMap<>();
    data.put("todayProduction", number(production, "today_production"));
    data.put("weekProduction", number(production, "week_production"));
    data.put("pendingInspection", number(production, "pending_inspection"));
    data.put(
        "planCompletionRate",
        expected == 0 ? 0 : Math.round(completedQuantity * 1000.0 / expected) / 10.0);
    data.put("planExpectedQuantity", expected);
    data.put("planCompletedQuantity", completedQuantity);
    data.put(
        "planStatus",
        Map.of(
            "notStarted", number(plans, "not_started"),
            "inProgress", number(plans, "in_progress"),
            "completed", number(plans, "completed")));
    data.put("dailyProduction", trend);
    data.put("generatedAt", java.time.LocalDateTime.now().toString());
    return Map.of("code", 0, "message", "success", "data", data);
  }

  private static long number(Map<String, Object> row, String key) {
    Object value = value(row, key);
    return value == null ? 0 : Long.parseLong(value.toString());
  }

  private static String text(Map<String, Object> row, String key) {
    return Objects.toString(value(row, key), "");
  }

  private static Object value(Map<String, Object> row, String key) {
    if (row == null) return null;
    Object value = row.get(key);
    if (value != null || !key.contains("_")) return value;
    String[] parts = key.split("_");
    StringBuilder camel = new StringBuilder(parts[0]);
    for (int i = 1; i < parts.length; i++) {
      if (!parts[i].isEmpty())
        camel.append(Character.toUpperCase(parts[i].charAt(0))).append(parts[i].substring(1));
    }
    return row.get(camel.toString());
  }
}
