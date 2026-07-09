package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
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
	if err := requireCanAuditInventoryFlow(h.Context); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.AuditInventoryFlow(mesCtx(h.Context), &rpcinventory.AuditInventoryFlowReq{
		Id:         req.GetId(),
		ApprovedBy: currentMESUserID(h.Context),
		Approved:   req.GetApproved(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
