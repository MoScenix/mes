package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	user "github.com/MoScenix/mes/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/mes/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
)

type UserLoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUserLoginService(Context context.Context, RequestContext *app.RequestContext) *UserLoginService {
	return &UserLoginService{RequestContext: RequestContext, Context: Context}
}

func (h *UserLoginService) Run(req *user.UserLoginRequest) (resp *user.BaseResponseLoginUserVO, err error) {
	res1, err := rpc.UserClient.Login(h.Context, &rpcuser.LoginReq{
		UserAccount:  req.UserAccount,
		UserPassword: req.UserPassword,
	})
	if err != nil {
		message := err.Error()
		if strings.Contains(message, "record not found") || strings.Contains(message, "bcrypt") {
			message = "账号或密码错误"
		}
		return &user.BaseResponseLoginUserVO{
			Code:    1,
			Message: message,
		}, nil
	}
	session := sessions.Default(h.RequestContext)
	session.Set(utils.UserIdKey, res1.UserId)
	session.Set(utils.UserRoleKey, res1.UserRole)
	if err := session.Save(); err != nil {
		return &user.BaseResponseLoginUserVO{
			Code:    1,
			Message: fmt.Sprintf("save session failed: %v", err),
		}, nil
	}
	res, err := rpc.UserClient.GetUser(h.Context, &rpcuser.GetUserReq{
		Id: int64(res1.UserId),
	})
	if err != nil {
		return &user.BaseResponseLoginUserVO{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	if res == nil {
		return &user.BaseResponseLoginUserVO{
			Code:    1,
			Message: "user service returned empty user",
		}, nil
	}
	return &user.BaseResponseLoginUserVO{
		Code:    0,
		Message: "success",
		Data: &user.LoginUserVO{
			Id:          res.Id,
			UserName:    res.UserName,
			UserAccount: res.UserAccount,
			UserAvatar:  res.UserAvatar,
			UserProfile: res.UserProfile,
			UserRole:    res.UserRole,
			CreateTime:  res.CreateTime,
			UpdateTime:  res.UpdateTime,
		},
	}, nil
}
