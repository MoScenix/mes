package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateProcessDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateProcessDraftService(Context context.Context, RequestContext *app.RequestContext) *UpdateProcessDraftService {
	return &UpdateProcessDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateProcessDraftService) Run(req *mes.UpdateProcessDraftRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.InventoryClient.GetProcess(mesCtx(h.Context), &rpcinventory.GetProcessReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateProcessDraft(h.Context, current.GetProcess()); err != nil {
		return mesBoolErr(err), nil
	}
	ownerUserID, err := requireSameUserOrAdmin(h.Context, req.GetOwnerUserId())
	if err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.UpdateProcessDraft(mesCtx(h.Context), &rpcinventory.UpdateProcessDraftReq{
		Id:          req.GetId(),
		OwnerUserId: ownerUserID,
		ItemId:      req.GetItemId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Items:       toRPCProcessItems(req.GetItems()),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
