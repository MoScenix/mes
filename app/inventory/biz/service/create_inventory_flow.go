package service

import "context"

type CreateInventoryFlowService struct {
	ctx context.Context
} // NewCreateInventoryFlowService new CreateInventoryFlowService
func NewCreateInventoryFlowService(ctx context.Context) *CreateInventoryFlowService {
	return &CreateInventoryFlowService{ctx: ctx}
}

// Run create note info
