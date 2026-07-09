package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListProcessService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListProcessService(Context context.Context, RequestContext *app.RequestContext) *ListProcessService {
	return &ListProcessService{RequestContext: RequestContext, Context: Context}
}

func (h *ListProcessService) Run(req *mes.ListProcessRequest) (resp *mes.BaseResponsePageProcessVO, err error) {
	ownerUserID := req.GetOwnerUserId()
	status := rpcinventory.DraftStatus(req.GetStatus())
	switch req.GetScope() {
	case mes.MesListScope_MES_LIST_SCOPE_MINE:
		currentUserID, err := requireBFFUserID(h.Context)
		if err != nil {
			return &mes.BaseResponsePageProcessVO{Code: 1, Message: err.Error()}, nil
		}
		ownerUserID = currentUserID
	case mes.MesListScope_MES_LIST_SCOPE_ALL:
		ownerUserID = 0
	default:
		if !bffIsAdmin(h.Context) {
			currentUserID, err := requireBFFUserID(h.Context)
			if err != nil {
				return &mes.BaseResponsePageProcessVO{Code: 1, Message: err.Error()}, nil
			}
			if ownerUserID != 0 && ownerUserID != currentUserID {
				return &mes.BaseResponsePageProcessVO{Code: 1, Message: errForbiddenAccess.Error()}, nil
			}
			if status == rpcinventory.DraftStatus_DRAFT_STATUS_DRAFT {
				ownerUserID = currentUserID
			}
			if status == rpcinventory.DraftStatus_DRAFT_STATUS_UNKNOWN && ownerUserID == 0 {
				status = rpcinventory.DraftStatus_DRAFT_STATUS_SUBMITTED
			}
		}
	}
	res, err := rpc.InventoryClient.ListProcess(mesCtx(h.Context), &rpcinventory.ListProcessReq{
		OwnerUserId:     ownerUserID,
		ItemId:          req.GetItemId(),
		Status:          status,
		PageNum:         req.GetPageNum(),
		PageSize:        req.GetPageSize(),
		SinceTime:       req.GetSinceTime(),
		RecentSeconds:   req.GetRecentSeconds(),
		CursorUpdatedAt: req.GetCursorUpdatedAt(),
		CursorId:        req.GetCursorId(),
		NamePrefix:      req.GetNamePrefix(),
		ItemNamePrefix:  req.GetItemNamePrefix(),
	})
	if err != nil {
		return &mes.BaseResponsePageProcessVO{Code: 1, Message: err.Error()}, nil
	}
	records := make([]*mes.ProcessVO, 0, len(res.GetProcessList()))
	for _, item := range res.GetProcessList() {
		records = append(records, toProcessVO(item))
	}
	page := pageProcess(records, req.GetPageNum(), req.GetPageSize(), res.GetTotal())
	page.HasMore = res.GetHasMore()
	page.NextCursorUpdatedAt = res.GetNextCursorUpdatedAt()
	page.NextCursorId = res.GetNextCursorId()
	return &mes.BaseResponsePageProcessVO{
		Code:    0,
		Message: "success",
		Data:    page,
	}, nil
}
