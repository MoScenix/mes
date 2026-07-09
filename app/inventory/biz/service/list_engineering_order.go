package service

import (
	"context"

	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type ListEngineeringOrderService struct {
	ctx context.Context
}

func NewListEngineeringOrderService(ctx context.Context) *ListEngineeringOrderService {
	return &ListEngineeringOrderService{ctx: ctx}
}

func (s *ListEngineeringOrderService) Run(req *inventory.ListEngineeringOrderReq) (*inventory.ListEngineeringOrderResp, error) {
	return runListEngineeringOrder(s.ctx, req)
}
