package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListInventoryFlowService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListInventoryFlowService(Context context.Context, RequestContext *app.RequestContext) *ListInventoryFlowService {
	return &ListInventoryFlowService{RequestContext: RequestContext, Context: Context}
}

func (h *ListInventoryFlowService) Run(req *mes.ListInventoryFlowRequest) (resp *mes.BaseResponsePageInventoryFlowVO, err error) {
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
		if req.GetIsTo() && req.GetFlowStatus() == mes.FlowStatus_FLOW_STATUS_DRAFT && !bffIsAdmin(h.Context) {
			return &mes.BaseResponsePageInventoryFlowVO{Code: 1, Message: errForbiddenAccess.Error()}, nil
		}
		var scopeErr error
		userID, scopeErr = scopedUserID(h.Context, req.GetUserId())
		if scopeErr != nil {
			return &mes.BaseResponsePageInventoryFlowVO{Code: 1, Message: scopeErr.Error()}, nil
		}
	}
	res, err := rpc.InventoryClient.ListInventoryFlow(mesCtx(h.Context), &rpcinventory.ListInventoryFlowReq{
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
		ItemUnitId:      req.GetItemUnitId(),
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
