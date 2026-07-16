CREATE TABLE IF NOT EXISTS `user` (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  name VARCHAR(50) NOT NULL,
  password_hash VARCHAR(100) NOT NULL,
  user_account VARCHAR(50) NOT NULL,
  user_avatar VARCHAR(100) NOT NULL DEFAULT '',
  user_profile VARCHAR(100) NULL DEFAULT '',
  user_role VARCHAR(32) NOT NULL DEFAULT 'worker',
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_account (user_account),
  KEY idx_user_deleted_at (deleted_at),
  KEY idx_user_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  name VARCHAR(100) NOT NULL,
  unit VARCHAR(20) NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  total_count BIGINT NOT NULL DEFAULT 0,
  in_stock_count BIGINT NOT NULL DEFAULT 0,
  reserved_count BIGINT NOT NULL DEFAULT 0,
  out_stock_count BIGINT NOT NULL DEFAULT 0,
  pending_count BIGINT NOT NULL DEFAULT 0,
  qualified_count BIGINT NOT NULL DEFAULT 0,
  unqualified_count BIGINT NOT NULL DEFAULT 0,
  available_count BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_item_name_prefix (name(64)),
  KEY idx_items_deleted_updated_id (deleted_at, updated_at DESC, id DESC),
  KEY idx_items_deleted_name_updated_id (deleted_at, name, updated_at DESC, id DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS processes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  item_id BIGINT UNSIGNED NOT NULL,
  owner_user_id BIGINT NOT NULL,
  name VARCHAR(100) NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  status INT NOT NULL DEFAULT 1,
  PRIMARY KEY (id),
  KEY idx_processes_deleted_updated_id (deleted_at, updated_at DESC, id DESC),
  KEY idx_processes_deleted_owner_updated_id (deleted_at, owner_user_id, updated_at DESC, id DESC),
  KEY idx_processes_deleted_owner_status_updated_id (deleted_at, owner_user_id, status, updated_at DESC, id DESC),
  KEY idx_processes_deleted_item_updated_id (deleted_at, item_id, updated_at DESC, id DESC),
  KEY idx_processes_deleted_status_updated_id (deleted_at, status, updated_at DESC, id DESC),
  KEY idx_processes_deleted_name_updated_id (deleted_at, name, updated_at DESC, id DESC),
  CONSTRAINT fk_processes_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS process_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  process_id BIGINT UNSIGNED NOT NULL,
  consume_item_id BIGINT UNSIGNED NOT NULL,
  quantity BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_process_items_process_consume_item (process_id, consume_item_id),
  KEY idx_process_items_deleted_process (deleted_at, process_id),
  KEY idx_process_items_deleted_consume_process (deleted_at, consume_item_id, process_id),
  KEY idx_process_items_consume_item_fk (consume_item_id),
  CONSTRAINT fk_processes_items FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_process_items_consume_item FOREIGN KEY (consume_item_id) REFERENCES items(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS engineering_orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  leader_user_id BIGINT NOT NULL,
  item_id BIGINT UNSIGNED NOT NULL,
  expected_quantity BIGINT NOT NULL DEFAULT 0,
  qualified_quantity BIGINT NOT NULL DEFAULT 0,
  produced_quantity BIGINT NOT NULL DEFAULT 0,
  description VARCHAR(255) NOT NULL DEFAULT '',
  process_id BIGINT UNSIGNED NOT NULL,
  unqualified_quantity BIGINT NOT NULL DEFAULT 0,
  status INT NOT NULL DEFAULT 1,
  name VARCHAR(100) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_engineering_orders_deleted_updated_id (deleted_at, updated_at DESC, id DESC),
  KEY idx_engineering_orders_deleted_leader_updated_id (deleted_at, leader_user_id, updated_at DESC, id DESC),
  KEY idx_engineering_orders_deleted_leader_status_updated_id (deleted_at, leader_user_id, status, updated_at DESC, id DESC),
  KEY idx_engineering_orders_deleted_item_updated_id (deleted_at, item_id, updated_at DESC, id DESC),
  KEY idx_engineering_orders_deleted_process_updated_id (deleted_at, process_id, updated_at DESC, id DESC),
  KEY idx_engineering_orders_deleted_process_status_updated_id (deleted_at, process_id, status, updated_at DESC, id DESC),
  KEY idx_engineering_orders_deleted_status_updated_id (deleted_at, status, updated_at DESC, id DESC),
  KEY idx_engineering_orders_deleted_name_updated_id (deleted_at, name, updated_at DESC, id DESC),
  CONSTRAINT fk_engineering_orders_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_engineering_orders_process FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS item_units (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  item_id BIGINT UNSIGNED NOT NULL,
  engineering_order_id BIGINT UNSIGNED NULL,
  stock_status INT NOT NULL DEFAULT 0,
  quality_status INT NOT NULL DEFAULT 0,
  description VARCHAR(255) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_item_units_deleted_updated_id (deleted_at, updated_at DESC, id DESC),
  KEY idx_item_units_deleted_item_updated_id (deleted_at, item_id, updated_at DESC, id DESC),
  KEY idx_item_units_deleted_engineering_updated_id (deleted_at, engineering_order_id, updated_at DESC, id DESC),
  KEY idx_item_units_deleted_stock_quality_updated_id (deleted_at, stock_status, quality_status, updated_at DESC, id DESC),
  KEY idx_item_units_deleted_item_stock_quality_updated_id (deleted_at, item_id, stock_status, quality_status, updated_at DESC, id DESC),
  KEY idx_item_units_deleted_engineering_quality (deleted_at, engineering_order_id, quality_status),
  CONSTRAINT fk_item_units_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_engineering_orders_item_units FOREIGN KEY (engineering_order_id) REFERENCES engineering_orders(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS inventory_flows (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  from_user_id BIGINT NOT NULL,
  to_user_id BIGINT NOT NULL,
  flow_type INT NOT NULL,
  business_type INT NOT NULL DEFAULT 1,
  flow_status INT NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  approved_by BIGINT NOT NULL DEFAULT 0,
  approved_at DATETIME(3) NULL,
  name VARCHAR(100) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_inventory_flows_deleted_updated_id (deleted_at, updated_at DESC, id DESC),
  KEY idx_inventory_flows_deleted_status_updated_id (deleted_at, flow_status, updated_at DESC, id DESC),
  KEY idx_inventory_flows_deleted_business_status_updated_id (deleted_at, business_type, flow_status, updated_at DESC, id DESC),
  KEY idx_inventory_flows_deleted_from_updated_id (deleted_at, from_user_id, updated_at DESC, id DESC),
  KEY idx_inventory_flows_deleted_to_updated_id (deleted_at, to_user_id, updated_at DESC, id DESC),
  KEY idx_inventory_flows_deleted_from_status_updated_id (deleted_at, from_user_id, flow_status, updated_at DESC, id DESC),
  KEY idx_inventory_flows_deleted_to_status_updated_id (deleted_at, to_user_id, flow_status, updated_at DESC, id DESC),
  KEY idx_inventory_flows_deleted_name_updated_id (deleted_at, name, updated_at DESC, id DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS inventory_flow_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  inventory_flow_id BIGINT UNSIGNED NOT NULL,
  item_id BIGINT UNSIGNED NOT NULL,
  apply_quantity BIGINT NOT NULL DEFAULT 0,
  finished_quantity BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_inventory_flow_items_flow_item (inventory_flow_id, item_id),
  KEY idx_inventory_flow_items_item_id (item_id),
  KEY idx_inventory_flow_items_deleted_flow_item (deleted_at, inventory_flow_id, item_id),
  KEY idx_inventory_flow_items_deleted_item_flow (deleted_at, item_id, inventory_flow_id),
  CONSTRAINT fk_inventory_flows_items FOREIGN KEY (inventory_flow_id) REFERENCES inventory_flows(id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_inventory_flow_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS inventory_flow_item_units (
  inventory_flow_id BIGINT UNSIGNED NOT NULL,
  item_unit_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (inventory_flow_id, item_unit_id),
  UNIQUE KEY uk_inventory_flow_item_units_flow_unit (inventory_flow_id, item_unit_id),
  KEY idx_inventory_flow_item_units_unit_flow (item_unit_id, inventory_flow_id),
  CONSTRAINT fk_inventory_flow_item_units_inventory_flow FOREIGN KEY (inventory_flow_id) REFERENCES inventory_flows(id),
  CONSTRAINT fk_inventory_flow_item_units_item_unit FOREIGN KEY (item_unit_id) REFERENCES item_units(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS work_order (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  from_user_id BIGINT NOT NULL,
  to_user_id BIGINT NOT NULL,
  description TEXT NOT NULL,
  status INT NOT NULL,
  read_status INT NOT NULL DEFAULT 1,
  name VARCHAR(100) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_work_order_deleted_from_updated_id (deleted_at, from_user_id, updated_at DESC, id DESC),
  KEY idx_work_order_deleted_to_updated_id (deleted_at, to_user_id, updated_at DESC, id DESC),
  KEY idx_work_order_deleted_from_status_updated_id (deleted_at, from_user_id, status, updated_at DESC, id DESC),
  KEY idx_work_order_deleted_to_status_updated_id (deleted_at, to_user_id, status, updated_at DESC, id DESC),
  KEY idx_work_order_deleted_to_read_updated_id (deleted_at, to_user_id, read_status, updated_at DESC, id DESC),
  KEY idx_work_order_deleted_to_status_read_updated_id (deleted_at, to_user_id, status, read_status, updated_at DESC, id DESC),
  KEY idx_work_order_deleted_name_updated_id (deleted_at, name, updated_at DESC, id DESC),
  KEY idx_work_order_deleted_from_name_updated_id (deleted_at, from_user_id, name, updated_at DESC, id DESC),
  KEY idx_work_order_deleted_to_name_updated_id (deleted_at, to_user_id, name, updated_at DESC, id DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS history (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  name VARCHAR(100) NULL,
  user_id BIGINT NULL,
  PRIMARY KEY (id),
  KEY idx_history_user_updated (user_id, updated_at, id),
  KEY idx_history_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS messages (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  deleted_at DATETIME(3) NULL,
  content TEXT NULL,
  role VARCHAR(50) NULL,
  is_file TINYINT(1) NULL DEFAULT 0,
  history_id BIGINT UNSIGNED NULL,
  user_id BIGINT NULL,
  PRIMARY KEY (id),
  KEY idx_messages_deleted_at (deleted_at),
  KEY idx_messages_history_time (history_id, created_at, id),
  KEY idx_messages_user_time (user_id, created_at, id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
