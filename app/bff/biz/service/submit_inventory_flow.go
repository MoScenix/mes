package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type SubmitInventoryFlowService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSubmitInventoryFlowService(Context context.Context, RequestContext *app.RequestContext) *SubmitInventoryFlowService {
	return &SubmitInventoryFlowService{RequestContext: RequestContext, Context: Context}
}

func (h *SubmitInventoryFlowService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(h.Context), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateInventoryFlowDraft(h.Context, current.GetInventoryFlow()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.SubmitInventoryFlow(mesCtx(h.Context), &rpcinventory.SubmitInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
