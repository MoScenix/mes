package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

type MarkWorkOrderReadService struct {
	ctx context.Context
} // NewMarkWorkOrderReadService new MarkWorkOrderReadService
func NewMarkWorkOrderReadService(ctx context.Context) *MarkWorkOrderReadService {
	return &MarkWorkOrderReadService{ctx: ctx}
}

// Run create note info
func (s *MarkWorkOrderReadService) Run(req *workorder.MarkWorkOrderReadReq) (resp *workorder.MarkWorkOrderReadResp, err error) {
	q, err := newWorkOrderQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	if err := q.MarkRead(req.GetId()); err != nil {
		return nil, err
	}

	return &workorder.MarkWorkOrderReadResp{Success: true}, nil
}
