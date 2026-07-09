package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type SubmitProcessService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSubmitProcessService(Context context.Context, RequestContext *app.RequestContext) *SubmitProcessService {
	return &SubmitProcessService{RequestContext: RequestContext, Context: Context}
}

func (h *SubmitProcessService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.InventoryClient.GetProcess(mesCtx(h.Context), &rpcinventory.GetProcessReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateProcessDraft(h.Context, current.GetProcess()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.SubmitProcess(mesCtx(h.Context), &rpcinventory.SubmitProcessReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
