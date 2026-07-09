package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"testing"
)

func TestSubmitProcess_Run(t *testing.T) {
	ctx := context.Background()
	s := NewSubmitProcessService(ctx)
	// init req and assert value

	req := &inventory.SubmitProcessReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
