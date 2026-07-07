package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"testing"
)

func TestUpdateItem_Run(t *testing.T) {
	ctx := context.Background()
	s := NewUpdateItemService(ctx)
	// init req and assert value

	req := &inventory.UpdateItemReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
