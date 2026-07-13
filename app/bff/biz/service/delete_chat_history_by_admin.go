package service

import (
	"context"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	lapp "github.com/MoScenix/mes/app/bff/hertz_gen/bff/app"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcapp "github.com/MoScenix/mes/rpc_gen/kitex_gen/app"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteChatHistoryByAdminService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteChatHistoryByAdminService(Context context.Context, RequestContext *app.RequestContext) *DeleteChatHistoryByAdminService {
	return &DeleteChatHistoryByAdminService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteChatHistoryByAdminService) Run(req *lapp.DeleteRequest) (resp *lapp.BaseResponseBoolean, err error) {
	res, err := rpc.AppClient.DeleteMessage(utils.WithIdentityMeta(h.Context), &rpcapp.DeleteMessageReq{
		Id: req.Id,
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
