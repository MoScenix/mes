package service

import (
	"context"

	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
)

type MarkWorkOrderReadService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewMarkWorkOrderReadService(Context context.Context, RequestContext *app.RequestContext) *MarkWorkOrderReadService {
	return &MarkWorkOrderReadService{RequestContext: RequestContext, Context: Context}
}

func (h *MarkWorkOrderReadService) Run(req *mes.DeleteRequest) (resp *mes.BaseResponseBoolean, err error) {
	return runMarkWorkOrderRead(h.Context, req)
}
