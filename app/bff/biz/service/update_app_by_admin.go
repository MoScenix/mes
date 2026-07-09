package service

import (
	"context"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	lapp "github.com/MoScenix/mes/app/bff/hertz_gen/bff/app"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcapp "github.com/MoScenix/mes/rpc_gen/kitex_gen/app"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateAppByAdminService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateAppByAdminService(Context context.Context, RequestContext *app.RequestContext) *UpdateAppByAdminService {
	return &UpdateAppByAdminService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateAppByAdminService) Run(req *lapp.AppAdminUpdateRequest) (resp *lapp.BaseResponseBoolean, err error) {
	ctx := utils.WithIdentityMeta(h.Context)
	res, err := rpc.AppClient.UpdateApp(ctx, &rpcapp.UpdateAppReq{
		Id:      req.Id,
		AppName: req.AppName,
	})
	if err != nil {
		return &lapp.BaseResponseBoolean{
			Code:    1,
			Message: err.Error(),
		}, err
	}
	return &lapp.BaseResponseBoolean{
		Code:    0,
		Message: "success",
		Data:    res.Success,
	}, nil
}
