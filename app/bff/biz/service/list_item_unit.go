package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
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
	return runListItemUnit(h.Context, req)
}
