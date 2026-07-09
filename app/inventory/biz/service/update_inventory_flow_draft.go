package service

import "context"

type UpdateInventoryFlowDraftService struct {
	ctx context.Context
} // NewUpdateInventoryFlowDraftService new UpdateInventoryFlowDraftService
func NewUpdateInventoryFlowDraftService(ctx context.Context) *UpdateInventoryFlowDraftService {
	return &UpdateInventoryFlowDraftService{ctx: ctx}
}

// Run create note info
