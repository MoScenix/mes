package service

import "context"

type ListItemService struct {
	ctx context.Context
} // NewListItemService new ListItemService
func NewListItemService(ctx context.Context) *ListItemService {
	return &ListItemService{ctx: ctx}
}

// Run create note info
