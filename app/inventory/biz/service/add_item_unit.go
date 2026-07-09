package service

import "context"

type AddItemUnitService struct {
	ctx context.Context
} // NewAddItemUnitService new AddItemUnitService
func NewAddItemUnitService(ctx context.Context) *AddItemUnitService {
	return &AddItemUnitService{ctx: ctx}
}

// Run create note info
