package service

import "context"

type UpdateProcessDraftService struct {
	ctx context.Context
} // NewUpdateProcessDraftService new UpdateProcessDraftService
func NewUpdateProcessDraftService(ctx context.Context) *UpdateProcessDraftService {
	return &UpdateProcessDraftService{ctx: ctx}
}

// Run create note info
