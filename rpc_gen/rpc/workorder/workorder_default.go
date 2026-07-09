package workorder

import (
	"context"
	workorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func CreateWorkOrder(ctx context.Context, req *workorder.CreateWorkOrderReq, callOptions ...callopt.Option) (resp *workorder.CreateWorkOrderResp, err error) {
	resp, err = defaultClient.CreateWorkOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CreateWorkOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateWorkOrderDraft(ctx context.Context, req *workorder.UpdateWorkOrderDraftReq, callOptions ...callopt.Option) (resp *workorder.UpdateWorkOrderDraftResp, err error) {
	resp, err = defaultClient.UpdateWorkOrderDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateWorkOrderDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteWorkOrderDraft(ctx context.Context, req *workorder.DeleteWorkOrderDraftReq, callOptions ...callopt.Option) (resp *workorder.DeleteWorkOrderDraftResp, err error) {
	resp, err = defaultClient.DeleteWorkOrderDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteWorkOrderDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SubmitWorkOrder(ctx context.Context, req *workorder.SubmitWorkOrderReq, callOptions ...callopt.Option) (resp *workorder.SubmitWorkOrderResp, err error) {
	resp, err = defaultClient.SubmitWorkOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SubmitWorkOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetWorkOrder(ctx context.Context, req *workorder.GetWorkOrderReq, callOptions ...callopt.Option) (resp *workorder.GetWorkOrderResp, err error) {
	resp, err = defaultClient.GetWorkOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetWorkOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListWorkOrder(ctx context.Context, req *workorder.ListWorkOrderReq, callOptions ...callopt.Option) (resp *workorder.ListWorkOrderResp, err error) {
	resp, err = defaultClient.ListWorkOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListWorkOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func MarkWorkOrderRead(ctx context.Context, req *workorder.MarkWorkOrderReadReq, callOptions ...callopt.Option) (resp *workorder.MarkWorkOrderReadResp, err error) {
	resp, err = defaultClient.MarkWorkOrderRead(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "MarkWorkOrderRead call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
