package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type UpdateInventoryFlowDraftService struct {
	ctx context.Context
} // NewUpdateInventoryFlowDraftService new UpdateInventoryFlowDraftService
func NewUpdateInventoryFlowDraftService(ctx context.Context) *UpdateInventoryFlowDraftService {
	return &UpdateInventoryFlowDraftService{ctx: ctx}
}

// Run create note info
func (s *UpdateInventoryFlowDraftService) Run(req *inventory.UpdateInventoryFlowDraftReq) (resp *inventory.UpdateInventoryFlowDraftResp, err error) {
	return runUpdateInventoryFlowDraft(s.ctx, req)
}
