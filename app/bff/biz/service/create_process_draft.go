package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type CreateProcessDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateProcessDraftService(Context context.Context, RequestContext *app.RequestContext) *CreateProcessDraftService {
	return &CreateProcessDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateProcessDraftService) Run(req *mes.CreateProcessDraftRequest) (resp *mes.BaseResponseLong, err error) {
	ownerUserID, err := requireSameUserOrAdmin(h.Context, req.GetOwnerUserId())
	if err != nil {
		return mesLongErr(err), nil
	}
	res, err := rpc.InventoryClient.CreateProcessDraft(mesCtx(h.Context), &rpcinventory.CreateProcessDraftReq{
		OwnerUserId: ownerUserID,
		ItemId:      req.GetItemId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Items:       toRPCProcessItems(req.GetItems()),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}
