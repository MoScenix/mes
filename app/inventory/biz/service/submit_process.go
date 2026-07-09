package service

import "context"

type SubmitProcessService struct {
	ctx context.Context
} // NewSubmitProcessService new SubmitProcessService
func NewSubmitProcessService(ctx context.Context) *SubmitProcessService {
	return &SubmitProcessService{ctx: ctx}
}

// Run create note info
