package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type SearchItemsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSearchItemsService(Context context.Context, RequestContext *app.RequestContext) *SearchItemsService {
	return &SearchItemsService{RequestContext: RequestContext, Context: Context}
}

func (h *SearchItemsService) Run(req *mes.SearchItemsRequest) (resp *mes.BaseResponsePageItemVO, err error) {
	return runListItem(h.Context, req.GetNamePrefix(), req.GetPageNum(), req.GetPageSize(), req.GetCursorName(), req.GetCursorId())
}
