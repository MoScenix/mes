package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"testing"
)

func TestMarkWorkOrderRead_Run(t *testing.T) {
	ctx := context.Background()
	s := NewMarkWorkOrderReadService(ctx)
	// init req and assert value

	req := &workorder.MarkWorkOrderReadReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
