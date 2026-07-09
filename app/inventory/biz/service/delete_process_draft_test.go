package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"testing"
)

func TestDeleteProcessDraft_Run(t *testing.T) {
	ctx := context.Background()
	s := NewDeleteProcessDraftService(ctx)
	// init req and assert value

	req := &inventory.DeleteProcessDraftReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
