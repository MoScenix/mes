package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcworkorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"github.com/cloudwego/hertz/pkg/app"
)

type SubmitWorkOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSubmitWorkOrderService(Context context.Context, RequestContext *app.RequestContext) *SubmitWorkOrderService {
	return &SubmitWorkOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *SubmitWorkOrderService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(h.Context), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateWorkOrderDraft(h.Context, current.GetWorkOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.WorkOrderClient.SubmitWorkOrder(mesCtx(h.Context), &rpcworkorder.SubmitWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
