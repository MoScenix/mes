package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"testing"
)

func TestUpdateWorkOrderDraft_Run(t *testing.T) {
	ctx := context.Background()
	s := NewUpdateWorkOrderDraftService(ctx)
	// init req and assert value

	req := &workorder.UpdateWorkOrderDraftReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
