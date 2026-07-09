package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type SubmitInventoryFlowService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSubmitInventoryFlowService(Context context.Context, RequestContext *app.RequestContext) *SubmitInventoryFlowService {
	return &SubmitInventoryFlowService{RequestContext: RequestContext, Context: Context}
}

func (h *SubmitInventoryFlowService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runSubmitInventoryFlow(h.Context, req)
}
