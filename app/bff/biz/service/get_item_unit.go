package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetItemUnitService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetItemUnitService(Context context.Context, RequestContext *app.RequestContext) *GetItemUnitService {
	return &GetItemUnitService{RequestContext: RequestContext, Context: Context}
}

func (h *GetItemUnitService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseItemUnitVO, err error) {
	res, err := rpc.InventoryClient.GetItemUnit(mesCtx(h.Context), &rpcinventory.GetItemUnitReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseItemUnitVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseItemUnitVO{Code: 0, Message: "success", Data: toItemUnitVO(res.GetItemUnit())}, nil
}
