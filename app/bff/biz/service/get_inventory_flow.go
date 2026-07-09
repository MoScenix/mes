package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetInventoryFlowService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetInventoryFlowService(Context context.Context, RequestContext *app.RequestContext) *GetInventoryFlowService {
	return &GetInventoryFlowService{RequestContext: RequestContext, Context: Context}
}

func (h *GetInventoryFlowService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseInventoryFlowVO, err error) {
	return runGetInventoryFlow(h.Context, req)
}
