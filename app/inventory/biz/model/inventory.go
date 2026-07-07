package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model

	Name        string `gorm:"type:varchar(100);not null;index:idx_item_name_prefix,length:64"`
	Unit        string `gorm:"type:varchar(20);not null"`
	Description string `gorm:"type:varchar(255);not null;default:''"`

	TotalCount       int64 `gorm:"not null;default:0"`
	InStockCount     int64 `gorm:"not null;default:0"`
	ReservedCount    int64 `gorm:"not null;default:0"`
	OutStockCount    int64 `gorm:"not null;default:0"`
	PendingCount     int64 `gorm:"not null;default:0"`
	QualifiedCount   int64 `gorm:"not null;default:0"`
	UnqualifiedCount int64 `gorm:"not null;default:0"`
	AvailableCount   int64 `gorm:"not null;default:0"`
}

type ItemUnit struct {
	gorm.Model

	ItemID        uint   `gorm:"not null;index:idx_item_unit_item_stock_quality,priority:1"`
	Item          Item   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	StockStatus   int32  `gorm:"not null;default:0;index:idx_item_unit_item_stock_quality,priority:2"`
	QualityStatus int32  `gorm:"not null;default:0;index:idx_item_unit_item_stock_quality,priority:3"`
	Description   string `gorm:"type:varchar(255);not null;default:''"`

	InventoryFlows []InventoryFlow `gorm:"many2many:inventory_flow_item_units;"`
}

type InventoryFlow struct {
	gorm.Model

	FromUserID  int64      `gorm:"not null;index:idx_inventory_flow_filter,priority:1"`
	ToUserID    int64      `gorm:"not null;index:idx_inventory_flow_filter,priority:2"`
	FlowType    int32      `gorm:"not null;index:idx_inventory_flow_filter,priority:3"`
	FlowStatus  int32      `gorm:"not null;index:idx_inventory_flow_filter,priority:4"`
	Description string     `gorm:"type:varchar(255);not null;default:''"`
	ApprovedBy  int64      `gorm:"not null;default:0"`
	ApprovedAt  *time.Time `gorm:"default:null"`

	Items     []InventoryFlowItem `gorm:"foreignKey:InventoryFlowID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ItemUnits []ItemUnit          `gorm:"many2many:inventory_flow_item_units;"`
}

type InventoryFlowItem struct {
	gorm.Model

	InventoryFlowID  uint  `gorm:"not null;uniqueIndex:idx_inventory_flow_item,priority:1"`
	ItemID           uint  `gorm:"not null;uniqueIndex:idx_inventory_flow_item,priority:2;index"`
	Item             Item  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ApplyQuantity    int64 `gorm:"not null;default:0"`
	FinishedQuantity int64 `gorm:"not null;default:0"`
}

type ItemQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewItemQuery(ctx context.Context, db *gorm.DB) *ItemQuery {
	return &ItemQuery{ctx: ctx, db: db}
}

func (q *ItemQuery) Create(item *Item) error {
	return q.db.WithContext(q.ctx).Create(item).Error
}

func (q *ItemQuery) Update(id uint, updates map[string]any) error {
	return q.db.WithContext(q.ctx).Model(&Item{}).Where("id = ?", id).Updates(updates).Error
}

func (q *ItemQuery) Get(id uint) (Item, error) {
	var item Item
	err := q.db.WithContext(q.ctx).First(&item, id).Error
	return item, err
}

func (q *ItemQuery) List(pageNum, pageSize int, namePrefix string) ([]Item, int64, error) {
	var items []Item
	var total int64
	db := q.db.WithContext(q.ctx).Model(&Item{})
	if namePrefix != "" {
		db = db.Where("name LIKE ?", namePrefix+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("id DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&items).Error
	return items, total, err
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

func (q *ItemUnitQuery) List(pageNum, pageSize int, itemID uint, stockStatus, qualityStatus int32) ([]ItemUnit, int64, error) {
	var units []ItemUnit
	var total int64
	db := q.db.WithContext(q.ctx).Model(&ItemUnit{})
	if itemID > 0 {
		db = db.Where("item_id = ?", itemID)
	}
	if stockStatus > 0 {
		db = db.Where("stock_status = ?", stockStatus)
	}
	if qualityStatus > 0 {
		db = db.Where("quality_status = ?", qualityStatus)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("id DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&units).Error
	return units, total, err
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

func (q *InventoryFlowQuery) List(pageNum, pageSize int, fromUserID, toUserID int64, flowType, flowStatus int32) ([]InventoryFlow, int64, error) {
	var flows []InventoryFlow
	var total int64
	db := q.db.WithContext(q.ctx).Model(&InventoryFlow{})
	if fromUserID > 0 {
		db = db.Where("from_user_id = ?", fromUserID)
	}
	if toUserID > 0 {
		db = db.Where("to_user_id = ?", toUserID)
	}
	if flowType > 0 {
		db = db.Where("flow_type = ?", flowType)
	}
	if flowStatus > 0 {
		db = db.Where("flow_status = ?", flowStatus)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Preload("Items.Item").
		Preload("ItemUnits").
		Order("id DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&flows).Error
	return flows, total, err
}

type ItemAggregate struct {
	TotalCount       int64
	InStockCount     int64
	ReservedCount    int64
	OutStockCount    int64
	PendingCount     int64
	QualifiedCount   int64
	UnqualifiedCount int64
	AvailableCount   int64
}

func RecalculateItemCounts(ctx context.Context, db *gorm.DB, itemIDs ...uint) error {
	seen := make(map[uint]struct{}, len(itemIDs))
	for _, itemID := range itemIDs {
		if itemID == 0 {
			continue
		}
		if _, ok := seen[itemID]; ok {
			continue
		}
		seen[itemID] = struct{}{}
		if err := recalculateOneItem(ctx, db, itemID); err != nil {
			return err
		}
	}
	return nil
}

func recalculateOneItem(ctx context.Context, db *gorm.DB, itemID uint) error {
	var agg ItemAggregate
	err := db.WithContext(ctx).Model(&ItemUnit{}).
		Select(`
			COUNT(*) AS total_count,
			COALESCE(SUM(CASE WHEN stock_status = ? THEN 1 ELSE 0 END), 0) AS in_stock_count,
			COALESCE(SUM(CASE WHEN stock_status = ? THEN 1 ELSE 0 END), 0) AS reserved_count,
			COALESCE(SUM(CASE WHEN stock_status = ? THEN 1 ELSE 0 END), 0) AS out_stock_count,
			COALESCE(SUM(CASE WHEN quality_status = ? THEN 1 ELSE 0 END), 0) AS pending_count,
			COALESCE(SUM(CASE WHEN quality_status = ? THEN 1 ELSE 0 END), 0) AS qualified_count,
			COALESCE(SUM(CASE WHEN quality_status = ? THEN 1 ELSE 0 END), 0) AS unqualified_count,
			COALESCE(SUM(CASE WHEN stock_status = ? AND quality_status = ? THEN 1 ELSE 0 END), 0) AS available_count
		`, 1, 2, 3, 1, 2, 3, 1, 2).
		Where("item_id = ?", itemID).
		Scan(&agg).Error
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Model(&Item{}).Where("id = ?", itemID).Updates(map[string]any{
		"total_count":       agg.TotalCount,
		"in_stock_count":    agg.InStockCount,
		"reserved_count":    agg.ReservedCount,
		"out_stock_count":   agg.OutStockCount,
		"pending_count":     agg.PendingCount,
		"qualified_count":   agg.QualifiedCount,
		"unqualified_count": agg.UnqualifiedCount,
		"available_count":   agg.AvailableCount,
	}).Error
}
