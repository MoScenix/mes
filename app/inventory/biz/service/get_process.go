package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type GetProcessService struct {
	ctx context.Context
} // NewGetProcessService new GetProcessService
func NewGetProcessService(ctx context.Context) *GetProcessService {
	return &GetProcessService{ctx: ctx}
}

// Run create note info
func (s *GetProcessService) Run(req *inventory.GetProcessReq) (resp *inventory.GetProcessResp, err error) {
	return runGetProcess(s.ctx, req)
}
