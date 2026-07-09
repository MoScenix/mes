package service

import (
	"context"

	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type CompleteInventoryFlowService struct {
	ctx context.Context
}

func NewCompleteInventoryFlowService(ctx context.Context) *CompleteInventoryFlowService {
	return &CompleteInventoryFlowService{ctx: ctx}
}

func (s *CompleteInventoryFlowService) Run(req *inventory.CompleteInventoryFlowReq) (resp *inventory.CompleteInventoryFlowResp, err error) {
	return runCompleteInventoryFlow(s.ctx, req)
}
