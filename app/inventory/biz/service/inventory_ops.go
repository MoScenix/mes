package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

func (s *AddItemService) Run(req *inventory.AddItemReq) (resp *inventory.AddItemResp, err error) {
	ctx := s.ctx
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

func (s *UpdateItemService) Run(req *inventory.UpdateItemReq) (resp *inventory.UpdateItemResp, err error) {
	ctx := s.ctx
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

func (s *GetItemService) Run(req *inventory.GetItemReq) (resp *inventory.GetItemResp, err error) {
	ctx := s.ctx
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

func (s *ListItemService) Run(req *inventory.ListItemReq) (resp *inventory.ListItemResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	_, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	var cursorID uint
	if req.GetCursorId() > 0 {
		cursorID, err = uintID(req.GetCursorId(), "cursor id")
		if err != nil {
			return nil, err
		}
	}
	cursorUpdatedAt, err := parseCursorTime(req.GetCursorUpdatedAt())
	if err != nil {
		return nil, err
	}
	items, hasMore, err := model.NewItemQuery(ctx, db).List(pageSize, req.GetNamePrefix(), cursorUpdatedAt, cursorID)
	if err != nil {
		return nil, err
	}
	resp = &inventory.ListItemResp{Total: int64(len(items)), HasMore: hasMore, ItemList: make([]*inventory.ItemInfo, 0, len(items))}
	for _, item := range items {
		resp.ItemList = append(resp.ItemList, itemInfo(item))
	}
	if len(items) > 0 {
		last := items[len(items)-1]
		resp.NextCursorUpdatedAt = formatTime(last.UpdatedAt)
		resp.NextCursorId = int64(last.ID)
	}
	return resp, nil
}

func (s *CreateProcessDraftService) Run(req *inventory.CreateProcessDraftReq) (resp *inventory.CreateProcessDraftResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	if req.GetOwnerUserId() <= 0 {
		return nil, errors.New("owner user id must be positive")
	}
	itemID, err := uintID(req.GetItemId(), "item id")
	if err != nil {
		return nil, err
	}
	processItems, consumeItemIDs, err := buildProcessItems(req.GetItems())
	if err != nil {
		return nil, err
	}
	process := &model.Process{
		ItemID:      itemID,
		OwnerUserID: req.GetOwnerUserId(),
		Name:        strings.TrimSpace(req.GetName()),
		Description: req.GetDescription(),
		Status:      model.DraftStatusDraft,
		Items:       processItems,
	}
	if process.Name == "" {
		return nil, errors.New("process name is required")
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.Item{}, itemID).Error; err != nil {
			return err
		}
		if err := ensureItemsExist(ctx, tx, consumeItemIDs); err != nil {
			return err
		}
		return model.NewProcessQuery(ctx, tx).Create(process)
	})
	if err != nil {
		return nil, err
	}
	return &inventory.CreateProcessDraftResp{Id: int64(process.ID)}, nil
}

func (s *UpdateProcessDraftService) Run(req *inventory.UpdateProcessDraftReq) (resp *inventory.UpdateProcessDraftResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "process id")
	if err != nil {
		return nil, err
	}
	if req.GetOwnerUserId() <= 0 {
		return nil, errors.New("owner user id must be positive")
	}
	itemID, err := uintID(req.GetItemId(), "item id")
	if err != nil {
		return nil, err
	}
	processItems, consumeItemIDs, err := buildProcessItems(req.GetItems())
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errors.New("process name is required")
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var process model.Process
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&process, id).Error; err != nil {
			return err
		}
		if process.Status != model.DraftStatusDraft {
			return errors.New("only draft process can be updated")
		}
		if err := tx.First(&model.Item{}, itemID).Error; err != nil {
			return err
		}
		if err := ensureItemsExist(ctx, tx, consumeItemIDs); err != nil {
			return err
		}
		if err := model.NewProcessQuery(ctx, tx).UpdateDraft(id, map[string]any{
			"owner_user_id": req.GetOwnerUserId(),
			"item_id":       itemID,
			"name":          name,
			"description":   req.GetDescription(),
		}); err != nil {
			return err
		}
		if err := tx.Unscoped().Where("process_id = ?", id).Delete(&model.ProcessItem{}).Error; err != nil {
			return err
		}
		for i := range processItems {
			processItems[i].ProcessID = id
		}
		return tx.Create(&processItems).Error
	})
	if err != nil {
		return nil, err
	}
	return &inventory.UpdateProcessDraftResp{Success: true}, nil
}

func (s *DeleteProcessDraftService) Run(req *inventory.DeleteProcessDraftReq) (resp *inventory.DeleteProcessDraftResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "process id")
	if err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var process model.Process
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&process, id).Error; err != nil {
			return err
		}
		if process.Status != model.DraftStatusDraft {
			return errors.New("only draft process can be deleted")
		}
		if err := tx.Unscoped().Where("process_id = ?", id).Delete(&model.ProcessItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&process).Error
	})
	if err != nil {
		return nil, err
	}
	return &inventory.DeleteProcessDraftResp{Success: true}, nil
}

func (s *SubmitProcessService) Run(req *inventory.SubmitProcessReq) (resp *inventory.SubmitProcessResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "process id")
	if err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.ProcessItem{}).Where("process_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("process items are required")
		}
		return model.NewProcessQuery(ctx, tx).SubmitDraft(id)
	})
	if err != nil {
		return nil, err
	}
	return &inventory.SubmitProcessResp{Success: true}, nil
}

func (s *GetProcessService) Run(req *inventory.GetProcessReq) (resp *inventory.GetProcessResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "process id")
	if err != nil {
		return nil, err
	}
	process, err := model.NewProcessQuery(ctx, db).Get(id, true)
	if err != nil {
		return nil, err
	}
	return &inventory.GetProcessResp{Process: processInfo(process, true)}, nil
}

func (s *ListProcessService) Run(req *inventory.ListProcessReq) (resp *inventory.ListProcessResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	_, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	var itemID uint
	if req.GetItemId() > 0 {
		itemID, err = uintID(req.GetItemId(), "item id")
		if err != nil {
			return nil, err
		}
	}
	sinceTime, err := parseSinceTime(req.GetSinceTime(), req.GetRecentSeconds())
	if err != nil {
		return nil, err
	}
	cursorUpdatedAt, err := parseCursorTime(req.GetCursorUpdatedAt())
	if err != nil {
		return nil, err
	}
	var cursorID uint
	if req.GetCursorId() > 0 {
		cursorID, err = uintID(req.GetCursorId(), "cursor id")
		if err != nil {
			return nil, err
		}
	}
	processes, hasMore, err := model.NewProcessQuery(ctx, db).List(pageSize, req.GetOwnerUserId(), itemID, int32(req.GetStatus()), req.GetNamePrefix(), req.GetItemNamePrefix(), sinceTime, cursorUpdatedAt, cursorID)
	if err != nil {
		return nil, err
	}
	resp = &inventory.ListProcessResp{Total: int64(len(processes)), HasMore: hasMore, ProcessList: make([]*inventory.ProcessInfo, 0, len(processes))}
	for _, process := range processes {
		resp.ProcessList = append(resp.ProcessList, processInfo(process, false))
	}
	if len(processes) > 0 {
		last := processes[len(processes)-1]
		resp.NextCursorUpdatedAt = formatTime(last.UpdatedAt)
		resp.NextCursorId = int64(last.ID)
	}
	return resp, nil
}

func (s *AddItemUnitService) Run(req *inventory.AddItemUnitReq) (resp *inventory.AddItemUnitResp, err error) {
	ctx := s.ctx
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
		StockStatus:   int32(inventory.StockStatus_STOCK_STATUS_OUT_STOCK),
		QualityStatus: int32(req.GetQualityStatus()),
		Description:   req.GetDescription(),
	}
	var engineeringOrderID uint
	if req.GetEngineeringOrderId() > 0 {
		engineeringOrderID, err = uintID(req.GetEngineeringOrderId(), "engineering order id")
		if err != nil {
			return nil, err
		}
		unit.EngineeringOrderID = &engineeringOrderID
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.Item{}, itemID).Error; err != nil {
			return err
		}
		if engineeringOrderID > 0 {
			if err := validateEngineeringOrderBinding(ctx, tx, engineeringOrderID, itemID, int32(req.GetQualityStatus())); err != nil {
				return err
			}
		}
		if err := model.NewItemUnitQuery(ctx, tx).Create(unit); err != nil {
			return err
		}
		if err := model.RecalculateItemCounts(ctx, tx, itemID); err != nil {
			return err
		}
		return model.RecalculateEngineeringOrderProducedQuantity(ctx, tx, engineeringOrderID)
	})
	if err != nil {
		return nil, err
	}
	return &inventory.AddItemUnitResp{Id: int64(unit.ID)}, nil
}

func (s *UpdateItemUnitStatusService) Run(req *inventory.UpdateItemUnitStatusReq) (resp *inventory.UpdateItemUnitStatusResp, err error) {
	ctx := s.ctx
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
		if err := model.RecalculateItemCounts(ctx, tx, unit.ItemID); err != nil {
			return err
		}
		if unit.EngineeringOrderID != nil {
			return model.RecalculateEngineeringOrderProducedQuantity(ctx, tx, *unit.EngineeringOrderID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &inventory.UpdateItemUnitStatusResp{Success: true}, nil
}

func (s *GetItemUnitService) Run(req *inventory.GetItemUnitReq) (resp *inventory.GetItemUnitResp, err error) {
	ctx := s.ctx
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

func (s *ListItemUnitService) Run(req *inventory.ListItemUnitReq) (resp *inventory.ListItemUnitResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	_, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	var itemID uint
	if req.GetItemId() > 0 {
		itemID, err = uintID(req.GetItemId(), "item id")
		if err != nil {
			return nil, err
		}
	}
	var engineeringOrderID uint
	if req.GetEngineeringOrderId() > 0 {
		engineeringOrderID, err = uintID(req.GetEngineeringOrderId(), "engineering order id")
		if err != nil {
			return nil, err
		}
	}
	var inventoryFlowID uint
	if req.GetInventoryFlowId() > 0 {
		inventoryFlowID, err = uintID(req.GetInventoryFlowId(), "inventory flow id")
		if err != nil {
			return nil, err
		}
	}
	var cursorID uint
	if req.GetCursorId() > 0 {
		cursorID, err = uintID(req.GetCursorId(), "cursor id")
		if err != nil {
			return nil, err
		}
	}
	cursorUpdatedAt, err := parseCursorTime(req.GetCursorUpdatedAt())
	if err != nil {
		return nil, err
	}
	units, hasMore, err := model.NewItemUnitQuery(ctx, db).List(pageSize, itemID, engineeringOrderID, inventoryFlowID, int32(req.GetStockStatus()), int32(req.GetQualityStatus()), req.GetItemNamePrefix(), cursorUpdatedAt, cursorID)
	if err != nil {
		return nil, err
	}
	resp = &inventory.ListItemUnitResp{Total: int64(len(units)), HasMore: hasMore, ItemUnitList: make([]*inventory.ItemUnitInfo, 0, len(units))}
	for _, unit := range units {
		resp.ItemUnitList = append(resp.ItemUnitList, itemUnitInfo(unit))
	}
	if len(units) > 0 {
		last := units[len(units)-1]
		resp.NextCursorUpdatedAt = formatTime(last.UpdatedAt)
		resp.NextCursorId = int64(last.ID)
	}
	return resp, nil
}

func (s *CreateInventoryFlowService) Run(req *inventory.CreateInventoryFlowReq) (resp *inventory.CreateInventoryFlowResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	if !validFlowType(req.GetFlowType()) {
		return nil, errors.New("invalid flow type")
	}
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errors.New("inventory flow name is required")
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
			Name:        name,
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

func (s *UpdateInventoryFlowDraftService) Run(req *inventory.UpdateInventoryFlowDraftReq) (resp *inventory.UpdateInventoryFlowDraftResp, err error) {
	ctx := s.ctx
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
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errors.New("inventory flow name is required")
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
			"name":         name,
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

func (s *DeleteInventoryFlowDraftService) Run(req *inventory.DeleteInventoryFlowDraftReq) (resp *inventory.DeleteInventoryFlowDraftResp, err error) {
	ctx := s.ctx
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

func (s *SubmitInventoryFlowService) Run(req *inventory.SubmitInventoryFlowReq) (resp *inventory.SubmitInventoryFlowResp, err error) {
	ctx := s.ctx
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

func (s *CompleteInventoryFlowService) Run(req *inventory.CompleteInventoryFlowReq) (resp *inventory.CompleteInventoryFlowResp, err error) {
	ctx := s.ctx
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
		if flow.FlowStatus != int32(inventory.FlowStatus_FLOW_STATUS_APPROVED) {
			return errors.New("only approved inventory flow can be completed")
		}
		if flow.FlowType != int32(inventory.FlowType_FLOW_TYPE_IN) && flow.FlowType != int32(inventory.FlowType_FLOW_TYPE_OUT) {
			return errors.New("invalid flow type")
		}
		details, itemIDs, _, err := loadFlowDetails(ctx, tx, flow.ID)
		if err != nil {
			return err
		}
		units, err := findUnitsByIDsForUpdate(ctx, tx, req.GetItemUnitIds())
		if err != nil {
			return err
		}
		unitByID := make(map[uint]model.ItemUnit, len(units))
		unitIDs := make([]uint, 0, len(units))
		for _, unit := range units {
			unitByID[unit.ID] = unit
			unitIDs = append(unitIDs, unit.ID)
		}
		var existingJoins []model.InventoryFlowItemUnit
		if err := tx.WithContext(ctx).
			Where("inventory_flow_id = ? AND item_unit_id IN ?", flow.ID, unitIDs).
			Find(&existingJoins).Error; err != nil {
			return err
		}
		if len(existingJoins) > 0 {
			return fmt.Errorf("item unit %d has already been completed in this flow", existingJoins[0].ItemUnitID)
		}
		detailByItem := make(map[uint]model.InventoryFlowItem, len(details))
		for _, detail := range details {
			detailByItem[detail.ItemID] = detail
		}
		finishedByItem := make(map[uint]int64, len(detailByItem))
		affectedItemIDs := append([]uint(nil), itemIDs...)
		for _, unit := range units {
			detail, ok := detailByItem[unit.ItemID]
			if !ok {
				return fmt.Errorf("item unit %d does not belong to flow items", unit.ID)
			}
			if detail.FinishedQuantity+finishedByItem[unit.ItemID]+1 > detail.ApplyQuantity {
				return fmt.Errorf("item %d item unit quantity exceeds apply quantity", unit.ItemID)
			}
			if flow.FlowType == int32(inventory.FlowType_FLOW_TYPE_IN) {
				if unit.StockStatus != int32(inventory.StockStatus_STOCK_STATUS_OUT_STOCK) {
					return fmt.Errorf("item unit %d is already in stock or reserved", unit.ID)
				}
				if unit.QualityStatus != int32(inventory.QualityStatus_QUALITY_STATUS_QUALIFIED) {
					return fmt.Errorf("item unit %d is not qualified", unit.ID)
				}
			} else {
				if unit.StockStatus != int32(inventory.StockStatus_STOCK_STATUS_IN_STOCK) {
					return fmt.Errorf("item unit %d is not in stock", unit.ID)
				}
				if unit.QualityStatus != int32(inventory.QualityStatus_QUALITY_STATUS_QUALIFIED) {
					return fmt.Errorf("item unit %d is not qualified", unit.ID)
				}
			}
			finishedByItem[unit.ItemID]++
			affectedItemIDs = append(affectedItemIDs, unit.ItemID)
		}
		joins := make([]model.InventoryFlowItemUnit, 0, len(units))
		for _, unit := range units {
			if _, ok := unitByID[unit.ID]; !ok {
				return fmt.Errorf("item unit %d is invalid", unit.ID)
			}
			joins = append(joins, model.InventoryFlowItemUnit{InventoryFlowID: flow.ID, ItemUnitID: unit.ID})
		}
		if err := tx.Create(&joins).Error; err != nil {
			return err
		}
		for itemID, quantity := range finishedByItem {
			detail := detailByItem[itemID]
			if err := tx.Model(&model.InventoryFlowItem{}).
				Where("id = ?", detail.ID).
				Update("finished_quantity", detail.FinishedQuantity+quantity).Error; err != nil {
				return err
			}
		}
		if flow.FlowType == int32(inventory.FlowType_FLOW_TYPE_IN) {
			if err := approveInFlow(ctx, tx, units); err != nil {
				return err
			}
		} else {
			if err := approveOutFlow(ctx, tx, units); err != nil {
				return err
			}
		}
		return model.RecalculateItemCounts(ctx, tx, affectedItemIDs...)
	})
	if err != nil {
		return nil, err
	}
	return &inventory.CompleteInventoryFlowResp{Success: true}, nil
}

func (s *AuditInventoryFlowService) Run(req *inventory.AuditInventoryFlowReq) (resp *inventory.AuditInventoryFlowResp, err error) {
	ctx := s.ctx
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

		if flow.FlowType == int32(inventory.FlowType_FLOW_TYPE_IN) {
			return tx.Model(&flow).Updates(map[string]any{
				"flow_status": int32(inventory.FlowStatus_FLOW_STATUS_APPROVED),
				"approved_by": req.GetApprovedBy(),
				"approved_at": &now,
			}).Error
		}
		if flow.FlowType != int32(inventory.FlowType_FLOW_TYPE_OUT) {
			return errors.New("invalid flow type")
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

func (s *GetInventoryFlowService) Run(req *inventory.GetInventoryFlowReq) (resp *inventory.GetInventoryFlowResp, err error) {
	ctx := s.ctx
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

func (s *ListInventoryFlowService) Run(req *inventory.ListInventoryFlowReq) (resp *inventory.ListInventoryFlowResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	scope := req.GetScope()
	filterUser := scope != inventory.ListScope_LIST_SCOPE_ALL && scope != inventory.ListScope_LIST_SCOPE_AUDIT
	if filterUser && req.GetUserId() <= 0 {
		return nil, errors.New("user id must be positive")
	}
	_, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	sinceTime, err := parseSinceTime(req.GetSinceTime(), req.GetRecentSeconds())
	if err != nil {
		return nil, err
	}
	cursorUpdatedAt, err := parseCursorTime(req.GetCursorUpdatedAt())
	if err != nil {
		return nil, err
	}
	var cursorID uint
	if req.GetCursorId() > 0 {
		cursorID, err = uintID(req.GetCursorId(), "cursor id")
		if err != nil {
			return nil, err
		}
	}
	var itemUnitID uint
	if req.GetItemUnitId() > 0 {
		itemUnitID, err = uintID(req.GetItemUnitId(), "item unit id")
		if err != nil {
			return nil, err
		}
	}
	flowStatus := int32(req.GetFlowStatus())
	if scope == inventory.ListScope_LIST_SCOPE_AUDIT && flowStatus <= 0 {
		flowStatus = int32(inventory.FlowStatus_FLOW_STATUS_SUBMITTED)
	}
	flows, hasMore, err := model.NewInventoryFlowQuery(ctx, db).List(
		pageSize,
		req.GetUserId(),
		req.GetIsTo(),
		filterUser,
		flowStatus,
		req.GetNamePrefix(),
		req.GetItemNamePrefix(),
		itemUnitID,
		sinceTime,
		cursorUpdatedAt,
		cursorID,
	)
	if err != nil {
		return nil, err
	}
	resp = &inventory.ListInventoryFlowResp{Total: int64(len(flows)), HasMore: hasMore, InventoryFlowList: make([]*inventory.InventoryFlowInfo, 0, len(flows))}
	for _, flow := range flows {
		resp.InventoryFlowList = append(resp.InventoryFlowList, flowInfo(flow))
	}
	if len(flows) > 0 {
		last := flows[len(flows)-1]
		resp.NextCursorUpdatedAt = formatTime(last.UpdatedAt)
		resp.NextCursorId = int64(last.ID)
	}
	return resp, nil
}

func (s *CreateEngineeringOrderDraftService) Run(req *inventory.CreateEngineeringOrderDraftReq) (resp *inventory.CreateEngineeringOrderDraftResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	processID, err := uintID(req.GetProcessId(), "process id")
	if err != nil {
		return nil, err
	}
	if req.GetLeaderUserId() <= 0 {
		return nil, errors.New("leader user id must be positive")
	}
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errors.New("engineering order name is required")
	}
	if err := validateEngineeringQuantities(req.GetExpectedQuantity(), 0); err != nil {
		return nil, err
	}
	order := &model.EngineeringOrder{
		LeaderUserID:     req.GetLeaderUserId(),
		ProcessID:        processID,
		Name:             name,
		ExpectedQuantity: req.GetExpectedQuantity(),
		Status:           model.DraftStatusDraft,
		Description:      req.GetDescription(),
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		process, err := getSubmittedProcessForEngineering(ctx, tx, processID, req.GetItemId())
		if err != nil {
			return err
		}
		order.ItemID = process.ItemID
		return model.NewEngineeringOrderQuery(ctx, tx).Create(order)
	})
	if err != nil {
		return nil, err
	}
	return &inventory.CreateEngineeringOrderDraftResp{Id: int64(order.ID)}, nil
}

func (s *UpdateEngineeringOrderDraftService) Run(req *inventory.UpdateEngineeringOrderDraftReq) (resp *inventory.UpdateEngineeringOrderDraftResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "engineering order id")
	if err != nil {
		return nil, err
	}
	processID, err := uintID(req.GetProcessId(), "process id")
	if err != nil {
		return nil, err
	}
	if req.GetLeaderUserId() <= 0 {
		return nil, errors.New("leader user id must be positive")
	}
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errors.New("engineering order name is required")
	}
	if err := validateEngineeringQuantities(req.GetExpectedQuantity(), 0); err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order model.EngineeringOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error; err != nil {
			return err
		}
		if order.Status != model.DraftStatusDraft {
			return errors.New("only draft engineering order can be updated")
		}
		process, err := getSubmittedProcessForEngineering(ctx, tx, processID, req.GetItemId())
		if err != nil {
			return err
		}
		return model.NewEngineeringOrderQuery(ctx, tx).Update(id, map[string]any{
			"leader_user_id":    req.GetLeaderUserId(),
			"process_id":        process.ID,
			"item_id":           process.ItemID,
			"name":              name,
			"expected_quantity": req.GetExpectedQuantity(),
			"description":       req.GetDescription(),
		})
	})
	if err != nil {
		return nil, err
	}
	return &inventory.UpdateEngineeringOrderDraftResp{Success: true}, nil
}

func (s *DeleteEngineeringOrderDraftService) Run(req *inventory.DeleteEngineeringOrderDraftReq) (resp *inventory.DeleteEngineeringOrderDraftResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "engineering order id")
	if err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order model.EngineeringOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error; err != nil {
			return err
		}
		if order.Status != model.DraftStatusDraft {
			return errors.New("only draft engineering order can be deleted")
		}
		return tx.Delete(&order).Error
	})
	if err != nil {
		return nil, err
	}
	return &inventory.DeleteEngineeringOrderDraftResp{Success: true}, nil
}

func (s *SubmitEngineeringOrderService) Run(req *inventory.SubmitEngineeringOrderReq) (resp *inventory.SubmitEngineeringOrderResp, err error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "engineering order id")
	if err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order model.EngineeringOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error; err != nil {
			return err
		}
		if order.Status != model.DraftStatusDraft {
			return errors.New("only draft engineering order can be submitted")
		}
		return tx.Model(&order).Update("status", model.DraftStatusSubmitted).Error
	})
	if err != nil {
		return nil, err
	}
	return &inventory.SubmitEngineeringOrderResp{Success: true}, nil
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

func buildProcessItems(reqItems []*inventory.ProcessItemReq) ([]model.ProcessItem, []uint, error) {
	if len(reqItems) == 0 {
		return nil, nil, errors.New("process items are required")
	}
	quantityByItem := make(map[uint]int64, len(reqItems))
	for _, reqItem := range reqItems {
		itemID, err := uintID(reqItem.GetConsumeItemId(), "consume item id")
		if err != nil {
			return nil, nil, err
		}
		if reqItem.GetQuantity() <= 0 {
			return nil, nil, errors.New("process item quantity must be positive")
		}
		quantityByItem[itemID] += reqItem.GetQuantity()
	}
	items := make([]model.ProcessItem, 0, len(quantityByItem))
	itemIDs := make([]uint, 0, len(quantityByItem))
	for itemID, quantity := range quantityByItem {
		itemIDs = append(itemIDs, itemID)
		items = append(items, model.ProcessItem{
			ConsumeItemID: itemID,
			Quantity:      quantity,
		})
	}
	return items, itemIDs, nil
}

func ensureItemsExist(ctx context.Context, db *gorm.DB, itemIDs []uint) error {
	if len(itemIDs) == 0 {
		return nil
	}
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

func findUnitsByIDsForUpdate(ctx context.Context, db *gorm.DB, ids []int64) ([]model.ItemUnit, error) {
	unitIDs, err := uniqueUintIDs(ids, "item unit id")
	if err != nil {
		return nil, err
	}
	if len(unitIDs) == 0 {
		return nil, errors.New("item units are required")
	}
	var units []model.ItemUnit
	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id IN ?", unitIDs).Find(&units).Error; err != nil {
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

func validateInFlow(units []model.ItemUnit, applyByItem, finishedByItem map[uint]int64) error {
	for _, unit := range units {
		if unit.StockStatus != int32(inventory.StockStatus_STOCK_STATUS_OUT_STOCK) {
			return fmt.Errorf("item unit %d is already in stock or reserved", unit.ID)
		}
		if unit.QualityStatus != int32(inventory.QualityStatus_QUALITY_STATUS_QUALIFIED) {
			return fmt.Errorf("item unit %d is not qualified", unit.ID)
		}
	}
	return validateExactFlowQuantities(applyByItem, finishedByItem)
}

func validateExactFlowQuantities(applyByItem, finishedByItem map[uint]int64) error {
	for itemID, applyQuantity := range applyByItem {
		if finishedByItem[itemID] != applyQuantity {
			return fmt.Errorf("item %d item unit quantity must equal apply quantity", itemID)
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

func validateEngineeringOrderBinding(ctx context.Context, db *gorm.DB, orderID, itemID uint, qualityStatus int32) error {
	if !validQualityStatus(inventory.QualityStatus(qualityStatus)) {
		return errors.New("invalid item unit quality status")
	}
	var order model.EngineeringOrder
	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, orderID).Error; err != nil {
		return err
	}
	if order.Status != model.DraftStatusSubmitted {
		return errors.New("item unit can only bind submitted engineering order")
	}
	if order.ItemID != itemID {
		return errors.New("item unit item does not match engineering order item")
	}
	if err := model.RecalculateEngineeringOrderProducedQuantity(ctx, db, orderID); err != nil {
		return err
	}
	if err := db.WithContext(ctx).First(&order, orderID).Error; err != nil {
		return err
	}
	if order.ProducedQuantity+1 > order.ExpectedQuantity {
		return errors.New("engineering order produced quantity exceeds expected quantity")
	}
	return nil
}

func validateEngineeringQuantities(expectedQuantity, producedQuantity int64) error {
	if expectedQuantity <= 0 {
		return errors.New("expected quantity must be positive")
	}
	if expectedQuantity < producedQuantity {
		return errors.New("expected quantity cannot be less than produced quantity")
	}
	return nil
}

func getSubmittedProcessForEngineering(ctx context.Context, db *gorm.DB, processID uint, reqItemID int64) (model.Process, error) {
	var process model.Process
	if err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&process, processID).Error; err != nil {
		return process, err
	}
	if process.Status != model.DraftStatusSubmitted {
		return process, errors.New("engineering order requires submitted process")
	}
	if reqItemID > 0 && uint(reqItemID) != process.ItemID {
		return process, errors.New("engineering order item does not match process item")
	}
	return process, nil
}

func validFlowType(status inventory.FlowType) bool {
	return status == inventory.FlowType_FLOW_TYPE_IN || status == inventory.FlowType_FLOW_TYPE_OUT
}

func validStockStatus(status inventory.StockStatus) bool {
	switch status {
	case inventory.StockStatus_STOCK_STATUS_UNKNOWN,
		inventory.StockStatus_STOCK_STATUS_IN_STOCK,
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
	info := &inventory.ItemUnitInfo{
		Id:            int64(unit.ID),
		ItemId:        int64(unit.ItemID),
		StockStatus:   inventory.StockStatus(unit.StockStatus),
		QualityStatus: inventory.QualityStatus(unit.QualityStatus),
		Description:   unit.Description,
		CreateTime:    formatTime(unit.CreatedAt),
		UpdateTime:    formatTime(unit.UpdatedAt),
	}
	if unit.EngineeringOrderID != nil {
		info.EngineeringOrderId = int64(*unit.EngineeringOrderID)
	}
	return info
}

func processItemInfo(item model.ProcessItem) *inventory.ProcessItemInfo {
	return &inventory.ProcessItemInfo{
		Id:            int64(item.ID),
		ProcessId:     int64(item.ProcessID),
		ConsumeItemId: int64(item.ConsumeItemID),
		Quantity:      item.Quantity,
		ConsumeItem:   itemInfo(item.ConsumeItem),
	}
}

func processInfo(process model.Process, withItems bool) *inventory.ProcessInfo {
	info := &inventory.ProcessInfo{
		Id:          int64(process.ID),
		ItemId:      int64(process.ItemID),
		OwnerUserId: process.OwnerUserID,
		Name:        process.Name,
		Description: process.Description,
		Status:      inventory.DraftStatus(process.Status),
		Item:        itemInfo(process.Item),
		CreateTime:  formatTime(process.CreatedAt),
		UpdateTime:  formatTime(process.UpdatedAt),
	}
	if withItems {
		info.Items = make([]*inventory.ProcessItemInfo, 0, len(process.Items))
		for _, item := range process.Items {
			info.Items = append(info.Items, processItemInfo(item))
		}
	}
	return info
}

func engineeringOrderInfo(order model.EngineeringOrder, withUnits bool) *inventory.EngineeringOrderInfo {
	info := &inventory.EngineeringOrderInfo{
		Id:                  int64(order.ID),
		LeaderUserId:        order.LeaderUserID,
		ItemId:              int64(order.ItemID),
		Item:                itemInfo(order.Item),
		Name:                order.Name,
		ExpectedQuantity:    order.ExpectedQuantity,
		QualifiedQuantity:   order.QualifiedQuantity,
		ProducedQuantity:    order.ProducedQuantity,
		Description:         order.Description,
		CreateTime:          formatTime(order.CreatedAt),
		UpdateTime:          formatTime(order.UpdatedAt),
		ProcessId:           int64(order.ProcessID),
		Process:             processInfo(order.Process, false),
		Status:              inventory.DraftStatus(order.Status),
		UnqualifiedQuantity: order.UnqualifiedQuantity,
	}
	if withUnits {
		info.ItemUnits = make([]*inventory.ItemUnitInfo, 0, len(order.ItemUnits))
		for _, unit := range order.ItemUnits {
			info.ItemUnits = append(info.ItemUnits, itemUnitInfo(unit))
		}
	}
	return info
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
		Name:        flow.Name,
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

func parseSinceTime(sinceTime string, recentSeconds int64) (*time.Time, error) {
	value := strings.TrimSpace(sinceTime)
	if value != "" {
		t, err := parseListTime(value)
		if err == nil {
			return t, nil
		}
		return nil, errors.New("sinceTime must use format 2006-01-02 15:04:05")
	}
	if recentSeconds > 0 {
		t := time.Now().Add(-time.Duration(recentSeconds) * time.Second)
		return &t, nil
	}
	return nil, nil
}

func parseCursorTime(cursorUpdatedAt string) (*time.Time, error) {
	value := strings.TrimSpace(cursorUpdatedAt)
	if value == "" {
		return nil, nil
	}
	t, err := parseListTime(value)
	if err != nil {
		return nil, errors.New("cursorUpdatedAt must use format 2006-01-02 15:04:05")
	}
	return t, nil
}

func parseListTime(value string) (*time.Time, error) {
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339} {
		t, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return &t, nil
		}
	}
	return nil, errors.New("invalid time")
}
