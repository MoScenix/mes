# MES Table Logic

This document records the current MES data model, table relationships, status rules, list-query rules, and index design direction.

## Common Rules

All main list pages use time cursor ordering:

```sql
ORDER BY updated_at DESC, id DESC
```

Cursor condition:

```sql
updated_at < :cursorUpdatedAt
OR (updated_at = :cursorUpdatedAt AND id < :cursorId)
```

All GORM models use soft delete through `deleted_at`. Main-table list indexes should start with `deleted_at`, then equality filters, then `updated_at`, then `id`.

Recommended index shape:

```text
deleted_at + equality filter columns + updated_at + id
```

Name prefix search uses:

```sql
name LIKE 'prefix%'
```

For prefix search plus time ordering, keep separate indexes for common access paths instead of trying to make one universal index.

## Tables

### items

Entity: material catalog row.

Fields:

| Field | Type | Meaning |
|---|---|---|
| id | uint | Primary key |
| created_at | time | Created time |
| updated_at | time | Updated time, list cursor |
| deleted_at | soft delete | Soft delete marker |
| name | varchar(100) | Material name |
| unit | varchar(20) | Unit |
| description | varchar(255) | Description |
| total_count | int64 | Total item-unit count |
| in_stock_count | int64 | Item units in stock |
| reserved_count | int64 | Reserved item units |
| out_stock_count | int64 | Out-of-stock item units |
| pending_count | int64 | Pending quality count |
| qualified_count | int64 | Qualified count |
| unqualified_count | int64 | Unqualified count |
| available_count | int64 | In stock and qualified count |

Relationships:

| Relation | Table | Type |
|---|---|---|
| item_units.item_id | item_units | One item has many item units |
| processes.item_id | processes | One item can be process output |
| process_items.consume_item_id | process_items | One item can be consumed by process |
| engineering_orders.item_id | engineering_orders | One item can be produced by engineering orders |
| inventory_flow_items.item_id | inventory_flow_items | One item can appear in flow requested items |

Current high-frequency queries:

| Query | Filter | Order |
|---|---|---|
| list materials | optional `namePrefix` | `updated_at DESC, id DESC` |
| search item picker | `namePrefix` | `updated_at DESC, id DESC` |
| aggregate recalculation | `id = ?` | none |

Important indexes:

| Purpose | Index columns |
|---|---|
| normal time list | `deleted_at, updated_at, id` |
| name prefix list | `deleted_at, name, updated_at, id` |
| id lookup | primary key |

### item_units

Entity: physical item unit.

Fields:

| Field | Type | Meaning |
|---|---|---|
| id | uint | Primary key |
| created_at | time | Created time |
| updated_at | time | Updated time, list cursor |
| deleted_at | soft delete | Soft delete marker |
| item_id | uint | Material id |
| engineering_order_id | nullable uint | Producing engineering order |
| stock_status | int32 | Stock status |
| quality_status | int32 | Quality status |
| description | varchar(255) | Description |

Enums:

| stock_status | Meaning |
|---|---|
| 0 | unknown |
| 1 | in stock |
| 2 | reserved |
| 3 | out of stock |

| quality_status | Meaning |
|---|---|
| 0 | unknown |
| 1 | pending |
| 2 | qualified |
| 3 | unqualified |

Relationships:

| Relation | Table | Type |
|---|---|---|
| item_units.item_id | items | Many units belong to one item |
| item_units.engineering_order_id | engineering_orders | Many units can be produced by one engineering order |
| inventory_flow_item_units.item_unit_id | inventory_flow_item_units | Many-to-many with inventory flows |

Current high-frequency queries:

| Query | Filter | Order |
|---|---|---|
| list item units | optional `itemId`, `engineeringOrderId`, `inventoryFlowId`, `stockStatus`, `qualityStatus`, `itemNamePrefix` | `item_units.updated_at DESC, item_units.id DESC` |
| item detail trace | `itemUnitId` through flow relation | flow time order |
| engineering order related units | `engineeringOrderId` | item-unit time order |
| flow related units | `inventoryFlowId` through join table | item-unit time order |
| item count recalculation | `item_id = ?` | none |
| engineering produced recalculation | `engineering_order_id = ?` | none |

Important indexes:

| Purpose | Index columns |
|---|---|
| normal time list | `deleted_at, updated_at, id` |
| by item | `deleted_at, item_id, updated_at, id` |
| by engineering order | `deleted_at, engineering_order_id, updated_at, id` |
| by stock and quality | `deleted_at, stock_status, quality_status, updated_at, id` |
| by stock only | `deleted_at, stock_status, updated_at, id` |
| by quality only | `deleted_at, quality_status, updated_at, id` |
| by item and statuses | `deleted_at, item_id, stock_status, quality_status, updated_at, id` |
| by engineering and statuses | `deleted_at, engineering_order_id, stock_status, quality_status, updated_at, id` |
| aggregate by item | `deleted_at, item_id` or `deleted_at, item_id, stock_status, quality_status` |
| aggregate by engineering order | `deleted_at, engineering_order_id` or `deleted_at, engineering_order_id, quality_status` |

### processes

Entity: process sheet. Process engineers manage these.

Fields:

| Field | Type | Meaning |
|---|---|---|
| id | uint | Primary key |
| created_at | time | Created time |
| updated_at | time | Updated time |
| deleted_at | soft delete | Soft delete marker |
| item_id | uint | Output item |
| owner_user_id | int64 | Process owner |
| name | varchar(100) | Process name |
| description | varchar(255) | Description |
| status | int32 | Draft status |

Enums:

| status | Meaning |
|---|---|
| 1 | draft |
| 2 | submitted |
| 3 | done |

Relationships:

| Relation | Table | Type |
|---|---|---|
| processes.item_id | items | Process outputs one item |
| process_items.process_id | process_items | One process has many consume items |
| engineering_orders.process_id | engineering_orders | One process can have many engineering orders |

Current high-frequency queries:

| Query | Filter | Order |
|---|---|---|
| process list | optional `ownerUserId`, `itemId`, `status`, `namePrefix`, `itemNamePrefix`, recent time | `processes.updated_at DESC, processes.id DESC` |
| process related engineering orders | `processId` through engineering order list | engineering order time order |

Important indexes:

| Purpose | Index columns |
|---|---|
| normal time list | `deleted_at, updated_at, id` |
| owner list | `deleted_at, owner_user_id, updated_at, id` |
| owner and status | `deleted_at, owner_user_id, status, updated_at, id` |
| item list | `deleted_at, item_id, updated_at, id` |
| item and status | `deleted_at, item_id, status, updated_at, id` |
| status list | `deleted_at, status, updated_at, id` |
| name prefix | `deleted_at, name, updated_at, id` |
| owner name prefix | `deleted_at, owner_user_id, name, updated_at, id` |

### process_items

Entity: process consume item line.

Fields:

| Field | Type | Meaning |
|---|---|---|
| id | uint | Primary key |
| created_at | time | Created time |
| updated_at | time | Updated time |
| deleted_at | soft delete | Soft delete marker |
| process_id | uint | Process id |
| consume_item_id | uint | Consumed item id |
| quantity | int64 | Required quantity |

Relationships:

| Relation | Table | Type |
|---|---|---|
| process_items.process_id | processes | Many consume lines belong to one process |
| process_items.consume_item_id | items | Consume material |

Important indexes:

| Purpose | Index columns |
|---|---|
| unique consume line | `process_id, consume_item_id` |
| delete or replace draft lines | `deleted_at, process_id` |
| reverse item lookup | `deleted_at, consume_item_id, process_id` |

### engineering_orders

Entity: engineering order. Engineers can see their own orders and process-related engineering orders.

Fields:

| Field | Type | Meaning |
|---|---|---|
| id | uint | Primary key |
| created_at | time | Created time |
| updated_at | time | Updated time |
| deleted_at | soft delete | Soft delete marker |
| leader_user_id | int64 | Owner / leader |
| process_id | uint | Related process |
| item_id | uint | Produced item |
| name | varchar(100) | Engineering order name |
| expected_quantity | int64 | Planned quantity |
| qualified_quantity | int64 | Qualified target or count field |
| unqualified_quantity | int64 | Unqualified count |
| produced_quantity | int64 | Produced unit count |
| status | int32 | Draft status |
| description | varchar(255) | Description |

Relationships:

| Relation | Table | Type |
|---|---|---|
| engineering_orders.process_id | processes | Many orders can use one process |
| engineering_orders.item_id | items | Engineering order produces one item |
| item_units.engineering_order_id | item_units | One order has many produced units |

Current high-frequency queries:

| Query | Filter | Order |
|---|---|---|
| engineering order list | optional `leaderUserId`, `itemId`, `processId`, `status`, `namePrefix`, `itemNamePrefix`, recent time | `engineering_orders.updated_at DESC, engineering_orders.id DESC` |
| process related orders | `processId` | time order |
| worker engineering list | optional `itemNamePrefix` | time order |

Important indexes:

| Purpose | Index columns |
|---|---|
| normal time list | `deleted_at, updated_at, id` |
| leader list | `deleted_at, leader_user_id, updated_at, id` |
| leader and status | `deleted_at, leader_user_id, status, updated_at, id` |
| item list | `deleted_at, item_id, updated_at, id` |
| item and status | `deleted_at, item_id, status, updated_at, id` |
| process list | `deleted_at, process_id, updated_at, id` |
| process and status | `deleted_at, process_id, status, updated_at, id` |
| item and process | `deleted_at, item_id, process_id, updated_at, id` |
| name prefix | `deleted_at, name, updated_at, id` |
| leader name prefix | `deleted_at, leader_user_id, name, updated_at, id` |
| process name prefix | `deleted_at, process_id, name, updated_at, id` |

### inventory_flows

Entity: inventory flow / circulation order. Used for inbound, outbound, approval, receiving, and scan operations.

Fields:

| Field | Type | Meaning |
|---|---|---|
| id | uint | Primary key |
| created_at | time | Created time |
| updated_at | time | Updated time |
| deleted_at | soft delete | Soft delete marker |
| from_user_id | int64 | Initiator |
| to_user_id | int64 | Receiver |
| flow_type | int32 | Inbound or outbound |
| flow_status | int32 | Flow status |
| name | varchar(100) | Flow name |
| description | varchar(255) | Description |
| approved_by | int64 | Approver |
| approved_at | nullable time | Approval time |

Enums:

| flow_type | Meaning |
|---|---|
| 1 | inbound |
| 2 | outbound |

| flow_status | Meaning |
|---|---|
| 1 | draft |
| 2 | submitted |
| 3 | approved |
| 4 | rejected |

Relationships:

| Relation | Table | Type |
|---|---|---|
| inventory_flow_items.inventory_flow_id | inventory_flow_items | One flow has many requested item lines |
| inventory_flow_item_units.inventory_flow_id | inventory_flow_item_units | One flow can bind many scanned item units |

Current high-frequency queries:

| Query | Filter | Order |
|---|---|---|
| my flow list | `from_user_id` or `to_user_id`, optional `flowStatus`, optional `itemNamePrefix` | `inventory_flows.updated_at DESC, inventory_flows.id DESC` |
| warehouse all flows | optional `flowStatus`, optional `itemNamePrefix` | time order |
| warehouse audit list | `flowStatus = submitted`, optional `itemNamePrefix` | time order |
| flow trace by item unit | `itemUnitId` through relation table | time order |
| related item units | `inventoryFlowId` through item-unit list | item-unit time order |

Important indexes:

| Purpose | Index columns |
|---|---|
| normal time list | `deleted_at, updated_at, id` |
| status list | `deleted_at, flow_status, updated_at, id` |
| initiator list | `deleted_at, from_user_id, updated_at, id` |
| receiver list | `deleted_at, to_user_id, updated_at, id` |
| initiator and status | `deleted_at, from_user_id, flow_status, updated_at, id` |
| receiver and status | `deleted_at, to_user_id, flow_status, updated_at, id` |
| name prefix | `deleted_at, name, updated_at, id` |
| initiator name prefix | `deleted_at, from_user_id, name, updated_at, id` |
| receiver name prefix | `deleted_at, to_user_id, name, updated_at, id` |
| status name prefix | `deleted_at, flow_status, name, updated_at, id` |

### inventory_flow_items

Entity: requested item line in an inventory flow.

Fields:

| Field | Type | Meaning |
|---|---|---|
| id | uint | Primary key |
| created_at | time | Created time |
| updated_at | time | Updated time |
| deleted_at | soft delete | Soft delete marker |
| inventory_flow_id | uint | Flow id |
| item_id | uint | Requested item |
| apply_quantity | int64 | Requested quantity |
| finished_quantity | int64 | Already operated quantity |

Relationships:

| Relation | Table | Type |
|---|---|---|
| inventory_flow_items.inventory_flow_id | inventory_flows | Many requested lines belong to one flow |
| inventory_flow_items.item_id | items | Requested material |

Important indexes:

| Purpose | Index columns |
|---|---|
| unique requested item per flow | `inventory_flow_id, item_id` |
| load flow details | `deleted_at, inventory_flow_id, item_id` |
| reverse item search | `deleted_at, item_id, inventory_flow_id` |

### inventory_flow_item_units

Entity: binding table between inventory flows and scanned item units.

Fields:

| Field | Type | Meaning |
|---|---|---|
| inventory_flow_id | uint | Flow id |
| item_unit_id | uint | Item unit id |

Relationships:

| Relation | Table | Type |
|---|---|---|
| inventory_flow_item_units.inventory_flow_id | inventory_flows | Flow side |
| inventory_flow_item_units.item_unit_id | item_units | Item-unit side |

Important indexes:

| Purpose | Index columns |
|---|---|
| unique binding | `inventory_flow_id, item_unit_id` |
| flow related units | `inventory_flow_id, item_unit_id` |
| item-unit trace | `item_unit_id, inventory_flow_id` |

### work_order

Entity: work order / message-like task.

Fields:

| Field | Type | Meaning |
|---|---|---|
| id | uint | Primary key |
| created_at | time | Created time |
| updated_at | time | Updated time |
| deleted_at | soft delete | Soft delete marker |
| from_user_id | int64 | Sender |
| to_user_id | int64 | Receiver |
| name | varchar(100) | Work order title |
| description | text | Content |
| status | int32 | Draft/submitted |
| read_status | int32 | Unread/read |

Enums:

| status | Meaning |
|---|---|
| 1 | draft |
| 2 | submitted |

| read_status | Meaning |
|---|---|
| 1 | unread |
| 2 | read |

Relationships:

Work orders currently do not have foreign-key model relations to inventory tables. They are related by user ids and UI navigation.

Current high-frequency queries:

| Query | Filter | Order |
|---|---|---|
| sent work orders | `from_user_id`, optional `status`, optional `namePrefix` | `updated_at DESC, id DESC` |
| inbox work orders | `to_user_id`, optional `status`, optional `isUnread`, optional `namePrefix` | time order |

Important indexes:

| Purpose | Index columns |
|---|---|
| sent time list | `deleted_at, from_user_id, updated_at, id` |
| inbox time list | `deleted_at, to_user_id, updated_at, id` |
| sent status list | `deleted_at, from_user_id, status, updated_at, id` |
| inbox status list | `deleted_at, to_user_id, status, updated_at, id` |
| sent unread/read list | `deleted_at, from_user_id, read_status, updated_at, id` |
| inbox unread/read list | `deleted_at, to_user_id, read_status, updated_at, id` |
| sent status + read | `deleted_at, from_user_id, status, read_status, updated_at, id` |
| inbox status + read | `deleted_at, to_user_id, status, read_status, updated_at, id` |
| global name prefix | `deleted_at, name, updated_at, id` |
| sent name prefix | `deleted_at, from_user_id, name, updated_at, id` |
| inbox name prefix | `deleted_at, to_user_id, name, updated_at, id` |
| sent status name prefix | `deleted_at, from_user_id, status, name, updated_at, id` |
| inbox status name prefix | `deleted_at, to_user_id, status, name, updated_at, id` |

## Status And Operation Logic

### Process

Draft lifecycle:

```text
draft -> submitted -> done
```

Current operations:

| Operation | Rule |
|---|---|
| create draft | create process and consume lines |
| update draft | only draft can update; replace consume lines |
| delete draft | only draft can delete; remove consume lines |
| submit | process must have at least one consume line |
| get | loads process meta and, for detail, consume lines |
| list | returns time-ordered dropdown/list rows |

### Engineering Order

Draft lifecycle:

```text
draft -> submitted -> done
```

Current operations:

| Operation | Rule |
|---|---|
| create draft | validate process and item |
| update draft | only draft can update |
| delete draft | only draft can delete |
| submit | only draft can submit |
| get | recalculates produced/qualified/unqualified quantities |
| list | time-ordered rows with optional process/item filters |

### Inventory Flow

Lifecycle:

```text
draft -> submitted -> approved / rejected
approved -> complete scan operation
```

Current operation logic:

| Operation | Rule |
|---|---|
| create draft | store flow meta, requested item lines, optional initial item-unit ids |
| update draft | only draft can update; replace requested lines and bindings |
| submit | only draft can submit |
| audit | submitted flow can be approved or rejected |
| complete | approved flow binds scanned item units and updates item-unit stock/quality |
| item progress | `inventory_flow_items.finished_quantity` is incremented as units are completed |
| related item units | item units are related only after scan/bind |

Important detail:

`InventoryFlow.items` is intrinsic to the flow and should be returned with flow detail/list for progress. `InventoryFlow.itemUnits` can be large and should be loaded by `ListItemUnit(inventoryFlowId)` when the UI needs the related unit dropdown/list.

### Item Unit

Current operation logic:

| Operation | Rule |
|---|---|
| add unit | create physical unit; optionally bind engineering order |
| update status | update stock and quality status |
| complete flow | bind unit to flow and update status |
| recalculate item counts | aggregate by `item_id` |
| recalculate engineering order counts | aggregate by `engineering_order_id` |

## Query Simplification Plan

The lowest-complexity list design is:

1. Keep every list time ordered.
2. Use many dedicated composite indexes for common filter combinations.
3. Avoid loading large associations in detail APIs unless the association is intrinsic and small.
4. For related dropdown/list views, call dedicated list endpoints with filter ids.
5. For item-name search across related tables, prefer two-step filtering:
   - first find matching `item.id` by `items.name LIKE prefix%`
   - then filter relation tables by `item_id`

