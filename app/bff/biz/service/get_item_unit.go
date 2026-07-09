package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
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
	return runGetItemUnit(h.Context, req)
}
