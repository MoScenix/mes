package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/MoScenix/mes/app/ai/conf"
	"github.com/MoScenix/mes/app/ai/infra"
	"github.com/MoScenix/mes/common/rpcmeta"
	inventorypb "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
	userpb "github.com/MoScenix/mes/rpc_gen/kitex_gen/user"
	workorderpb "github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ListWorkOrdersInput struct {
	Limit      int64  `json:"limit,omitempty" jsonschema:"description=Maximum number of work orders to return, default 30."`
	NamePrefix string `json:"name_prefix,omitempty" jsonschema:"description=Work order name prefix keyword."`
	UserID     int64  `json:"user_id,omitempty" jsonschema:"description=Admin-only user id to query. Non-admin values are ignored and current operator is used."`
	IsTo       bool   `json:"is_to,omitempty" jsonschema:"description=true lists work orders assigned to the user; false lists work orders submitted by the user."`
	Unread     bool   `json:"unread,omitempty" jsonschema:"description=true returns only unread work orders."`
	Status     string `json:"status,omitempty" jsonschema:"description=Optional status: draft or submitted. Drafts are only visible to their creator."`
}

type MarkWorkOrderReadInput struct {
	ID int64 `json:"id" jsonschema:"description=Work order id."`
}

type CreateWorkOrderDraftInput struct {
	Name        string `json:"name" jsonschema:"description=Required work order name."`
	ToUserID    int64  `json:"to_user_id" jsonschema:"description=Recipient user id."`
	Description string `json:"description" jsonschema:"description=Work order description."`
	FromUserID  int64  `json:"from_user_id,omitempty" jsonschema:"description=Admin-only submitter override. Non-admin values are ignored and current operator is used."`
}

type SearchUsersInput struct {
	ID       int64  `json:"id,omitempty" jsonschema:"description=Exact user id. If set, returns this user only."`
	Name     string `json:"name,omitempty" jsonschema:"description=User name prefix or keyword. Use this before ask_user when the user gives a name such as root instead of a numeric id."`
	Account  string `json:"account,omitempty" jsonschema:"description=User account prefix or keyword. Use this before ask_user when the user gives an account such as root instead of a numeric id."`
	Role     string `json:"role,omitempty" jsonschema:"description=Optional role filter such as leader, purchase, worker, process_engineer, warehouse_admin, sales, admin."`
	PageSize int64  `json:"page_size,omitempty" jsonschema:"description=Maximum users to return, default 10."`
}

type safeUserInfo struct {
	ID      int64  `json:"id"`
	Account string `json:"account,omitempty"`
	Name    string `json:"name,omitempty"`
	Role    string `json:"role,omitempty"`
}

type searchUsersOutput struct {
	Users []safeUserInfo `json:"users"`
	Total int64          `json:"total,omitempty"`
	Note  string         `json:"note,omitempty"`
}

type UpdateWorkOrderDraftInput struct {
	ID          int64  `json:"id" jsonschema:"description=Work order draft id."`
	Name        string `json:"name" jsonschema:"description=Required updated work order name."`
	ToUserID    int64  `json:"to_user_id" jsonschema:"description=Recipient user id."`
	Description string `json:"description" jsonschema:"description=Updated work order description."`
	FromUserID  int64  `json:"from_user_id,omitempty" jsonschema:"description=Admin-only submitter override. Non-admin values are ignored and current operator is used."`
}

type CreateEngineeringOrderDraftInput struct {
	Name              string `json:"name" jsonschema:"description=Required engineering order name."`
	LeaderUserID      int64  `json:"leader_user_id,omitempty" jsonschema:"description=Admin-only leader override. Non-admin values are ignored and current operator is used."`
	ProcessID         int64  `json:"process_id,omitempty" jsonschema:"description=Process id for the planned production draft."`
	ItemID            int64  `json:"item_id" jsonschema:"description=Produced item definition id."`
	ExpectedQuantity  int64  `json:"expected_quantity" jsonschema:"description=Planned production quantity."`
	QualifiedQuantity int64  `json:"qualified_quantity,omitempty" jsonschema:"description=Planned qualified quantity."`
	Description       string `json:"description,omitempty" jsonschema:"description=Engineering order draft description."`
}

type UpdateEngineeringOrderDraftInput struct {
	ID                int64  `json:"id" jsonschema:"description=Engineering order draft id."`
	Name              string `json:"name" jsonschema:"description=Required updated engineering order name."`
	LeaderUserID      int64  `json:"leader_user_id,omitempty" jsonschema:"description=Admin-only leader override. Non-admin values are ignored and current operator is used."`
	ProcessID         int64  `json:"process_id,omitempty" jsonschema:"description=Process id for the planned production draft."`
	ItemID            int64  `json:"item_id" jsonschema:"description=Produced item definition id."`
	ExpectedQuantity  int64  `json:"expected_quantity" jsonschema:"description=Planned production quantity."`
	QualifiedQuantity int64  `json:"qualified_quantity,omitempty" jsonschema:"description=Planned qualified quantity."`
	Description       string `json:"description,omitempty" jsonschema:"description=Engineering order draft description."`
}

type ListEngineeringOrdersInput struct {
	PageNum      int64  `json:"page_num,omitempty" jsonschema:"description=Page number, default 1."`
	PageSize     int64  `json:"page_size,omitempty" jsonschema:"description=Page size, default 30."`
	NamePrefix   string `json:"name_prefix,omitempty" jsonschema:"description=Engineering order name prefix keyword."`
	LeaderUserID int64  `json:"leader_user_id,omitempty" jsonschema:"description=Admin-only leader filter. Non-admin values are ignored and current operator is used."`
	ItemID       int64  `json:"item_id,omitempty" jsonschema:"description=Produced item definition id filter."`
}

type GetEngineeringOrderInput struct {
	ID int64 `json:"id" jsonschema:"description=Engineering order id."`
}

type InventoryFlowItemInput struct {
	ItemID        int64 `json:"item_id" jsonschema:"description=Item definition id."`
	ApplyQuantity int64 `json:"apply_quantity" jsonschema:"description=Requested item quantity."`
}

type CreateInventoryFlowDraftInput struct {
	Name        string                   `json:"name" jsonschema:"description=Required inventory flow name."`
	ToUserID    int64                    `json:"to_user_id" jsonschema:"description=Recipient user id."`
	FlowType    string                   `json:"flow_type" jsonschema:"description=Inventory flow type: in or out."`
	Description string                   `json:"description,omitempty" jsonschema:"description=Inventory flow description."`
	Items       []InventoryFlowItemInput `json:"items,omitempty" jsonschema:"description=Item quantities requested by item definition."`
	ItemUnitIDs []int64                  `json:"item_unit_ids,omitempty" jsonschema:"description=Concrete item unit ids, usually for outbound flows."`
	FromUserID  int64                    `json:"from_user_id,omitempty" jsonschema:"description=Admin-only submitter override. Non-admin values are ignored and current operator is used."`
}

type ListInventoryFlowsInput struct {
	Limit      int64  `json:"limit,omitempty" jsonschema:"description=Maximum number of inventory flows to return, default 30."`
	NamePrefix string `json:"name_prefix,omitempty" jsonschema:"description=Inventory flow name prefix keyword."`
	Scope      string `json:"scope,omitempty" jsonschema:"description=Scope: to_me, from_me, or submitted_by_me. Default submitted_by_me."`
	FlowStatus string `json:"flow_status,omitempty" jsonschema:"description=Optional status: draft, submitted, approved, rejected."`
	UserID     int64  `json:"user_id,omitempty" jsonschema:"description=Admin-only user id to query. Non-admin values are ignored and current operator is used."`
}

type GetInventoryFlowInput struct {
	ID int64 `json:"id" jsonschema:"description=Inventory flow id."`
}

type SearchItemsInput struct {
	NamePrefix string `json:"name_prefix,omitempty" jsonschema:"description=Item name prefix keyword."`
	PageNum    int64  `json:"page_num,omitempty" jsonschema:"description=Page number, default 1."`
	PageSize   int64  `json:"page_size,omitempty" jsonschema:"description=Page size, default 30."`
}

type GetItemInput struct {
	ID int64 `json:"id" jsonschema:"description=Item definition id."`
}

type ListItemUnitsInput struct {
	ItemID        int64  `json:"item_id,omitempty" jsonschema:"description=Item definition id filter."`
	StockStatus   string `json:"stock_status,omitempty" jsonschema:"description=Optional stock status: in_stock, reserved, out_stock."`
	QualityStatus string `json:"quality_status,omitempty" jsonschema:"description=Optional quality status: pending, qualified, unqualified."`
	PageNum       int64  `json:"page_num,omitempty" jsonschema:"description=Page number, default 1."`
	PageSize      int64  `json:"page_size,omitempty" jsonschema:"description=Page size, default 30."`
}

type ListPendingInventoryFlowsInput struct {
	Limit      int64  `json:"limit,omitempty" jsonschema:"description=Maximum number of pending inventory flows to return, default 30."`
	NamePrefix string `json:"name_prefix,omitempty" jsonschema:"description=Inventory flow name prefix keyword."`
	UserID     int64  `json:"user_id,omitempty" jsonschema:"description=Admin-only warehouse user id. Non-admin values are ignored and current operator is used."`
}

type InventoryCheckInput struct {
	NamePrefix    string `json:"name_prefix,omitempty" jsonschema:"description=Item name prefix keyword."`
	ItemID        int64  `json:"item_id,omitempty" jsonschema:"description=Optional item definition id for item units."`
	StockStatus   string `json:"stock_status,omitempty" jsonschema:"description=Optional item unit stock status: in_stock, reserved, out_stock."`
	QualityStatus string `json:"quality_status,omitempty" jsonschema:"description=Optional item unit quality status: pending, qualified, unqualified."`
	PageNum       int64  `json:"page_num,omitempty" jsonschema:"description=Page number, default 1."`
	PageSize      int64  `json:"page_size,omitempty" jsonschema:"description=Page size, default 30."`
}

var nonAIToolNames = map[string]bool{
	"submit_work_order":              true,
	"submit_process":                 true,
	"submit_engineering_order":       true,
	"submit_inventory_flow":          true,
	"audit_inventory_flow":           true,
	"update_item_unit_status":        true,
	"add_item":                       true,
	"update_item":                    true,
	"delete_work_order_draft":        true,
	"delete_process_draft":           true,
	"delete_engineering_order":       true,
	"delete_engineering_order_draft": true,
	"delete_inventory_flow":          true,
	"delete_inventory_flow_draft":    true,
}

func NewMESTools(ctx context.Context, baseTools ...tool.BaseTool) ([]tool.BaseTool, error) {
	cfg := conf.GetConf().AITools
	role := operatorRole(ctx, cfg.RoleAliases)
	toolNames := toolNamesForRole(role, cfg)
	builders := mesToolBuilders(baseTools...)

	tools := make([]tool.BaseTool, 0, len(toolNames))
	for _, name := range toolNames {
		builder, ok := builders[name]
		if !ok {
			continue
		}
		t, err := builder()
		if err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}
	return tools, nil
}

func operatorRole(ctx context.Context, aliases map[string]string) string {
	return rpcmeta.NormalizeRole(rpcmeta.FromContext(ctx).OperatorRole, aliases)
}

func toolNamesForRole(role string, cfg conf.AITools) []string {
	groups := cfg.RoleGroups[role]
	if rpcmeta.IsAdmin(role) && len(groups) == 0 {
		groups = []string{"common", "workorder", "engineering_order", "inventory_flow", "item", "warehouse_admin"}
	}
	seen := map[string]bool{}
	var names []string
	for _, group := range groups {
		for _, name := range cfg.ToolGroups[group] {
			if name == "" || seen[name] || nonAIToolNames[name] {
				continue
			}
			seen[name] = true
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

type toolBuilder func() (tool.BaseTool, error)

func mesToolBuilders(baseTools ...tool.BaseTool) map[string]toolBuilder {
	builders := map[string]toolBuilder{}
	for _, t := range baseTools {
		if t == nil {
			continue
		}
		info, err := t.Info(context.Background())
		if err != nil || info == nil || info.Name == "" {
			continue
		}
		baseTool := t
		builders[info.Name] = func() (tool.BaseTool, error) { return baseTool, nil }
	}
	builders["list_work_orders"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[ListWorkOrdersInput, string]("list_work_orders", "Get the latest time-ordered work orders as dropdown/list data. Use limit for the number of records, default 30.", runListWorkOrders)
	}
	builders["mark_work_order_read"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[MarkWorkOrderReadInput, string]("mark_work_order_read", "Mark a work order as read.", runMarkWorkOrderRead)
	}
	builders["search_users"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[SearchUsersInput, string]("search_users", "Search assignable MES users by id, name, account, or role. If the user provides a username/account like root but not a numeric id, call this tool before asking the user for an id. Returns only safe fields: id, account, name, role.", runSearchUsers)
	}
	builders["create_work_order_draft"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[CreateWorkOrderDraftInput, string]("create_work_order_draft", "Create a work order draft.", runCreateWorkOrderDraft)
	}
	builders["update_work_order_draft"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[UpdateWorkOrderDraftInput, string]("update_work_order_draft", "Update an existing work order draft.", runUpdateWorkOrderDraft)
	}
	builders["create_engineering_order_draft"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[CreateEngineeringOrderDraftInput, string]("create_engineering_order_draft", "Create an engineering order draft for planned production output.", runCreateEngineeringOrderDraft)
	}
	builders["update_engineering_order_draft"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[UpdateEngineeringOrderDraftInput, string]("update_engineering_order_draft", "Update an existing engineering order draft.", runUpdateEngineeringOrderDraft)
	}
	builders["list_engineering_orders"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[ListEngineeringOrdersInput, string]("list_engineering_orders", "List engineering orders.", runListEngineeringOrders)
	}
	builders["get_engineering_order"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[GetEngineeringOrderInput, string]("get_engineering_order", "Get engineering order details.", runGetEngineeringOrder)
	}
	builders["create_inventory_flow_draft"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[CreateInventoryFlowDraftInput, string]("create_inventory_flow_draft", "Create an inbound or outbound inventory flow draft.", runCreateInventoryFlowDraft)
	}
	builders["list_inventory_flows"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[ListInventoryFlowsInput, string]("list_inventory_flows", "Get the latest time-ordered inventory flows as dropdown/list data by scope/status. Use limit for the number of records, default 30.", runListInventoryFlows)
	}
	builders["get_inventory_flow"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[GetInventoryFlowInput, string]("get_inventory_flow", "Get inventory flow details.", runGetInventoryFlow)
	}
	builders["search_items"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[SearchItemsInput, string]("search_items", "Search item definitions/material types. Count fields describe concrete item units: total_count is all units of this item type, in_stock_count is physically in stock, reserved_count is reserved, out_stock_count is already out of stock, pending/qualified/unqualified are quality counts, and available_count is usable stock for new outbound requests. Do not treat total_count as available stock.", runSearchItems)
	}
	builders["get_item"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[GetItemInput, string]("get_item", "Get item definition/material type details.", runGetItem)
	}
	builders["list_item_units"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[ListItemUnitsInput, string]("list_item_units", "List concrete item units in inventory.", runListItemUnits)
	}
	builders["list_pending_inventory_flows"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[ListPendingInventoryFlowsInput, string]("list_pending_inventory_flows", "Get the latest time-ordered submitted inventory flows pending warehouse processing. Use limit for the number of records, default 30.", runListPendingInventoryFlows)
	}
	builders["inventory_check"] = func() (tool.BaseTool, error) {
		return toolutils.InferTool[InventoryCheckInput, string]("inventory_check", "Read-only inventory check for item stock and item unit details.", runInventoryCheck)
	}
	return builders
}

func runListWorkOrders(ctx context.Context, input ListWorkOrdersInput) (string, error) {
	client, err := infra.WorkOrderClient()
	if err != nil {
		return "", err
	}
	userID, err := effectiveUserID(ctx, input.UserID)
	if err != nil {
		return "", err
	}
	status := parseWorkOrderStatus(input.Status)
	isTo := input.IsTo
	if status == workorderpb.WorkOrderStatus_WORK_ORDER_STATUS_DRAFT {
		isTo = false
	}
	resp, err := client.ListWorkOrder(ctx, &workorderpb.ListWorkOrderReq{
		PageNum:    1,
		PageSize:   pageSize(input.Limit),
		Id:         userID,
		IsTo:       isTo,
		IsUnread:   input.Unread,
		NamePrefix: strings.TrimSpace(input.NamePrefix),
		Status:     status,
	})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runMarkWorkOrderRead(ctx context.Context, input MarkWorkOrderReadInput) (string, error) {
	client, err := infra.WorkOrderClient()
	if err != nil {
		return "", err
	}
	current, err := client.GetWorkOrder(ctx, &workorderpb.GetWorkOrderReq{Id: input.ID})
	if err != nil {
		return "", err
	}
	if err := requireAICanMarkWorkOrderRead(ctx, current.GetWorkOrder()); err != nil {
		return "", err
	}
	resp, err := client.MarkWorkOrderRead(ctx, &workorderpb.MarkWorkOrderReadReq{Id: input.ID})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runSearchUsers(ctx context.Context, input SearchUsersInput) (string, error) {
	client, err := infra.UserClient()
	if err != nil {
		return "", err
	}
	if input.ID > 0 {
		user, err := client.GetUser(ctx, &userpb.GetUserReq{Id: input.ID})
		if err != nil {
			return "", err
		}
		if roleMatches(input.Role, user.GetUserRole()) {
			return marshalJSON(searchUsersOutput{
				Users: []safeUserInfo{safeUser(user)},
				Total: 1,
				Note:  "password and private credential fields are intentionally omitted",
			})
		}
		return marshalJSON(searchUsersOutput{Users: []safeUserInfo{}, Total: 0})
	}
	resp, err := client.ListUser(ctx, &userpb.ListUserReq{
		PageNum:  1,
		PageSize: pageSize(input.PageSize),
		UserName: strings.TrimSpace(input.Name),
		Account:  strings.TrimSpace(input.Account),
	})
	if err != nil {
		return "", err
	}
	users := make([]safeUserInfo, 0, len(resp.GetUserList()))
	for _, user := range resp.GetUserList() {
		if !roleMatches(input.Role, user.GetUserRole()) {
			continue
		}
		users = append(users, safeUser(user))
	}
	return marshalJSON(searchUsersOutput{
		Users: users,
		Total: resp.GetTotal(),
		Note:  "password and private credential fields are intentionally omitted",
	})
}

func runCreateWorkOrderDraft(ctx context.Context, input CreateWorkOrderDraftInput) (string, error) {
	client, err := infra.WorkOrderClient()
	if err != nil {
		return "", err
	}
	fromUserID, err := effectiveUserID(ctx, input.FromUserID)
	if err != nil {
		return "", err
	}
	resp, err := client.CreateWorkOrder(ctx, &workorderpb.CreateWorkOrderReq{
		FromUserId:  fromUserID,
		ToUserId:    input.ToUserID,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		return "", err
	}
	return marshalDraftCreated("work_order", resp.GetId())
}

func runUpdateWorkOrderDraft(ctx context.Context, input UpdateWorkOrderDraftInput) (string, error) {
	client, err := infra.WorkOrderClient()
	if err != nil {
		return "", err
	}
	current, err := client.GetWorkOrder(ctx, &workorderpb.GetWorkOrderReq{Id: input.ID})
	if err != nil {
		return "", err
	}
	if err := requireAICanUpdateWorkOrderDraft(ctx, current.GetWorkOrder()); err != nil {
		return "", err
	}
	fromUserID, err := effectiveUserID(ctx, input.FromUserID)
	if err != nil {
		return "", err
	}
	resp, err := client.UpdateWorkOrderDraft(ctx, &workorderpb.UpdateWorkOrderDraftReq{
		Id:          input.ID,
		FromUserId:  fromUserID,
		ToUserId:    input.ToUserID,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runCreateEngineeringOrderDraft(ctx context.Context, input CreateEngineeringOrderDraftInput) (string, error) {
	if !rpcmeta.CanCreateEngineeringOrder(operatorRole(ctx, conf.GetConf().AITools.RoleAliases)) {
		return "", fmt.Errorf("create_engineering_order_draft requires leader or admin role")
	}
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	leaderUserID, err := effectiveUserID(ctx, input.LeaderUserID)
	if err != nil {
		return "", err
	}
	resp, err := client.CreateEngineeringOrderDraft(ctx, &inventorypb.CreateEngineeringOrderDraftReq{
		LeaderUserId:      leaderUserID,
		ProcessId:         input.ProcessID,
		ItemId:            input.ItemID,
		Name:              input.Name,
		ExpectedQuantity:  input.ExpectedQuantity,
		QualifiedQuantity: input.QualifiedQuantity,
		Description:       input.Description,
	})
	if err != nil {
		return "", err
	}
	return marshalDraftCreated("engineering_order", resp.GetId())
}

func runUpdateEngineeringOrderDraft(ctx context.Context, input UpdateEngineeringOrderDraftInput) (string, error) {
	if !rpcmeta.CanCreateEngineeringOrder(operatorRole(ctx, conf.GetConf().AITools.RoleAliases)) {
		return "", fmt.Errorf("update_engineering_order_draft requires leader or admin role")
	}
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	current, err := client.GetEngineeringOrder(ctx, &inventorypb.GetEngineeringOrderReq{Id: input.ID})
	if err != nil {
		return "", err
	}
	if err := requireAICanUpdateEngineeringOrderDraft(ctx, current.GetEngineeringOrder()); err != nil {
		return "", err
	}
	leaderUserID, err := effectiveUserID(ctx, input.LeaderUserID)
	if err != nil {
		return "", err
	}
	_, err = client.UpdateEngineeringOrderDraft(ctx, &inventorypb.UpdateEngineeringOrderDraftReq{
		Id:                input.ID,
		LeaderUserId:      leaderUserID,
		ProcessId:         input.ProcessID,
		ItemId:            input.ItemID,
		Name:              input.Name,
		ExpectedQuantity:  input.ExpectedQuantity,
		QualifiedQuantity: input.QualifiedQuantity,
		Description:       input.Description,
	})
	if err != nil {
		return "", err
	}
	return marshalJSON(map[string]any{
		"id":      input.ID,
		"status":  "draft_updated",
		"message": "草稿已更新",
	})
}

func runListEngineeringOrders(ctx context.Context, input ListEngineeringOrdersInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	leaderUserID, err := effectiveOptionalUserID(ctx, input.LeaderUserID)
	if err != nil {
		return "", err
	}
	resp, err := client.ListEngineeringOrder(ctx, &inventorypb.ListEngineeringOrderReq{
		LeaderUserId: leaderUserID,
		ItemId:       input.ItemID,
		PageNum:      pageNum(input.PageNum),
		PageSize:     pageSize(input.PageSize),
		NamePrefix:   strings.TrimSpace(input.NamePrefix),
	})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runGetEngineeringOrder(ctx context.Context, input GetEngineeringOrderInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	resp, err := client.GetEngineeringOrder(ctx, &inventorypb.GetEngineeringOrderReq{Id: input.ID})
	if err != nil {
		return "", err
	}
	if err := requireAICanViewEngineeringOrder(ctx, resp.GetEngineeringOrder()); err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runCreateInventoryFlowDraft(ctx context.Context, input CreateInventoryFlowDraftInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	fromUserID, err := effectiveUserID(ctx, input.FromUserID)
	if err != nil {
		return "", err
	}
	resp, err := client.CreateInventoryFlow(ctx, &inventorypb.CreateInventoryFlowReq{
		FromUserId:  fromUserID,
		ToUserId:    input.ToUserID,
		FlowType:    parseFlowType(input.FlowType),
		Name:        input.Name,
		Description: input.Description,
		Items:       inventoryFlowItems(input.Items),
		ItemUnitIds: input.ItemUnitIDs,
	})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runListInventoryFlows(ctx context.Context, input ListInventoryFlowsInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	userID, err := effectiveUserID(ctx, input.UserID)
	if err != nil {
		return "", err
	}
	isTo := isInventoryFlowToScope(input.Scope)
	flowStatus := parseFlowStatus(input.FlowStatus)
	resp, err := client.ListInventoryFlow(ctx, &inventorypb.ListInventoryFlowReq{
		UserId:     userID,
		IsTo:       isTo,
		FlowStatus: flowStatus,
		NamePrefix: strings.TrimSpace(input.NamePrefix),
		PageNum:    1,
		PageSize:   pageSize(input.Limit),
	})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runGetInventoryFlow(ctx context.Context, input GetInventoryFlowInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	resp, err := client.GetInventoryFlow(ctx, &inventorypb.GetInventoryFlowReq{Id: input.ID})
	if err != nil {
		return "", err
	}
	if err := requireAICanViewInventoryFlow(ctx, resp.GetInventoryFlow()); err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runSearchItems(ctx context.Context, input SearchItemsInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	resp, err := client.ListItem(ctx, &inventorypb.ListItemReq{
		PageNum:    pageNum(input.PageNum),
		PageSize:   pageSize(input.PageSize),
		NamePrefix: input.NamePrefix,
	})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runGetItem(ctx context.Context, input GetItemInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	resp, err := client.GetItem(ctx, &inventorypb.GetItemReq{Id: input.ID})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runListItemUnits(ctx context.Context, input ListItemUnitsInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	resp, err := client.ListItemUnit(ctx, &inventorypb.ListItemUnitReq{
		PageNum:       pageNum(input.PageNum),
		PageSize:      pageSize(input.PageSize),
		ItemId:        input.ItemID,
		StockStatus:   parseStockStatus(input.StockStatus),
		QualityStatus: parseQualityStatus(input.QualityStatus),
	})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runListPendingInventoryFlows(ctx context.Context, input ListPendingInventoryFlowsInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	userID, err := effectiveUserID(ctx, input.UserID)
	if err != nil {
		return "", err
	}
	resp, err := client.ListInventoryFlow(ctx, &inventorypb.ListInventoryFlowReq{
		UserId:     userID,
		IsTo:       true,
		FlowStatus: inventorypb.FlowStatus_FLOW_STATUS_SUBMITTED,
		NamePrefix: strings.TrimSpace(input.NamePrefix),
		PageNum:    1,
		PageSize:   pageSize(input.Limit),
	})
	if err != nil {
		return "", err
	}
	return marshalProto(resp)
}

func runInventoryCheck(ctx context.Context, input InventoryCheckInput) (string, error) {
	client, err := infra.InventoryClient()
	if err != nil {
		return "", err
	}
	items, err := client.ListItem(ctx, &inventorypb.ListItemReq{
		PageNum:    pageNum(input.PageNum),
		PageSize:   pageSize(input.PageSize),
		NamePrefix: input.NamePrefix,
	})
	if err != nil {
		return "", err
	}
	units, err := client.ListItemUnit(ctx, &inventorypb.ListItemUnitReq{
		PageNum:       pageNum(input.PageNum),
		PageSize:      pageSize(input.PageSize),
		ItemId:        input.ItemID,
		StockStatus:   parseStockStatus(input.StockStatus),
		QualityStatus: parseQualityStatus(input.QualityStatus),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"items":%s,"itemUnits":%s}`, mustMarshalProto(items), mustMarshalProto(units)), nil
}

func marshalProto(msg proto.Message) (string, error) {
	data, err := protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: true}.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func marshalJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func marshalDraftCreated(kind string, id int64) (string, error) {
	return marshalJSON(map[string]any{
		"id":         id,
		"kind":       kind,
		"status":     "draft",
		"visible_to": "creator",
		"message":    "草稿已创建，仅创建人可见；可以在自己的发件/草稿列表中继续提交或更新。",
	})
}

func safeUser(user *userpb.GetUserResp) safeUserInfo {
	if user == nil {
		return safeUserInfo{}
	}
	return safeUserInfo{
		ID:      user.GetId(),
		Account: user.GetUserAccount(),
		Name:    user.GetUserName(),
		Role:    rpcmeta.NormalizeRole(user.GetUserRole(), conf.GetConf().AITools.RoleAliases),
	}
}

func roleMatches(filter string, userRole string) bool {
	normalized := rpcmeta.NormalizeRole(filter, conf.GetConf().AITools.RoleAliases)
	if strings.TrimSpace(filter) == "" {
		return true
	}
	return rpcmeta.NormalizeRole(userRole, conf.GetConf().AITools.RoleAliases) == normalized
}

func mustMarshalProto(msg proto.Message) string {
	data, err := marshalProto(msg)
	if err != nil {
		return "{}"
	}
	return data
}

func effectiveUserID(ctx context.Context, inputUserID int64) (int64, error) {
	if inputUserID > 0 && rpcmeta.IsAdmin(operatorRole(ctx, conf.GetConf().AITools.RoleAliases)) {
		return inputUserID, nil
	}
	return operatorID(ctx)
}

func effectiveOptionalUserID(ctx context.Context, inputUserID int64) (int64, error) {
	if rpcmeta.IsAdmin(operatorRole(ctx, conf.GetConf().AITools.RoleAliases)) {
		return inputUserID, nil
	}
	return operatorID(ctx)
}

func operatorID(ctx context.Context) (int64, error) {
	if id, ok := rpcmeta.OperatorIDFromContext(ctx); ok && id > 0 {
		return id, nil
	}
	raw := strings.TrimSpace(rpcmeta.FromContext(ctx).OperatorID)
	if raw == "" {
		return 0, fmt.Errorf("operator id is missing from rpc meta")
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("operator id %q is invalid", raw)
	}
	return id, nil
}

func aiIsAdmin(ctx context.Context) bool {
	return rpcmeta.IsAdmin(operatorRole(ctx, conf.GetConf().AITools.RoleAliases))
}

func requireAICanViewWorkOrder(ctx context.Context, order *workorderpb.WorkOrderInfo) error {
	if order == nil {
		return fmt.Errorf("work order not found")
	}
	userID, err := operatorID(ctx)
	if err != nil {
		return err
	}
	if order.GetStatus() == workorderpb.WorkOrderStatus_WORK_ORDER_STATUS_DRAFT {
		if aiIsAdmin(ctx) || order.GetFromUserId() == userID {
			return nil
		}
		return fmt.Errorf("forbidden: no permission")
	}
	if aiIsAdmin(ctx) || order.GetFromUserId() == userID || order.GetToUserId() == userID {
		return nil
	}
	return fmt.Errorf("forbidden: no permission")
}

func requireAICanUpdateWorkOrderDraft(ctx context.Context, order *workorderpb.WorkOrderInfo) error {
	if order == nil {
		return fmt.Errorf("work order not found")
	}
	userID, err := operatorID(ctx)
	if err != nil {
		return err
	}
	if aiIsAdmin(ctx) || order.GetFromUserId() == userID {
		return nil
	}
	return fmt.Errorf("forbidden: no permission")
}

func requireAICanMarkWorkOrderRead(ctx context.Context, order *workorderpb.WorkOrderInfo) error {
	if order == nil {
		return fmt.Errorf("work order not found")
	}
	userID, err := operatorID(ctx)
	if err != nil {
		return err
	}
	if aiIsAdmin(ctx) || order.GetToUserId() == userID {
		return nil
	}
	return fmt.Errorf("forbidden: no permission")
}

func requireAICanViewEngineeringOrder(ctx context.Context, order *inventorypb.EngineeringOrderInfo) error {
	if order == nil {
		return fmt.Errorf("engineering order not found")
	}
	if order.GetStatus() != inventorypb.DraftStatus_DRAFT_STATUS_DRAFT {
		return nil
	}
	return requireAICanUpdateEngineeringOrderDraft(ctx, order)
}

func requireAICanUpdateEngineeringOrderDraft(ctx context.Context, order *inventorypb.EngineeringOrderInfo) error {
	if order == nil {
		return fmt.Errorf("engineering order not found")
	}
	userID, err := operatorID(ctx)
	if err != nil {
		return err
	}
	if aiIsAdmin(ctx) || order.GetLeaderUserId() == userID {
		return nil
	}
	return fmt.Errorf("forbidden: no permission")
}

func requireAICanViewInventoryFlow(ctx context.Context, flow *inventorypb.InventoryFlowInfo) error {
	if flow == nil {
		return fmt.Errorf("inventory flow not found")
	}
	userID, err := operatorID(ctx)
	if err != nil {
		return err
	}
	if flow.GetFlowStatus() == inventorypb.FlowStatus_FLOW_STATUS_DRAFT {
		if aiIsAdmin(ctx) || flow.GetFromUserId() == userID {
			return nil
		}
		return fmt.Errorf("forbidden: no permission")
	}
	if aiIsAdmin(ctx) || flow.GetFromUserId() == userID || flow.GetToUserId() == userID {
		return nil
	}
	return fmt.Errorf("forbidden: no permission")
}

func pageNum(v int64) int64 {
	if v <= 0 {
		return 1
	}
	return v
}

func pageSize(v int64) int64 {
	if v <= 0 {
		return 30
	}
	if v > 100 {
		return 100
	}
	return v
}

func inventoryFlowItems(items []InventoryFlowItemInput) []*inventorypb.InventoryFlowItemReq {
	out := make([]*inventorypb.InventoryFlowItemReq, 0, len(items))
	for _, item := range items {
		out = append(out, &inventorypb.InventoryFlowItemReq{
			ItemId:        item.ItemID,
			ApplyQuantity: item.ApplyQuantity,
		})
	}
	return out
}

func parseFlowType(value string) inventorypb.FlowType {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "in", "入库", "flow_type_in":
		return inventorypb.FlowType_FLOW_TYPE_IN
	case "out", "出库", "flow_type_out":
		return inventorypb.FlowType_FLOW_TYPE_OUT
	default:
		return inventorypb.FlowType_FLOW_TYPE_UNKNOWN
	}
}

func parseFlowStatus(value string) inventorypb.FlowStatus {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "draft", "草稿", "flow_status_draft":
		return inventorypb.FlowStatus_FLOW_STATUS_DRAFT
	case "submitted", "pending", "待审核", "已提交", "flow_status_submitted":
		return inventorypb.FlowStatus_FLOW_STATUS_SUBMITTED
	case "approved", "accepted", "通过", "已通过", "flow_status_approved":
		return inventorypb.FlowStatus_FLOW_STATUS_APPROVED
	case "rejected", "拒绝", "已拒绝", "flow_status_rejected":
		return inventorypb.FlowStatus_FLOW_STATUS_REJECTED
	default:
		return inventorypb.FlowStatus_FLOW_STATUS_UNKNOWN
	}
}

func parseWorkOrderStatus(value string) workorderpb.WorkOrderStatus {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "draft", "草稿", "work_order_status_draft":
		return workorderpb.WorkOrderStatus_WORK_ORDER_STATUS_DRAFT
	case "submitted", "已提交", "work_order_status_submitted":
		return workorderpb.WorkOrderStatus_WORK_ORDER_STATUS_SUBMITTED
	default:
		return workorderpb.WorkOrderStatus_WORK_ORDER_STATUS_UNKNOWN
	}
}

func parseStockStatus(value string) inventorypb.StockStatus {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "in_stock", "in", "在库", "stock_status_in_stock":
		return inventorypb.StockStatus_STOCK_STATUS_IN_STOCK
	case "reserved", "已预留", "stock_status_reserved":
		return inventorypb.StockStatus_STOCK_STATUS_RESERVED
	case "out_stock", "out", "出库", "stock_status_out_stock":
		return inventorypb.StockStatus_STOCK_STATUS_OUT_STOCK
	default:
		return inventorypb.StockStatus_STOCK_STATUS_UNKNOWN
	}
}

func parseQualityStatus(value string) inventorypb.QualityStatus {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "pending", "待检", "quality_status_pending":
		return inventorypb.QualityStatus_QUALITY_STATUS_PENDING
	case "qualified", "合格", "quality_status_qualified":
		return inventorypb.QualityStatus_QUALITY_STATUS_QUALIFIED
	case "unqualified", "不合格", "quality_status_unqualified":
		return inventorypb.QualityStatus_QUALITY_STATUS_UNQUALIFIED
	default:
		return inventorypb.QualityStatus_QUALITY_STATUS_UNKNOWN
	}
}

func isInventoryFlowToScope(scope string) bool {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "to_me", "to", "assigned_to_me", "待我处理":
		return true
	default:
		return false
	}
}
