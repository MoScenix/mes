package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type CreateInventoryFlowService struct {
	ctx context.Context
} // NewCreateInventoryFlowService new CreateInventoryFlowService
func NewCreateInventoryFlowService(ctx context.Context) *CreateInventoryFlowService {
	return &CreateInventoryFlowService{ctx: ctx}
}

// Run create note info
func (s *CreateInventoryFlowService) Run(req *inventory.CreateInventoryFlowReq) (resp *inventory.CreateInventoryFlowResp, err error) {
	return runCreateInventoryFlow(s.ctx, req)
}
