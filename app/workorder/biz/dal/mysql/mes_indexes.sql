CREATE INDEX idx_work_order_deleted_from_updated_id ON work_order (deleted_at, from_user_id, updated_at DESC, id DESC);
CREATE INDEX idx_work_order_deleted_to_updated_id ON work_order (deleted_at, to_user_id, updated_at DESC, id DESC);
CREATE INDEX idx_work_order_deleted_from_status_updated_id ON work_order (deleted_at, from_user_id, status, updated_at DESC, id DESC);
CREATE INDEX idx_work_order_deleted_to_status_updated_id ON work_order (deleted_at, to_user_id, status, updated_at DESC, id DESC);
CREATE INDEX idx_work_order_deleted_to_read_updated_id ON work_order (deleted_at, to_user_id, read_status, updated_at DESC, id DESC);
CREATE INDEX idx_work_order_deleted_to_status_read_updated_id ON work_order (deleted_at, to_user_id, status, read_status, updated_at DESC, id DESC);
CREATE INDEX idx_work_order_deleted_name_updated_id ON work_order (deleted_at, name, updated_at DESC, id DESC);
CREATE INDEX idx_work_order_deleted_from_name_updated_id ON work_order (deleted_at, from_user_id, name, updated_at DESC, id DESC);
CREATE INDEX idx_work_order_deleted_to_name_updated_id ON work_order (deleted_at, to_user_id, name, updated_at DESC, id DESC);
