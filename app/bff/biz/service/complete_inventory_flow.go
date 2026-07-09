package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type CompleteInventoryFlowService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCompleteInventoryFlowService(Context context.Context, RequestContext *app.RequestContext) *CompleteInventoryFlowService {
	return &CompleteInventoryFlowService{RequestContext: RequestContext, Context: Context}
}

func (h *CompleteInventoryFlowService) Run(req *mes.CompleteInventoryFlowRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(h.Context), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanViewInventoryFlow(h.Context, current.GetInventoryFlow()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.CompleteInventoryFlow(mesCtx(h.Context), &rpcinventory.CompleteInventoryFlowReq{
		Id:          req.GetId(),
		ItemUnitIds: req.GetItemUnitIds(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
