package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type ListItemUnitService struct {
	ctx context.Context
} // NewListItemUnitService new ListItemUnitService
func NewListItemUnitService(ctx context.Context) *ListItemUnitService {
	return &ListItemUnitService{ctx: ctx}
}

// Run create note info
func (s *ListItemUnitService) Run(req *inventory.ListItemUnitReq) (resp *inventory.ListItemUnitResp, err error) {
	return runListItemUnit(s.ctx, req)
}
