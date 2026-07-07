package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type UpdateItemUnitStatusService struct {
	ctx context.Context
} // NewUpdateItemUnitStatusService new UpdateItemUnitStatusService
func NewUpdateItemUnitStatusService(ctx context.Context) *UpdateItemUnitStatusService {
	return &UpdateItemUnitStatusService{ctx: ctx}
}

// Run create note info
func (s *UpdateItemUnitStatusService) Run(req *inventory.UpdateItemUnitStatusReq) (resp *inventory.UpdateItemUnitStatusResp, err error) {
	return runUpdateItemUnitStatus(s.ctx, req)
}
