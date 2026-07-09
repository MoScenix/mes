package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type EngineeringOrder struct {
	ID        uint `gorm:"primarykey;index:idx_engineering_order_updated_id,priority:3;index:idx_engineering_order_leader_updated_id,priority:4;index:idx_engineering_order_leader_status_updated_id,priority:5;index:idx_engineering_order_status_updated_id,priority:4;index:idx_engineering_order_item_updated_id,priority:4;index:idx_engineering_order_item_status_updated_id,priority:5;index:idx_engineering_order_process_updated_id,priority:4;index:idx_engineering_order_process_status_updated_id,priority:5;index:idx_engineering_order_item_process_updated_id,priority:5;index:idx_engineering_order_name_id,priority:3"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index:idx_engineering_order_updated_id,priority:2;index:idx_engineering_order_leader_updated_id,priority:3;index:idx_engineering_order_leader_status_updated_id,priority:4;index:idx_engineering_order_status_updated_id,priority:3;index:idx_engineering_order_item_updated_id,priority:3;index:idx_engineering_order_item_status_updated_id,priority:4;index:idx_engineering_order_process_updated_id,priority:3;index:idx_engineering_order_process_status_updated_id,priority:4;index:idx_engineering_order_item_process_updated_id,priority:4"`
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_engineering_order_updated_id,priority:1;index:idx_engineering_order_leader_updated_id,priority:1;index:idx_engineering_order_leader_status_updated_id,priority:1;index:idx_engineering_order_status_updated_id,priority:1;index:idx_engineering_order_item_updated_id,priority:1;index:idx_engineering_order_item_status_updated_id,priority:1;index:idx_engineering_order_process_updated_id,priority:1;index:idx_engineering_order_process_status_updated_id,priority:1;index:idx_engineering_order_item_process_updated_id,priority:1;index:idx_engineering_order_name_id,priority:1"`

	LeaderUserID        int64   `gorm:"not null;index:idx_engineering_order_leader_updated_id,priority:2;index:idx_engineering_order_leader_status_updated_id,priority:2"`
	ProcessID           uint    `gorm:"not null;index:idx_engineering_order_process_updated_id,priority:2;index:idx_engineering_order_process_status_updated_id,priority:2;index:idx_engineering_order_item_process_updated_id,priority:3"`
	Process             Process `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ItemID              uint    `gorm:"not null;index:idx_engineering_order_item_updated_id,priority:2;index:idx_engineering_order_item_status_updated_id,priority:2;index:idx_engineering_order_item_process_updated_id,priority:2"`
	Item                Item    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name                string  `gorm:"type:varchar(100);not null;default:'';index:idx_engineering_order_name_id,priority:2,length:64"`
	ExpectedQuantity    int64   `gorm:"not null;default:0"`
	QualifiedQuantity   int64   `gorm:"not null;default:0"`
	UnqualifiedQuantity int64   `gorm:"not null;default:0"`
	ProducedQuantity    int64   `gorm:"not null;default:0"`
	Status              int32   `gorm:"not null;default:1;index:idx_engineering_order_leader_status_updated_id,priority:3;index:idx_engineering_order_status_updated_id,priority:2;index:idx_engineering_order_process_status_updated_id,priority:3;index:idx_engineering_order_item_status_updated_id,priority:3"`
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

func (q *EngineeringOrderQuery) List(pageSize int, leaderUserID int64, itemID uint, processID uint, status int32, namePrefix string, itemNamePrefix string, sinceTime *time.Time, cursorUpdatedAt *time.Time, cursorID uint) ([]EngineeringOrder, bool, error) {
	var orders []EngineeringOrder
	db := q.db.WithContext(q.ctx).Model(&EngineeringOrder{})
	if itemNamePrefix != "" {
		db = db.Joins("JOIN items ON items.id = engineering_orders.item_id AND items.deleted_at IS NULL").
			Where("items.name LIKE ?", itemNamePrefix+"%")
	}
	if leaderUserID > 0 {
		db = db.Where("engineering_orders.leader_user_id = ?", leaderUserID)
	}
	if itemID > 0 {
		db = db.Where("engineering_orders.item_id = ?", itemID)
	}
	if processID > 0 {
		db = db.Where("engineering_orders.process_id = ?", processID)
	}
	if status > 0 {
		db = db.Where("engineering_orders.status = ?", status)
	}
	if namePrefix != "" {
		db = db.Where("engineering_orders.name LIKE ?", namePrefix+"%")
	}
	if sinceTime != nil {
		db = db.Where("engineering_orders.updated_at > ?", *sinceTime)
	}
	if cursorUpdatedAt != nil && cursorID > 0 {
		db = db.Where("(engineering_orders.updated_at < ? OR (engineering_orders.updated_at = ? AND engineering_orders.id < ?))", *cursorUpdatedAt, *cursorUpdatedAt, cursorID)
	}
	err := db.Preload("Item").Preload("Process").
		Order("engineering_orders.updated_at DESC, engineering_orders.id DESC").
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
