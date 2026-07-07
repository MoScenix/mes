package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"testing"
)

func TestListWorkOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewListWorkOrderService(ctx)
	// init req and assert value

	req := &workorder.ListWorkOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
