package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

type UpdateWorkOrderDraftService struct {
	ctx context.Context
} // NewUpdateWorkOrderDraftService new UpdateWorkOrderDraftService
func NewUpdateWorkOrderDraftService(ctx context.Context) *UpdateWorkOrderDraftService {
	return &UpdateWorkOrderDraftService{ctx: ctx}
}

// Run create note info
func (s *UpdateWorkOrderDraftService) Run(req *workorder.UpdateWorkOrderDraftReq) (resp *workorder.UpdateWorkOrderDraftResp, err error) {
	q, err := newWorkOrderQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	if err := q.UpdateDraft(req.GetId(), req.GetFromUserId(), req.GetToUserId(), req.GetDescription()); err != nil {
		return nil, err
	}

	return &workorder.UpdateWorkOrderDraftResp{Success: true}, nil
}
