package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type DeleteProcessDraftService struct {
	ctx context.Context
} // NewDeleteProcessDraftService new DeleteProcessDraftService
func NewDeleteProcessDraftService(ctx context.Context) *DeleteProcessDraftService {
	return &DeleteProcessDraftService{ctx: ctx}
}

// Run create note info
func (s *DeleteProcessDraftService) Run(req *inventory.DeleteProcessDraftReq) (resp *inventory.DeleteProcessDraftResp, err error) {
	return runDeleteProcessDraft(s.ctx, req)
}
