package service

import "context"

type GetItemUnitService struct {
	ctx context.Context
} // NewGetItemUnitService new GetItemUnitService
func NewGetItemUnitService(ctx context.Context) *GetItemUnitService {
	return &GetItemUnitService{ctx: ctx}
}

// Run create note info
