package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

type GetWorkOrderService struct {
	ctx context.Context
} // NewGetWorkOrderService new GetWorkOrderService
func NewGetWorkOrderService(ctx context.Context) *GetWorkOrderService {
	return &GetWorkOrderService{ctx: ctx}
}

// Run create note info
func (s *GetWorkOrderService) Run(req *workorder.GetWorkOrderReq) (resp *workorder.GetWorkOrderResp, err error) {
	q, err := newWorkOrderQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	order, err := q.GetWorkOrder(req.GetId())
	if err != nil {
		return nil, err
	}

	return &workorder.GetWorkOrderResp{WorkOrder: toWorkOrderInfo(order)}, nil
}
