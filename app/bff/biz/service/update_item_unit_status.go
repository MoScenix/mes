package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateItemUnitStatusService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateItemUnitStatusService(Context context.Context, RequestContext *app.RequestContext) *UpdateItemUnitStatusService {
	return &UpdateItemUnitStatusService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateItemUnitStatusService) Run(req *mes.UpdateItemUnitStatusRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runUpdateItemUnitStatus(h.Context, req)
}
