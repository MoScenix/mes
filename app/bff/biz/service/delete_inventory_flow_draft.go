package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteInventoryFlowDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteInventoryFlowDraftService(Context context.Context, RequestContext *app.RequestContext) *DeleteInventoryFlowDraftService {
	return &DeleteInventoryFlowDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteInventoryFlowDraftService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(h.Context), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateInventoryFlowDraft(h.Context, current.GetInventoryFlow()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.DeleteInventoryFlowDraft(mesCtx(h.Context), &rpcinventory.DeleteInventoryFlowDraftReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
