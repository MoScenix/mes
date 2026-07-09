package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type CreateInventoryFlowDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateInventoryFlowDraftService(Context context.Context, RequestContext *app.RequestContext) *CreateInventoryFlowDraftService {
	return &CreateInventoryFlowDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateInventoryFlowDraftService) Run(req *mes.CreateInventoryFlowDraftRequest) (resp *mes.BaseResponseLong, err error) {
	return runCreateInventoryFlowDraft(h.Context, req)
}
