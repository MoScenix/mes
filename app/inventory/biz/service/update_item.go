package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type UpdateItemService struct {
	ctx context.Context
} // NewUpdateItemService new UpdateItemService
func NewUpdateItemService(ctx context.Context) *UpdateItemService {
	return &UpdateItemService{ctx: ctx}
}

// Run create note info
func (s *UpdateItemService) Run(req *inventory.UpdateItemReq) (resp *inventory.UpdateItemResp, err error) {
	return runUpdateItem(s.ctx, req)
}
