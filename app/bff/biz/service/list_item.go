package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListItemService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListItemService(Context context.Context, RequestContext *app.RequestContext) *ListItemService {
	return &ListItemService{RequestContext: RequestContext, Context: Context}
}

func (h *ListItemService) Run(req *mes.ListItemRequest) (resp *mes.BaseResponsePageItemVO, err error) {
	res, err := rpc.InventoryClient.ListItem(mesCtx(h.Context), &rpcinventory.ListItemReq{
		PageNum:         req.GetPageNum(),
		PageSize:        req.GetPageSize(),
		NamePrefix:      req.GetNamePrefix(),
		CursorUpdatedAt: req.GetCursorUpdatedAt(),
		CursorId:        req.GetCursorId(),
	})
	if err != nil {
		return &mes.BaseResponsePageItemVO{Code: 1, Message: err.Error()}, nil
	}
	records := make([]*mes.ItemVO, 0, len(res.GetItemList()))
	for _, item := range res.GetItemList() {
		records = append(records, toItemVO(item))
	}
	page := pageItem(records, req.GetPageNum(), req.GetPageSize(), res.GetTotal())
	page.HasMore = res.GetHasMore()
	page.NextCursorUpdatedAt = res.GetNextCursorUpdatedAt()
	page.NextCursorId = res.GetNextCursorId()
	return &mes.BaseResponsePageItemVO{
		Code:    0,
		Message: "success",
		Data:    page,
	}, nil
}
