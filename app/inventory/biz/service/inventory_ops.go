package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MoScenix/mes/app/inventory/biz/dal/mysql"
	"github.com/MoScenix/mes/app/inventory/biz/model"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

func runAddItem(ctx context.Context, req *inventory.AddItemReq) (*inventory.AddItemResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	item := &model.Item{
		Name:        req.GetName(),
		Unit:        req.GetUnit(),
		Description: req.GetDescription(),
	}
	if err := model.NewItemQuery(ctx, db).Create(item); err != nil {
		return nil, err
	}
	return &inventory.AddItemResp{Id: int64(item.ID)}, nil
}

func runUpdateItem(ctx context.Context, req *inventory.UpdateItemReq) (*inventory.UpdateItemResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "item id")
	if err != nil {
		return nil, err
	}
	if _, err := model.NewItemQuery(ctx, db).Get(id); err != nil {
		return nil, err
	}
	err = model.NewItemQuery(ctx, db).Update(id, map[string]any{
		"name":        req.GetName(),
		"unit":        req.GetUnit(),
		"description": req.GetDescription(),
	})
	if err != nil {
		return nil, err
	}
	return &inventory.UpdateItemResp{Success: true}, nil
}

func runGetItem(ctx context.Context, req *inventory.GetItemReq) (*inventory.GetItemResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "item id")
	if err != nil {
		return nil, err
	}
	item, err := model.NewItemQuery(ctx, db).Get(id)
	if err != nil {
		return nil, err
	}
	return &inventory.GetItemResp{Item: itemInfo(item)}, nil
}

func runListItem(ctx context.Context, req *inventory.ListItemReq) (*inventory.ListItemResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	pageNum, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	items, total, err := model.NewItemQuery(ctx, db).List(pageNum, pageSize, req.GetNamePrefix())
	if err != nil {
		return nil, err
	}
	resp := &inventory.ListItemResp{Total: total, ItemList: make([]*inventory.ItemInfo, 0, len(items))}
	for _, item := range items {
		resp.ItemList = append(resp.ItemList, itemInfo(item))
	}
	return resp, nil
}

func runAddItemUnit(ctx context.Context, req *inventory.AddItemUnitReq) (*inventory.AddItemUnitResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	itemID, err := uintID(req.GetItemId(), "item id")
	if err != nil {
		return nil, err
	}
	if !validStockStatus(req.GetStockStatus()) || !validQualityStatus(req.GetQualityStatus()) {
		return nil, errors.New("invalid item unit status")
	}
	unit := &model.ItemUnit{
		ItemID:        itemID,
		StockStatus:   int32(req.GetStockStatus()),
		QualityStatus: int32(req.GetQualityStatus()),
		Description:   req.GetDescription(),
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.Item{}, itemID).Error; err != nil {
			return err
		}
		if err := model.NewItemUnitQuery(ctx, tx).Create(unit); err != nil {
			return err
		}
		return model.RecalculateItemCounts(ctx, tx, itemID)
	})
	if err != nil {
		return nil, err
	}
	return &inventory.AddItemUnitResp{Id: int64(unit.ID)}, nil
}

func runUpdateItemUnitStatus(ctx context.Context, req *inventory.UpdateItemUnitStatusReq) (*inventory.UpdateItemUnitStatusResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "item unit id")
	if err != nil {
		return nil, err
	}
	if !validStockStatus(req.GetStockStatus()) || !validQualityStatus(req.GetQualityStatus()) {
		return nil, errors.New("invalid item unit status")
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var unit model.ItemUnit
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&unit, id).Error; err != nil {
			return err
		}
		if err := model.NewItemUnitQuery(ctx, tx).UpdateStatus(id, int32(req.GetStockStatus()), int32(req.GetQualityStatus())); err != nil {
			return err
		}
		return model.RecalculateItemCounts(ctx, tx, unit.ItemID)
	})
	if err != nil {
		return nil, err
	}
	return &inventory.UpdateItemUnitStatusResp{Success: true}, nil
}

func runGetItemUnit(ctx context.Context, req *inventory.GetItemUnitReq) (*inventory.GetItemUnitResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "item unit id")
	if err != nil {
		return nil, err
	}
	unit, err := model.NewItemUnitQuery(ctx, db).Get(id)
	if err != nil {
		return nil, err
	}
	return &inventory.GetItemUnitResp{ItemUnit: itemUnitInfo(unit)}, nil
}

func runListItemUnit(ctx context.Context, req *inventory.ListItemUnitReq) (*inventory.ListItemUnitResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	pageNum, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	var itemID uint
	if req.GetItemId() > 0 {
		itemID, err = uintID(req.GetItemId(), "item id")
		if err != nil {
			return nil, err
		}
	}
	units, total, err := model.NewItemUnitQuery(ctx, db).List(pageNum, pageSize, itemID, int32(req.GetStockStatus()), int32(req.GetQualityStatus()))
	if err != nil {
		return nil, err
	}
	resp := &inventory.ListItemUnitResp{Total: total, ItemUnitList: make([]*inventory.ItemUnitInfo, 0, len(units))}
	for _, unit := range units {
		resp.ItemUnitList = append(resp.ItemUnitList, itemUnitInfo(unit))
	}
	return resp, nil
}

func runCreateInventoryFlow(ctx context.Context, req *inventory.CreateInventoryFlowReq) (*inventory.CreateInventoryFlowResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	if !validFlowType(req.GetFlowType()) {
		return nil, errors.New("invalid flow type")
	}
	var flowID uint
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		items, itemIDs, err := buildFlowItems(req.GetItems())
		if err != nil {
			return err
		}
		if err := ensureItemsExist(ctx, tx, itemIDs); err != nil {
			return err
		}
		units, err := findUnitsByIDs(ctx, tx, req.GetItemUnitIds())
		if err != nil {
			return err
		}
		flow := model.InventoryFlow{
			FromUserID:  req.GetFromUserId(),
			ToUserID:    req.GetToUserId(),
			FlowType:    int32(req.GetFlowType()),
			FlowStatus:  int32(inventory.FlowStatus_FLOW_STATUS_DRAFT),
			Description: req.GetDescription(),
			Items:       items,
		}
		if err := tx.WithContext(ctx).Create(&flow).Error; err != nil {
			return err
		}
		if len(units) > 0 {
			if err := tx.Model(&flow).Association("ItemUnits").Replace(&units); err != nil {
				return err
			}
		}
		flowID = flow.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &inventory.CreateInventoryFlowResp{Id: int64(flowID)}, nil
}

func runUpdateInventoryFlowDraft(ctx context.Context, req *inventory.UpdateInventoryFlowDraftReq) (*inventory.UpdateInventoryFlowDraftResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "inventory flow id")
	if err != nil {
		return nil, err
	}
	if !validFlowType(req.GetFlowType()) {
		return nil, errors.New("invalid flow type")
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var flow model.InventoryFlow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&flow, id).Error; err != nil {
			return err
		}
		if flow.FlowStatus != int32(inventory.FlowStatus_FLOW_STATUS_DRAFT) {
			return errors.New("only draft inventory flow can be updated")
		}
		items, itemIDs, err := buildFlowItems(req.GetItems())
		if err != nil {
			return err
		}
		if err := ensureItemsExist(ctx, tx, itemIDs); err != nil {
			return err
		}
		units, err := findUnitsByIDs(ctx, tx, req.GetItemUnitIds())
		if err != nil {
			return err
		}
		if err := tx.Model(&flow).Updates(map[string]any{
			"from_user_id": req.GetFromUserId(),
			"to_user_id":   req.GetToUserId(),
			"flow_type":    int32(req.GetFlowType()),
			"description":  req.GetDescription(),
		}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("inventory_flow_id = ?", flow.ID).Delete(&model.InventoryFlowItem{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].InventoryFlowID = flow.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		return tx.Model(&flow).Association("ItemUnits").Replace(&units)
	})
	if err != nil {
		return nil, err
	}
	return &inventory.UpdateInventoryFlowDraftResp{Success: true}, nil
}

func runDeleteInventoryFlowDraft(ctx context.Context, req *inventory.DeleteInventoryFlowDraftReq) (*inventory.DeleteInventoryFlowDraftResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "inventory flow id")
	if err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var flow model.InventoryFlow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&flow, id).Error; err != nil {
			return err
		}
		if flow.FlowStatus != int32(inventory.FlowStatus_FLOW_STATUS_DRAFT) {
			return errors.New("only draft inventory flow can be deleted")
		}
		if err := tx.Model(&flow).Association("ItemUnits").Clear(); err != nil {
			return err
		}
		if err := tx.Unscoped().Where("inventory_flow_id = ?", flow.ID).Delete(&model.InventoryFlowItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&flow).Error
	})
	if err != nil {
		return nil, err
	}
	return &inventory.DeleteInventoryFlowDraftResp{Success: true}, nil
}

func runSubmitInventoryFlow(ctx context.Context, req *inventory.SubmitInventoryFlowReq) (*inventory.SubmitInventoryFlowResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "inventory flow id")
	if err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var flow model.InventoryFlow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&flow, id).Error; err != nil {
			return err
		}
		if flow.FlowStatus != int32(inventory.FlowStatus_FLOW_STATUS_DRAFT) {
			return errors.New("only draft inventory flow can be submitted")
		}
		return tx.Model(&flow).Update("flow_status", int32(inventory.FlowStatus_FLOW_STATUS_SUBMITTED)).Error
	})
	if err != nil {
		return nil, err
	}
	return &inventory.SubmitInventoryFlowResp{Success: true}, nil
}

func runAuditInventoryFlow(ctx context.Context, req *inventory.AuditInventoryFlowReq) (*inventory.AuditInventoryFlowResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "inventory flow id")
	if err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var flow model.InventoryFlow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&flow, id).Error; err != nil {
			return err
		}
		if flow.FlowStatus != int32(inventory.FlowStatus_FLOW_STATUS_SUBMITTED) {
			return errors.New("only submitted inventory flow can be audited")
		}
		now := time.Now()
		if !req.GetApproved() {
			return tx.Model(&flow).Updates(map[string]any{
				"flow_status": int32(inventory.FlowStatus_FLOW_STATUS_REJECTED),
				"approved_by": req.GetApprovedBy(),
				"approved_at": &now,
			}).Error
		}

		details, itemIDs, applyByItem, err := loadFlowDetails(ctx, tx, flow.ID)
		if err != nil {
			return err
		}
		units, err := loadFlowUnitsForUpdate(ctx, tx, flow.ID)
		if err != nil {
			return err
		}
		finishedByItem := make(map[uint]int64, len(applyByItem))
		affectedItemIDs := append([]uint(nil), itemIDs...)
		for _, unit := range units {
			if _, ok := applyByItem[unit.ItemID]; !ok {
				return fmt.Errorf("item unit %d does not belong to flow items", unit.ID)
			}
			finishedByItem[unit.ItemID]++
			affectedItemIDs = append(affectedItemIDs, unit.ItemID)
		}

		switch flow.FlowType {
		case int32(inventory.FlowType_FLOW_TYPE_IN):
			if err := approveInFlow(ctx, tx, units); err != nil {
				return err
			}
		case int32(inventory.FlowType_FLOW_TYPE_OUT):
			if err := validateOutFlow(units, applyByItem, finishedByItem); err != nil {
				return err
			}
			if err := approveOutFlow(ctx, tx, units); err != nil {
				return err
			}
		default:
			return errors.New("invalid flow type")
		}

		for _, detail := range details {
			if err := tx.Model(&model.InventoryFlowItem{}).
				Where("id = ?", detail.ID).
				Update("finished_quantity", finishedByItem[detail.ItemID]).Error; err != nil {
				return err
			}
		}
		if err := model.RecalculateItemCounts(ctx, tx, affectedItemIDs...); err != nil {
			return err
		}
		return tx.Model(&flow).Updates(map[string]any{
			"flow_status": int32(inventory.FlowStatus_FLOW_STATUS_APPROVED),
			"approved_by": req.GetApprovedBy(),
			"approved_at": &now,
		}).Error
	})
	if err != nil {
		return nil, err
	}
	return &inventory.AuditInventoryFlowResp{Success: true}, nil
}

func runGetInventoryFlow(ctx context.Context, req *inventory.GetInventoryFlowReq) (*inventory.GetInventoryFlowResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "inventory flow id")
	if err != nil {
		return nil, err
	}
	flow, err := model.NewInventoryFlowQuery(ctx, db).Get(id)
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryFlowResp{InventoryFlow: flowInfo(flow)}, nil
}

func runListInventoryFlow(ctx context.Context, req *inventory.ListInventoryFlowReq) (*inventory.ListInventoryFlowResp, error) {
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	pageNum, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	flows, total, err := model.NewInventoryFlowQuery(ctx, db).List(
		pageNum,
		pageSize,
		req.GetFromUserId(),
		req.GetToUserId(),
		int32(req.GetFlowType()),
		int32(req.GetFlowStatus()),
	)
	if err != nil {
		return nil, err
	}
	resp := &inventory.ListInventoryFlowResp{Total: total, InventoryFlowList: make([]*inventory.InventoryFlowInfo, 0, len(flows))}
	for _, flow := range flows {
		resp.InventoryFlowList = append(resp.InventoryFlowList, flowInfo(flow))
	}
	return resp, nil
}

func inventoryDB() (*gorm.DB, error) {
	if mysql.DB == nil {
		return nil, errors.New("inventory mysql is not initialized")
	}
	return mysql.DB, nil
}

func normalizePage(pageNum, pageSize int64) (int, int) {
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return int(pageNum), int(pageSize)
}

func uintID(id int64, name string) (uint, error) {
	if id <= 0 {
		return 0, fmt.Errorf("%s must be positive", name)
	}
	return uint(id), nil
}

func buildFlowItems(reqItems []*inventory.InventoryFlowItemReq) ([]model.InventoryFlowItem, []uint, error) {
	if len(reqItems) == 0 {
		return nil, nil, errors.New("inventory flow items are required")
	}
	quantityByItem := make(map[uint]int64, len(reqItems))
	for _, reqItem := range reqItems {
		itemID, err := uintID(reqItem.GetItemId(), "item id")
		if err != nil {
			return nil, nil, err
		}
		if reqItem.GetApplyQuantity() <= 0 {
			return nil, nil, errors.New("apply quantity must be positive")
		}
		quantityByItem[itemID] += reqItem.GetApplyQuantity()
	}
	items := make([]model.InventoryFlowItem, 0, len(quantityByItem))
	itemIDs := make([]uint, 0, len(quantityByItem))
	for itemID, quantity := range quantityByItem {
		itemIDs = append(itemIDs, itemID)
		items = append(items, model.InventoryFlowItem{
			ItemID:        itemID,
			ApplyQuantity: quantity,
		})
	}
	return items, itemIDs, nil
}

func ensureItemsExist(ctx context.Context, db *gorm.DB, itemIDs []uint) error {
	var count int64
	if err := db.WithContext(ctx).Model(&model.Item{}).Where("id IN ?", itemIDs).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(itemIDs)) {
		return errors.New("inventory flow contains unknown item")
	}
	return nil
}

func findUnitsByIDs(ctx context.Context, db *gorm.DB, ids []int64) ([]model.ItemUnit, error) {
	unitIDs, err := uniqueUintIDs(ids, "item unit id")
	if err != nil {
		return nil, err
	}
	if len(unitIDs) == 0 {
		return nil, nil
	}
	var units []model.ItemUnit
	if err := db.WithContext(ctx).Where("id IN ?", unitIDs).Find(&units).Error; err != nil {
		return nil, err
	}
	if len(units) != len(unitIDs) {
		return nil, errors.New("inventory flow contains unknown item unit")
	}
	return units, nil
}

func uniqueUintIDs(ids []int64, name string) ([]uint, error) {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		uintID, err := uintID(id, name)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[uintID]; ok {
			continue
		}
		seen[uintID] = struct{}{}
		result = append(result, uintID)
	}
	return result, nil
}

func loadFlowDetails(ctx context.Context, db *gorm.DB, flowID uint) ([]model.InventoryFlowItem, []uint, map[uint]int64, error) {
	var details []model.InventoryFlowItem
	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("inventory_flow_id = ?", flowID).
		Find(&details).Error; err != nil {
		return nil, nil, nil, err
	}
	if len(details) == 0 {
		return nil, nil, nil, errors.New("inventory flow items are required")
	}
	itemIDs := make([]uint, 0, len(details))
	applyByItem := make(map[uint]int64, len(details))
	for _, detail := range details {
		itemIDs = append(itemIDs, detail.ItemID)
		applyByItem[detail.ItemID] = detail.ApplyQuantity
	}
	return details, itemIDs, applyByItem, nil
}

func loadFlowUnitsForUpdate(ctx context.Context, db *gorm.DB, flowID uint) ([]model.ItemUnit, error) {
	var units []model.ItemUnit
	err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&model.ItemUnit{}).
		Joins("JOIN inventory_flow_item_units ON inventory_flow_item_units.item_unit_id = item_units.id").
		Where("inventory_flow_item_units.inventory_flow_id = ?", flowID).
		Find(&units).Error
	return units, err
}

func approveInFlow(ctx context.Context, db *gorm.DB, units []model.ItemUnit) error {
	for _, unit := range units {
		qualityStatus := unit.QualityStatus
		if qualityStatus == int32(inventory.QualityStatus_QUALITY_STATUS_UNKNOWN) {
			qualityStatus = int32(inventory.QualityStatus_QUALITY_STATUS_PENDING)
		}
		if err := db.WithContext(ctx).Model(&model.ItemUnit{}).
			Where("id = ?", unit.ID).
			Updates(map[string]any{
				"stock_status":   int32(inventory.StockStatus_STOCK_STATUS_IN_STOCK),
				"quality_status": qualityStatus,
			}).Error; err != nil {
			return err
		}
	}
	return nil
}

func validateOutFlow(units []model.ItemUnit, applyByItem, finishedByItem map[uint]int64) error {
	for _, unit := range units {
		if unit.StockStatus != int32(inventory.StockStatus_STOCK_STATUS_IN_STOCK) ||
			unit.QualityStatus != int32(inventory.QualityStatus_QUALITY_STATUS_QUALIFIED) {
			return fmt.Errorf("item unit %d is not available for out stock", unit.ID)
		}
	}
	for itemID, finishedQuantity := range finishedByItem {
		if finishedQuantity > applyByItem[itemID] {
			return fmt.Errorf("item %d item unit quantity exceeds apply quantity", itemID)
		}
	}
	return nil
}

func approveOutFlow(ctx context.Context, db *gorm.DB, units []model.ItemUnit) error {
	if len(units) == 0 {
		return nil
	}
	ids := make([]uint, 0, len(units))
	for _, unit := range units {
		ids = append(ids, unit.ID)
	}
	return db.WithContext(ctx).Model(&model.ItemUnit{}).
		Where("id IN ?", ids).
		Update("stock_status", int32(inventory.StockStatus_STOCK_STATUS_OUT_STOCK)).Error
}

func validFlowType(status inventory.FlowType) bool {
	return status == inventory.FlowType_FLOW_TYPE_IN || status == inventory.FlowType_FLOW_TYPE_OUT
}

func validStockStatus(status inventory.StockStatus) bool {
	switch status {
	case inventory.StockStatus_STOCK_STATUS_UNKNOWN,
		inventory.StockStatus_STOCK_STATUS_IN_STOCK,
		inventory.StockStatus_STOCK_STATUS_RESERVED,
		inventory.StockStatus_STOCK_STATUS_OUT_STOCK:
		return true
	default:
		return false
	}
}

func validQualityStatus(status inventory.QualityStatus) bool {
	switch status {
	case inventory.QualityStatus_QUALITY_STATUS_UNKNOWN,
		inventory.QualityStatus_QUALITY_STATUS_PENDING,
		inventory.QualityStatus_QUALITY_STATUS_QUALIFIED,
		inventory.QualityStatus_QUALITY_STATUS_UNQUALIFIED:
		return true
	default:
		return false
	}
}

func itemInfo(item model.Item) *inventory.ItemInfo {
	return &inventory.ItemInfo{
		Id:               int64(item.ID),
		Name:             item.Name,
		Unit:             item.Unit,
		Description:      item.Description,
		TotalCount:       item.TotalCount,
		InStockCount:     item.InStockCount,
		ReservedCount:    item.ReservedCount,
		OutStockCount:    item.OutStockCount,
		PendingCount:     item.PendingCount,
		QualifiedCount:   item.QualifiedCount,
		UnqualifiedCount: item.UnqualifiedCount,
		AvailableCount:   item.AvailableCount,
		CreateTime:       formatTime(item.CreatedAt),
		UpdateTime:       formatTime(item.UpdatedAt),
	}
}

func itemUnitInfo(unit model.ItemUnit) *inventory.ItemUnitInfo {
	return &inventory.ItemUnitInfo{
		Id:            int64(unit.ID),
		ItemId:        int64(unit.ItemID),
		StockStatus:   inventory.StockStatus(unit.StockStatus),
		QualityStatus: inventory.QualityStatus(unit.QualityStatus),
		Description:   unit.Description,
		CreateTime:    formatTime(unit.CreatedAt),
		UpdateTime:    formatTime(unit.UpdatedAt),
	}
}

func flowItemInfo(item model.InventoryFlowItem) *inventory.InventoryFlowItemInfo {
	return &inventory.InventoryFlowItemInfo{
		Id:               int64(item.ID),
		InventoryFlowId:  int64(item.InventoryFlowID),
		ItemId:           int64(item.ItemID),
		ApplyQuantity:    item.ApplyQuantity,
		FinishedQuantity: item.FinishedQuantity,
		Item:             itemInfo(item.Item),
	}
}

func flowInfo(flow model.InventoryFlow) *inventory.InventoryFlowInfo {
	resp := &inventory.InventoryFlowInfo{
		Id:          int64(flow.ID),
		FromUserId:  flow.FromUserID,
		ToUserId:    flow.ToUserID,
		FlowType:    inventory.FlowType(flow.FlowType),
		FlowStatus:  inventory.FlowStatus(flow.FlowStatus),
		Description: flow.Description,
		ApprovedBy:  flow.ApprovedBy,
		CreateTime:  formatTime(flow.CreatedAt),
		UpdateTime:  formatTime(flow.UpdatedAt),
		Items:       make([]*inventory.InventoryFlowItemInfo, 0, len(flow.Items)),
		ItemUnits:   make([]*inventory.ItemUnitInfo, 0, len(flow.ItemUnits)),
	}
	if flow.ApprovedAt != nil {
		resp.ApprovedAt = formatTime(*flow.ApprovedAt)
	}
	for _, item := range flow.Items {
		resp.Items = append(resp.Items, flowItemInfo(item))
	}
	for _, unit := range flow.ItemUnits {
		resp.ItemUnits = append(resp.ItemUnits, itemUnitInfo(unit))
	}
	return resp
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
