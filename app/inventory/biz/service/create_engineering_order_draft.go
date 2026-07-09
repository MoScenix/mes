package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type CreateEngineeringOrderDraftService struct {
	ctx context.Context
} // NewCreateEngineeringOrderDraftService new CreateEngineeringOrderDraftService
func NewCreateEngineeringOrderDraftService(ctx context.Context) *CreateEngineeringOrderDraftService {
	return &CreateEngineeringOrderDraftService{ctx: ctx}
}

// Run create note info
func (s *CreateEngineeringOrderDraftService) Run(req *inventory.CreateEngineeringOrderDraftReq) (resp *inventory.CreateEngineeringOrderDraftResp, err error) {
	return runCreateEngineeringOrderDraft(s.ctx, req)
}
