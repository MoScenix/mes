package service

import (
	"context"

	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type GetEngineeringOrderService struct {
	ctx context.Context
}

func NewGetEngineeringOrderService(ctx context.Context) *GetEngineeringOrderService {
	return &GetEngineeringOrderService{ctx: ctx}
}

func (s *GetEngineeringOrderService) Run(req *inventory.GetEngineeringOrderReq) (*inventory.GetEngineeringOrderResp, error) {
	return runGetEngineeringOrder(s.ctx, req)
}
