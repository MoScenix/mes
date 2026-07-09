package service

import (
	"context"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
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

func runCreateWorkOrderDraft(ctx context.Context, req *mes.CreateWorkOrderDraftRequest) (*mes.BaseResponseLong, error) {
	currentUserID, err := requireBFFUserID(ctx)
	if err != nil {
		return mesLongErr(err), nil
	}
	res, err := rpc.WorkOrderClient.CreateWorkOrder(mesCtx(ctx), &rpcworkorder.CreateWorkOrderReq{
		FromUserId:  currentUserID,
		ToUserId:    req.GetToUserId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}

func runUpdateWorkOrderDraft(ctx context.Context, req *mes.UpdateWorkOrderDraftRequest) (*mes.BaseResponseBoolean, error) {
	currentUserID, err := requireBFFUserID(ctx)
	if err != nil {
		return mesBoolErr(err), nil
	}
	current, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(ctx), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateWorkOrderDraft(ctx, current.GetWorkOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.WorkOrderClient.UpdateWorkOrderDraft(mesCtx(ctx), &rpcworkorder.UpdateWorkOrderDraftReq{
		Id:          req.GetId(),
		FromUserId:  currentUserID,
		ToUserId:    req.GetToUserId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runDeleteWorkOrderDraft(ctx context.Context, req *mes.DeleteRequest) (*mes.BaseResponseBoolean, error) {
	current, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(ctx), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateWorkOrderDraft(ctx, current.GetWorkOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.WorkOrderClient.DeleteWorkOrderDraft(mesCtx(ctx), &rpcworkorder.DeleteWorkOrderDraftReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runSubmitWorkOrder(ctx context.Context, req *mes.DeleteRequest) (*mes.BaseResponseBoolean, error) {
	current, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(ctx), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateWorkOrderDraft(ctx, current.GetWorkOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.WorkOrderClient.SubmitWorkOrder(mesCtx(ctx), &rpcworkorder.SubmitWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runGetWorkOrder(ctx context.Context, req *mes.GetByIdRequest) (*mes.BaseResponseWorkOrderVO, error) {
	res, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(ctx), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseWorkOrderVO{Code: 1, Message: err.Error()}, nil
	}
	if err = requireCanViewWorkOrder(ctx, res.GetWorkOrder()); err != nil {
		return &mes.BaseResponseWorkOrderVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseWorkOrderVO{Code: 0, Message: "success", Data: toWorkOrderVO(res.GetWorkOrder())}, nil
}

func runListWorkOrder(ctx context.Context, req *mes.ListWorkOrderRequest) (*mes.BaseResponsePageWorkOrderVO, error) {
	isTo := req.GetIsTo()
	status := rpcworkorder.WorkOrderStatus(req.GetStatus())
	if status == rpcworkorder.WorkOrderStatus_WORK_ORDER_STATUS_DRAFT {
		isTo = false
	}
	userID, err := scopedUserID(ctx, req.GetId())
	if err != nil {
		return &mes.BaseResponsePageWorkOrderVO{Code: 1, Message: err.Error()}, nil
	}
	res, err := rpc.WorkOrderClient.ListWorkOrder(mesCtx(ctx), &rpcworkorder.ListWorkOrderReq{
		PageNum:         req.GetPageNum(),
		PageSize:        req.GetPageSize(),
		Id:              userID,
		IsTo:            isTo,
		IsUnread:        req.GetIsUnread(),
		SinceTime:       req.GetSinceTime(),
		RecentSeconds:   req.GetRecentSeconds(),
		CursorUpdatedAt: req.GetCursorUpdatedAt(),
		CursorId:        req.GetCursorId(),
		NamePrefix:      req.GetNamePrefix(),
		Status:          status,
	})
	if err != nil {
		return &mes.BaseResponsePageWorkOrderVO{Code: 1, Message: err.Error()}, nil
	}
	records := make([]*mes.WorkOrderVO, 0, len(res.GetWorkOrderList()))
	for _, item := range res.GetWorkOrderList() {
		records = append(records, toWorkOrderVO(item))
	}
	page := pageWorkOrder(records, req.GetPageNum(), req.GetPageSize(), res.GetTotal())
	page.HasMore = res.GetHasMore()
	page.NextCursorUpdatedAt = res.GetNextCursorUpdatedAt()
	page.NextCursorId = res.GetNextCursorId()
	return &mes.BaseResponsePageWorkOrderVO{
		Code:    0,
		Message: "success",
		Data:    page,
	}, nil
}

func runMarkWorkOrderRead(ctx context.Context, req *mes.DeleteRequest) (*mes.BaseResponseBoolean, error) {
	current, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(ctx), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanMarkWorkOrderRead(ctx, current.GetWorkOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.WorkOrderClient.MarkWorkOrderRead(mesCtx(ctx), &rpcworkorder.MarkWorkOrderReadReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runAddItem(ctx context.Context, req *mes.AddItemRequest) (*mes.BaseResponseLong, error) {
	res, err := rpc.InventoryClient.AddItem(mesCtx(ctx), &rpcinventory.AddItemReq{
		Name:        req.GetName(),
		Unit:        req.GetUnit(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}

func runUpdateItem(ctx context.Context, req *mes.UpdateItemRequest) (*mes.BaseResponseBoolean, error) {
	res, err := rpc.InventoryClient.UpdateItem(mesCtx(ctx), &rpcinventory.UpdateItemReq{
		Id:          req.GetId(),
		Name:        req.GetName(),
		Unit:        req.GetUnit(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runGetItem(ctx context.Context, req *mes.GetByIdRequest) (*mes.BaseResponseItemVO, error) {
	res, err := rpc.InventoryClient.GetItem(mesCtx(ctx), &rpcinventory.GetItemReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseItemVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseItemVO{Code: 0, Message: "success", Data: toItemVO(res.GetItem())}, nil
}

func runListItem(ctx context.Context, namePrefix string, pageNum int64, pageSize int64, cursorName string, cursorID int64) (*mes.BaseResponsePageItemVO, error) {
	res, err := rpc.InventoryClient.ListItem(mesCtx(ctx), &rpcinventory.ListItemReq{
		PageNum:    pageNum,
		PageSize:   pageSize,
		NamePrefix: namePrefix,
		CursorName: cursorName,
		CursorId:   cursorID,
	})
	if err != nil {
		return &mes.BaseResponsePageItemVO{Code: 1, Message: err.Error()}, nil
	}
	records := make([]*mes.ItemVO, 0, len(res.GetItemList()))
	for _, item := range res.GetItemList() {
		records = append(records, toItemVO(item))
	}
	page := pageItem(records, pageNum, pageSize, res.GetTotal())
	page.HasMore = res.GetHasMore()
	page.NextCursorName = res.GetNextCursorName()
	page.NextCursorId = res.GetNextCursorId()
	return &mes.BaseResponsePageItemVO{
		Code:    0,
		Message: "success",
		Data:    page,
	}, nil
}

func runAddItemUnit(ctx context.Context, req *mes.AddItemUnitRequest) (*mes.BaseResponseLong, error) {
	res, err := rpc.InventoryClient.AddItemUnit(mesCtx(ctx), &rpcinventory.AddItemUnitReq{
		ItemId:             req.GetItemId(),
		StockStatus:        rpcinventory.StockStatus(req.GetStockStatus()),
		QualityStatus:      rpcinventory.QualityStatus(req.GetQualityStatus()),
		Description:        req.GetDescription(),
		EngineeringOrderId: req.GetEngineeringOrderId(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}

func runUpdateItemUnitStatus(ctx context.Context, req *mes.UpdateItemUnitStatusRequest) (*mes.BaseResponseBoolean, error) {
	res, err := rpc.InventoryClient.UpdateItemUnitStatus(mesCtx(ctx), &rpcinventory.UpdateItemUnitStatusReq{
		Id:            req.GetId(),
		StockStatus:   rpcinventory.StockStatus(req.GetStockStatus()),
		QualityStatus: rpcinventory.QualityStatus(req.GetQualityStatus()),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runGetItemUnit(ctx context.Context, req *mes.GetByIdRequest) (*mes.BaseResponseItemUnitVO, error) {
	res, err := rpc.InventoryClient.GetItemUnit(mesCtx(ctx), &rpcinventory.GetItemUnitReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseItemUnitVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseItemUnitVO{Code: 0, Message: "success", Data: toItemUnitVO(res.GetItemUnit())}, nil
}

func runListItemUnit(ctx context.Context, req *mes.ListItemUnitRequest) (*mes.BaseResponsePageItemUnitVO, error) {
	res, err := rpc.InventoryClient.ListItemUnit(mesCtx(ctx), &rpcinventory.ListItemUnitReq{
		PageNum:            req.GetPageNum(),
		PageSize:           req.GetPageSize(),
		ItemId:             req.GetItemId(),
		StockStatus:        rpcinventory.StockStatus(req.GetStockStatus()),
		QualityStatus:      rpcinventory.QualityStatus(req.GetQualityStatus()),
		EngineeringOrderId: req.GetEngineeringOrderId(),
		CursorId:           req.GetCursorId(),
		ItemNamePrefix:     req.GetItemNamePrefix(),
	})
	if err != nil {
		return &mes.BaseResponsePageItemUnitVO{Code: 1, Message: err.Error()}, nil
	}
	records := make([]*mes.ItemUnitVO, 0, len(res.GetItemUnitList()))
	for _, item := range res.GetItemUnitList() {
		records = append(records, toItemUnitVO(item))
	}
	page := pageItemUnit(records, req.GetPageNum(), req.GetPageSize(), res.GetTotal())
	page.HasMore = res.GetHasMore()
	page.NextCursorId = res.GetNextCursorId()
	return &mes.BaseResponsePageItemUnitVO{
		Code:    0,
		Message: "success",
		Data:    page,
	}, nil
}

func runCreateInventoryFlowDraft(ctx context.Context, req *mes.CreateInventoryFlowDraftRequest) (*mes.BaseResponseLong, error) {
	currentUserID, err := requireBFFUserID(ctx)
	if err != nil {
		return mesLongErr(err), nil
	}
	res, err := rpc.InventoryClient.CreateInventoryFlow(mesCtx(ctx), &rpcinventory.CreateInventoryFlowReq{
		FromUserId:  currentUserID,
		ToUserId:    req.GetToUserId(),
		FlowType:    rpcinventory.FlowType(req.GetFlowType()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Items:       toRPCInventoryFlowItems(req.GetItems()),
		ItemUnitIds: req.GetItemUnitIds(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}

func runUpdateInventoryFlowDraft(ctx context.Context, req *mes.UpdateInventoryFlowDraftRequest) (*mes.BaseResponseBoolean, error) {
	currentUserID, err := requireBFFUserID(ctx)
	if err != nil {
		return mesBoolErr(err), nil
	}
	current, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(ctx), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateInventoryFlowDraft(ctx, current.GetInventoryFlow()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.UpdateInventoryFlowDraft(mesCtx(ctx), &rpcinventory.UpdateInventoryFlowDraftReq{
		Id:          req.GetId(),
		FromUserId:  currentUserID,
		ToUserId:    req.GetToUserId(),
		FlowType:    rpcinventory.FlowType(req.GetFlowType()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Items:       toRPCInventoryFlowItems(req.GetItems()),
		ItemUnitIds: req.GetItemUnitIds(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runDeleteInventoryFlowDraft(ctx context.Context, req *mes.DeleteRequest) (*mes.BaseResponseBoolean, error) {
	current, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(ctx), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateInventoryFlowDraft(ctx, current.GetInventoryFlow()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.DeleteInventoryFlowDraft(mesCtx(ctx), &rpcinventory.DeleteInventoryFlowDraftReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runSubmitInventoryFlow(ctx context.Context, req *mes.DeleteRequest) (*mes.BaseResponseBoolean, error) {
	current, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(ctx), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateInventoryFlowDraft(ctx, current.GetInventoryFlow()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.SubmitInventoryFlow(mesCtx(ctx), &rpcinventory.SubmitInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runAuditInventoryFlow(ctx context.Context, req *mes.AuditInventoryFlowRequest) (*mes.BaseResponseBoolean, error) {
	if err := requireCanAuditInventoryFlow(ctx); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.AuditInventoryFlow(mesCtx(ctx), &rpcinventory.AuditInventoryFlowReq{
		Id:         req.GetId(),
		ApprovedBy: currentMESUserID(ctx),
		Approved:   req.GetApproved(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}

func runGetInventoryFlow(ctx context.Context, req *mes.GetByIdRequest) (*mes.BaseResponseInventoryFlowVO, error) {
	res, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(ctx), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseInventoryFlowVO{Code: 1, Message: err.Error()}, nil
	}
	if err = requireCanViewInventoryFlow(ctx, res.GetInventoryFlow()); err != nil {
		return &mes.BaseResponseInventoryFlowVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseInventoryFlowVO{Code: 0, Message: "success", Data: toInventoryFlowVO(res.GetInventoryFlow())}, nil
}

func runListInventoryFlow(ctx context.Context, req *mes.ListInventoryFlowRequest) (*mes.BaseResponsePageInventoryFlowVO, error) {
	scope := req.GetScope()
	userID := req.GetUserId()
	isTo := req.GetIsTo()
	flowStatus := rpcinventory.FlowStatus(req.GetFlowStatus())
	switch scope {
	case mes.MesListScope_MES_LIST_SCOPE_ALL:
		userID = 0
		isTo = false
	case mes.MesListScope_MES_LIST_SCOPE_AUDIT:
		userID = 0
		isTo = false
		if flowStatus == rpcinventory.FlowStatus_FLOW_STATUS_UNKNOWN {
			flowStatus = rpcinventory.FlowStatus_FLOW_STATUS_SUBMITTED
		}
	default:
		if req.GetIsTo() && req.GetFlowStatus() == mes.FlowStatus_FLOW_STATUS_DRAFT && !bffIsAdmin(ctx) {
			return &mes.BaseResponsePageInventoryFlowVO{Code: 1, Message: errForbiddenAccess.Error()}, nil
		}
		var err error
		userID, err = scopedUserID(ctx, req.GetUserId())
		if err != nil {
			return &mes.BaseResponsePageInventoryFlowVO{Code: 1, Message: err.Error()}, nil
		}
	}
	res, err := rpc.InventoryClient.ListInventoryFlow(mesCtx(ctx), &rpcinventory.ListInventoryFlowReq{
		UserId:          userID,
		IsTo:            isTo,
		FlowStatus:      flowStatus,
		PageNum:         req.GetPageNum(),
		PageSize:        req.GetPageSize(),
		SinceTime:       req.GetSinceTime(),
		RecentSeconds:   req.GetRecentSeconds(),
		CursorUpdatedAt: req.GetCursorUpdatedAt(),
		CursorId:        req.GetCursorId(),
		NamePrefix:      req.GetNamePrefix(),
		ItemNamePrefix:  req.GetItemNamePrefix(),
		Scope:           toInventoryListScope(scope),
	})
	if err != nil {
		return &mes.BaseResponsePageInventoryFlowVO{Code: 1, Message: err.Error()}, nil
	}
	records := make([]*mes.InventoryFlowVO, 0, len(res.GetInventoryFlowList()))
	for _, item := range res.GetInventoryFlowList() {
		records = append(records, toInventoryFlowVO(item))
	}
	page := pageInventoryFlow(records, req.GetPageNum(), req.GetPageSize(), res.GetTotal())
	page.HasMore = res.GetHasMore()
	page.NextCursorUpdatedAt = res.GetNextCursorUpdatedAt()
	page.NextCursorId = res.GetNextCursorId()
	return &mes.BaseResponsePageInventoryFlowVO{
		Code:    0,
		Message: "success",
		Data:    page,
	}, nil
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
