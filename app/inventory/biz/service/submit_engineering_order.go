package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type SubmitEngineeringOrderService struct {
	ctx context.Context
} // NewSubmitEngineeringOrderService new SubmitEngineeringOrderService
func NewSubmitEngineeringOrderService(ctx context.Context) *SubmitEngineeringOrderService {
	return &SubmitEngineeringOrderService{ctx: ctx}
}

// Run create note info
func (s *SubmitEngineeringOrderService) Run(req *inventory.SubmitEngineeringOrderReq) (resp *inventory.SubmitEngineeringOrderResp, err error) {
	return runSubmitEngineeringOrder(s.ctx, req)
}
