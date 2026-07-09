package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetWorkOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetWorkOrderService(Context context.Context, RequestContext *app.RequestContext) *GetWorkOrderService {
	return &GetWorkOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *GetWorkOrderService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseWorkOrderVO, err error) {
	return runGetWorkOrder(h.Context, req)
}
