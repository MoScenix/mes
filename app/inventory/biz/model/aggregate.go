package model

import (
	"context"

	"gorm.io/gorm"
)

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

func RecalculateEngineeringOrderProducedQuantity(ctx context.Context, db *gorm.DB, orderIDs ...uint) error {
	seen := make(map[uint]struct{}, len(orderIDs))
	for _, orderID := range orderIDs {
		if orderID == 0 {
			continue
		}
		if _, ok := seen[orderID]; ok {
			continue
		}
		seen[orderID] = struct{}{}
		var agg struct {
			ProducedQuantity    int64
			QualifiedQuantity   int64
			UnqualifiedQuantity int64
		}
		if err := db.WithContext(ctx).Model(&ItemUnit{}).
			Select(`
				COUNT(*) AS produced_quantity,
				COALESCE(SUM(CASE WHEN quality_status = ? THEN 1 ELSE 0 END), 0) AS qualified_quantity,
				COALESCE(SUM(CASE WHEN quality_status = ? THEN 1 ELSE 0 END), 0) AS unqualified_quantity
			`, 2, 3).
			Where("engineering_order_id = ?", orderID).
			Scan(&agg).Error; err != nil {
			return err
		}
		if err := db.WithContext(ctx).Model(&EngineeringOrder{}).Where("id = ?", orderID).Updates(map[string]any{
			"produced_quantity":    agg.ProducedQuantity,
			"qualified_quantity":   agg.QualifiedQuantity,
			"unqualified_quantity": agg.UnqualifiedQuantity,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
