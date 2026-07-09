CREATE INDEX idx_items_deleted_updated_id ON items (deleted_at, updated_at DESC, id DESC);
CREATE INDEX idx_items_deleted_name_updated_id ON items (deleted_at, name, updated_at DESC, id DESC);

CREATE INDEX idx_item_units_deleted_updated_id ON item_units (deleted_at, updated_at DESC, id DESC);
CREATE INDEX idx_item_units_deleted_item_updated_id ON item_units (deleted_at, item_id, updated_at DESC, id DESC);
CREATE INDEX idx_item_units_deleted_engineering_updated_id ON item_units (deleted_at, engineering_order_id, updated_at DESC, id DESC);
CREATE INDEX idx_item_units_deleted_stock_quality_updated_id ON item_units (deleted_at, stock_status, quality_status, updated_at DESC, id DESC);
CREATE INDEX idx_item_units_deleted_item_stock_quality_updated_id ON item_units (deleted_at, item_id, stock_status, quality_status, updated_at DESC, id DESC);
CREATE INDEX idx_item_units_deleted_engineering_quality ON item_units (deleted_at, engineering_order_id, quality_status);

CREATE INDEX idx_processes_deleted_updated_id ON processes (deleted_at, updated_at DESC, id DESC);
CREATE INDEX idx_processes_deleted_owner_updated_id ON processes (deleted_at, owner_user_id, updated_at DESC, id DESC);
CREATE INDEX idx_processes_deleted_owner_status_updated_id ON processes (deleted_at, owner_user_id, status, updated_at DESC, id DESC);
CREATE INDEX idx_processes_deleted_item_updated_id ON processes (deleted_at, item_id, updated_at DESC, id DESC);
CREATE INDEX idx_processes_deleted_status_updated_id ON processes (deleted_at, status, updated_at DESC, id DESC);
CREATE INDEX idx_processes_deleted_name_updated_id ON processes (deleted_at, name, updated_at DESC, id DESC);

CREATE UNIQUE INDEX uk_process_items_process_consume_item ON process_items (process_id, consume_item_id);
CREATE INDEX idx_process_items_deleted_process ON process_items (deleted_at, process_id);
CREATE INDEX idx_process_items_deleted_consume_process ON process_items (deleted_at, consume_item_id, process_id);

CREATE INDEX idx_engineering_orders_deleted_updated_id ON engineering_orders (deleted_at, updated_at DESC, id DESC);
CREATE INDEX idx_engineering_orders_deleted_leader_updated_id ON engineering_orders (deleted_at, leader_user_id, updated_at DESC, id DESC);
CREATE INDEX idx_engineering_orders_deleted_leader_status_updated_id ON engineering_orders (deleted_at, leader_user_id, status, updated_at DESC, id DESC);
CREATE INDEX idx_engineering_orders_deleted_item_updated_id ON engineering_orders (deleted_at, item_id, updated_at DESC, id DESC);
CREATE INDEX idx_engineering_orders_deleted_process_updated_id ON engineering_orders (deleted_at, process_id, updated_at DESC, id DESC);
CREATE INDEX idx_engineering_orders_deleted_process_status_updated_id ON engineering_orders (deleted_at, process_id, status, updated_at DESC, id DESC);
CREATE INDEX idx_engineering_orders_deleted_status_updated_id ON engineering_orders (deleted_at, status, updated_at DESC, id DESC);
CREATE INDEX idx_engineering_orders_deleted_name_updated_id ON engineering_orders (deleted_at, name, updated_at DESC, id DESC);

CREATE INDEX idx_inventory_flows_deleted_updated_id ON inventory_flows (deleted_at, updated_at DESC, id DESC);
CREATE INDEX idx_inventory_flows_deleted_status_updated_id ON inventory_flows (deleted_at, flow_status, updated_at DESC, id DESC);
CREATE INDEX idx_inventory_flows_deleted_from_updated_id ON inventory_flows (deleted_at, from_user_id, updated_at DESC, id DESC);
CREATE INDEX idx_inventory_flows_deleted_to_updated_id ON inventory_flows (deleted_at, to_user_id, updated_at DESC, id DESC);
CREATE INDEX idx_inventory_flows_deleted_from_status_updated_id ON inventory_flows (deleted_at, from_user_id, flow_status, updated_at DESC, id DESC);
CREATE INDEX idx_inventory_flows_deleted_to_status_updated_id ON inventory_flows (deleted_at, to_user_id, flow_status, updated_at DESC, id DESC);
CREATE INDEX idx_inventory_flows_deleted_name_updated_id ON inventory_flows (deleted_at, name, updated_at DESC, id DESC);

CREATE UNIQUE INDEX uk_inventory_flow_items_flow_item ON inventory_flow_items (inventory_flow_id, item_id);
CREATE INDEX idx_inventory_flow_items_deleted_flow_item ON inventory_flow_items (deleted_at, inventory_flow_id, item_id);
CREATE INDEX idx_inventory_flow_items_deleted_item_flow ON inventory_flow_items (deleted_at, item_id, inventory_flow_id);

CREATE UNIQUE INDEX uk_inventory_flow_item_units_flow_unit ON inventory_flow_item_units (inventory_flow_id, item_unit_id);
CREATE INDEX idx_inventory_flow_item_units_unit_flow ON inventory_flow_item_units (item_unit_id, inventory_flow_id);
