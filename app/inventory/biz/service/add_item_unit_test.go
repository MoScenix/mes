package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"testing"
)

func TestAddItemUnit_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAddItemUnitService(ctx)
	// init req and assert value

	req := &inventory.AddItemUnitReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
