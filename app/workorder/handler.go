package main

import (
	"context"
	"github.com/MoScenix/mes/app/workorder/biz/service"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

// WorkOrderServiceImpl implements the last service interface defined in the IDL.
type WorkOrderServiceImpl struct{}

// CreateWorkOrder implements the WorkOrderServiceImpl interface.
func (s *WorkOrderServiceImpl) CreateWorkOrder(ctx context.Context, req *workorder.CreateWorkOrderReq) (resp *workorder.CreateWorkOrderResp, err error) {
	resp, err = service.NewCreateWorkOrderService(ctx).Run(req)

	return resp, err
}

// UpdateWorkOrderDraft implements the WorkOrderServiceImpl interface.
func (s *WorkOrderServiceImpl) UpdateWorkOrderDraft(ctx context.Context, req *workorder.UpdateWorkOrderDraftReq) (resp *workorder.UpdateWorkOrderDraftResp, err error) {
	resp, err = service.NewUpdateWorkOrderDraftService(ctx).Run(req)

	return resp, err
}

// DeleteWorkOrderDraft implements the WorkOrderServiceImpl interface.
func (s *WorkOrderServiceImpl) DeleteWorkOrderDraft(ctx context.Context, req *workorder.DeleteWorkOrderDraftReq) (resp *workorder.DeleteWorkOrderDraftResp, err error) {
	resp, err = service.NewDeleteWorkOrderDraftService(ctx).Run(req)

	return resp, err
}

// SubmitWorkOrder implements the WorkOrderServiceImpl interface.
func (s *WorkOrderServiceImpl) SubmitWorkOrder(ctx context.Context, req *workorder.SubmitWorkOrderReq) (resp *workorder.SubmitWorkOrderResp, err error) {
	resp, err = service.NewSubmitWorkOrderService(ctx).Run(req)

	return resp, err
}

// GetWorkOrder implements the WorkOrderServiceImpl interface.
func (s *WorkOrderServiceImpl) GetWorkOrder(ctx context.Context, req *workorder.GetWorkOrderReq) (resp *workorder.GetWorkOrderResp, err error) {
	resp, err = service.NewGetWorkOrderService(ctx).Run(req)

	return resp, err
}

// ListWorkOrder implements the WorkOrderServiceImpl interface.
func (s *WorkOrderServiceImpl) ListWorkOrder(ctx context.Context, req *workorder.ListWorkOrderReq) (resp *workorder.ListWorkOrderResp, err error) {
	resp, err = service.NewListWorkOrderService(ctx).Run(req)

	return resp, err
}

// MarkWorkOrderRead implements the WorkOrderServiceImpl interface.
func (s *WorkOrderServiceImpl) MarkWorkOrderRead(ctx context.Context, req *workorder.MarkWorkOrderReadReq) (resp *workorder.MarkWorkOrderReadResp, err error) {
	resp, err = service.NewMarkWorkOrderReadService(ctx).Run(req)

	return resp, err
}
