package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type AuditInventoryFlowService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAuditInventoryFlowService(Context context.Context, RequestContext *app.RequestContext) *AuditInventoryFlowService {
	return &AuditInventoryFlowService{RequestContext: RequestContext, Context: Context}
}

func (h *AuditInventoryFlowService) Run(req *mes.AuditInventoryFlowRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runAuditInventoryFlow(h.Context, req)
}
