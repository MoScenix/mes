package service

import (
	"context"
	"errors"
	"strings"
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
		Name:        order.Name,
		Description: order.Description,
		Status:      workorder.WorkOrderStatus(order.Status),
		ReadStatus:  workorder.WorkOrderReadStatus(order.ReadStatus),
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

func parseSinceTime(sinceTime string, recentSeconds int64) (*time.Time, error) {
	value := strings.TrimSpace(sinceTime)
	if value != "" {
		t, err := parseListTime(value)
		if err == nil {
			return t, nil
		}
		return nil, errors.New("sinceTime must use format 2006-01-02 15:04:05")
	}
	if recentSeconds > 0 {
		t := time.Now().Add(-time.Duration(recentSeconds) * time.Second)
		return &t, nil
	}
	return nil, nil
}

func parseCursorTime(cursorUpdatedAt string) (*time.Time, error) {
	value := strings.TrimSpace(cursorUpdatedAt)
	if value == "" {
		return nil, nil
	}
	t, err := parseListTime(value)
	if err != nil {
		return nil, errors.New("cursorUpdatedAt must use format 2006-01-02 15:04:05")
	}
	return t, nil
}

func parseListTime(value string) (*time.Time, error) {
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339} {
		t, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return &t, nil
		}
	}
	return nil, errors.New("invalid time")
}
