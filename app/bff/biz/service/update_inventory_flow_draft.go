package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
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
	currentUserID, err := requireBFFUserID(h.Context)
	if err != nil {
		return mesBoolErr(err), nil
	}
	current, err := rpc.InventoryClient.GetInventoryFlow(mesCtx(h.Context), &rpcinventory.GetInventoryFlowReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateInventoryFlowDraft(h.Context, current.GetInventoryFlow()); err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.UpdateInventoryFlowDraft(mesCtx(h.Context), &rpcinventory.UpdateInventoryFlowDraftReq{
		Id:          req.GetId(),
		FromUserId:  currentUserID,
		ToUserId:    req.GetToUserId(),
		FlowType:    rpcinventory.FlowType(req.GetFlowType()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Items:       toRPCInventoryFlowItems(req.GetItems()),
		ItemUnitIds: req.GetItemUnitIds(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
