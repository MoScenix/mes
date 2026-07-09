package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddItemService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddItemService(Context context.Context, RequestContext *app.RequestContext) *AddItemService {
	return &AddItemService{RequestContext: RequestContext, Context: Context}
}

func (h *AddItemService) Run(req *mes.AddItemRequest) (resp *mes.BaseResponseLong, err error) {
	res, err := rpc.InventoryClient.AddItem(mesCtx(h.Context), &rpcinventory.AddItemReq{
		Name:        req.GetName(),
		Unit:        req.GetUnit(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}
