package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteProcessDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteProcessDraftService(Context context.Context, RequestContext *app.RequestContext) *DeleteProcessDraftService {
	return &DeleteProcessDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteProcessDraftService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.InventoryClient.GetProcess(mesCtx(h.Context), &rpcinventory.GetProcessReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateProcessDraft(h.Context, current.GetProcess()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.DeleteProcessDraft(mesCtx(h.Context), &rpcinventory.DeleteProcessDraftReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
