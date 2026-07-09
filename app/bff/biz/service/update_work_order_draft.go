package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcworkorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
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
	currentUserID, err := requireBFFUserID(h.Context)
	if err != nil {
		return mesBoolErr(err), nil
	}
	current, err := rpc.WorkOrderClient.GetWorkOrder(mesCtx(h.Context), &rpcworkorder.GetWorkOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateWorkOrderDraft(h.Context, current.GetWorkOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.WorkOrderClient.UpdateWorkOrderDraft(mesCtx(h.Context), &rpcworkorder.UpdateWorkOrderDraftReq{
		Id:          req.GetId(),
		FromUserId:  currentUserID,
		ToUserId:    req.GetToUserId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
