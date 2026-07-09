package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	"github.com/MoScenix/mes/common/rpcmeta"
	rpcapp "github.com/MoScenix/mes/rpc_gen/kitex_gen/app"
	rpcinventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	rpcworkorder "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
)

var errForbiddenAccess = errors.New("forbidden: no permission")

func requireBFFUserID(ctx context.Context) (int64, error) {
	userID, ok := utils.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return 0, utils.ErrUnauthorizedUserID
	}
	return userID, nil
}

func bffUserRole(ctx context.Context) string {
	role, _ := ctx.Value(utils.UserRoleKey).(string)
	return role
}

func bffIsAdmin(ctx context.Context) bool {
	return rpcmeta.IsAdmin(bffUserRole(ctx))
}

func scopedUserID(ctx context.Context, requested int64) (int64, error) {
	currentUserID, err := requireBFFUserID(ctx)
	if err != nil {
		return 0, err
	}
	if bffIsAdmin(ctx) {
		return requested, nil
	}
	if requested != 0 && requested != currentUserID {
		return 0, errForbiddenAccess
	}
	return currentUserID, nil
}

func requireSameUserOrAdmin(ctx context.Context, requested int64) (int64, error) {
	currentUserID, err := requireBFFUserID(ctx)
	if err != nil {
		return 0, err
	}
	if bffIsAdmin(ctx) {
		if requested != 0 {
			return requested, nil
		}
		return currentUserID, nil
	}
	if requested != 0 && requested != currentUserID {
		return 0, errForbiddenAccess
	}
	return currentUserID, nil
}

func requireAppOwnerOrAdmin(ctx context.Context, appID int64) error {
	if appID <= 0 {
		return fmt.Errorf("appId is required")
	}
	currentUserID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	res, err := rpc.AppClient.GetApp(utils.WithIdentityMeta(ctx), &rpcapp.GetAppReq{Id: appID})
	if err != nil {
		return err
	}
	appInfo := res.GetApp()
	if appInfo == nil {
		return fmt.Errorf("app not found")
	}
	if bffIsAdmin(ctx) || appInfo.GetUserId() == currentUserID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanViewWorkOrder(ctx context.Context, order *rpcworkorder.WorkOrderInfo) error {
	if order == nil {
		return fmt.Errorf("work order not found")
	}
	userID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	if order.GetStatus() == rpcworkorder.WorkOrderStatus_WORK_ORDER_STATUS_DRAFT {
		if bffIsAdmin(ctx) || order.GetFromUserId() == userID {
			return nil
		}
		return errForbiddenAccess
	}
	if bffIsAdmin(ctx) || order.GetFromUserId() == userID || order.GetToUserId() == userID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanUpdateWorkOrderDraft(ctx context.Context, order *rpcworkorder.WorkOrderInfo) error {
	if order == nil {
		return fmt.Errorf("work order not found")
	}
	userID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	if bffIsAdmin(ctx) || order.GetFromUserId() == userID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanMarkWorkOrderRead(ctx context.Context, order *rpcworkorder.WorkOrderInfo) error {
	if order == nil {
		return fmt.Errorf("work order not found")
	}
	userID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	if bffIsAdmin(ctx) || order.GetToUserId() == userID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanViewEngineeringOrder(ctx context.Context, order *rpcinventory.EngineeringOrderInfo) error {
	if order == nil {
		return fmt.Errorf("engineering order not found")
	}
	if order.GetStatus() != rpcinventory.DraftStatus_DRAFT_STATUS_DRAFT {
		return nil
	}
	userID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	if bffIsAdmin(ctx) || order.GetLeaderUserId() == userID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanUpdateEngineeringOrder(ctx context.Context, order *rpcinventory.EngineeringOrderInfo) error {
	return requireCanViewEngineeringOrder(ctx, order)
}

func requireCanViewProcess(ctx context.Context, process *rpcinventory.ProcessInfo) error {
	if process == nil {
		return fmt.Errorf("process not found")
	}
	if process.GetStatus() == rpcinventory.DraftStatus_DRAFT_STATUS_SUBMITTED || process.GetStatus() == rpcinventory.DraftStatus_DRAFT_STATUS_DONE {
		return nil
	}
	userID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	if bffIsAdmin(ctx) || process.GetOwnerUserId() == userID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanUpdateProcessDraft(ctx context.Context, process *rpcinventory.ProcessInfo) error {
	if process == nil {
		return fmt.Errorf("process not found")
	}
	userID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	if bffIsAdmin(ctx) || process.GetOwnerUserId() == userID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanViewInventoryFlow(ctx context.Context, flow *rpcinventory.InventoryFlowInfo) error {
	if flow == nil {
		return fmt.Errorf("inventory flow not found")
	}
	userID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	if flow.GetFlowStatus() == rpcinventory.FlowStatus_FLOW_STATUS_DRAFT {
		if bffIsAdmin(ctx) || flow.GetFromUserId() == userID {
			return nil
		}
		return errForbiddenAccess
	}
	if bffIsAdmin(ctx) || flow.GetFromUserId() == userID || flow.GetToUserId() == userID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanUpdateInventoryFlowDraft(ctx context.Context, flow *rpcinventory.InventoryFlowInfo) error {
	if flow == nil {
		return fmt.Errorf("inventory flow not found")
	}
	userID, err := requireBFFUserID(ctx)
	if err != nil {
		return err
	}
	if bffIsAdmin(ctx) || flow.GetFromUserId() == userID {
		return nil
	}
	return errForbiddenAccess
}

func requireCanAuditInventoryFlow(ctx context.Context) error {
	if _, err := requireBFFUserID(ctx); err != nil {
		return err
	}
	if rpcmeta.CanAuditInventoryFlow(bffUserRole(ctx)) {
		return nil
	}
	return errForbiddenAccess
}
