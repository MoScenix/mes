package service

import "context"

type UpdateItemUnitStatusService struct {
	ctx context.Context
} // NewUpdateItemUnitStatusService new UpdateItemUnitStatusService
func NewUpdateItemUnitStatusService(ctx context.Context) *UpdateItemUnitStatusService {
	return &UpdateItemUnitStatusService{ctx: ctx}
}

// Run create note info
