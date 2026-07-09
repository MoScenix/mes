package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcworkorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
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
	currentUserID, err := requireBFFUserID(h.Context)
	if err != nil {
		return mesLongErr(err), nil
	}
	res, err := rpc.WorkOrderClient.CreateWorkOrder(mesCtx(h.Context), &rpcworkorder.CreateWorkOrderReq{
		FromUserId:  currentUserID,
		ToUserId:    req.GetToUserId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return mesLongErr(err), nil
	}
	return mesLong(res.GetId()), nil
}
