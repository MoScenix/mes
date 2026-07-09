package service

import "context"

type UpdateEngineeringOrderDraftService struct {
	ctx context.Context
} // NewUpdateEngineeringOrderDraftService new UpdateEngineeringOrderDraftService
func NewUpdateEngineeringOrderDraftService(ctx context.Context) *UpdateEngineeringOrderDraftService {
	return &UpdateEngineeringOrderDraftService{ctx: ctx}
}

// Run create note info
