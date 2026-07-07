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

	pageNum, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	orders, total, err := q.ListWorkOrder(pageNum, pageSize, model.WorkOrderListFilter{
		FromUserID: req.GetFromUserId(),
		ToUserID:   req.GetToUserId(),
		Status:     int32(req.GetStatus()),
	})
	if err != nil {
		return nil, err
	}

	resp = &workorder.ListWorkOrderResp{
		Total: total,
	}
	for _, order := range orders {
		resp.WorkOrderList = append(resp.WorkOrderList, toWorkOrderInfo(order))
	}
	return resp, nil
}
