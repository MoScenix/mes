package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetProcessService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetProcessService(Context context.Context, RequestContext *app.RequestContext) *GetProcessService {
	return &GetProcessService{RequestContext: RequestContext, Context: Context}
}

func (h *GetProcessService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseProcessVO, err error) {
	res, err := rpc.InventoryClient.GetProcess(mesCtx(h.Context), &rpcinventory.GetProcessReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseProcessVO{Code: 1, Message: err.Error()}, nil
	}
	if err = requireCanViewProcess(h.Context, res.GetProcess()); err != nil {
		return &mes.BaseResponseProcessVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseProcessVO{Code: 0, Message: "success", Data: toProcessVO(res.GetProcess())}, nil
}
