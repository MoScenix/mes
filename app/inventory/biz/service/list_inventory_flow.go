package service

import "context"

type ListInventoryFlowService struct {
	ctx context.Context
} // NewListInventoryFlowService new ListInventoryFlowService
func NewListInventoryFlowService(ctx context.Context) *ListInventoryFlowService {
	return &ListInventoryFlowService{ctx: ctx}
}

// Run create note info
