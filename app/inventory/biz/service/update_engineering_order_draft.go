package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type UpdateEngineeringOrderDraftService struct {
	ctx context.Context
} // NewUpdateEngineeringOrderDraftService new UpdateEngineeringOrderDraftService
func NewUpdateEngineeringOrderDraftService(ctx context.Context) *UpdateEngineeringOrderDraftService {
	return &UpdateEngineeringOrderDraftService{ctx: ctx}
}

// Run create note info
func (s *UpdateEngineeringOrderDraftService) Run(req *inventory.UpdateEngineeringOrderDraftReq) (resp *inventory.UpdateEngineeringOrderDraftResp, err error) {
	return runUpdateEngineeringOrderDraft(s.ctx, req)
}
