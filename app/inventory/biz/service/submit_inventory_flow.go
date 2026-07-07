package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type SubmitInventoryFlowService struct {
	ctx context.Context
} // NewSubmitInventoryFlowService new SubmitInventoryFlowService
func NewSubmitInventoryFlowService(ctx context.Context) *SubmitInventoryFlowService {
	return &SubmitInventoryFlowService{ctx: ctx}
}

// Run create note info
func (s *SubmitInventoryFlowService) Run(req *inventory.SubmitInventoryFlowReq) (resp *inventory.SubmitInventoryFlowResp, err error) {
	return runSubmitInventoryFlow(s.ctx, req)
}
