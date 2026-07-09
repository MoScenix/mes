package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcworkorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListWorkOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListWorkOrderService(Context context.Context, RequestContext *app.RequestContext) *ListWorkOrderService {
	return &ListWorkOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *ListWorkOrderService) Run(req *mes.ListWorkOrderRequest) (resp *mes.BaseResponsePageWorkOrderVO, err error) {
	isTo := req.GetIsTo()
	status := rpcworkorder.WorkOrderStatus(req.GetStatus())
	if status == rpcworkorder.WorkOrderStatus_WORK_ORDER_STATUS_DRAFT {
		isTo = false
	}
	userID, err := scopedUserID(h.Context, req.GetId())
	if err != nil {
		return &mes.BaseResponsePageWorkOrderVO{Code: 1, Message: err.Error()}, nil
	}
	res, err := rpc.WorkOrderClient.ListWorkOrder(mesCtx(h.Context), &rpcworkorder.ListWorkOrderReq{
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
