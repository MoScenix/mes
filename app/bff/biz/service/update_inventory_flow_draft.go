package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateInventoryFlowDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateInventoryFlowDraftService(Context context.Context, RequestContext *app.RequestContext) *UpdateInventoryFlowDraftService {
	return &UpdateInventoryFlowDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateInventoryFlowDraftService) Run(req *mes.UpdateInventoryFlowDraftRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runUpdateInventoryFlowDraft(h.Context, req)
}
