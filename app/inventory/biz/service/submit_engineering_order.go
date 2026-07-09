package service

import "context"

type SubmitEngineeringOrderService struct {
	ctx context.Context
} // NewSubmitEngineeringOrderService new SubmitEngineeringOrderService
func NewSubmitEngineeringOrderService(ctx context.Context) *SubmitEngineeringOrderService {
	return &SubmitEngineeringOrderService{ctx: ctx}
}

// Run create note info
