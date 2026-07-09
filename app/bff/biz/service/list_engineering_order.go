package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListEngineeringOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListEngineeringOrderService(Context context.Context, RequestContext *app.RequestContext) *ListEngineeringOrderService {
	return &ListEngineeringOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *ListEngineeringOrderService) Run(req *mes.ListEngineeringOrderRequest) (resp *mes.BaseResponsePageEngineeringOrderVO, err error) {
	leaderUserID := req.GetLeaderUserId()
	status := rpcinventory.DraftStatus(req.GetStatus())
	switch req.GetScope() {
	case mes.MesListScope_MES_LIST_SCOPE_MINE:
		currentUserID, err := requireBFFUserID(h.Context)
		if err != nil {
			return &mes.BaseResponsePageEngineeringOrderVO{Code: 1, Message: err.Error()}, nil
		}
		leaderUserID = currentUserID
	case mes.MesListScope_MES_LIST_SCOPE_ALL, mes.MesListScope_MES_LIST_SCOPE_BY_PROCESS:
		leaderUserID = 0
	default:
		if !bffIsAdmin(h.Context) {
			currentUserID, err := requireBFFUserID(h.Context)
			if err != nil {
				return &mes.BaseResponsePageEngineeringOrderVO{Code: 1, Message: err.Error()}, nil
			}
			if leaderUserID != 0 && leaderUserID != currentUserID {
				return &mes.BaseResponsePageEngineeringOrderVO{Code: 1, Message: errForbiddenAccess.Error()}, nil
			}
			if status == rpcinventory.DraftStatus_DRAFT_STATUS_DRAFT {
				leaderUserID = currentUserID
			}
			if status == rpcinventory.DraftStatus_DRAFT_STATUS_UNKNOWN && leaderUserID == 0 {
				status = rpcinventory.DraftStatus_DRAFT_STATUS_SUBMITTED
			}
		}
	}
	res, err := rpc.InventoryClient.ListEngineeringOrder(mesCtx(h.Context), &rpcinventory.ListEngineeringOrderReq{
		LeaderUserId:    leaderUserID,
		ItemId:          req.GetItemId(),
		PageNum:         req.GetPageNum(),
		PageSize:        req.GetPageSize(),
		ProcessId:       req.GetProcessId(),
		Status:          status,
		SinceTime:       req.GetSinceTime(),
		RecentSeconds:   req.GetRecentSeconds(),
		CursorUpdatedAt: req.GetCursorUpdatedAt(),
		CursorId:        req.GetCursorId(),
		NamePrefix:      req.GetNamePrefix(),
		ItemNamePrefix:  req.GetItemNamePrefix(),
	})
	if err != nil {
		return &mes.BaseResponsePageEngineeringOrderVO{Code: 1, Message: err.Error()}, nil
	}
	records := make([]*mes.EngineeringOrderVO, 0, len(res.GetEngineeringOrderList()))
	for _, item := range res.GetEngineeringOrderList() {
		records = append(records, toEngineeringOrderVO(item))
	}
	page := pageEngineeringOrder(records, req.GetPageNum(), req.GetPageSize(), res.GetTotal())
	page.HasMore = res.GetHasMore()
	page.NextCursorUpdatedAt = res.GetNextCursorUpdatedAt()
	page.NextCursorId = res.GetNextCursorId()
	return &mes.BaseResponsePageEngineeringOrderVO{
		Code:    0,
		Message: "success",
		Data:    page,
	}, nil
}
