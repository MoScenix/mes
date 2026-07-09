package service

import "context"

type CreateProcessDraftService struct {
	ctx context.Context
} // NewCreateProcessDraftService new CreateProcessDraftService
func NewCreateProcessDraftService(ctx context.Context) *CreateProcessDraftService {
	return &CreateProcessDraftService{ctx: ctx}
}

// Run create note info
