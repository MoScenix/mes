package service

import (
	"context"

	"github.com/MoScenix/mes/app/inventory/biz/model"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type ListEngineeringOrderService struct {
	ctx context.Context
}

func NewListEngineeringOrderService(ctx context.Context) *ListEngineeringOrderService {
	return &ListEngineeringOrderService{ctx: ctx}
}

func (s *ListEngineeringOrderService) Run(req *inventory.ListEngineeringOrderReq) (*inventory.ListEngineeringOrderResp, error) {
	ctx := s.ctx
	db, err := inventoryDB()
	if err != nil {
		return nil, err
	}
	_, pageSize := normalizePage(req.GetPageNum(), req.GetPageSize())
	var itemID uint
	if req.GetItemId() > 0 {
		itemID, err = uintID(req.GetItemId(), "item id")
		if err != nil {
			return nil, err
		}
	}
	var processID uint
	if req.GetProcessId() > 0 {
		processID, err = uintID(req.GetProcessId(), "process id")
		if err != nil {
			return nil, err
		}
	}
	sinceTime, err := parseSinceTime(req.GetSinceTime(), req.GetRecentSeconds())
	if err != nil {
		return nil, err
	}
	cursorUpdatedAt, err := parseCursorTime(req.GetCursorUpdatedAt())
	if err != nil {
		return nil, err
	}
	var cursorID uint
	if req.GetCursorId() > 0 {
		cursorID, err = uintID(req.GetCursorId(), "cursor id")
		if err != nil {
			return nil, err
		}
	}
	orders, hasMore, err := model.NewEngineeringOrderQuery(ctx, db).List(pageSize, req.GetLeaderUserId(), itemID, processID, int32(req.GetStatus()), req.GetNamePrefix(), req.GetItemNamePrefix(), sinceTime, cursorUpdatedAt, cursorID)
	if err != nil {
		return nil, err
	}
	resp := &inventory.ListEngineeringOrderResp{Total: int64(len(orders)), HasMore: hasMore, EngineeringOrderList: make([]*inventory.EngineeringOrderInfo, 0, len(orders))}
	for _, order := range orders {
		resp.EngineeringOrderList = append(resp.EngineeringOrderList, engineeringOrderInfo(order, false))
	}
	if len(orders) > 0 {
		last := orders[len(orders)-1]
		resp.NextCursorUpdatedAt = formatTime(last.UpdatedAt)
		resp.NextCursorId = int64(last.ID)
	}
	return resp, nil
}
