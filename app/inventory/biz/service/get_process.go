package service

import "context"

type GetProcessService struct {
	ctx context.Context
} // NewGetProcessService new GetProcessService
func NewGetProcessService(ctx context.Context) *GetProcessService {
	return &GetProcessService{ctx: ctx}
}

// Run create note info
