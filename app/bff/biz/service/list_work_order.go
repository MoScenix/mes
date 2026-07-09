package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListWorkOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListWorkOrderService(Context context.Context, RequestContext *app.RequestContext) *ListWorkOrderService {
	return &ListWorkOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *ListWorkOrderService) Run(req *mes.ListWorkOrderRequest) (resp *mes.BaseResponsePageWorkOrderVO, err error) {
	return runListWorkOrder(h.Context, req)
}
