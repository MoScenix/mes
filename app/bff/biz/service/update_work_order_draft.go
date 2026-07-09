package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateWorkOrderDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateWorkOrderDraftService(Context context.Context, RequestContext *app.RequestContext) *UpdateWorkOrderDraftService {
	return &UpdateWorkOrderDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateWorkOrderDraftService) Run(req *mes.UpdateWorkOrderDraftRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runUpdateWorkOrderDraft(h.Context, req)
}
