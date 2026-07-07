package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type GetItemService struct {
	ctx context.Context
} // NewGetItemService new GetItemService
func NewGetItemService(ctx context.Context) *GetItemService {
	return &GetItemService{ctx: ctx}
}

// Run create note info
func (s *GetItemService) Run(req *inventory.GetItemReq) (resp *inventory.GetItemResp, err error) {
	return runGetItem(s.ctx, req)
}
