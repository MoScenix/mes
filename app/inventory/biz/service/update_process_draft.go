package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type UpdateProcessDraftService struct {
	ctx context.Context
} // NewUpdateProcessDraftService new UpdateProcessDraftService
func NewUpdateProcessDraftService(ctx context.Context) *UpdateProcessDraftService {
	return &UpdateProcessDraftService{ctx: ctx}
}

// Run create note info
func (s *UpdateProcessDraftService) Run(req *inventory.UpdateProcessDraftReq) (resp *inventory.UpdateProcessDraftResp, err error) {
	return runUpdateProcessDraft(s.ctx, req)
}
