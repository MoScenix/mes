package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type SubmitProcessService struct {
	ctx context.Context
} // NewSubmitProcessService new SubmitProcessService
func NewSubmitProcessService(ctx context.Context) *SubmitProcessService {
	return &SubmitProcessService{ctx: ctx}
}

// Run create note info
func (s *SubmitProcessService) Run(req *inventory.SubmitProcessReq) (resp *inventory.SubmitProcessResp, err error) {
	return runSubmitProcess(s.ctx, req)
}
