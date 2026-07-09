package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListItemUnitService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListItemUnitService(Context context.Context, RequestContext *app.RequestContext) *ListItemUnitService {
	return &ListItemUnitService{RequestContext: RequestContext, Context: Context}
}

func (h *ListItemUnitService) Run(req *mes.ListItemUnitRequest) (resp *mes.BaseResponsePageItemUnitVO, err error) {
	res, err := rpc.InventoryClient.ListItemUnit(mesCtx(h.Context), &rpcinventory.ListItemUnitReq{
		PageNum:            req.GetPageNum(),
		PageSize:           req.GetPageSize(),
		ItemId:             req.GetItemId(),
		StockStatus:        rpcinventory.StockStatus(req.GetStockStatus()),
		QualityStatus:      rpcinventory.QualityStatus(req.GetQualityStatus()),
		EngineeringOrderId: req.GetEngineeringOrderId(),
		CursorId:           req.GetCursorId(),
		ItemNamePrefix:     req.GetItemNamePrefix(),
		InventoryFlowId:    req.GetInventoryFlowId(),
		CursorUpdatedAt:    req.GetCursorUpdatedAt(),
	})
	if err != nil {
		return &mes.BaseResponsePageItemUnitVO{Code: 1, Message: err.Error()}, nil
	}
	records := make([]*mes.ItemUnitVO, 0, len(res.GetItemUnitList()))
	for _, item := range res.GetItemUnitList() {
		records = append(records, toItemUnitVO(item))
	}
	page := pageItemUnit(records, req.GetPageNum(), req.GetPageSize(), res.GetTotal())
	page.HasMore = res.GetHasMore()
	page.NextCursorUpdatedAt = res.GetNextCursorUpdatedAt()
	page.NextCursorId = res.GetNextCursorId()
	return &mes.BaseResponsePageItemUnitVO{
		Code:    0,
		Message: "success",
		Data:    page,
	}, nil
}
