package main

import (
	"context"
	"github.com/MoScenix/mes/app/inventory/biz/service"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

// InventoryServiceImpl implements the last service interface defined in the IDL.
type InventoryServiceImpl struct{}

// AddItem implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) AddItem(ctx context.Context, req *inventory.AddItemReq) (resp *inventory.AddItemResp, err error) {
	resp, err = service.NewAddItemService(ctx).Run(req)

	return resp, err
}

// UpdateItem implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) UpdateItem(ctx context.Context, req *inventory.UpdateItemReq) (resp *inventory.UpdateItemResp, err error) {
	resp, err = service.NewUpdateItemService(ctx).Run(req)

	return resp, err
}

// GetItem implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) GetItem(ctx context.Context, req *inventory.GetItemReq) (resp *inventory.GetItemResp, err error) {
	resp, err = service.NewGetItemService(ctx).Run(req)

	return resp, err
}

// ListItem implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) ListItem(ctx context.Context, req *inventory.ListItemReq) (resp *inventory.ListItemResp, err error) {
	resp, err = service.NewListItemService(ctx).Run(req)

	return resp, err
}

// AddItemUnit implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) AddItemUnit(ctx context.Context, req *inventory.AddItemUnitReq) (resp *inventory.AddItemUnitResp, err error) {
	resp, err = service.NewAddItemUnitService(ctx).Run(req)

	return resp, err
}

// UpdateItemUnitStatus implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) UpdateItemUnitStatus(ctx context.Context, req *inventory.UpdateItemUnitStatusReq) (resp *inventory.UpdateItemUnitStatusResp, err error) {
	resp, err = service.NewUpdateItemUnitStatusService(ctx).Run(req)

	return resp, err
}

// GetItemUnit implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) GetItemUnit(ctx context.Context, req *inventory.GetItemUnitReq) (resp *inventory.GetItemUnitResp, err error) {
	resp, err = service.NewGetItemUnitService(ctx).Run(req)

	return resp, err
}

// ListItemUnit implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) ListItemUnit(ctx context.Context, req *inventory.ListItemUnitReq) (resp *inventory.ListItemUnitResp, err error) {
	resp, err = service.NewListItemUnitService(ctx).Run(req)

	return resp, err
}

// CreateInventoryFlow implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) CreateInventoryFlow(ctx context.Context, req *inventory.CreateInventoryFlowReq) (resp *inventory.CreateInventoryFlowResp, err error) {
	resp, err = service.NewCreateInventoryFlowService(ctx).Run(req)

	return resp, err
}

// UpdateInventoryFlowDraft implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) UpdateInventoryFlowDraft(ctx context.Context, req *inventory.UpdateInventoryFlowDraftReq) (resp *inventory.UpdateInventoryFlowDraftResp, err error) {
	resp, err = service.NewUpdateInventoryFlowDraftService(ctx).Run(req)

	return resp, err
}

// DeleteInventoryFlowDraft implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) DeleteInventoryFlowDraft(ctx context.Context, req *inventory.DeleteInventoryFlowDraftReq) (resp *inventory.DeleteInventoryFlowDraftResp, err error) {
	resp, err = service.NewDeleteInventoryFlowDraftService(ctx).Run(req)

	return resp, err
}

// SubmitInventoryFlow implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) SubmitInventoryFlow(ctx context.Context, req *inventory.SubmitInventoryFlowReq) (resp *inventory.SubmitInventoryFlowResp, err error) {
	resp, err = service.NewSubmitInventoryFlowService(ctx).Run(req)

	return resp, err
}

// CompleteInventoryFlow implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) CompleteInventoryFlow(ctx context.Context, req *inventory.CompleteInventoryFlowReq) (resp *inventory.CompleteInventoryFlowResp, err error) {
	resp, err = service.NewCompleteInventoryFlowService(ctx).Run(req)

	return resp, err
}

// AuditInventoryFlow implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) AuditInventoryFlow(ctx context.Context, req *inventory.AuditInventoryFlowReq) (resp *inventory.AuditInventoryFlowResp, err error) {
	resp, err = service.NewAuditInventoryFlowService(ctx).Run(req)

	return resp, err
}

// GetInventoryFlow implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) GetInventoryFlow(ctx context.Context, req *inventory.GetInventoryFlowReq) (resp *inventory.GetInventoryFlowResp, err error) {
	resp, err = service.NewGetInventoryFlowService(ctx).Run(req)

	return resp, err
}

// ListInventoryFlow implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) ListInventoryFlow(ctx context.Context, req *inventory.ListInventoryFlowReq) (resp *inventory.ListInventoryFlowResp, err error) {
	resp, err = service.NewListInventoryFlowService(ctx).Run(req)

	return resp, err
}

// GetEngineeringOrder implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) GetEngineeringOrder(ctx context.Context, req *inventory.GetEngineeringOrderReq) (resp *inventory.GetEngineeringOrderResp, err error) {
	resp, err = service.NewGetEngineeringOrderService(ctx).Run(req)

	return resp, err
}

// ListEngineeringOrder implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) ListEngineeringOrder(ctx context.Context, req *inventory.ListEngineeringOrderReq) (resp *inventory.ListEngineeringOrderResp, err error) {
	resp, err = service.NewListEngineeringOrderService(ctx).Run(req)

	return resp, err
}

// CreateProcessDraft implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) CreateProcessDraft(ctx context.Context, req *inventory.CreateProcessDraftReq) (resp *inventory.CreateProcessDraftResp, err error) {
	resp, err = service.NewCreateProcessDraftService(ctx).Run(req)

	return resp, err
}

// UpdateProcessDraft implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) UpdateProcessDraft(ctx context.Context, req *inventory.UpdateProcessDraftReq) (resp *inventory.UpdateProcessDraftResp, err error) {
	resp, err = service.NewUpdateProcessDraftService(ctx).Run(req)

	return resp, err
}

// DeleteProcessDraft implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) DeleteProcessDraft(ctx context.Context, req *inventory.DeleteProcessDraftReq) (resp *inventory.DeleteProcessDraftResp, err error) {
	resp, err = service.NewDeleteProcessDraftService(ctx).Run(req)

	return resp, err
}

// SubmitProcess implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) SubmitProcess(ctx context.Context, req *inventory.SubmitProcessReq) (resp *inventory.SubmitProcessResp, err error) {
	resp, err = service.NewSubmitProcessService(ctx).Run(req)

	return resp, err
}

// GetProcess implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) GetProcess(ctx context.Context, req *inventory.GetProcessReq) (resp *inventory.GetProcessResp, err error) {
	resp, err = service.NewGetProcessService(ctx).Run(req)

	return resp, err
}

// ListProcess implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) ListProcess(ctx context.Context, req *inventory.ListProcessReq) (resp *inventory.ListProcessResp, err error) {
	resp, err = service.NewListProcessService(ctx).Run(req)

	return resp, err
}

// CreateEngineeringOrderDraft implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) CreateEngineeringOrderDraft(ctx context.Context, req *inventory.CreateEngineeringOrderDraftReq) (resp *inventory.CreateEngineeringOrderDraftResp, err error) {
	resp, err = service.NewCreateEngineeringOrderDraftService(ctx).Run(req)

	return resp, err
}

// UpdateEngineeringOrderDraft implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) UpdateEngineeringOrderDraft(ctx context.Context, req *inventory.UpdateEngineeringOrderDraftReq) (resp *inventory.UpdateEngineeringOrderDraftResp, err error) {
	resp, err = service.NewUpdateEngineeringOrderDraftService(ctx).Run(req)

	return resp, err
}

// DeleteEngineeringOrderDraft implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) DeleteEngineeringOrderDraft(ctx context.Context, req *inventory.DeleteEngineeringOrderDraftReq) (resp *inventory.DeleteEngineeringOrderDraftResp, err error) {
	resp, err = service.NewDeleteEngineeringOrderDraftService(ctx).Run(req)

	return resp, err
}

// SubmitEngineeringOrder implements the InventoryServiceImpl interface.
func (s *InventoryServiceImpl) SubmitEngineeringOrder(ctx context.Context, req *inventory.SubmitEngineeringOrderReq) (resp *inventory.SubmitEngineeringOrderResp, err error) {
	resp, err = service.NewSubmitEngineeringOrderService(ctx).Run(req)

	return resp, err
}
