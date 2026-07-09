package service

import "context"

type UpdateItemService struct {
	ctx context.Context
} // NewUpdateItemService new UpdateItemService
func NewUpdateItemService(ctx context.Context) *UpdateItemService {
	return &UpdateItemService{ctx: ctx}
}

// Run create note info
