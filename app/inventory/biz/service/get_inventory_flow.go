package service

import "context"

type GetInventoryFlowService struct {
	ctx context.Context
} // NewGetInventoryFlowService new GetInventoryFlowService
func NewGetInventoryFlowService(ctx context.Context) *GetInventoryFlowService {
	return &GetInventoryFlowService{ctx: ctx}
}

// Run create note info
