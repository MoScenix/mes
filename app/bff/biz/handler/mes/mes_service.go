package mes

import (
	"context"

	"github.com/MoScenix/mes/app/bff/biz/service"
	"github.com/MoScenix/mes/app/bff/biz/utils"
	mes "github.com/MoScenix/mes/app/bff/hertz_gen/bff/mes"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CreateWorkOrderDraft .
// @router /mes/work-order/draft/create [POST]
func CreateWorkOrderDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.CreateWorkOrderDraftRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseLong{}
	resp, err = service.NewCreateWorkOrderDraftService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateWorkOrderDraft .
// @router /mes/work-order/draft/update [POST]
func UpdateWorkOrderDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.UpdateWorkOrderDraftRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewUpdateWorkOrderDraftService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// DeleteWorkOrderDraft .
// @router /mes/work-order/draft/delete [POST]
func DeleteWorkOrderDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewDeleteWorkOrderDraftService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// SubmitWorkOrder .
// @router /mes/work-order/submit [POST]
func SubmitWorkOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewSubmitWorkOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetWorkOrder .
// @router /mes/work-order/get [GET]
func GetWorkOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.GetByIdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseWorkOrderVO{}
	resp, err = service.NewGetWorkOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListWorkOrder .
// @router /mes/work-order/list [POST]
func ListWorkOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.ListWorkOrderRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponsePageWorkOrderVO{}
	resp, err = service.NewListWorkOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// MarkWorkOrderRead .
// @router /mes/work-order/read [POST]
func MarkWorkOrderRead(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewMarkWorkOrderReadService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CreateEngineeringOrder .
// @router /mes/engineering-order/create [POST]
func CreateEngineeringOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.CreateEngineeringOrderRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseLong{}
	resp, err = service.NewCreateEngineeringOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateEngineeringOrder .
// @router /mes/engineering-order/update [POST]
func UpdateEngineeringOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.UpdateEngineeringOrderRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewUpdateEngineeringOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetEngineeringOrder .
// @router /mes/engineering-order/get [GET]
func GetEngineeringOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.GetByIdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseEngineeringOrderVO{}
	resp, err = service.NewGetEngineeringOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListEngineeringOrder .
// @router /mes/engineering-order/list [POST]
func ListEngineeringOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.ListEngineeringOrderRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponsePageEngineeringOrderVO{}
	resp, err = service.NewListEngineeringOrderService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AddItem .
// @router /mes/item/add [POST]
func AddItem(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.AddItemRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseLong{}
	resp, err = service.NewAddItemService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateItem .
// @router /mes/item/update [POST]
func UpdateItem(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.UpdateItemRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewUpdateItemService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetItem .
// @router /mes/item/get [GET]
func GetItem(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.GetByIdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseItemVO{}
	resp, err = service.NewGetItemService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListItem .
// @router /mes/item/list [POST]
func ListItem(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.ListItemRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponsePageItemVO{}
	resp, err = service.NewListItemService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// SearchItems .
// @router /mes/item/search [GET]
func SearchItems(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.SearchItemsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponsePageItemVO{}
	resp, err = service.NewSearchItemsService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AddItemUnit .
// @router /mes/item-unit/add [POST]
func AddItemUnit(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.AddItemUnitRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseLong{}
	resp, err = service.NewAddItemUnitService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateItemUnitStatus .
// @router /mes/item-unit/status/update [POST]
func UpdateItemUnitStatus(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.UpdateItemUnitStatusRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewUpdateItemUnitStatusService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetItemUnit .
// @router /mes/item-unit/get [GET]
func GetItemUnit(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.GetByIdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseItemUnitVO{}
	resp, err = service.NewGetItemUnitService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListItemUnit .
// @router /mes/item-unit/list [POST]
func ListItemUnit(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.ListItemUnitRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponsePageItemUnitVO{}
	resp, err = service.NewListItemUnitService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CreateInventoryFlowDraft .
// @router /mes/inventory-flow/draft/create [POST]
func CreateInventoryFlowDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.CreateInventoryFlowDraftRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseLong{}
	resp, err = service.NewCreateInventoryFlowDraftService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateInventoryFlowDraft .
// @router /mes/inventory-flow/draft/update [POST]
func UpdateInventoryFlowDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.UpdateInventoryFlowDraftRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewUpdateInventoryFlowDraftService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// DeleteInventoryFlowDraft .
// @router /mes/inventory-flow/draft/delete [POST]
func DeleteInventoryFlowDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewDeleteInventoryFlowDraftService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// SubmitInventoryFlow .
// @router /mes/inventory-flow/submit [POST]
func SubmitInventoryFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewSubmitInventoryFlowService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// AuditInventoryFlow .
// @router /mes/inventory-flow/audit [POST]
func AuditInventoryFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.AuditInventoryFlowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseBoolean{}
	resp, err = service.NewAuditInventoryFlowService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetInventoryFlow .
// @router /mes/inventory-flow/get [GET]
func GetInventoryFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.GetByIdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponseInventoryFlowVO{}
	resp, err = service.NewGetInventoryFlowService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListInventoryFlow .
// @router /mes/inventory-flow/list [POST]
func ListInventoryFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.ListInventoryFlowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &mes.BaseResponsePageInventoryFlowVO{}
	resp, err = service.NewListInventoryFlowService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CreateProcessDraft .
// @router /mes/process/draft/create [POST]
func CreateProcessDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.CreateProcessDraftRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewCreateProcessDraftService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateProcessDraft .
// @router /mes/process/draft/update [POST]
func UpdateProcessDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.UpdateProcessDraftRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewUpdateProcessDraftService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// DeleteProcessDraft .
// @router /mes/process/draft/delete [POST]
func DeleteProcessDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewDeleteProcessDraftService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// SubmitProcess .
// @router /mes/process/submit [POST]
func SubmitProcess(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewSubmitProcessService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetProcess .
// @router /mes/process/get [GET]
func GetProcess(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.GetByIdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewGetProcessService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListProcess .
// @router /mes/process/list [POST]
func ListProcess(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.ListProcessRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewListProcessService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CreateEngineeringOrderDraft .
// @router /mes/engineering-order/draft/create [POST]
func CreateEngineeringOrderDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.CreateEngineeringOrderRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewCreateEngineeringOrderDraftService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateEngineeringOrderDraft .
// @router /mes/engineering-order/draft/update [POST]
func UpdateEngineeringOrderDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.UpdateEngineeringOrderRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewUpdateEngineeringOrderDraftService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// DeleteEngineeringOrderDraft .
// @router /mes/engineering-order/draft/delete [POST]
func DeleteEngineeringOrderDraft(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewDeleteEngineeringOrderDraftService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// SubmitEngineeringOrder .
// @router /mes/engineering-order/submit [POST]
func SubmitEngineeringOrder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewSubmitEngineeringOrderService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CompleteInventoryFlow .
// @router /mes/inventory-flow/complete [POST]
func CompleteInventoryFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mes.CompleteInventoryFlowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewCompleteInventoryFlowService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
