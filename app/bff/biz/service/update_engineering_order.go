package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateEngineeringOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateEngineeringOrderService(Context context.Context, RequestContext *app.RequestContext) *UpdateEngineeringOrderService {
	return &UpdateEngineeringOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateEngineeringOrderService) Run(req *mes.UpdateEngineeringOrderRequest) (resp *mes.BaseResponseBoolean, err error) {
	current, err := rpc.InventoryClient.GetEngineeringOrder(mesCtx(h.Context), &rpcinventory.GetEngineeringOrderReq{Id: req.GetId()})
	if err != nil {
		return mesBoolErr(err), nil
	}
	if err = requireCanUpdateEngineeringOrder(h.Context, current.GetEngineeringOrder()); err != nil {
		return mesBoolErr(err), nil
	}
	leaderUserID, err := requireSameUserOrAdmin(h.Context, req.GetLeaderUserId())
	if err != nil {
		return mesBoolErr(err), nil
	}
	res, err := rpc.InventoryClient.UpdateEngineeringOrderDraft(mesCtx(h.Context), &rpcinventory.UpdateEngineeringOrderDraftReq{
		Id:                req.GetId(),
		LeaderUserId:      leaderUserID,
		ItemId:            req.GetItemId(),
		Name:              req.GetName(),
		ExpectedQuantity:  req.GetExpectedQuantity(),
		QualifiedQuantity: req.GetQualifiedQuantity(),
		Description:       req.GetDescription(),
		ProcessId:         req.GetProcessId(),
	})
	if err != nil {
		return mesBoolErr(err), nil
	}
	return mesBool(res.GetSuccess()), nil
}
