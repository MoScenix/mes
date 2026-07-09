package service

import (
	"context"
	"time"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	lapp "github.com/MoScenix/mes/app/bff/hertz_gen/bff/app"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcapp "github.com/MoScenix/mes/rpc_gen/kitex_gen/app"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListAppChatHistoryService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListAppChatHistoryService(Context context.Context, RequestContext *app.RequestContext) *ListAppChatHistoryService {
	return &ListAppChatHistoryService{RequestContext: RequestContext, Context: Context}
}

func (h *ListAppChatHistoryService) Run(req *lapp.ListAppChatHistoryRequest) (resp *lapp.BaseResponsePageChatHistory, err error) {
	if err = requireAppOwnerOrAdmin(h.Context, req.GetAppId()); err != nil {
		return &lapp.BaseResponsePageChatHistory{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	q := rpc.AppClient
	res, err := q.ListAppMessage(utils.WithIdentityMeta(h.Context), &rpcapp.ListAppMessageReq{
		AppId:          req.AppId,
		PageSize:       20,
		LastCreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return &lapp.BaseResponsePageChatHistory{
			Code:    0,
			Message: err.Error(),
		}, err
	}
	resp = &lapp.BaseResponsePageChatHistory{
		Code:    0,
		Message: "success",
		Data: &lapp.PageChatHistory{
			Records: []*lapp.ChatHistory{},
		},
	}
	for _, v := range res.MessageList {
		resp.Data.Records = append(resp.Data.Records, &lapp.ChatHistory{
			Id:          v.Id,
			Message:     v.Content,
			MessageType: v.Role,
			AppId:       v.AppId,
			UserId:      v.UserId,
			CreateTime:  v.CreateTime,
			UpdateTime:  v.UpdateTime,
			IsDelete:    v.IsDelete,
			IsFile:      v.IsFile,
		})
	}
	return resp, nil
}
