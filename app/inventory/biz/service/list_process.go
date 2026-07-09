package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type ListProcessService struct {
	ctx context.Context
} // NewListProcessService new ListProcessService
func NewListProcessService(ctx context.Context) *ListProcessService {
	return &ListProcessService{ctx: ctx}
}

// Run create note info
func (s *ListProcessService) Run(req *inventory.ListProcessReq) (resp *inventory.ListProcessResp, err error) {
	return runListProcess(s.ctx, req)
}
