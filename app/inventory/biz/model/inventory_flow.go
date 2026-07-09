package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type InventoryFlow struct {
	ID        uint `gorm:"primarykey;index:idx_inventory_flow_updated_id,priority:3;index:idx_inventory_flow_status_updated_id,priority:4;index:idx_inventory_flow_from_updated_id,priority:4;index:idx_inventory_flow_to_updated_id,priority:4;index:idx_inventory_flow_from_status_updated_id,priority:5;index:idx_inventory_flow_to_status_updated_id,priority:5;index:idx_inventory_flow_name_id,priority:3"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index:idx_inventory_flow_updated_id,priority:2;index:idx_inventory_flow_status_updated_id,priority:3;index:idx_inventory_flow_from_updated_id,priority:3;index:idx_inventory_flow_to_updated_id,priority:3;index:idx_inventory_flow_from_status_updated_id,priority:4;index:idx_inventory_flow_to_status_updated_id,priority:4"`
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_inventory_flow_updated_id,priority:1;index:idx_inventory_flow_status_updated_id,priority:1;index:idx_inventory_flow_from_updated_id,priority:1;index:idx_inventory_flow_to_updated_id,priority:1;index:idx_inventory_flow_from_status_updated_id,priority:1;index:idx_inventory_flow_to_status_updated_id,priority:1;index:idx_inventory_flow_name_id,priority:1"`

	FromUserID  int64      `gorm:"not null;index:idx_inventory_flow_from_updated_id,priority:2;index:idx_inventory_flow_from_status_updated_id,priority:2"`
	ToUserID    int64      `gorm:"not null;index:idx_inventory_flow_to_updated_id,priority:2;index:idx_inventory_flow_to_status_updated_id,priority:2"`
	FlowType    int32      `gorm:"not null"`
	FlowStatus  int32      `gorm:"not null;index:idx_inventory_flow_status_updated_id,priority:2;index:idx_inventory_flow_from_status_updated_id,priority:3;index:idx_inventory_flow_to_status_updated_id,priority:3"`
	Name        string     `gorm:"type:varchar(100);not null;default:'';index:idx_inventory_flow_name_id,priority:2,length:64"`
	Description string     `gorm:"type:varchar(255);not null;default:''"`
	ApprovedBy  int64      `gorm:"not null;default:0"`
	ApprovedAt  *time.Time `gorm:"default:null"`

	Items     []InventoryFlowItem `gorm:"foreignKey:InventoryFlowID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ItemUnits []ItemUnit          `gorm:"many2many:inventory_flow_item_units;"`
}

type InventoryFlowItem struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	InventoryFlowID  uint  `gorm:"not null;uniqueIndex:idx_inventory_flow_item,priority:1;index:idx_inventory_flow_item_reverse,priority:2"`
	ItemID           uint  `gorm:"not null;uniqueIndex:idx_inventory_flow_item,priority:2;index:idx_inventory_flow_item_reverse,priority:1"`
	Item             Item  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ApplyQuantity    int64 `gorm:"not null;default:0"`
	FinishedQuantity int64 `gorm:"not null;default:0"`
}

type InventoryFlowItemUnit struct {
	InventoryFlowID uint `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_inventory_flow_unit,priority:1;index:idx_inventory_flow_unit_reverse,priority:2"`
	ItemUnitID      uint `gorm:"primaryKey;autoIncrement:false;uniqueIndex:idx_inventory_flow_unit,priority:2;index:idx_inventory_flow_unit_reverse,priority:1"`
}

type InventoryFlowQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewInventoryFlowQuery(ctx context.Context, db *gorm.DB) *InventoryFlowQuery {
	return &InventoryFlowQuery{ctx: ctx, db: db}
}

func (q *InventoryFlowQuery) Get(id uint) (InventoryFlow, error) {
	var flow InventoryFlow
	err := q.db.WithContext(q.ctx).
		Preload("Items.Item").
		Preload("ItemUnits").
		First(&flow, id).Error
	return flow, err
}

func (q *InventoryFlowQuery) List(pageSize int, userID int64, isTo bool, filterUser bool, flowStatus int32, namePrefix string, itemNamePrefix string, sinceTime *time.Time, cursorUpdatedAt *time.Time, cursorID uint) ([]InventoryFlow, bool, error) {
	var flows []InventoryFlow
	db := q.db.WithContext(q.ctx).Model(&InventoryFlow{})
	if filterUser {
		if isTo {
			db = db.Where("inventory_flows.to_user_id = ?", userID)
		} else {
			db = db.Where("inventory_flows.from_user_id = ?", userID)
		}
	}
	if flowStatus > 0 {
		db = db.Where("inventory_flows.flow_status = ?", flowStatus)
	} else {
		db = db.Where("inventory_flows.flow_status <> ?", int32(1))
	}
	if namePrefix != "" {
		db = db.Where("inventory_flows.name LIKE ?", namePrefix+"%")
	}
	if itemNamePrefix != "" {
		itemNameLike := itemNamePrefix + "%"
		db = db.Where(`EXISTS (
			SELECT 1
			FROM inventory_flow_items
			JOIN items ON items.id = inventory_flow_items.item_id AND items.deleted_at IS NULL
			WHERE inventory_flow_items.inventory_flow_id = inventory_flows.id
				AND inventory_flow_items.deleted_at IS NULL
				AND items.name LIKE ?
		) OR EXISTS (
			SELECT 1
			FROM inventory_flow_item_units
			JOIN item_units ON item_units.id = inventory_flow_item_units.item_unit_id AND item_units.deleted_at IS NULL
			JOIN items ON items.id = item_units.item_id AND items.deleted_at IS NULL
			WHERE inventory_flow_item_units.inventory_flow_id = inventory_flows.id
				AND items.name LIKE ?
		)`, itemNameLike, itemNameLike)
	}
	if sinceTime != nil {
		db = db.Where("inventory_flows.updated_at > ?", *sinceTime)
	}
	if cursorUpdatedAt != nil && cursorID > 0 {
		db = db.Where("(inventory_flows.updated_at < ? OR (inventory_flows.updated_at = ? AND inventory_flows.id < ?))", *cursorUpdatedAt, *cursorUpdatedAt, cursorID)
	}
	err := db.Preload("Items.Item").
		Order("inventory_flows.updated_at DESC, inventory_flows.id DESC").
		Limit(pageSize + 1).
		Find(&flows).Error
	if err != nil {
		return nil, false, err
	}
	hasMore := len(flows) > pageSize
	if hasMore {
		flows = flows[:pageSize]
	}
	return flows, hasMore, nil
}
