package service

import "context"

type DeleteProcessDraftService struct {
	ctx context.Context
} // NewDeleteProcessDraftService new DeleteProcessDraftService
func NewDeleteProcessDraftService(ctx context.Context) *DeleteProcessDraftService {
	return &DeleteProcessDraftService{ctx: ctx}
}

// Run create note info
