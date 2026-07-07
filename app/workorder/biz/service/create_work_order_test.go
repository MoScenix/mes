package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"testing"
)

func TestCreateWorkOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCreateWorkOrderService(ctx)
	// init req and assert value

	req := &workorder.CreateWorkOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
