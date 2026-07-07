package workorder

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"

	"github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder/workorderservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() workorderservice.Client
	Service() string
	CreateWorkOrder(ctx context.Context, Req *workorder.CreateWorkOrderReq, callOptions ...callopt.Option) (r *workorder.CreateWorkOrderResp, err error)
	UpdateWorkOrderDraft(ctx context.Context, Req *workorder.UpdateWorkOrderDraftReq, callOptions ...callopt.Option) (r *workorder.UpdateWorkOrderDraftResp, err error)
	DeleteWorkOrderDraft(ctx context.Context, Req *workorder.DeleteWorkOrderDraftReq, callOptions ...callopt.Option) (r *workorder.DeleteWorkOrderDraftResp, err error)
	SubmitWorkOrder(ctx context.Context, Req *workorder.SubmitWorkOrderReq, callOptions ...callopt.Option) (r *workorder.SubmitWorkOrderResp, err error)
	GetWorkOrder(ctx context.Context, Req *workorder.GetWorkOrderReq, callOptions ...callopt.Option) (r *workorder.GetWorkOrderResp, err error)
	ListWorkOrder(ctx context.Context, Req *workorder.ListWorkOrderReq, callOptions ...callopt.Option) (r *workorder.ListWorkOrderResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := workorderservice.NewClient(dstService, opts...)
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
	kitexClient workorderservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() workorderservice.Client {
	return c.kitexClient
}

func (c *clientImpl) CreateWorkOrder(ctx context.Context, Req *workorder.CreateWorkOrderReq, callOptions ...callopt.Option) (r *workorder.CreateWorkOrderResp, err error) {
	return c.kitexClient.CreateWorkOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateWorkOrderDraft(ctx context.Context, Req *workorder.UpdateWorkOrderDraftReq, callOptions ...callopt.Option) (r *workorder.UpdateWorkOrderDraftResp, err error) {
	return c.kitexClient.UpdateWorkOrderDraft(ctx, Req, callOptions...)
}

func (c *clientImpl) DeleteWorkOrderDraft(ctx context.Context, Req *workorder.DeleteWorkOrderDraftReq, callOptions ...callopt.Option) (r *workorder.DeleteWorkOrderDraftResp, err error) {
	return c.kitexClient.DeleteWorkOrderDraft(ctx, Req, callOptions...)
}

func (c *clientImpl) SubmitWorkOrder(ctx context.Context, Req *workorder.SubmitWorkOrderReq, callOptions ...callopt.Option) (r *workorder.SubmitWorkOrderResp, err error) {
	return c.kitexClient.SubmitWorkOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) GetWorkOrder(ctx context.Context, Req *workorder.GetWorkOrderReq, callOptions ...callopt.Option) (r *workorder.GetWorkOrderResp, err error) {
	return c.kitexClient.GetWorkOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) ListWorkOrder(ctx context.Context, Req *workorder.ListWorkOrderReq, callOptions ...callopt.Option) (r *workorder.ListWorkOrderResp, err error) {
	return c.kitexClient.ListWorkOrder(ctx, Req, callOptions...)
}
