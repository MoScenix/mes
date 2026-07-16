package com.team10.mes.config;

import java.util.List;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Component;

@Component
public class DatabaseIndexInitializer implements ApplicationRunner {
  private static final List<Index> INDEXES =
      List.of(
          index("items", "idx_items_deleted_updated_id", "deleted_at, updated_at DESC, id DESC"),
          index("items", "idx_items_deleted_name_updated_id", "deleted_at, name, updated_at DESC, id DESC"),
          index("item_units", "idx_item_units_deleted_updated_id", "deleted_at, updated_at DESC, id DESC"),
          index("item_units", "idx_item_units_deleted_item_updated_id", "deleted_at, item_id, updated_at DESC, id DESC"),
          index("item_units", "idx_item_units_deleted_engineering_updated_id", "deleted_at, engineering_order_id, updated_at DESC, id DESC"),
          index("item_units", "idx_item_units_deleted_stock_quality_updated_id", "deleted_at, stock_status, quality_status, updated_at DESC, id DESC"),
          index("item_units", "idx_item_units_deleted_item_stock_quality_updated_id", "deleted_at, item_id, stock_status, quality_status, updated_at DESC, id DESC"),
          index("item_units", "idx_item_units_deleted_engineering_quality", "deleted_at, engineering_order_id, quality_status"),
          index("processes", "idx_processes_deleted_updated_id", "deleted_at, updated_at DESC, id DESC"),
          index("processes", "idx_processes_deleted_owner_updated_id", "deleted_at, owner_user_id, updated_at DESC, id DESC"),
          index("processes", "idx_processes_deleted_owner_status_updated_id", "deleted_at, owner_user_id, status, updated_at DESC, id DESC"),
          index("processes", "idx_processes_deleted_item_updated_id", "deleted_at, item_id, updated_at DESC, id DESC"),
          index("processes", "idx_processes_deleted_status_updated_id", "deleted_at, status, updated_at DESC, id DESC"),
          index("processes", "idx_processes_deleted_name_updated_id", "deleted_at, name, updated_at DESC, id DESC"),
          unique("process_items", "uk_process_items_process_consume_item", "process_id, consume_item_id"),
          index("process_items", "idx_process_items_deleted_process", "deleted_at, process_id"),
          index("process_items", "idx_process_items_deleted_consume_process", "deleted_at, consume_item_id, process_id"),
          index("engineering_orders", "idx_engineering_orders_deleted_updated_id", "deleted_at, updated_at DESC, id DESC"),
          index("engineering_orders", "idx_engineering_orders_deleted_leader_updated_id", "deleted_at, leader_user_id, updated_at DESC, id DESC"),
          index("engineering_orders", "idx_engineering_orders_deleted_leader_status_updated_id", "deleted_at, leader_user_id, status, updated_at DESC, id DESC"),
          index("engineering_orders", "idx_engineering_orders_deleted_item_updated_id", "deleted_at, item_id, updated_at DESC, id DESC"),
          index("engineering_orders", "idx_engineering_orders_deleted_process_updated_id", "deleted_at, process_id, updated_at DESC, id DESC"),
          index("engineering_orders", "idx_engineering_orders_deleted_process_status_updated_id", "deleted_at, process_id, status, updated_at DESC, id DESC"),
          index("engineering_orders", "idx_engineering_orders_deleted_status_updated_id", "deleted_at, status, updated_at DESC, id DESC"),
          index("engineering_orders", "idx_engineering_orders_deleted_name_updated_id", "deleted_at, name, updated_at DESC, id DESC"),
          index("inventory_flows", "idx_inventory_flows_deleted_updated_id", "deleted_at, updated_at DESC, id DESC"),
          index("inventory_flows", "idx_inventory_flows_deleted_status_updated_id", "deleted_at, flow_status, updated_at DESC, id DESC"),
          index("inventory_flows", "idx_inventory_flows_deleted_business_status_updated_id", "deleted_at, business_type, flow_status, updated_at DESC, id DESC"),
          index("inventory_flows", "idx_inventory_flows_deleted_from_updated_id", "deleted_at, from_user_id, updated_at DESC, id DESC"),
          index("inventory_flows", "idx_inventory_flows_deleted_to_updated_id", "deleted_at, to_user_id, updated_at DESC, id DESC"),
          index("inventory_flows", "idx_inventory_flows_deleted_from_status_updated_id", "deleted_at, from_user_id, flow_status, updated_at DESC, id DESC"),
          index("inventory_flows", "idx_inventory_flows_deleted_to_status_updated_id", "deleted_at, to_user_id, flow_status, updated_at DESC, id DESC"),
          index("inventory_flows", "idx_inventory_flows_deleted_name_updated_id", "deleted_at, name, updated_at DESC, id DESC"),
          unique("inventory_flow_items", "uk_inventory_flow_items_flow_item", "inventory_flow_id, item_id"),
          index("inventory_flow_items", "idx_inventory_flow_items_deleted_flow_item", "deleted_at, inventory_flow_id, item_id"),
          index("inventory_flow_items", "idx_inventory_flow_items_deleted_item_flow", "deleted_at, item_id, inventory_flow_id"),
          unique("inventory_flow_item_units", "uk_inventory_flow_item_units_flow_unit", "inventory_flow_id, item_unit_id"),
          index("inventory_flow_item_units", "idx_inventory_flow_item_units_unit_flow", "item_unit_id, inventory_flow_id"),
          index("work_order", "idx_work_order_deleted_from_updated_id", "deleted_at, from_user_id, updated_at DESC, id DESC"),
          index("work_order", "idx_work_order_deleted_to_updated_id", "deleted_at, to_user_id, updated_at DESC, id DESC"),
          index("work_order", "idx_work_order_deleted_from_status_updated_id", "deleted_at, from_user_id, status, updated_at DESC, id DESC"),
          index("work_order", "idx_work_order_deleted_to_status_updated_id", "deleted_at, to_user_id, status, updated_at DESC, id DESC"),
          index("work_order", "idx_work_order_deleted_to_read_updated_id", "deleted_at, to_user_id, read_status, updated_at DESC, id DESC"),
          index("work_order", "idx_work_order_deleted_to_status_read_updated_id", "deleted_at, to_user_id, status, read_status, updated_at DESC, id DESC"),
          index("work_order", "idx_work_order_deleted_name_updated_id", "deleted_at, name, updated_at DESC, id DESC"),
          index("work_order", "idx_work_order_deleted_from_name_updated_id", "deleted_at, from_user_id, name, updated_at DESC, id DESC"),
          index("work_order", "idx_work_order_deleted_to_name_updated_id", "deleted_at, to_user_id, name, updated_at DESC, id DESC"));

  private final JdbcTemplate jdbc;

  public DatabaseIndexInitializer(JdbcTemplate jdbc) {
    this.jdbc = jdbc;
  }

  @Override
  public void run(ApplicationArguments args) {
    ensureColumn(
        "inventory_flows",
        "business_type",
        "ALTER TABLE inventory_flows ADD COLUMN business_type INT NOT NULL DEFAULT 1 AFTER flow_type");
    jdbc.update(
        "UPDATE inventory_flows SET business_type=2 WHERE flow_type=2 AND business_type=1");
    for (Index index : INDEXES) {
      Integer count =
          jdbc.queryForObject(
              "SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema=DATABASE() AND table_name=? AND index_name=?",
              Integer.class,
              index.table(),
              index.name());
      if (count != null && count > 0) continue;
      jdbc.execute(index.sql());
    }
  }

  private void ensureColumn(String table, String column, String sql) {
    Integer count =
        jdbc.queryForObject(
            "SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name=? AND column_name=?",
            Integer.class,
            table,
            column);
    if (count == null || count == 0) jdbc.execute(sql);
  }

  private static Index index(String table, String name, String columns) {
    return new Index(table, name, false, columns);
  }

  private static Index unique(String table, String name, String columns) {
    return new Index(table, name, true, columns);
  }

  private record Index(String table, String name, boolean unique, String columns) {
    private String sql() {
      return "CREATE "
          + (unique ? "UNIQUE " : "")
          + "INDEX `"
          + name
          + "` ON `"
          + table
          + "` ("
          + columns
          + ")";
    }
  }
}
