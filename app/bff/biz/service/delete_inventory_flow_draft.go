package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteInventoryFlowDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteInventoryFlowDraftService(Context context.Context, RequestContext *app.RequestContext) *DeleteInventoryFlowDraftService {
	return &DeleteInventoryFlowDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteInventoryFlowDraftService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runDeleteInventoryFlowDraft(h.Context, req)
}
