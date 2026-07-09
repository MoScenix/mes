package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/hertz/pkg/app"
)

type CreateEngineeringOrderDraftService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateEngineeringOrderDraftService(Context context.Context, RequestContext *app.RequestContext) *CreateEngineeringOrderDraftService {
	return &CreateEngineeringOrderDraftService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateEngineeringOrderDraftService) Run(req *mes.CreateEngineeringOrderRequest) (resp *mes.BaseResponseLong, err error) {
	leaderUserID, err := requireSameUserOrAdmin(h.Context, req.GetLeaderUserId())
	if err != nil {
		return mesLongErr(err), nil
	}
	res, err := rpc.InventoryClient.CreateEngineeringOrderDraft(mesCtx(h.Context), &rpcinventory.CreateEngineeringOrderDraftReq{
		LeaderUserId:      leaderUserID,
		ItemId:            req.GetItemId(),
		Name:              req.GetName(),
		ExpectedQuantity:  req.GetExpectedQuantity(),
		QualifiedQuantity: req.GetQualifiedQuantity(),
		Description:       req.GetDescription(),
		ProcessId:         req.GetProcessId(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}
