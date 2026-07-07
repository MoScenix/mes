package service

import (
	"context"
	"errors"
	"time"

	"github.com/MoScenix/mes/app/workorder/biz/dal/mysql"
	"github.com/MoScenix/mes/app/workorder/biz/model"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

var errDatabaseNotInitialized = errors.New("workorder mysql database is not initialized")

func newWorkOrderQuery(ctx context.Context) (*model.WorkOrderQuery, error) {
	if mysql.DB == nil {
		return nil, errDatabaseNotInitialized
	}
	return model.NewWorkOrderQuery(ctx, mysql.DB), nil
}

func toWorkOrderInfo(order model.WorkOrder) *workorder.WorkOrderInfo {
	return &workorder.WorkOrderInfo{
		Id:          int64(order.ID),
		FromUserId:  order.FromUserID,
		ToUserId:    order.ToUserID,
		Description: order.Description,
		Status:      workorder.WorkOrderStatus(order.Status),
		CreateTime:  order.CreatedAt.Format(time.RFC3339),
		UpdateTime:  order.UpdatedAt.Format(time.RFC3339),
	}
}

func normalizePage(pageNum int64, pageSize int64) (int64, int64) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return pageNum, pageSize
}
