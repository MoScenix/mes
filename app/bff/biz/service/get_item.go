package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetItemService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetItemService(Context context.Context, RequestContext *app.RequestContext) *GetItemService {
	return &GetItemService{RequestContext: RequestContext, Context: Context}
}

func (h *GetItemService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseItemVO, err error) {
	res, err := rpc.InventoryClient.GetItem(mesCtx(h.Context), &rpcinventory.GetItemReq{Id: req.GetId()})
	if err != nil {
		return &mes.BaseResponseItemVO{Code: 1, Message: err.Error()}, nil
	}
	return &mes.BaseResponseItemVO{Code: 0, Message: "success", Data: toItemVO(res.GetItem())}, nil
}
