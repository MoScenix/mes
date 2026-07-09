package service

import "context"

type CompleteInventoryFlowService struct {
	ctx context.Context
}

func NewCompleteInventoryFlowService(ctx context.Context) *CompleteInventoryFlowService {
	return &CompleteInventoryFlowService{ctx: ctx}
}
