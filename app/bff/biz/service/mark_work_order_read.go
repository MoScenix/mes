package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcworkorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"github.com/cloudwego/hertz/pkg/app"
)

type MarkWorkOrderReadService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewMarkWorkOrderReadService(Context context.Context, RequestContext *app.RequestContext) *MarkWorkOrderReadService {
	return &MarkWorkOrderReadService{RequestContext: RequestContext, Context: Context}
}

func (h *MarkWorkOrderReadService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(h.Context), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanMarkWorkOrderRead(h.Context, current.GetWorkOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.WorkOrderClient.MarkWorkOrderRead(mesCtx(h.Context), &rpcworkorder.MarkWorkOrderReadReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
