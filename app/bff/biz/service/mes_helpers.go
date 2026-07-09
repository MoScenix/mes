package service

import (
	"context"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	rpcworkorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

func mesCtx(ctx context.Context) context.Context {
	return utils.WithIdentityMeta(ctx)
}

func currentMESUserID(ctx context.Context) int64 {
	userID, _ := utils.UserIDFromContext(ctx)
	return userID
}

func toInventoryListScope(scope mes.MesListScope) rpcinventory.ListScope {
	switch scope {
	case mes.MesListScope_MES_LIST_SCOPE_MINE:
		return rpcinventory.ListScope_LIST_SCOPE_MINE
	case mes.MesListScope_MES_LIST_SCOPE_ALL:
		return rpcinventory.ListScope_LIST_SCOPE_ALL
	case mes.MesListScope_MES_LIST_SCOPE_AUDIT:
		return rpcinventory.ListScope_LIST_SCOPE_AUDIT
	case mes.MesListScope_MES_LIST_SCOPE_BY_PROCESS:
		return rpcinventory.ListScope_LIST_SCOPE_BY_PROCESS
	default:
		return rpcinventory.ListScope_LIST_SCOPE_UNSPECIFIED
	}
}

func mesBool(data bool) *mes.BaseResponseBoolean {
	return &mes.BaseResponseBoolean{Code: 0, Data: data, Message: "success"}
}

func mesBoolErr(err error) *mes.BaseResponseBoolean {
	return &mes.BaseResponseBoolean{Code: 1, Message: err.Error()}
}

func mesLong(data int64) *mes.BaseResponseLong {
	return &mes.BaseResponseLong{Code: 0, Data: data, Message: "success"}
}

func mesLongErr(err error) *mes.BaseResponseLong {
	return &mes.BaseResponseLong{Code: 1, Message: err.Error()}
}

func toWorkOrderVO(info *rpcworkorder.WorkOrderInfo) *mes.WorkOrderVO {
	if info == nil {
		return nil
	}
	return &mes.WorkOrderVO{
		Id:          info.GetId(),
		FromUserId:  info.GetFromUserId(),
		ToUserId:    info.GetToUserId(),
		Name:        info.GetName(),
		Description: info.GetDescription(),
		Status:      mes.WorkOrderStatus(info.GetStatus()),
		CreateTime:  info.GetCreateTime(),
		UpdateTime:  info.GetUpdateTime(),
		ReadStatus:  mes.WorkOrderReadStatus(info.GetReadStatus()),
	}
}

func toItemVO(info *rpcinventory.ItemInfo) *mes.ItemVO {
	if info == nil {
		return nil
	}
	return &mes.ItemVO{
		Id:               info.GetId(),
		Name:             info.GetName(),
		Unit:             info.GetUnit(),
		Description:      info.GetDescription(),
		TotalCount:       info.GetTotalCount(),
		InStockCount:     info.GetInStockCount(),
		ReservedCount:    info.GetReservedCount(),
		OutStockCount:    info.GetOutStockCount(),
		PendingCount:     info.GetPendingCount(),
		QualifiedCount:   info.GetQualifiedCount(),
		UnqualifiedCount: info.GetUnqualifiedCount(),
		AvailableCount:   info.GetAvailableCount(),
		CreateTime:       info.GetCreateTime(),
		UpdateTime:       info.GetUpdateTime(),
	}
}

func toItemUnitVO(info *rpcinventory.ItemUnitInfo) *mes.ItemUnitVO {
	if info == nil {
		return nil
	}
	return &mes.ItemUnitVO{
		Id:                 info.GetId(),
		ItemId:             info.GetItemId(),
		StockStatus:        mes.StockStatus(info.GetStockStatus()),
		QualityStatus:      mes.QualityStatus(info.GetQualityStatus()),
		Description:        info.GetDescription(),
		CreateTime:         info.GetCreateTime(),
		UpdateTime:         info.GetUpdateTime(),
		EngineeringOrderId: info.GetEngineeringOrderId(),
	}
}

func toProcessItemVO(info *rpcinventory.ProcessItemInfo) *mes.ProcessItemVO {
	if info == nil {
		return nil
	}
	return &mes.ProcessItemVO{
		Id:            info.GetId(),
		ProcessId:     info.GetProcessId(),
		ConsumeItemId: info.GetConsumeItemId(),
		Quantity:      info.GetQuantity(),
		ConsumeItem:   toItemVO(info.GetConsumeItem()),
	}
}

func toProcessVO(info *rpcinventory.ProcessInfo) *mes.ProcessVO {
	if info == nil {
		return nil
	}
	items := make([]*mes.ProcessItemVO, 0, len(info.GetItems()))
	for _, item := range info.GetItems() {
		items = append(items, toProcessItemVO(item))
	}
	return &mes.ProcessVO{
		Id:          info.GetId(),
		ItemId:      info.GetItemId(),
		OwnerUserId: info.GetOwnerUserId(),
		Name:        info.GetName(),
		Description: info.GetDescription(),
		Status:      mes.DraftStatus(info.GetStatus()),
		Item:        toItemVO(info.GetItem()),
		Items:       items,
		CreateTime:  info.GetCreateTime(),
		UpdateTime:  info.GetUpdateTime(),
	}
}

func toEngineeringOrderVO(info *rpcinventory.EngineeringOrderInfo) *mes.EngineeringOrderVO {
	if info == nil {
		return nil
	}
	itemUnits := make([]*mes.ItemUnitVO, 0, len(info.GetItemUnits()))
	for _, item := range info.GetItemUnits() {
		itemUnits = append(itemUnits, toItemUnitVO(item))
	}
	return &mes.EngineeringOrderVO{
		Id:                  info.GetId(),
		LeaderUserId:        info.GetLeaderUserId(),
		ItemId:              info.GetItemId(),
		Item:                toItemVO(info.GetItem()),
		Name:                info.GetName(),
		ExpectedQuantity:    info.GetExpectedQuantity(),
		QualifiedQuantity:   info.GetQualifiedQuantity(),
		ProducedQuantity:    info.GetProducedQuantity(),
		Description:         info.GetDescription(),
		ItemUnits:           itemUnits,
		CreateTime:          info.GetCreateTime(),
		UpdateTime:          info.GetUpdateTime(),
		ProcessId:           info.GetProcessId(),
		Process:             toProcessVO(info.GetProcess()),
		Status:              mes.DraftStatus(info.GetStatus()),
		UnqualifiedQuantity: info.GetUnqualifiedQuantity(),
	}
}

func toInventoryFlowVO(info *rpcinventory.InventoryFlowInfo) *mes.InventoryFlowVO {
	if info == nil {
		return nil
	}
	items := make([]*mes.InventoryFlowItemVO, 0, len(info.GetItems()))
	for _, item := range info.GetItems() {
		items = append(items, toInventoryFlowItemVO(item))
	}
	itemUnits := make([]*mes.ItemUnitVO, 0, len(info.GetItemUnits()))
	for _, item := range info.GetItemUnits() {
		itemUnits = append(itemUnits, toItemUnitVO(item))
	}
	return &mes.InventoryFlowVO{
		Id:          info.GetId(),
		FromUserId:  info.GetFromUserId(),
		ToUserId:    info.GetToUserId(),
		FlowType:    mes.FlowType(info.GetFlowType()),
		FlowStatus:  mes.FlowStatus(info.GetFlowStatus()),
		Name:        info.GetName(),
		Description: info.GetDescription(),
		ApprovedBy:  info.GetApprovedBy(),
		ApprovedAt:  info.GetApprovedAt(),
		Items:       items,
		ItemUnits:   itemUnits,
		CreateTime:  info.GetCreateTime(),
		UpdateTime:  info.GetUpdateTime(),
	}
}

func toInventoryFlowItemVO(info *rpcinventory.InventoryFlowItemInfo) *mes.InventoryFlowItemVO {
	if info == nil {
		return nil
	}
	return &mes.InventoryFlowItemVO{
		Id:               info.GetId(),
		InventoryFlowId:  info.GetInventoryFlowId(),
		ItemId:           info.GetItemId(),
		ApplyQuantity:    info.GetApplyQuantity(),
		FinishedQuantity: info.GetFinishedQuantity(),
		Item:             toItemVO(info.GetItem()),
	}
}

func toRPCInventoryFlowItems(items []*mes.InventoryFlowItemRequest) []*rpcinventory.InventoryFlowItemReq {
	result := make([]*rpcinventory.InventoryFlowItemReq, 0, len(items))
	for _, item := range items {
		result = append(result, &rpcinventory.InventoryFlowItemReq{
			ItemId:        item.GetItemId(),
			ApplyQuantity: item.GetApplyQuantity(),
		})
	}
	return result
}

func toRPCProcessItems(items []*mes.ProcessItemRequest) []*rpcinventory.ProcessItemReq {
	result := make([]*rpcinventory.ProcessItemReq, 0, len(items))
	for _, item := range items {
		result = append(result, &rpcinventory.ProcessItemReq{
			ConsumeItemId: item.GetConsumeItemId(),
			Quantity:      item.GetQuantity(),
		})
	}
	return result
}

func pageTotal(pageSize int64, total int64) int64 {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}

func pageWorkOrder(records []*mes.WorkOrderVO, pageNum int64, pageSize int64, total int64) *mes.PageWorkOrderVO {
	return &mes.PageWorkOrderVO{Records: records, PageNumber: pageNum, PageSize: pageSize, TotalPage: pageTotal(pageSize, total), TotalRow: total}
}

func pageItem(records []*mes.ItemVO, pageNum int64, pageSize int64, total int64) *mes.PageItemVO {
	return &mes.PageItemVO{Records: records, PageNumber: pageNum, PageSize: pageSize, TotalPage: pageTotal(pageSize, total), TotalRow: total}
}

func pageItemUnit(records []*mes.ItemUnitVO, pageNum int64, pageSize int64, total int64) *mes.PageItemUnitVO {
	return &mes.PageItemUnitVO{Records: records, PageNumber: pageNum, PageSize: pageSize, TotalPage: pageTotal(pageSize, total), TotalRow: total}
}

func pageProcess(records []*mes.ProcessVO, pageNum int64, pageSize int64, total int64) *mes.PageProcessVO {
	return &mes.PageProcessVO{Records: records, PageNumber: pageNum, PageSize: pageSize, TotalPage: pageTotal(pageSize, total), TotalRow: total}
}

func pageEngineeringOrder(records []*mes.EngineeringOrderVO, pageNum int64, pageSize int64, total int64) *mes.PageEngineeringOrderVO {
	return &mes.PageEngineeringOrderVO{Records: records, PageNumber: pageNum, PageSize: pageSize, TotalPage: pageTotal(pageSize, total), TotalRow: total}
}

func pageInventoryFlow(records []*mes.InventoryFlowVO, pageNum int64, pageSize int64, total int64) *mes.PageInventoryFlowVO {
	return &mes.PageInventoryFlowVO{Records: records, PageNumber: pageNum, PageSize: pageSize, TotalPage: pageTotal(pageSize, total), TotalRow: total}
}
