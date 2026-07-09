package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type EngineeringOrder struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	LeaderUserID        int64   `gorm:"not null"`
	ProcessID           uint    `gorm:"not null"`
	Process             Process `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ItemID              uint    `gorm:"not null"`
	Item                Item    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name                string  `gorm:"type:varchar(100);not null;default:''"`
	ExpectedQuantity    int64   `gorm:"not null;default:0"`
	QualifiedQuantity   int64   `gorm:"not null;default:0"`
	UnqualifiedQuantity int64   `gorm:"not null;default:0"`
	ProducedQuantity    int64   `gorm:"not null;default:0"`
	Status              int32   `gorm:"not null;default:1"`
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
