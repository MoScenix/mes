package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type EngineeringOrder struct {
	ID        uint `gorm:"primarykey;index:idx_engineering_order_leader_updated_id,priority:4;index:idx_engineering_order_leader_status_updated_id,priority:5;index:idx_engineering_order_process_updated_id,priority:4;index:idx_engineering_order_item_updated_id,priority:4;index:idx_engineering_order_status_updated_id,priority:4;index:idx_engineering_order_updated_id,priority:3;index:idx_engineering_order_process_status_updated_id,priority:5;index:idx_engineering_order_item_status_updated_id,priority:5;index:idx_engineering_order_leader_process_updated_id,priority:5;index:idx_engineering_order_leader_process_status_updated_id,priority:6;index:idx_engineering_order_leader_item_updated_id,priority:5;index:idx_engineering_order_leader_item_status_updated_id,priority:6;index:idx_engineering_order_name_id,priority:3;index:idx_engineering_order_leader_name_id,priority:4;index:idx_engineering_order_item_name_id,priority:4;index:idx_engineering_order_process_name_id,priority:4;index:idx_engineering_order_leader_status_name_id,priority:5;index:idx_engineering_order_item_status_name_id,priority:5;index:idx_engineering_order_process_status_name_id,priority:5;index:idx_engineering_order_item_process_updated_id,priority:5;index:idx_engineering_order_item_process_status_updated_id,priority:6"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index:idx_engineering_order_leader_updated_id,priority:3;index:idx_engineering_order_leader_status_updated_id,priority:4;index:idx_engineering_order_process_updated_id,priority:3;index:idx_engineering_order_item_updated_id,priority:3;index:idx_engineering_order_status_updated_id,priority:3;index:idx_engineering_order_updated_id,priority:2;index:idx_engineering_order_process_status_updated_id,priority:4;index:idx_engineering_order_item_status_updated_id,priority:4;index:idx_engineering_order_leader_process_updated_id,priority:4;index:idx_engineering_order_leader_process_status_updated_id,priority:5;index:idx_engineering_order_leader_item_updated_id,priority:4;index:idx_engineering_order_leader_item_status_updated_id,priority:5;index:idx_engineering_order_item_process_updated_id,priority:4;index:idx_engineering_order_item_process_status_updated_id,priority:5"`
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_engineering_order_leader_updated_id,priority:1;index:idx_engineering_order_leader_status_updated_id,priority:1;index:idx_engineering_order_process_updated_id,priority:1;index:idx_engineering_order_item_updated_id,priority:1;index:idx_engineering_order_status_updated_id,priority:1;index:idx_engineering_order_updated_id,priority:1;index:idx_engineering_order_process_status_updated_id,priority:1;index:idx_engineering_order_item_status_updated_id,priority:1;index:idx_engineering_order_leader_process_updated_id,priority:1;index:idx_engineering_order_leader_process_status_updated_id,priority:1;index:idx_engineering_order_leader_item_updated_id,priority:1;index:idx_engineering_order_leader_item_status_updated_id,priority:1;index:idx_engineering_order_name_id,priority:1;index:idx_engineering_order_leader_name_id,priority:1;index:idx_engineering_order_item_name_id,priority:1;index:idx_engineering_order_process_name_id,priority:1;index:idx_engineering_order_leader_status_name_id,priority:1;index:idx_engineering_order_item_status_name_id,priority:1;index:idx_engineering_order_process_status_name_id,priority:1;index:idx_engineering_order_item_process_updated_id,priority:1;index:idx_engineering_order_item_process_status_updated_id,priority:1"`

	LeaderUserID        int64   `gorm:"not null;index:idx_engineering_order_leader_updated_id,priority:2;index:idx_engineering_order_leader_status_updated_id,priority:2;index:idx_engineering_order_leader_process_updated_id,priority:2;index:idx_engineering_order_leader_process_status_updated_id,priority:2;index:idx_engineering_order_leader_item_updated_id,priority:2;index:idx_engineering_order_leader_item_status_updated_id,priority:2;index:idx_engineering_order_leader_name_id,priority:2;index:idx_engineering_order_leader_status_name_id,priority:2"`
	ProcessID           uint    `gorm:"not null;index:idx_engineering_order_process_updated_id,priority:2;index:idx_engineering_order_process_status_updated_id,priority:2;index:idx_engineering_order_leader_process_updated_id,priority:3;index:idx_engineering_order_leader_process_status_updated_id,priority:3;index:idx_engineering_order_process_name_id,priority:2;index:idx_engineering_order_process_status_name_id,priority:2;index:idx_engineering_order_item_process_updated_id,priority:3;index:idx_engineering_order_item_process_status_updated_id,priority:3"`
	Process             Process `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ItemID              uint    `gorm:"not null;index:idx_engineering_order_item_updated_id,priority:2;index:idx_engineering_order_item_status_updated_id,priority:2;index:idx_engineering_order_leader_item_updated_id,priority:3;index:idx_engineering_order_leader_item_status_updated_id,priority:3;index:idx_engineering_order_item_name_id,priority:2;index:idx_engineering_order_item_status_name_id,priority:2;index:idx_engineering_order_item_process_updated_id,priority:2;index:idx_engineering_order_item_process_status_updated_id,priority:2"`
	Item                Item    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name                string  `gorm:"type:varchar(100);not null;default:'';index:idx_engineering_order_name_id,priority:2,length:64;index:idx_engineering_order_leader_name_id,priority:3,length:64;index:idx_engineering_order_item_name_id,priority:3,length:64;index:idx_engineering_order_process_name_id,priority:3,length:64;index:idx_engineering_order_leader_status_name_id,priority:4,length:64;index:idx_engineering_order_item_status_name_id,priority:4,length:64;index:idx_engineering_order_process_status_name_id,priority:4,length:64"`
	ExpectedQuantity    int64   `gorm:"not null;default:0"`
	QualifiedQuantity   int64   `gorm:"not null;default:0"`
	UnqualifiedQuantity int64   `gorm:"not null;default:0"`
	ProducedQuantity    int64   `gorm:"not null;default:0"`
	Status              int32   `gorm:"not null;default:1;index:idx_engineering_order_leader_status_updated_id,priority:3;index:idx_engineering_order_status_updated_id,priority:2;index:idx_engineering_order_process_status_updated_id,priority:3;index:idx_engineering_order_item_status_updated_id,priority:3;index:idx_engineering_order_leader_process_status_updated_id,priority:4;index:idx_engineering_order_leader_item_status_updated_id,priority:4;index:idx_engineering_order_leader_status_name_id,priority:3;index:idx_engineering_order_item_status_name_id,priority:3;index:idx_engineering_order_process_status_name_id,priority:3;index:idx_engineering_order_item_process_status_updated_id,priority:4"`
	Description         string  `gorm:"type:varchar(255);not null;default:''"`

	ItemUnits []ItemUnit `gorm:"foreignKey:EngineeringOrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type EngineeringOrderQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewEngineeringOrderQuery(ctx context.Context, db *gorm.DB) *EngineeringOrderQuery {
	return &EngineeringOrderQuery{ctx: ctx, db: db}
}

func (q *EngineeringOrderQuery) Create(order *EngineeringOrder) error {
	return q.db.WithContext(q.ctx).Create(order).Error
}

func (q *EngineeringOrderQuery) Update(id uint, updates map[string]any) error {
	return q.db.WithContext(q.ctx).Model(&EngineeringOrder{}).Where("id = ?", id).Updates(updates).Error
}

func (q *EngineeringOrderQuery) Get(id uint, withUnits bool) (EngineeringOrder, error) {
	var order EngineeringOrder
	db := q.db.WithContext(q.ctx).Preload("Item").Preload("Process")
	if withUnits {
		db = db.Preload("ItemUnits")
	}
	err := db.First(&order, id).Error
	return order, err
}

func (q *EngineeringOrderQuery) List(pageSize int, leaderUserID int64, itemID uint, processID uint, status int32, namePrefix string, sinceTime *time.Time, cursorUpdatedAt *time.Time, cursorID uint) ([]EngineeringOrder, bool, error) {
	var orders []EngineeringOrder
	db := q.db.WithContext(q.ctx).Model(&EngineeringOrder{})
	if leaderUserID > 0 {
		db = db.Where("leader_user_id = ?", leaderUserID)
	}
	if itemID > 0 {
		db = db.Where("item_id = ?", itemID)
	}
	if processID > 0 {
		db = db.Where("process_id = ?", processID)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
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
	err := db.Preload("Item").Preload("Process").
		Order("updated_at DESC, id DESC").
		Limit(pageSize + 1).
		Find(&orders).Error
	if err != nil {
		return nil, false, err
	}
	hasMore := len(orders) > pageSize
	if hasMore {
		orders = orders[:pageSize]
	}
	return orders, hasMore, nil
}
