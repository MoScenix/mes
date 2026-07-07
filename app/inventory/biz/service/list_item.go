package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type ListItemService struct {
	ctx context.Context
} // NewListItemService new ListItemService
func NewListItemService(ctx context.Context) *ListItemService {
	return &ListItemService{ctx: ctx}
}

// Run create note info
func (s *ListItemService) Run(req *inventory.ListItemReq) (resp *inventory.ListItemResp, err error) {
	return runListItem(s.ctx, req)
}
