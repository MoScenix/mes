package service

import (
	"context"

	"github.com/MoScenix/mes/app/workorder/biz/model"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

type ListWorkOrderService struct {
	ctx context.Context
} // NewListWorkOrderService new ListWorkOrderService
func NewListWorkOrderService(ctx context.Context) *ListWorkOrderService {
	return &ListWorkOrderService{ctx: ctx}
}

// Run create note info
func (s *ListWorkOrderService) Run(req *workorder.ListWorkOrderReq) (resp *workorder.ListWorkOrderResp, err error) {
	q, err := newWorkOrderQuery(s.ctx)
	if err != nil {
		return nil, err
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
	isTo := req.GetIsTo()
	status := int32(req.GetStatus())
	if status == model.WorkOrderStatusDraft {
		isTo = false
	}
	orders, hasMore, err := q.ListWorkOrderByEmployee(req.GetId(), pageSize, isTo, req.GetIsUnread(), status, req.GetNamePrefix(), sinceTime, cursorUpdatedAt, req.GetCursorId())
	if err != nil {
		return nil, err
	}

	resp = &workorder.ListWorkOrderResp{Total: int64(len(orders)), HasMore: hasMore}
	for _, order := range orders {
		resp.WorkOrderList = append(resp.WorkOrderList, toWorkOrderInfo(order))
	}
	if len(orders) > 0 {
		last := orders[len(orders)-1]
		resp.NextCursorUpdatedAt = last.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.NextCursorId = int64(last.ID)
	}
	return resp, nil
}
