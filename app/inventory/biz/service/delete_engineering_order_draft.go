package service

import "context"

type DeleteEngineeringOrderDraftService struct {
	ctx context.Context
} // NewDeleteEngineeringOrderDraftService new DeleteEngineeringOrderDraftService
func NewDeleteEngineeringOrderDraftService(ctx context.Context) *DeleteEngineeringOrderDraftService {
	return &DeleteEngineeringOrderDraftService{ctx: ctx}
}

// Run create note info
