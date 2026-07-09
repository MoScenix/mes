package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
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
	return runAddItemUnit(h.Context, req)
}
