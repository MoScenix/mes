package service

import "context"

type CreateEngineeringOrderDraftService struct {
	ctx context.Context
} // NewCreateEngineeringOrderDraftService new CreateEngineeringOrderDraftService
func NewCreateEngineeringOrderDraftService(ctx context.Context) *CreateEngineeringOrderDraftService {
	return &CreateEngineeringOrderDraftService{ctx: ctx}
}

// Run create note info
