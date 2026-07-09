package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
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
	currentUserID, err := requireBFFUserID(h.Context)
	if err != nil {
		return mesLongErr(err), nil
	}
	res, err := rpc.InventoryClient.CreateInventoryFlow(mesCtx(h.Context), &rpcinventory.CreateInventoryFlowReq{
		FromUserId:  currentUserID,
		ToUserId:    req.GetToUserId(),
		FlowType:    rpcinventory.FlowType(req.GetFlowType()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Items:       toRPCInventoryFlowItems(req.GetItems()),
		ItemUnitIds: req.GetItemUnitIds(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}
