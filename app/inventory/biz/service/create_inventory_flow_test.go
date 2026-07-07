package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"testing"
)

func TestCreateInventoryFlow_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCreateInventoryFlowService(ctx)
	// init req and assert value

	req := &inventory.CreateInventoryFlowReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
