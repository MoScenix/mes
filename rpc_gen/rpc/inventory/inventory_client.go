package inventory

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"

	"github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory/inventoryservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() inventoryservice.Client
	Service() string
	AddItem(ctx context.Context, Req *inventory.AddItemReq, callOptions ...callopt.Option) (r *inventory.AddItemResp, err error)
	UpdateItem(ctx context.Context, Req *inventory.UpdateItemReq, callOptions ...callopt.Option) (r *inventory.UpdateItemResp, err error)
	GetItem(ctx context.Context, Req *inventory.GetItemReq, callOptions ...callopt.Option) (r *inventory.GetItemResp, err error)
	ListItem(ctx context.Context, Req *inventory.ListItemReq, callOptions ...callopt.Option) (r *inventory.ListItemResp, err error)
	AddItemUnit(ctx context.Context, Req *inventory.AddItemUnitReq, callOptions ...callopt.Option) (r *inventory.AddItemUnitResp, err error)
	UpdateItemUnitStatus(ctx context.Context, Req *inventory.UpdateItemUnitStatusReq, callOptions ...callopt.Option) (r *inventory.UpdateItemUnitStatusResp, err error)
	GetItemUnit(ctx context.Context, Req *inventory.GetItemUnitReq, callOptions ...callopt.Option) (r *inventory.GetItemUnitResp, err error)
	ListItemUnit(ctx context.Context, Req *inventory.ListItemUnitReq, callOptions ...callopt.Option) (r *inventory.ListItemUnitResp, err error)
	CreateInventoryFlow(ctx context.Context, Req *inventory.CreateInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.CreateInventoryFlowResp, err error)
	UpdateInventoryFlowDraft(ctx context.Context, Req *inventory.UpdateInventoryFlowDraftReq, callOptions ...callopt.Option) (r *inventory.UpdateInventoryFlowDraftResp, err error)
	DeleteInventoryFlowDraft(ctx context.Context, Req *inventory.DeleteInventoryFlowDraftReq, callOptions ...callopt.Option) (r *inventory.DeleteInventoryFlowDraftResp, err error)
	SubmitInventoryFlow(ctx context.Context, Req *inventory.SubmitInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.SubmitInventoryFlowResp, err error)
	AuditInventoryFlow(ctx context.Context, Req *inventory.AuditInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.AuditInventoryFlowResp, err error)
	GetInventoryFlow(ctx context.Context, Req *inventory.GetInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.GetInventoryFlowResp, err error)
	ListInventoryFlow(ctx context.Context, Req *inventory.ListInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.ListInventoryFlowResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := inventoryservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient inventoryservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() inventoryservice.Client {
	return c.kitexClient
}

func (c *clientImpl) AddItem(ctx context.Context, Req *inventory.AddItemReq, callOptions ...callopt.Option) (r *inventory.AddItemResp, err error) {
	return c.kitexClient.AddItem(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateItem(ctx context.Context, Req *inventory.UpdateItemReq, callOptions ...callopt.Option) (r *inventory.UpdateItemResp, err error) {
	return c.kitexClient.UpdateItem(ctx, Req, callOptions...)
}

func (c *clientImpl) GetItem(ctx context.Context, Req *inventory.GetItemReq, callOptions ...callopt.Option) (r *inventory.GetItemResp, err error) {
	return c.kitexClient.GetItem(ctx, Req, callOptions...)
}

func (c *clientImpl) ListItem(ctx context.Context, Req *inventory.ListItemReq, callOptions ...callopt.Option) (r *inventory.ListItemResp, err error) {
	return c.kitexClient.ListItem(ctx, Req, callOptions...)
}

func (c *clientImpl) AddItemUnit(ctx context.Context, Req *inventory.AddItemUnitReq, callOptions ...callopt.Option) (r *inventory.AddItemUnitResp, err error) {
	return c.kitexClient.AddItemUnit(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateItemUnitStatus(ctx context.Context, Req *inventory.UpdateItemUnitStatusReq, callOptions ...callopt.Option) (r *inventory.UpdateItemUnitStatusResp, err error) {
	return c.kitexClient.UpdateItemUnitStatus(ctx, Req, callOptions...)
}

func (c *clientImpl) GetItemUnit(ctx context.Context, Req *inventory.GetItemUnitReq, callOptions ...callopt.Option) (r *inventory.GetItemUnitResp, err error) {
	return c.kitexClient.GetItemUnit(ctx, Req, callOptions...)
}

func (c *clientImpl) ListItemUnit(ctx context.Context, Req *inventory.ListItemUnitReq, callOptions ...callopt.Option) (r *inventory.ListItemUnitResp, err error) {
	return c.kitexClient.ListItemUnit(ctx, Req, callOptions...)
}

func (c *clientImpl) CreateInventoryFlow(ctx context.Context, Req *inventory.CreateInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.CreateInventoryFlowResp, err error) {
	return c.kitexClient.CreateInventoryFlow(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateInventoryFlowDraft(ctx context.Context, Req *inventory.UpdateInventoryFlowDraftReq, callOptions ...callopt.Option) (r *inventory.UpdateInventoryFlowDraftResp, err error) {
	return c.kitexClient.UpdateInventoryFlowDraft(ctx, Req, callOptions...)
}

func (c *clientImpl) DeleteInventoryFlowDraft(ctx context.Context, Req *inventory.DeleteInventoryFlowDraftReq, callOptions ...callopt.Option) (r *inventory.DeleteInventoryFlowDraftResp, err error) {
	return c.kitexClient.DeleteInventoryFlowDraft(ctx, Req, callOptions...)
}

func (c *clientImpl) SubmitInventoryFlow(ctx context.Context, Req *inventory.SubmitInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.SubmitInventoryFlowResp, err error) {
	return c.kitexClient.SubmitInventoryFlow(ctx, Req, callOptions...)
}

func (c *clientImpl) AuditInventoryFlow(ctx context.Context, Req *inventory.AuditInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.AuditInventoryFlowResp, err error) {
	return c.kitexClient.AuditInventoryFlow(ctx, Req, callOptions...)
}

func (c *clientImpl) GetInventoryFlow(ctx context.Context, Req *inventory.GetInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.GetInventoryFlowResp, err error) {
	return c.kitexClient.GetInventoryFlow(ctx, Req, callOptions...)
}

func (c *clientImpl) ListInventoryFlow(ctx context.Context, Req *inventory.ListInventoryFlowReq, callOptions ...callopt.Option) (r *inventory.ListInventoryFlowResp, err error) {
	return c.kitexClient.ListInventoryFlow(ctx, Req, callOptions...)
}
