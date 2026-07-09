package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteWorkOrderDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteWorkOrderDraftService(Context context.Context, RequestContext *app.RequestContext) *DeleteWorkOrderDraftService {
	return &DeleteWorkOrderDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteWorkOrderDraftService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runDeleteWorkOrderDraft(h.Context, req)
}
