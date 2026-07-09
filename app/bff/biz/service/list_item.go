package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
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
	return runListItem(h.Context, req.GetNamePrefix(), req.GetPageNum(), req.GetPageSize(), req.GetCursorName(), req.GetCursorId())
}
