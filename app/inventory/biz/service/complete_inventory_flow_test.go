package service

import (
	"context"
	"testing"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

func TestCompleteInventoryFlow_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCompleteInventoryFlowService(ctx)
	// init req and assert value

	req := &inventory.CompleteInventoryFlowReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
