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
