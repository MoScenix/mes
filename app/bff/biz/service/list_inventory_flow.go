package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListInventoryFlowService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListInventoryFlowService(Context context.Context, RequestContext *app.RequestContext) *ListInventoryFlowService {
	return &ListInventoryFlowService{RequestContext: RequestContext, Context: Context}
}

func (h *ListInventoryFlowService) Run(req *mes.ListInventoryFlowRequest) (resp *mes.BaseResponsePageInventoryFlowVO, err error) {
	return runListInventoryFlow(h.Context, req)
}
