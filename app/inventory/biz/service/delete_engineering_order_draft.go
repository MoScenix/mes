package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type DeleteEngineeringOrderDraftService struct {
	ctx context.Context
} // NewDeleteEngineeringOrderDraftService new DeleteEngineeringOrderDraftService
func NewDeleteEngineeringOrderDraftService(ctx context.Context) *DeleteEngineeringOrderDraftService {
	return &DeleteEngineeringOrderDraftService{ctx: ctx}
}

// Run create note info
func (s *DeleteEngineeringOrderDraftService) Run(req *inventory.DeleteEngineeringOrderDraftReq) (resp *inventory.DeleteEngineeringOrderDraftResp, err error) {
	return runDeleteEngineeringOrderDraft(s.ctx, req)
}
