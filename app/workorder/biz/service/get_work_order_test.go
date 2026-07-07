package service

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"testing"
)

func TestGetWorkOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetWorkOrderService(ctx)
	// init req and assert value

	req := &workorder.GetWorkOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
