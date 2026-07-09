package service

import (
	"context"

	"github.com/MoScenix/mes/app/inventory/biz/model"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type GetEngineeringOrderService struct {
	ctx context.Context
}

func NewGetEngineeringOrderService(ctx context.Context) *GetEngineeringOrderService {
	return &GetEngineeringOrderService{ctx: ctx}
}

func (s *GetEngineeringOrderService) Run(req *inventory.GetEngineeringOrderReq) (*inventory.GetEngineeringOrderResp, error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	id, err := uintID(req.GetId(), "engineering order id")
	if err != nil {
		return nil, err
	}
	if err := model.RecalculateEngineeringOrderProducedQuantity(ctx, db, id); err != nil {
		return nil, err
	}
	order, err := model.NewEngineeringOrderQuery(ctx, db).Get(id, true)
	if err != nil {
		return nil, err
	}
	return &inventory.GetEngineeringOrderResp{EngineeringOrder: engineeringOrderInfo(order, false)}, nil
}
