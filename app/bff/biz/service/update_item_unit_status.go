package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateItemUnitStatusService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateItemUnitStatusService(Context context.Context, RequestContext *app.RequestContext) *UpdateItemUnitStatusService {
	return &UpdateItemUnitStatusService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateItemUnitStatusService) Run(req *mes.UpdateItemUnitStatusRequest) (resp *mes.BaseResponseBoolean, err error) {
	res, err := rpc.InventoryClient.UpdateItemUnitStatus(mesCtx(h.Context), &rpcinventory.UpdateItemUnitStatusReq{
		Id:            req.GetId(),
		StockStatus:   rpcinventory.StockStatus(req.GetStockStatus()),
		QualityStatus: rpcinventory.QualityStatus(req.GetQualityStatus()),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
