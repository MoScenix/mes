package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddItemUnitService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddItemUnitService(Context context.Context, RequestContext *app.RequestContext) *AddItemUnitService {
	return &AddItemUnitService{RequestContext: RequestContext, Context: Context}
}

func (h *AddItemUnitService) Run(req *mes.AddItemUnitRequest) (resp *mes.BaseResponseLong, err error) {
	res, err := rpc.InventoryClient.AddItemUnit(mesCtx(h.Context), &rpcinventory.AddItemUnitReq{
		ItemId:             req.GetItemId(),
		StockStatus:        rpcinventory.StockStatus(req.GetStockStatus()),
		QualityStatus:      rpcinventory.QualityStatus(req.GetQualityStatus()),
		Description:        req.GetDescription(),
		EngineeringOrderId: req.GetEngineeringOrderId(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}
