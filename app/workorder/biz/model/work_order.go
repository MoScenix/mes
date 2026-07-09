package model

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	WorkOrderStatusDraft     int32 = 1
	WorkOrderStatusSubmitted int32 = 2

	WorkOrderReadStatusUnread int32 = 1
	WorkOrderReadStatusRead   int32 = 2
)

var ErrDraftRequired = errors.New("work order is not a draft or does not exist")

type WorkOrder struct {
	ID        uint `gorm:"primarykey;index:idx_from_user_updated_id,priority:4;index:idx_to_user_updated_id,priority:4;index:idx_from_user_status_updated_id,priority:5;index:idx_to_user_status_updated_id,priority:5;index:idx_from_user_read_status_updated_id,priority:5;index:idx_to_user_read_status_updated_id,priority:5;index:idx_work_order_name_id,priority:3;index:idx_work_order_from_name_id,priority:4;index:idx_work_order_to_name_id,priority:4;index:idx_from_user_status_read_updated_id,priority:6;index:idx_to_user_status_read_updated_id,priority:6;index:idx_work_order_from_status_name_id,priority:5;index:idx_work_order_to_status_name_id,priority:5"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index:idx_from_user_updated_id,priority:3;index:idx_to_user_updated_id,priority:3;index:idx_from_user_status_updated_id,priority:4;index:idx_to_user_status_updated_id,priority:4;index:idx_from_user_read_status_updated_id,priority:4;index:idx_to_user_read_status_updated_id,priority:4;index:idx_from_user_status_read_updated_id,priority:5;index:idx_to_user_status_read_updated_id,priority:5"`
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_from_user_updated_id,priority:1;index:idx_to_user_updated_id,priority:1;index:idx_from_user_status_updated_id,priority:1;index:idx_to_user_status_updated_id,priority:1;index:idx_from_user_read_status_updated_id,priority:1;index:idx_to_user_read_status_updated_id,priority:1;index:idx_work_order_name_id,priority:1;index:idx_work_order_from_name_id,priority:1;index:idx_work_order_to_name_id,priority:1;index:idx_from_user_status_read_updated_id,priority:1;index:idx_to_user_status_read_updated_id,priority:1;index:idx_work_order_from_status_name_id,priority:1;index:idx_work_order_to_status_name_id,priority:1"`

	FromUserID  int64  `gorm:"column:from_user_id;not null;index:idx_from_user_id;index:idx_from_user_updated_id,priority:2;index:idx_from_user_status_updated_id,priority:2;index:idx_from_user_read_status_updated_id,priority:2;index:idx_work_order_from_name_id,priority:2;index:idx_from_user_status_read_updated_id,priority:2;index:idx_work_order_from_status_name_id,priority:2"`
	ToUserID    int64  `gorm:"column:to_user_id;not null;index:idx_to_user_id;index:idx_to_user_updated_id,priority:2;index:idx_to_user_status_updated_id,priority:2;index:idx_to_user_read_status_updated_id,priority:2;index:idx_work_order_to_name_id,priority:2;index:idx_to_user_status_read_updated_id,priority:2;index:idx_work_order_to_status_name_id,priority:2"`
	Name        string `gorm:"type:varchar(100);not null;default:'';index:idx_work_order_name_id,priority:2,length:64;index:idx_work_order_from_name_id,priority:3,length:64;index:idx_work_order_to_name_id,priority:3,length:64;index:idx_work_order_from_status_name_id,priority:4,length:64;index:idx_work_order_to_status_name_id,priority:4,length:64"`
	Description string `gorm:"type:text;not null"`
	Status      int32  `gorm:"not null;index:idx_from_user_status_updated_id,priority:3;index:idx_to_user_status_updated_id,priority:3;index:idx_from_user_status_read_updated_id,priority:3;index:idx_to_user_status_read_updated_id,priority:3;index:idx_work_order_from_status_name_id,priority:3;index:idx_work_order_to_status_name_id,priority:3"`
	ReadStatus  int32  `gorm:"column:read_status;not null;default:1;index:idx_from_user_read_status_updated_id,priority:3;index:idx_to_user_read_status_updated_id,priority:3;index:idx_from_user_status_read_updated_id,priority:4;index:idx_to_user_status_read_updated_id,priority:4"`
}

func (WorkOrder) TableName() string {
	return "work_order"
}

type WorkOrderQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewWorkOrderQuery(ctx context.Context, db *gorm.DB) *WorkOrderQuery {
	return &WorkOrderQuery{
		ctx: ctx,
		db:  db,
	}
}

func (q *WorkOrderQuery) CreateWorkOrder(order WorkOrder) (WorkOrder, error) {
	err := q.db.WithContext(q.ctx).Model(&WorkOrder{}).Create(&order).Error
	return order, err
}

func (q *WorkOrderQuery) GetWorkOrder(id int64) (WorkOrder, error) {
	order := WorkOrder{}
	err := q.db.WithContext(q.ctx).Model(&WorkOrder{}).Where("id = ?", id).First(&order).Error
	return order, err
}

func (q *WorkOrderQuery) UpdateDraft(id int64, fromUserID int64, toUserID int64, name string, description string) error {
	result := q.db.WithContext(q.ctx).Model(&WorkOrder{}).
		Where("id = ? AND status = ?", id, WorkOrderStatusDraft).
		Updates(map[string]interface{}{
			"from_user_id": fromUserID,
			"to_user_id":   toUserID,
			"name":         name,
			"description":  description,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDraftRequired
	}
	return nil
}

func (q *WorkOrderQuery) DeleteDraft(id int64) error {
	result := q.db.WithContext(q.ctx).Where("id = ? AND status = ?", id, WorkOrderStatusDraft).Delete(&WorkOrder{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDraftRequired
	}
	return nil
}

func (q *WorkOrderQuery) SubmitDraft(id int64) error {
	result := q.db.WithContext(q.ctx).Model(&WorkOrder{}).
		Where("id = ? AND status = ?", id, WorkOrderStatusDraft).
		Update("status", WorkOrderStatusSubmitted)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDraftRequired
	}
	return nil
}

func (q *WorkOrderQuery) MarkRead(id int64) error {
	result := q.db.WithContext(q.ctx).Model(&WorkOrder{}).
		Where("id = ?", id).
		Update("read_status", WorkOrderReadStatusRead)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (q *WorkOrderQuery) ListWorkOrderByEmployee(employeeID int64, pageSize int64, isTo bool, onlyUnread bool, status int32, namePrefix string, sinceTime *time.Time, cursorUpdatedAt *time.Time, cursorID int64) ([]WorkOrder, bool, error) {
	db := q.db.WithContext(q.ctx).Model(&WorkOrder{})
	if isTo {
		db = db.Where("to_user_id = ?", employeeID)
	} else {
		db = db.Where("from_user_id = ?", employeeID)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	} else if isTo {
		db = db.Where("status <> ?", WorkOrderStatusDraft)
	}
	if onlyUnread {
		db = db.Where("read_status = ?", WorkOrderReadStatusUnread)
	}
	if namePrefix != "" {
		db = db.Where("name LIKE ?", namePrefix+"%")
	}
	if sinceTime != nil {
		db = db.Where("updated_at > ?", *sinceTime)
	}
	if cursorUpdatedAt != nil && cursorID > 0 {
		db = db.Where("(updated_at < ? OR (updated_at = ? AND id < ?))", *cursorUpdatedAt, *cursorUpdatedAt, cursorID)
	}

	var orders []WorkOrder
	err := db.Order("updated_at DESC, id DESC").
		Limit(int(pageSize + 1)).
		Find(&orders).Error
	if err != nil {
		return nil, false, err
	}
	hasMore := int64(len(orders)) > pageSize
	if hasMore {
		orders = orders[:pageSize]
	}
	return orders, hasMore, nil
}
