package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetInventoryFlowService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetInventoryFlowService(Context context.Context, RequestContext *app.RequestContext) *GetInventoryFlowService {
	return &GetInventoryFlowService{RequestContext: RequestContext, Context: Context}
}

func (h *GetInventoryFlowService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseInventoryFlowVO, err error) {
	res, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(h.Context), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseInventoryFlowVO{Code: 1, Message: err.Error()}, nil
	}
	if err = requireCanViewInventoryFlow(h.Context, res.GetInventoryFlow()); err != nil {
		return &mes.BaseResponseInventoryFlowVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseInventoryFlowVO{Code: 0, Message: "success", Data: toInventoryFlowVO(res.GetInventoryFlow())}, nil
}
