package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateItemService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateItemService(Context context.Context, RequestContext *app.RequestContext) *UpdateItemService {
	return &UpdateItemService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateItemService) Run(req *mes.UpdateItemRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runUpdateItem(h.Context, req)
}
