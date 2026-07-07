package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type ListInventoryFlowService struct {
	ctx context.Context
} // NewListInventoryFlowService new ListInventoryFlowService
func NewListInventoryFlowService(ctx context.Context) *ListInventoryFlowService {
	return &ListInventoryFlowService{ctx: ctx}
}

// Run create note info
func (s *ListInventoryFlowService) Run(req *inventory.ListInventoryFlowReq) (resp *inventory.ListInventoryFlowResp, err error) {
	return runListInventoryFlow(s.ctx, req)
}
