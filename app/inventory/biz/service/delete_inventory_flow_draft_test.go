package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"testing"
)

func TestDeleteInventoryFlowDraft_Run(t *testing.T) {
	ctx := context.Background()
	s := NewDeleteInventoryFlowDraftService(ctx)
	// init req and assert value

	req := &inventory.DeleteInventoryFlowDraftReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
