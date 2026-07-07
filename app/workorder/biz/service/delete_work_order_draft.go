package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

type DeleteWorkOrderDraftService struct {
	ctx context.Context
} // NewDeleteWorkOrderDraftService new DeleteWorkOrderDraftService
func NewDeleteWorkOrderDraftService(ctx context.Context) *DeleteWorkOrderDraftService {
	return &DeleteWorkOrderDraftService{ctx: ctx}
}

// Run create note info
func (s *DeleteWorkOrderDraftService) Run(req *workorder.DeleteWorkOrderDraftReq) (resp *workorder.DeleteWorkOrderDraftResp, err error) {
	q, err := newWorkOrderQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	if err := q.DeleteDraft(req.GetId()); err != nil {
		return nil, err
	}

	return &workorder.DeleteWorkOrderDraftResp{Success: true}, nil
}
