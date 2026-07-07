package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type GetInventoryFlowService struct {
	ctx context.Context
} // NewGetInventoryFlowService new GetInventoryFlowService
func NewGetInventoryFlowService(ctx context.Context) *GetInventoryFlowService {
	return &GetInventoryFlowService{ctx: ctx}
}

// Run create note info
func (s *GetInventoryFlowService) Run(req *inventory.GetInventoryFlowReq) (resp *inventory.GetInventoryFlowResp, err error) {
	return runGetInventoryFlow(s.ctx, req)
}
