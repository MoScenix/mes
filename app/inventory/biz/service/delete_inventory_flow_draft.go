package service

import "context"

type DeleteInventoryFlowDraftService struct {
	ctx context.Context
} // NewDeleteInventoryFlowDraftService new DeleteInventoryFlowDraftService
func NewDeleteInventoryFlowDraftService(ctx context.Context) *DeleteInventoryFlowDraftService {
	return &DeleteInventoryFlowDraftService{ctx: ctx}
}

// Run create note info
