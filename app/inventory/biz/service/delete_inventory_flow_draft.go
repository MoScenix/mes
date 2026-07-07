package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type DeleteInventoryFlowDraftService struct {
	ctx context.Context
} // NewDeleteInventoryFlowDraftService new DeleteInventoryFlowDraftService
func NewDeleteInventoryFlowDraftService(ctx context.Context) *DeleteInventoryFlowDraftService {
	return &DeleteInventoryFlowDraftService{ctx: ctx}
}

// Run create note info
func (s *DeleteInventoryFlowDraftService) Run(req *inventory.DeleteInventoryFlowDraftReq) (resp *inventory.DeleteInventoryFlowDraftResp, err error) {
	return runDeleteInventoryFlowDraft(s.ctx, req)
}
