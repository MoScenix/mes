package model

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

const (
	WorkOrderStatusDraft     int32 = 1
	WorkOrderStatusSubmitted int32 = 2
)

var ErrDraftRequired = errors.New("work order is not a draft or does not exist")

type WorkOrder struct {
	gorm.Model

	FromUserID  int64  `gorm:"column:from_user_id;not null;index;index:idx_from_user_status"`
	ToUserID    int64  `gorm:"column:to_user_id;not null;index;index:idx_to_user_status"`
	Description string `gorm:"type:text;not null"`
	Status      int32  `gorm:"not null;index;index:idx_from_user_status;index:idx_to_user_status"`
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

func (q *WorkOrderQuery) UpdateDraft(id int64, fromUserID int64, toUserID int64, description string) error {
	result := q.db.WithContext(q.ctx).Model(&WorkOrder{}).
		Where("id = ? AND status = ?", id, WorkOrderStatusDraft).
		Updates(map[string]interface{}{
			"from_user_id": fromUserID,
			"to_user_id":   toUserID,
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

type WorkOrderListFilter struct {
	FromUserID int64
	ToUserID   int64
	Status     int32
}

func (q *WorkOrderQuery) ListWorkOrder(pageNum int64, pageSize int64, filter WorkOrderListFilter) ([]WorkOrder, int64, error) {
	db := q.db.WithContext(q.ctx).Model(&WorkOrder{})
	if filter.FromUserID > 0 {
		db = db.Where("from_user_id = ?", filter.FromUserID)
	}
	if filter.ToUserID > 0 {
		db = db.Where("to_user_id = ?", filter.ToUserID)
	}
	if filter.Status > 0 {
		db = db.Where("status = ?", filter.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []WorkOrder
	err := db.Order("id DESC").
		Limit(int(pageSize)).
		Offset(int((pageNum - 1) * pageSize)).
		Find(&orders).Error
	return orders, total, err
}
