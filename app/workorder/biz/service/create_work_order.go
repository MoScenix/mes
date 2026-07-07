package service

import (
	"context"

	"github.com/MoScenix/mes/app/workorder/biz/model"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

type CreateWorkOrderService struct {
	ctx context.Context
} // NewCreateWorkOrderService new CreateWorkOrderService
func NewCreateWorkOrderService(ctx context.Context) *CreateWorkOrderService {
	return &CreateWorkOrderService{ctx: ctx}
}

// Run create note info
func (s *CreateWorkOrderService) Run(req *workorder.CreateWorkOrderReq) (resp *workorder.CreateWorkOrderResp, err error) {
	q, err := newWorkOrderQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	order, err := q.CreateWorkOrder(model.WorkOrder{
		FromUserID:  req.GetFromUserId(),
		ToUserID:    req.GetToUserId(),
		Description: req.GetDescription(),
		Status:      model.WorkOrderStatusDraft,
	})
	if err != nil {
		return nil, err
	}

	return &workorder.CreateWorkOrderResp{Id: int64(order.ID)}, nil
}
