package inventory

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func AddItem(ctx context.Context, req *inventory.AddItemReq, callOptions ...callopt.Option) (resp *inventory.AddItemResp, err error) {
	resp, err = defaultClient.AddItem(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddItem call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateItem(ctx context.Context, req *inventory.UpdateItemReq, callOptions ...callopt.Option) (resp *inventory.UpdateItemResp, err error) {
	resp, err = defaultClient.UpdateItem(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateItem call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetItem(ctx context.Context, req *inventory.GetItemReq, callOptions ...callopt.Option) (resp *inventory.GetItemResp, err error) {
	resp, err = defaultClient.GetItem(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetItem call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListItem(ctx context.Context, req *inventory.ListItemReq, callOptions ...callopt.Option) (resp *inventory.ListItemResp, err error) {
	resp, err = defaultClient.ListItem(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListItem call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AddItemUnit(ctx context.Context, req *inventory.AddItemUnitReq, callOptions ...callopt.Option) (resp *inventory.AddItemUnitResp, err error) {
	resp, err = defaultClient.AddItemUnit(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddItemUnit call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateItemUnitStatus(ctx context.Context, req *inventory.UpdateItemUnitStatusReq, callOptions ...callopt.Option) (resp *inventory.UpdateItemUnitStatusResp, err error) {
	resp, err = defaultClient.UpdateItemUnitStatus(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateItemUnitStatus call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetItemUnit(ctx context.Context, req *inventory.GetItemUnitReq, callOptions ...callopt.Option) (resp *inventory.GetItemUnitResp, err error) {
	resp, err = defaultClient.GetItemUnit(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetItemUnit call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListItemUnit(ctx context.Context, req *inventory.ListItemUnitReq, callOptions ...callopt.Option) (resp *inventory.ListItemUnitResp, err error) {
	resp, err = defaultClient.ListItemUnit(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListItemUnit call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func CreateInventoryFlow(ctx context.Context, req *inventory.CreateInventoryFlowReq, callOptions ...callopt.Option) (resp *inventory.CreateInventoryFlowResp, err error) {
	resp, err = defaultClient.CreateInventoryFlow(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CreateInventoryFlow call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateInventoryFlowDraft(ctx context.Context, req *inventory.UpdateInventoryFlowDraftReq, callOptions ...callopt.Option) (resp *inventory.UpdateInventoryFlowDraftResp, err error) {
	resp, err = defaultClient.UpdateInventoryFlowDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateInventoryFlowDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteInventoryFlowDraft(ctx context.Context, req *inventory.DeleteInventoryFlowDraftReq, callOptions ...callopt.Option) (resp *inventory.DeleteInventoryFlowDraftResp, err error) {
	resp, err = defaultClient.DeleteInventoryFlowDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteInventoryFlowDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SubmitInventoryFlow(ctx context.Context, req *inventory.SubmitInventoryFlowReq, callOptions ...callopt.Option) (resp *inventory.SubmitInventoryFlowResp, err error) {
	resp, err = defaultClient.SubmitInventoryFlow(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SubmitInventoryFlow call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AuditInventoryFlow(ctx context.Context, req *inventory.AuditInventoryFlowReq, callOptions ...callopt.Option) (resp *inventory.AuditInventoryFlowResp, err error) {
	resp, err = defaultClient.AuditInventoryFlow(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AuditInventoryFlow call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetInventoryFlow(ctx context.Context, req *inventory.GetInventoryFlowReq, callOptions ...callopt.Option) (resp *inventory.GetInventoryFlowResp, err error) {
	resp, err = defaultClient.GetInventoryFlow(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetInventoryFlow call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListInventoryFlow(ctx context.Context, req *inventory.ListInventoryFlowReq, callOptions ...callopt.Option) (resp *inventory.ListInventoryFlowResp, err error) {
	resp, err = defaultClient.ListInventoryFlow(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListInventoryFlow call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
