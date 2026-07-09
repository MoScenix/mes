package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type SubmitEngineeringOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSubmitEngineeringOrderService(Context context.Context, RequestContext *app.RequestContext) *SubmitEngineeringOrderService {
	return &SubmitEngineeringOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *SubmitEngineeringOrderService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.InventoryClient.GetEngineeringOrder(mesCtx(h.Context), &rpcinventory.GetEngineeringOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateEngineeringOrder(h.Context, current.GetEngineeringOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.SubmitEngineeringOrder(mesCtx(h.Context), &rpcinventory.SubmitEngineeringOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
