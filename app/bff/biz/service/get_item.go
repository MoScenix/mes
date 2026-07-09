package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetItemService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetItemService(Context context.Context, RequestContext *app.RequestContext) *GetItemService {
	return &GetItemService{RequestContext: RequestContext, Context: Context}
}

func (h *GetItemService) Run(req *mes.GetByIdRequest) (resp *mes.BaseResponseItemVO, err error) {
	return runGetItem(h.Context, req)
}
