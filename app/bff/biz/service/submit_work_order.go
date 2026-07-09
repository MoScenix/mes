package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type SubmitWorkOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSubmitWorkOrderService(Context context.Context, RequestContext *app.RequestContext) *SubmitWorkOrderService {
	return &SubmitWorkOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *SubmitWorkOrderService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runSubmitWorkOrder(h.Context, req)
}
