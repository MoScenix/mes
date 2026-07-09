# MES Filter Interfaces

This document records the current MES HTTP interfaces, internal RPC interfaces, list filters, and response shape.

## Shared List Contract

All list responses are dropdown/list style and time ordered.

Request fields used by most lists:

| Field | Meaning |
|---|---|
| pageNum | Compatibility field; current UI mainly uses cursor flow |
| pageSize | Number of rows to return |
| cursorUpdatedAt | Cursor updated time from previous response |
| cursorId | Cursor id from previous response |
| sinceTime | Optional absolute lower bound |
| recentSeconds | Optional relative lower bound |
| namePrefix | Prefix search on entity name |
| itemNamePrefix | Prefix search on related item name |

Response fields used by every list page:

| Field | Meaning |
|---|---|
| records | Current rows |
| hasMore | Whether another cursor request exists |
| nextCursorUpdatedAt | Cursor updated time for next request |
| nextCursorId | Cursor id for next request |

List ordering:

```text
updated_at DESC, id DESC
```

## Scope Rules

`MesListScope`:

| Value | Meaning |
|---|---|
| MES_LIST_SCOPE_UNSPECIFIED | Default behavior |
| MES_LIST_SCOPE_MINE | Current user's own rows |
| MES_LIST_SCOPE_ALL | All rows |
| MES_LIST_SCOPE_AUDIT | Rows pending audit |
| MES_LIST_SCOPE_BY_PROCESS | Rows related to process |

Current note:

The code currently passes scope/filter values through the BFF. This is not a complete permission system yet. The intended permission model is:

| Role | Visible data |
|---|---|
| engineer | own engineering orders; all process sheets through intelligent/search list; engineering order related process and item-unit lists |
| process engineer | own process sheets; process-related engineering order lists |
| warehouse | all inventory flows; audit submitted flows |
| other user | own inventory flows and own work orders |

## BFF HTTP Interfaces

Source: `idl/bff/mes_bff.proto`.

### Work Order

| Method | Path | Request | Response |
|---|---|---|---|
| POST | `/mes/work-order/draft/create` | `CreateWorkOrderDraftRequest` | `BaseResponseLong` |
| POST | `/mes/work-order/draft/update` | `UpdateWorkOrderDraftRequest` | `BaseResponseBoolean` |
| POST | `/mes/work-order/draft/delete` | `DeleteRequest` | `BaseResponseBoolean` |
| POST | `/mes/work-order/submit` | `DeleteRequest` | `BaseResponseBoolean` |
| GET | `/mes/work-order/get` | `GetByIdRequest` | `BaseResponseWorkOrderVO` |
| POST | `/mes/work-order/list` | `ListWorkOrderRequest` | `BaseResponsePageWorkOrderVO` |
| POST | `/mes/work-order/read` | `DeleteRequest` | `BaseResponseBoolean` |

#### CreateWorkOrderDraftRequest

| Field | Meaning |
|---|---|
| toUserId | Receiver |
| description | Content |
| name | Title |

#### UpdateWorkOrderDraftRequest

| Field | Meaning |
|---|---|
| id | Work order id |
| toUserId | Receiver |
| description | Content |
| name | Title |

#### ListWorkOrderRequest

| Field | Meaning |
|---|---|
| pageNum | Compatibility page number |
| pageSize | Row count |
| id | User id filter |
| isTo | true = inbox, false = sent |
| isUnread | unread only |
| sinceTime | Absolute lower time |
| recentSeconds | Relative lower time |
| cursorUpdatedAt | Time cursor |
| cursorId | Id cursor |
| namePrefix | Work order title prefix |
| status | Draft/submitted |
| scope | Mine/all/audit style scope |

### Process

| Method | Path | Request | Response |
|---|---|---|---|
| POST | `/mes/process/draft/create` | `CreateProcessDraftRequest` | `BaseResponseLong` |
| POST | `/mes/process/draft/update` | `UpdateProcessDraftRequest` | `BaseResponseBoolean` |
| POST | `/mes/process/draft/delete` | `DeleteRequest` | `BaseResponseBoolean` |
| POST | `/mes/process/submit` | `DeleteRequest` | `BaseResponseBoolean` |
| GET | `/mes/process/get` | `GetByIdRequest` | `BaseResponseProcessVO` |
| POST | `/mes/process/list` | `ListProcessRequest` | `BaseResponsePageProcessVO` |

#### CreateProcessDraftRequest

| Field | Meaning |
|---|---|
| ownerUserId | Owner user id |
| itemId | Output item id |
| name | Process name |
| description | Description |
| items | Consume item lines |

#### ProcessItemRequest

| Field | Meaning |
|---|---|
| consumeItemId | Consumed material id |
| quantity | Quantity |

#### UpdateProcessDraftRequest

| Field | Meaning |
|---|---|
| id | Process id |
| ownerUserId | Owner user id |
| itemId | Output item id |
| name | Process name |
| description | Description |
| items | Consume item lines |

#### ListProcessRequest

| Field | Meaning |
|---|---|
| ownerUserId | Owner filter |
| itemId | Output item filter |
| status | Draft/submitted/done |
| pageNum | Compatibility page number |
| pageSize | Row count |
| sinceTime | Absolute lower time |
| recentSeconds | Relative lower time |
| cursorUpdatedAt | Time cursor |
| cursorId | Id cursor |
| namePrefix | Process name prefix |
| itemNamePrefix | Output item name prefix |
| scope | Mine/all/by process style scope |

### Engineering Order

| Method | Path | Request | Response |
|---|---|---|---|
| POST | `/mes/engineering-order/draft/create` | `CreateEngineeringOrderRequest` | `BaseResponseLong` |
| POST | `/mes/engineering-order/draft/update` | `UpdateEngineeringOrderRequest` | `BaseResponseBoolean` |
| POST | `/mes/engineering-order/draft/delete` | `DeleteRequest` | `BaseResponseBoolean` |
| POST | `/mes/engineering-order/submit` | `DeleteRequest` | `BaseResponseBoolean` |
| GET | `/mes/engineering-order/get` | `GetByIdRequest` | `BaseResponseEngineeringOrderVO` |
| POST | `/mes/engineering-order/list` | `ListEngineeringOrderRequest` | `BaseResponsePageEngineeringOrderVO` |

#### CreateEngineeringOrderRequest

| Field | Meaning |
|---|---|
| leaderUserId | Engineering leader |
| itemId | Produced item |
| expectedQuantity | Expected production quantity |
| qualifiedQuantity | Qualified quantity field |
| description | Description |
| processId | Related process |
| name | Engineering order name |

#### UpdateEngineeringOrderRequest

| Field | Meaning |
|---|---|
| id | Engineering order id |
| leaderUserId | Engineering leader |
| itemId | Produced item |
| expectedQuantity | Expected production quantity |
| qualifiedQuantity | Qualified quantity field |
| description | Description |
| processId | Related process |
| name | Engineering order name |

#### ListEngineeringOrderRequest

| Field | Meaning |
|---|---|
| leaderUserId | Owner filter |
| itemId | Produced item filter |
| pageNum | Compatibility page number |
| pageSize | Row count |
| processId | Related process filter |
| status | Draft/submitted/done |
| sinceTime | Absolute lower time |
| recentSeconds | Relative lower time |
| cursorUpdatedAt | Time cursor |
| cursorId | Id cursor |
| namePrefix | Engineering order name prefix |
| itemNamePrefix | Produced item name prefix |
| scope | Mine/all/by process style scope |

### Item

| Method | Path | Request | Response |
|---|---|---|---|
| POST | `/mes/item/add` | `AddItemRequest` | `BaseResponseLong` |
| POST | `/mes/item/update` | `UpdateItemRequest` | `BaseResponseBoolean` |
| GET | `/mes/item/get` | `GetByIdRequest` | `BaseResponseItemVO` |
| POST | `/mes/item/list` | `ListItemRequest` | `BaseResponsePageItemVO` |
| GET | `/mes/item/search` | `SearchItemsRequest` | `BaseResponsePageItemVO` |

#### AddItemRequest

| Field | Meaning |
|---|---|
| name | Material name |
| unit | Unit |
| description | Description |

#### UpdateItemRequest

| Field | Meaning |
|---|---|
| id | Item id |
| name | Material name |
| unit | Unit |
| description | Description |

#### ListItemRequest / SearchItemsRequest

| Field | Meaning |
|---|---|
| pageNum | Compatibility page number |
| pageSize | Row count |
| namePrefix | Material name prefix |
| cursorUpdatedAt | Time cursor |
| cursorId | Id cursor |

### Item Unit

| Method | Path | Request | Response |
|---|---|---|---|
| POST | `/mes/item-unit/add` | `AddItemUnitRequest` | `BaseResponseLong` |
| POST | `/mes/item-unit/status/update` | `UpdateItemUnitStatusRequest` | `BaseResponseBoolean` |
| GET | `/mes/item-unit/get` | `GetByIdRequest` | `BaseResponseItemUnitVO` |
| POST | `/mes/item-unit/list` | `ListItemUnitRequest` | `BaseResponsePageItemUnitVO` |

#### AddItemUnitRequest

| Field | Meaning |
|---|---|
| itemId | Material id |
| stockStatus | Initial stock status |
| qualityStatus | Initial quality status |
| description | Description |
| engineeringOrderId | Producing engineering order |

#### UpdateItemUnitStatusRequest

| Field | Meaning |
|---|---|
| id | Item unit id |
| stockStatus | New stock status |
| qualityStatus | New quality status |

#### ListItemUnitRequest

| Field | Meaning |
|---|---|
| pageNum | Compatibility page number |
| pageSize | Row count |
| itemId | Material filter |
| stockStatus | Stock status filter |
| qualityStatus | Quality status filter |
| engineeringOrderId | Related engineering order filter |
| cursorId | Id cursor |
| itemNamePrefix | Material name prefix |
| scope | Scope |
| inventoryFlowId | Related flow filter; only units scanned/bound to this flow |
| cursorUpdatedAt | Time cursor |

Important behavior:

`inventoryFlowId` filters through `inventory_flow_item_units`. Item units are related to a flow only after scan/bind, so a flow can have requested `items` but zero related `itemUnits`.

### Inventory Flow

| Method | Path | Request | Response |
|---|---|---|---|
| POST | `/mes/inventory-flow/draft/create` | `CreateInventoryFlowDraftRequest` | `BaseResponseLong` |
| POST | `/mes/inventory-flow/draft/update` | `UpdateInventoryFlowDraftRequest` | `BaseResponseBoolean` |
| POST | `/mes/inventory-flow/draft/delete` | `DeleteRequest` | `BaseResponseBoolean` |
| POST | `/mes/inventory-flow/submit` | `DeleteRequest` | `BaseResponseBoolean` |
| POST | `/mes/inventory-flow/complete` | `CompleteInventoryFlowRequest` | `BaseResponseBoolean` |
| POST | `/mes/inventory-flow/audit` | `AuditInventoryFlowRequest` | `BaseResponseBoolean` |
| GET | `/mes/inventory-flow/get` | `GetByIdRequest` | `BaseResponseInventoryFlowVO` |
| POST | `/mes/inventory-flow/list` | `ListInventoryFlowRequest` | `BaseResponsePageInventoryFlowVO` |

#### InventoryFlowItemRequest

| Field | Meaning |
|---|---|
| itemId | Requested material |
| applyQuantity | Requested quantity |

#### CreateInventoryFlowDraftRequest

| Field | Meaning |
|---|---|
| toUserId | Receiver |
| flowType | Inbound/outbound |
| description | Description |
| items | Requested material lines |
| itemUnitIds | Initially bound item units |
| name | Flow name |

#### UpdateInventoryFlowDraftRequest

| Field | Meaning |
|---|---|
| id | Flow id |
| toUserId | Receiver |
| flowType | Inbound/outbound |
| description | Description |
| items | Requested material lines |
| itemUnitIds | Bound item units |
| name | Flow name |

#### AuditInventoryFlowRequest

| Field | Meaning |
|---|---|
| id | Flow id |
| approved | true = approve, false = reject |

#### CompleteInventoryFlowRequest

| Field | Meaning |
|---|---|
| id | Flow id |
| itemUnitIds | Scanned item units to bind/operate |

#### ListInventoryFlowRequest

| Field | Meaning |
|---|---|
| userId | User filter |
| isTo | true = receiver side, false = initiator side |
| flowStatus | Flow status filter |
| pageNum | Compatibility page number |
| pageSize | Row count |
| sinceTime | Absolute lower time |
| recentSeconds | Relative lower time |
| cursorUpdatedAt | Time cursor |
| cursorId | Id cursor |
| namePrefix | Flow name prefix |
| itemNamePrefix | Requested/bound item name prefix |
| scope | Mine/all/audit |
| itemUnitId | Trace flows related to one item unit |

## VO Response Fields

### WorkOrderVO

| Field | Meaning |
|---|---|
| id | Work order id |
| fromUserId | Sender |
| toUserId | Receiver |
| description | Content |
| status | Draft/submitted |
| createTime | Created time |
| updateTime | Updated time |
| readStatus | Unread/read |
| name | Title |

### ItemVO

| Field | Meaning |
|---|---|
| id | Item id |
| name | Material name |
| unit | Unit |
| description | Description |
| totalCount | Total units |
| inStockCount | In-stock units |
| reservedCount | Reserved units |
| outStockCount | Out-of-stock units |
| pendingCount | Pending quality units |
| qualifiedCount | Qualified units |
| unqualifiedCount | Unqualified units |
| availableCount | In-stock and qualified units |
| createTime | Created time |
| updateTime | Updated time |

### ItemUnitVO

| Field | Meaning |
|---|---|
| id | Item unit id |
| itemId | Material id |
| stockStatus | Stock status |
| qualityStatus | Quality status |
| description | Description |
| createTime | Created time |
| updateTime | Updated time |
| engineeringOrderId | Producing engineering order |

### ProcessVO

| Field | Meaning |
|---|---|
| id | Process id |
| itemId | Output item id |
| ownerUserId | Owner |
| name | Process name |
| description | Description |
| status | Draft/submitted/done |
| item | Output item meta |
| items | Consume item lines |
| createTime | Created time |
| updateTime | Updated time |

### EngineeringOrderVO

| Field | Meaning |
|---|---|
| id | Engineering order id |
| leaderUserId | Leader |
| itemId | Produced item id |
| item | Produced item meta |
| expectedQuantity | Expected quantity |
| qualifiedQuantity | Qualified quantity |
| producedQuantity | Produced quantity |
| description | Description |
| itemUnits | Related item units; should be loaded by list endpoint for large lists |
| createTime | Created time |
| updateTime | Updated time |
| processId | Related process id |
| process | Related process meta |
| status | Draft/submitted/done |
| unqualifiedQuantity | Unqualified quantity |
| name | Engineering order name |

### InventoryFlowVO

| Field | Meaning |
|---|---|
| id | Flow id |
| fromUserId | Initiator |
| toUserId | Receiver |
| flowType | Inbound/outbound |
| flowStatus | Draft/submitted/approved/rejected |
| description | Description |
| approvedBy | Approver |
| approvedAt | Approval time |
| items | Requested material lines and progress |
| itemUnits | Bound/scanned item units; large lists should use `ListItemUnit(inventoryFlowId)` |
| createTime | Created time |
| updateTime | Updated time |
| name | Flow name |

### InventoryFlowItemVO

| Field | Meaning |
|---|---|
| id | Flow item line id |
| inventoryFlowId | Flow id |
| itemId | Material id |
| applyQuantity | Requested quantity |
| finishedQuantity | Completed/scanned quantity |
| item | Material meta |

## Internal RPC Interfaces

### InventoryService

| RPC | Request | Response |
|---|---|---|
| AddItem | AddItemReq | AddItemResp |
| UpdateItem | UpdateItemReq | UpdateItemResp |
| GetItem | GetItemReq | GetItemResp |
| ListItem | ListItemReq | ListItemResp |
| CreateProcessDraft | CreateProcessDraftReq | CreateProcessDraftResp |
| UpdateProcessDraft | UpdateProcessDraftReq | UpdateProcessDraftResp |
| DeleteProcessDraft | DeleteProcessDraftReq | DeleteProcessDraftResp |
| SubmitProcess | SubmitProcessReq | SubmitProcessResp |
| GetProcess | GetProcessReq | GetProcessResp |
| ListProcess | ListProcessReq | ListProcessResp |
| AddItemUnit | AddItemUnitReq | AddItemUnitResp |
| UpdateItemUnitStatus | UpdateItemUnitStatusReq | UpdateItemUnitStatusResp |
| GetItemUnit | GetItemUnitReq | GetItemUnitResp |
| ListItemUnit | ListItemUnitReq | ListItemUnitResp |
| CreateInventoryFlow | CreateInventoryFlowReq | CreateInventoryFlowResp |
| UpdateInventoryFlowDraft | UpdateInventoryFlowDraftReq | UpdateInventoryFlowDraftResp |
| DeleteInventoryFlowDraft | DeleteInventoryFlowDraftReq | DeleteInventoryFlowDraftResp |
| SubmitInventoryFlow | SubmitInventoryFlowReq | SubmitInventoryFlowResp |
| CompleteInventoryFlow | CompleteInventoryFlowReq | CompleteInventoryFlowResp |
| AuditInventoryFlow | AuditInventoryFlowReq | AuditInventoryFlowResp |
| GetInventoryFlow | GetInventoryFlowReq | GetInventoryFlowResp |
| ListInventoryFlow | ListInventoryFlowReq | ListInventoryFlowResp |
| CreateEngineeringOrderDraft | CreateEngineeringOrderDraftReq | CreateEngineeringOrderDraftResp |
| UpdateEngineeringOrderDraft | UpdateEngineeringOrderDraftReq | UpdateEngineeringOrderDraftResp |
| DeleteEngineeringOrderDraft | DeleteEngineeringOrderDraftReq | DeleteEngineeringOrderDraftResp |
| SubmitEngineeringOrder | SubmitEngineeringOrderReq | SubmitEngineeringOrderResp |
| GetEngineeringOrder | GetEngineeringOrderReq | GetEngineeringOrderResp |
| ListEngineeringOrder | ListEngineeringOrderReq | ListEngineeringOrderResp |

### WorkOrderService

| RPC | Request | Response |
|---|---|---|
| CreateWorkOrder | CreateWorkOrderReq | CreateWorkOrderResp |
| UpdateWorkOrderDraft | UpdateWorkOrderDraftReq | UpdateWorkOrderDraftResp |
| DeleteWorkOrderDraft | DeleteWorkOrderDraftReq | DeleteWorkOrderDraftResp |
| SubmitWorkOrder | SubmitWorkOrderReq | SubmitWorkOrderResp |
| GetWorkOrder | GetWorkOrderReq | GetWorkOrderResp |
| ListWorkOrder | ListWorkOrderReq | ListWorkOrderResp |
| MarkWorkOrderRead | MarkWorkOrderReadReq | MarkWorkOrderReadResp |

## Filter To Table Mapping

| Interface | Filter | Table path |
|---|---|---|
| ListItem | `namePrefix` | `items.name LIKE prefix%` |
| ListItemUnit | `itemId` | `item_units.item_id` |
| ListItemUnit | `engineeringOrderId` | `item_units.engineering_order_id` |
| ListItemUnit | `inventoryFlowId` | `inventory_flow_item_units.inventory_flow_id -> item_units.id` |
| ListItemUnit | `stockStatus` | `item_units.stock_status` |
| ListItemUnit | `qualityStatus` | `item_units.quality_status` |
| ListItemUnit | `itemNamePrefix` | `item_units.item_id -> items.name LIKE prefix%` |
| ListProcess | `ownerUserId` | `processes.owner_user_id` |
| ListProcess | `itemId` | `processes.item_id` |
| ListProcess | `status` | `processes.status` |
| ListProcess | `namePrefix` | `processes.name LIKE prefix%` |
| ListProcess | `itemNamePrefix` | `processes.item_id -> items.name LIKE prefix%` |
| ListEngineeringOrder | `leaderUserId` | `engineering_orders.leader_user_id` |
| ListEngineeringOrder | `itemId` | `engineering_orders.item_id` |
| ListEngineeringOrder | `processId` | `engineering_orders.process_id` |
| ListEngineeringOrder | `status` | `engineering_orders.status` |
| ListEngineeringOrder | `namePrefix` | `engineering_orders.name LIKE prefix%` |
| ListEngineeringOrder | `itemNamePrefix` | `engineering_orders.item_id -> items.name LIKE prefix%` |
| ListInventoryFlow | `userId + isTo=false` | `inventory_flows.from_user_id` |
| ListInventoryFlow | `userId + isTo=true` | `inventory_flows.to_user_id` |
| ListInventoryFlow | `scope=ALL` | no user filter |
| ListInventoryFlow | `scope=AUDIT` | no user filter, default `flow_status=submitted` |
| ListInventoryFlow | `flowStatus` | `inventory_flows.flow_status` |
| ListInventoryFlow | `namePrefix` | `inventory_flows.name LIKE prefix%` |
| ListInventoryFlow | `itemNamePrefix` | `inventory_flow_items/item_units -> items.name LIKE prefix%` |
| ListInventoryFlow | `itemUnitId` | `inventory_flow_item_units.item_unit_id -> inventory_flows.id` |
| ListWorkOrder | `id + isTo=false` | `work_order.from_user_id` |
| ListWorkOrder | `id + isTo=true` | `work_order.to_user_id` |
| ListWorkOrder | `isUnread` | `work_order.read_status=unread` |
| ListWorkOrder | `status` | `work_order.status` |
| ListWorkOrder | `namePrefix` | `work_order.name LIKE prefix%` |

