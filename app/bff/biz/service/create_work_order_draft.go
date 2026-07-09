package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type CreateWorkOrderDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateWorkOrderDraftService(Context context.Context, RequestContext *app.RequestContext) *CreateWorkOrderDraftService {
	return &CreateWorkOrderDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateWorkOrderDraftService) Run(req *mes.CreateWorkOrderDraftRequest) (resp *mes.BaseResponseLong, err error) {
	return runCreateWorkOrderDraft(h.Context, req)
}
