package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

type SubmitWorkOrderService struct {
	ctx context.Context
} // NewSubmitWorkOrderService new SubmitWorkOrderService
func NewSubmitWorkOrderService(ctx context.Context) *SubmitWorkOrderService {
	return &SubmitWorkOrderService{ctx: ctx}
}

// Run create note info
func (s *SubmitWorkOrderService) Run(req *workorder.SubmitWorkOrderReq) (resp *workorder.SubmitWorkOrderResp, err error) {
	q, err := newWorkOrderQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	if err := q.SubmitDraft(req.GetId()); err != nil {
		return nil, err
	}

	return &workorder.SubmitWorkOrderResp{Success: true}, nil
}
