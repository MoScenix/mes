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

func CreateProcessDraft(ctx context.Context, req *inventory.CreateProcessDraftReq, callOptions ...callopt.Option) (resp *inventory.CreateProcessDraftResp, err error) {
	resp, err = defaultClient.CreateProcessDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CreateProcessDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateProcessDraft(ctx context.Context, req *inventory.UpdateProcessDraftReq, callOptions ...callopt.Option) (resp *inventory.UpdateProcessDraftResp, err error) {
	resp, err = defaultClient.UpdateProcessDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateProcessDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteProcessDraft(ctx context.Context, req *inventory.DeleteProcessDraftReq, callOptions ...callopt.Option) (resp *inventory.DeleteProcessDraftResp, err error) {
	resp, err = defaultClient.DeleteProcessDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteProcessDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SubmitProcess(ctx context.Context, req *inventory.SubmitProcessReq, callOptions ...callopt.Option) (resp *inventory.SubmitProcessResp, err error) {
	resp, err = defaultClient.SubmitProcess(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SubmitProcess call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetProcess(ctx context.Context, req *inventory.GetProcessReq, callOptions ...callopt.Option) (resp *inventory.GetProcessResp, err error) {
	resp, err = defaultClient.GetProcess(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetProcess call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListProcess(ctx context.Context, req *inventory.ListProcessReq, callOptions ...callopt.Option) (resp *inventory.ListProcessResp, err error) {
	resp, err = defaultClient.ListProcess(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListProcess call failed,err =%+v", err)
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

func CompleteInventoryFlow(ctx context.Context, req *inventory.CompleteInventoryFlowReq, callOptions ...callopt.Option) (resp *inventory.CompleteInventoryFlowResp, err error) {
	resp, err = defaultClient.CompleteInventoryFlow(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CompleteInventoryFlow call failed,err =%+v", err)
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

func CreateEngineeringOrderDraft(ctx context.Context, req *inventory.CreateEngineeringOrderDraftReq, callOptions ...callopt.Option) (resp *inventory.CreateEngineeringOrderDraftResp, err error) {
	resp, err = defaultClient.CreateEngineeringOrderDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CreateEngineeringOrderDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateEngineeringOrderDraft(ctx context.Context, req *inventory.UpdateEngineeringOrderDraftReq, callOptions ...callopt.Option) (resp *inventory.UpdateEngineeringOrderDraftResp, err error) {
	resp, err = defaultClient.UpdateEngineeringOrderDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateEngineeringOrderDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteEngineeringOrderDraft(ctx context.Context, req *inventory.DeleteEngineeringOrderDraftReq, callOptions ...callopt.Option) (resp *inventory.DeleteEngineeringOrderDraftResp, err error) {
	resp, err = defaultClient.DeleteEngineeringOrderDraft(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteEngineeringOrderDraft call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SubmitEngineeringOrder(ctx context.Context, req *inventory.SubmitEngineeringOrderReq, callOptions ...callopt.Option) (resp *inventory.SubmitEngineeringOrderResp, err error) {
	resp, err = defaultClient.SubmitEngineeringOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SubmitEngineeringOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetEngineeringOrder(ctx context.Context, req *inventory.GetEngineeringOrderReq, callOptions ...callopt.Option) (resp *inventory.GetEngineeringOrderResp, err error) {
	resp, err = defaultClient.GetEngineeringOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetEngineeringOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListEngineeringOrder(ctx context.Context, req *inventory.ListEngineeringOrderReq, callOptions ...callopt.Option) (resp *inventory.ListEngineeringOrderResp, err error) {
	resp, err = defaultClient.ListEngineeringOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListEngineeringOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
