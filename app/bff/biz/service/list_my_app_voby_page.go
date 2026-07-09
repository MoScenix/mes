package service

import (
	"context"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	lapp "github.com/MoScenix/mes/app/bff/hertz_gen/bff/app"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcapp "github.com/MoScenix/mes/rpc_gen/kitex_gen/app"
	rpcuser "github.com/MoScenix/mes/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListMyAppVOByPageService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListMyAppVOByPageService(Context context.Context, RequestContext *app.RequestContext) *ListMyAppVOByPageService {
	return &ListMyAppVOByPageService{RequestContext: RequestContext, Context: Context}
}

func (h *ListMyAppVOByPageService) Run(req *lapp.AppQueryRequest) (resp *lapp.BaseResponsePageAppVO, err error) {
	userID, err := scopedUserID(h.Context, req.GetUserId())
	if err != nil {
		return &lapp.BaseResponsePageAppVO{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	req.UserId = userID
	res, err := rpc.AppClient.ListApp(utils.WithIdentityMeta(h.Context), &rpcapp.ListAppReq{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		UserId:   req.UserId,
	})
	if err != nil {
		return &lapp.BaseResponsePageAppVO{
			Code:    1,
			Message: err.Error(),
		}, err
	}
	resp = &lapp.BaseResponsePageAppVO{
		Code:    0,
		Message: "success",
		Data: &lapp.PageAppVO{
			TotalPage:  res.Total,
			PageSize:   req.PageSize,
			PageNumber: req.PageNum,
			Records:    []*lapp.AppVO{},
			TotalRow:   int64(len(res.AppList)),
		},
	}
	q := rpc.UserClient
	for _, app := range res.AppList {
		r, err := q.GetUser(h.Context, &rpcuser.GetUserReq{
			Id: app.UserId,
		})
		if err != nil {
			return &lapp.BaseResponsePageAppVO{
				Code:    1,
				Message: err.Error(),
			}, err
		}
		resp.Data.Records = append(resp.Data.Records, &lapp.AppVO{
			Id:         app.Id,
			AppName:    app.AppName,
			CreateTime: app.CreateTime,
			UpdateTime: app.UpdateTime,
			UserId:     app.UserId,
			User: &lapp.UserVO{
				Id:          r.Id,
				UserName:    r.UserName,
				UserAccount: r.UserAccount,
				UserAvatar:  r.UserAvatar,
				UserProfile: r.UserProfile,
				UserRole:    r.UserRole,
				CreateTime:  r.CreateTime,
			},
		})
	}
	return resp, nil
}
