package service

import "context"

type ListProcessService struct {
	ctx context.Context
} // NewListProcessService new ListProcessService
func NewListProcessService(ctx context.Context) *ListProcessService {
	return &ListProcessService{ctx: ctx}
}

// Run create note info
