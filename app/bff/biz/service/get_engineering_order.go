package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetEngineeringOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetEngineeringOrderService(Context context.Context, RequestContext *app.RequestContext) *GetEngineeringOrderService {
	return &GetEngineeringOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *GetEngineeringOrderService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseEngineeringOrderVO, err error) {
	res, err := rpc.InventoryClient.GetEngineeringOrder(mesCtx(h.Context), &rpcinventory.GetEngineeringOrderReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseEngineeringOrderVO{Code: 1, Message: err.Error()}, nil
	}
	if err = requireCanViewEngineeringOrder(h.Context, res.GetEngineeringOrder()); err != nil {
		return &mes.BaseResponseEngineeringOrderVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseEngineeringOrderVO{Code: 0, Message: "success", Data: toEngineeringOrderVO(res.GetEngineeringOrder())}, nil
}
