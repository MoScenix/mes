package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ItemUnit struct {
	ID        uint `gorm:"primarykey;index:idx_item_unit_item_id,priority:3;index:idx_item_unit_engineering_id,priority:3;index:idx_item_unit_stock_quality_id,priority:4;index:idx_item_unit_item_stock_quality_id,priority:5;index:idx_item_unit_engineering_stock_quality_id,priority:5"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_item_unit_item_id,priority:1;index:idx_item_unit_engineering_id,priority:1;index:idx_item_unit_stock_quality_id,priority:1;index:idx_item_unit_item_stock_quality_id,priority:1;index:idx_item_unit_engineering_stock_quality_id,priority:1"`

	ItemID             uint              `gorm:"not null;index:idx_item_unit_item_id,priority:2;index:idx_item_unit_item_stock_quality_id,priority:2"`
	Item               Item              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	EngineeringOrderID *uint             `gorm:"index:idx_item_unit_engineering_id,priority:2;index:idx_item_unit_engineering_stock_quality_id,priority:2"`
	EngineeringOrder   *EngineeringOrder `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StockStatus        int32             `gorm:"not null;default:0;index:idx_item_unit_stock_quality_id,priority:2;index:idx_item_unit_item_stock_quality_id,priority:3;index:idx_item_unit_engineering_stock_quality_id,priority:3"`
	QualityStatus      int32             `gorm:"not null;default:0;index:idx_item_unit_stock_quality_id,priority:3;index:idx_item_unit_item_stock_quality_id,priority:4;index:idx_item_unit_engineering_stock_quality_id,priority:4"`
	Description        string            `gorm:"type:varchar(255);not null;default:''"`

	InventoryFlows []InventoryFlow `gorm:"many2many:inventory_flow_item_units;"`
}

type ItemUnitQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewItemUnitQuery(ctx context.Context, db *gorm.DB) *ItemUnitQuery {
	return &ItemUnitQuery{ctx: ctx, db: db}
}

func (q *ItemUnitQuery) Create(unit *ItemUnit) error {
	return q.db.WithContext(q.ctx).Create(unit).Error
}

func (q *ItemUnitQuery) UpdateStatus(id uint, stockStatus, qualityStatus int32) error {
	return q.db.WithContext(q.ctx).Model(&ItemUnit{}).Where("id = ?", id).Updates(map[string]any{
		"stock_status":   stockStatus,
		"quality_status": qualityStatus,
	}).Error
}

func (q *ItemUnitQuery) Get(id uint) (ItemUnit, error) {
	var unit ItemUnit
	err := q.db.WithContext(q.ctx).First(&unit, id).Error
	return unit, err
}

func (q *ItemUnitQuery) List(pageSize int, itemID uint, engineeringOrderID uint, stockStatus, qualityStatus int32, itemNamePrefix string, cursorID uint) ([]ItemUnit, bool, error) {
	var units []ItemUnit
	db := q.db.WithContext(q.ctx).Model(&ItemUnit{})
	if itemNamePrefix != "" {
		db = db.Joins("JOIN items ON items.id = item_units.item_id AND items.deleted_at IS NULL").
			Where("items.name LIKE ?", itemNamePrefix+"%")
	}
	if itemID > 0 {
		db = db.Where("item_units.item_id = ?", itemID)
	}
	if engineeringOrderID > 0 {
		db = db.Where("item_units.engineering_order_id = ?", engineeringOrderID)
	}
	if stockStatus > 0 {
		db = db.Where("item_units.stock_status = ?", stockStatus)
	}
	if qualityStatus > 0 {
		db = db.Where("item_units.quality_status = ?", qualityStatus)
	}
	if cursorID > 0 {
		db = db.Where("item_units.id < ?", cursorID)
	}
	err := db.Order("item_units.id DESC").Limit(pageSize + 1).Find(&units).Error
	if err != nil {
		return nil, false, err
	}
	hasMore := len(units) > pageSize
	if hasMore {
		units = units[:pageSize]
	}
	return units, hasMore, nil
}
