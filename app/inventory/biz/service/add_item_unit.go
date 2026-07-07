package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type AddItemUnitService struct {
	ctx context.Context
} // NewAddItemUnitService new AddItemUnitService
func NewAddItemUnitService(ctx context.Context) *AddItemUnitService {
	return &AddItemUnitService{ctx: ctx}
}

// Run create note info
func (s *AddItemUnitService) Run(req *inventory.AddItemUnitReq) (resp *inventory.AddItemUnitResp, err error) {
	return runAddItemUnit(s.ctx, req)
}
