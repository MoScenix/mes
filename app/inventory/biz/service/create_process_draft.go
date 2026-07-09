package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type CreateProcessDraftService struct {
	ctx context.Context
} // NewCreateProcessDraftService new CreateProcessDraftService
func NewCreateProcessDraftService(ctx context.Context) *CreateProcessDraftService {
	return &CreateProcessDraftService{ctx: ctx}
}

// Run create note info
func (s *CreateProcessDraftService) Run(req *inventory.CreateProcessDraftReq) (resp *inventory.CreateProcessDraftResp, err error) {
	return runCreateProcessDraft(s.ctx, req)
}
