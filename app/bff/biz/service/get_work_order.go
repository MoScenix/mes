package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcworkorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetWorkOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetWorkOrderService(Context context.Context, RequestContext *app.RequestContext) *GetWorkOrderService {
	return &GetWorkOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *GetWorkOrderService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseWorkOrderVO, err error) {
	res, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(h.Context), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseWorkOrderVO{Code: 1, Message: err.Error()}, nil
	}
	if err = requireCanViewWorkOrder(h.Context, res.GetWorkOrder()); err != nil {
		return &mes.BaseResponseWorkOrderVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseWorkOrderVO{Code: 0, Message: "success", Data: toWorkOrderVO(res.GetWorkOrder())}, nil
}
