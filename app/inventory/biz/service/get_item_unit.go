package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type GetItemUnitService struct {
	ctx context.Context
} // NewGetItemUnitService new GetItemUnitService
func NewGetItemUnitService(ctx context.Context) *GetItemUnitService {
	return &GetItemUnitService{ctx: ctx}
}

// Run create note info
func (s *GetItemUnitService) Run(req *inventory.GetItemUnitReq) (resp *inventory.GetItemUnitResp, err error) {
	return runGetItemUnit(s.ctx, req)
}
