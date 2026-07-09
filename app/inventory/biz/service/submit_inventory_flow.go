package service

import "context"

type SubmitInventoryFlowService struct {
	ctx context.Context
} // NewSubmitInventoryFlowService new SubmitInventoryFlowService
func NewSubmitInventoryFlowService(ctx context.Context) *SubmitInventoryFlowService {
	return &SubmitInventoryFlowService{ctx: ctx}
}

// Run create note info
