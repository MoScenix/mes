package service

import "context"

type ListItemUnitService struct {
	ctx context.Context
} // NewListItemUnitService new ListItemUnitService
func NewListItemUnitService(ctx context.Context) *ListItemUnitService {
	return &ListItemUnitService{ctx: ctx}
}

// Run create note info
