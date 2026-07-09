package service

import (
	"context"
	"os"
	"path/filepath"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	"github.com/MoScenix/mes/app/bff/conf"
	user "github.com/MoScenix/mes/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	"github.com/MoScenix/mes/common/rpcmeta"
	rpcuser "github.com/MoScenix/mes/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserService {
	return &UpdateUserService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserService) Run(req *user.UserUpdateRequest) (resp *user.BaseResponseBoolean, err error) {
	currentUserID, ok := utils.UserIDFromContext(h.Context)
	if !ok {
		return &user.BaseResponseBoolean{
			Code:    2,
			Message: "用户未登录",
			Data:    false,
		}, nil
	}
	currentRole, _ := h.Context.Value(utils.UserRoleKey).(string)
	if req.Id != 0 {
		if currentUserID != req.Id && !rpcmeta.IsAdmin(currentRole) {
			return &user.BaseResponseBoolean{
				Code:    2,
				Message: "用户id不一致",
				Data:    false,
			}, nil
		}
	} else {
		req.Id = currentUserID
	}
	avatar, err := h.RequestContext.FormFile("avatar")
	if avatar != nil && err == nil {
		avatarPath := conf.AvatarPath(req.Id)
		req.UserAvatar = conf.AvatarURL(req.Id)
		os.MkdirAll(filepath.Dir(avatarPath), os.ModePerm)
		h.RequestContext.SaveUploadedFile(avatar, avatarPath)
	}
	_, err = rpc.UserClient.Update(h.Context, &rpcuser.UpdateReq{
		Id:          req.Id,
		UserName:    req.UserName,
		UserAvatar:  req.UserAvatar,
		UserProfile: req.UserProfile,
		UserRole:    req.UserRole,
	})
	if err != nil {
		return &user.BaseResponseBoolean{
			Code:    1,
			Message: "更新用户信息失败",
			Data:    false,
		}, err
	}
	return &user.BaseResponseBoolean{
		Code:    0,
		Message: "success",
		Data:    true,
	}, err
}
