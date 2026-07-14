// @ts-ignore
/* eslint-disable */
import request from '@/request'

const jsonHeaders = { 'Content-Type': 'application/json' }

function withLimit<T extends { pageNum?: number; pageSize?: number; limit?: number }>(body: T) {
  const { limit, ...rest } = body
  return {
    pageNum: rest.pageNum ?? 1,
    pageSize: rest.pageSize ?? limit ?? 30,
    ...rest,
  }
}

export enum FlowType {
  Unknown = 0,
  In = 1,
  Out = 2,
}

export enum FlowStatus {
  Unknown = 0,
  Draft = 1,
  Submitted = 2,
  Approved = 3,
  Rejected = 4,
}

export enum StockStatus {
  Unknown = 0,
  InStock = 1,
  Reserved = 2,
  OutStock = 3,
}

export enum QualityStatus {
  Unknown = 0,
  Pending = 1,
  Qualified = 2,
  Unqualified = 3,
}

export enum WorkOrderStatus {
  Unknown = 0,
  Draft = 1,
  Submitted = 2,
}

export enum DraftStatus {
  Unknown = 0,
  Draft = 1,
  Submitted = 2,
  Done = 3,
}

export enum MesListScope {
  Unspecified = 0,
  Mine = 1,
  All = 2,
  Audit = 3,
  ByProcess = 4,
}

export type BaseResponse<T> = {
  code?: number
  data?: T
  message?: string
}

export type PageResult<T> = {
  records?: T[]
  pageNumber?: number
  pageSize?: number
  totalPage?: number
  totalRow?: number
  hasMore?: boolean
  nextCursorUpdatedAt?: string
  nextCursorId?: number
}

export type ItemVO = {
  id?: number
  name?: string
  unit?: string
  description?: string
  totalCount?: number
  inStockCount?: number
  reservedCount?: number
  outStockCount?: number
  pendingCount?: number
  qualifiedCount?: number
  unqualifiedCount?: number
  availableCount?: number
  createTime?: string
  updateTime?: string
}

export type ItemUnitVO = {
  id?: number
  itemId?: number
  stockStatus?: StockStatus
  qualityStatus?: QualityStatus
  description?: string
  createTime?: string
  updateTime?: string
  engineeringOrderId?: number
}

export type ProcessItemRequest = {
  consumeItemId?: number
  quantity?: number
}

export type ProcessItemVO = {
  id?: number
  processId?: number
  consumeItemId?: number
  quantity?: number
  consumeItem?: ItemVO
}

export type ProcessVO = {
  id?: number
  itemId?: number
  ownerUserId?: number
  name?: string
  description?: string
  status?: DraftStatus
  item?: ItemVO
  items?: ProcessItemVO[]
  createTime?: string
  updateTime?: string
}

export type InventoryFlowItemRequest = {
  itemId?: number
  applyQuantity?: number
}

export type InventoryFlowItemVO = {
  id?: number
  inventoryFlowId?: number
  itemId?: number
  applyQuantity?: number
  finishedQuantity?: number
  item?: ItemVO
}

export type InventoryFlowVO = {
  id?: number
  name?: string
  fromUserId?: number
  toUserId?: number
  flowType?: FlowType
  flowStatus?: FlowStatus
  description?: string
  approvedBy?: number
  approvedAt?: string
  items?: InventoryFlowItemVO[]
  itemUnits?: ItemUnitVO[]
  createTime?: string
  updateTime?: string
}

export type EngineeringOrderVO = {
  id?: number
  name?: string
  leaderUserId?: number
  itemId?: number
  item?: ItemVO
  expectedQuantity?: number
  qualifiedQuantity?: number
  producedQuantity?: number
  description?: string
  itemUnits?: ItemUnitVO[]
  createTime?: string
  updateTime?: string
  processId?: number
  process?: ProcessVO
  status?: DraftStatus
  unqualifiedQuantity?: number
}

export type WorkOrderVO = {
  id?: number
  name?: string
  fromUserId?: number
  toUserId?: number
  description?: string
  status?: WorkOrderStatus
  createTime?: string
  updateTime?: string
  readStatus?: API.WorkOrderReadStatus
}

export type CreateInventoryFlowDraftRequest = {
  name?: string
  toUserId?: number
  flowType?: FlowType
  description?: string
  items?: InventoryFlowItemRequest[]
  itemUnitIds?: number[]
}

export type UpdateInventoryFlowDraftRequest = CreateInventoryFlowDraftRequest & {
  id?: number
}

export type CompleteInventoryFlowRequest = {
  id?: number
  itemUnitIds?: number[]
}

export type ListInventoryFlowRequest = {
  userId?: number
  isTo?: boolean
  flowStatus?: FlowStatus
  namePrefix?: string
  itemNamePrefix?: string
  scope?: MesListScope
  itemUnitId?: number
  pageNum?: number
  pageSize?: number
  sinceTime?: string
  recentSeconds?: number
  cursorUpdatedAt?: string
  cursorId?: number
}

export type CreateEngineeringOrderRequest = {
  name?: string
  leaderUserId?: number
  itemId?: number
  expectedQuantity?: number
  qualifiedQuantity?: number
  description?: string
  processId?: number
}

export type UpdateEngineeringOrderRequest = CreateEngineeringOrderRequest & {
  id?: number
}

export type ListEngineeringOrderRequest = {
  leaderUserId?: number
  itemId?: number
  namePrefix?: string
  itemNamePrefix?: string
  scope?: MesListScope
  pageNum?: number
  pageSize?: number
  processId?: number
  status?: DraftStatus
  sinceTime?: string
  recentSeconds?: number
  cursorUpdatedAt?: string
  cursorId?: number
}

export type CreateProcessDraftRequest = {
  ownerUserId?: number
  itemId?: number
  name?: string
  description?: string
  items?: ProcessItemRequest[]
}

export type UpdateProcessDraftRequest = CreateProcessDraftRequest & {
  id?: number
}

export type ListProcessRequest = {
  ownerUserId?: number
  itemId?: number
  namePrefix?: string
  itemNamePrefix?: string
  scope?: MesListScope
  status?: DraftStatus
  pageNum?: number
  pageSize?: number
  limit?: number
  sinceTime?: string
  recentSeconds?: number
  cursorUpdatedAt?: string
  cursorId?: number
}

export type CreateWorkOrderDraftRequest = {
  name?: string
  toUserId?: number
  description?: string
}

export type UpdateWorkOrderDraftRequest = CreateWorkOrderDraftRequest & {
  id?: number
}

export type ListWorkOrderRequest = {
  pageNum?: number
  pageSize?: number
  limit?: number
  id?: number
  namePrefix?: string
  scope?: MesListScope
  isTo?: boolean
  isUnread?: boolean
  status?: WorkOrderStatus
  sinceTime?: string
  recentSeconds?: number
  cursorUpdatedAt?: string
  cursorId?: number
}

export type AddItemRequest = {
  name?: string
  unit?: string
  description?: string
}

export type UpdateItemRequest = AddItemRequest & {
  id?: number
}

export type ListItemRequest = {
  pageNum?: number
  pageSize?: number
  limit?: number
  namePrefix?: string
  cursorUpdatedAt?: string
  cursorId?: number
}

export type SearchItemsRequest = ListItemRequest

export type ListItemUnitRequest = {
  pageNum?: number
  pageSize?: number
  limit?: number
  itemId?: number
  itemNamePrefix?: string
  scope?: MesListScope
  stockStatus?: StockStatus
  qualityStatus?: QualityStatus
  engineeringOrderId?: number
  inventoryFlowId?: number
  cursorUpdatedAt?: string
  cursorId?: number
}

export type AddItemUnitRequest = {
  itemId?: number
  stockStatus?: StockStatus
  qualityStatus?: QualityStatus
  description?: string
  engineeringOrderId?: number
}

export type UpdateItemUnitStatusRequest = {
  id?: number
  stockStatus?: StockStatus
  qualityStatus?: QualityStatus
}

export type MesInventoryFlow = InventoryFlowVO
export type MesItemUnit = ItemUnitVO

export const FLOW_TYPE_IN = FlowType.In
export const FLOW_TYPE_OUT = FlowType.Out
export const FLOW_STATUS_SUBMITTED = FlowStatus.Submitted
export const FLOW_STATUS_APPROVED = FlowStatus.Approved
export const STOCK_STATUS_IN_STOCK = StockStatus.InStock
export const STOCK_STATUS_RESERVED = StockStatus.Reserved
export const STOCK_STATUS_OUT_STOCK = StockStatus.OutStock
export const QUALITY_STATUS_PENDING = QualityStatus.Pending
export const QUALITY_STATUS_QUALIFIED = QualityStatus.Qualified
export const QUALITY_STATUS_UNQUALIFIED = QualityStatus.Unqualified

export async function createInventoryFlowDraft(body: CreateInventoryFlowDraftRequest) {
  return request<BaseResponse<number>>('/mes/inventory-flow/draft/create', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function updateInventoryFlowDraft(body: UpdateInventoryFlowDraftRequest) {
  return request<BaseResponse<boolean>>('/mes/inventory-flow/draft/update', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function deleteInventoryFlowDraft(body: { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/inventory-flow/draft/delete', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function submitInventoryFlow(id: number | { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/inventory-flow/submit', {
    method: 'POST',
    headers: jsonHeaders,
    data: typeof id === 'number' ? { id } : id,
  })
}

export async function completeInventoryFlow(body: CompleteInventoryFlowRequest) {
  return request<BaseResponse<boolean>>('/mes/inventory-flow/complete', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function auditInventoryFlow(
  id: number | { id?: number; approved?: boolean },
  approved?: boolean,
) {
  return request<BaseResponse<boolean>>('/mes/inventory-flow/audit', {
    method: 'POST',
    headers: jsonHeaders,
    data: typeof id === 'number' ? { id, approved } : id,
  })
}

export async function getInventoryFlow(params: { id: number }) {
  return request<BaseResponse<InventoryFlowVO>>('/mes/inventory-flow/get', {
    method: 'GET',
    params,
  })
}

export async function listInventoryFlow(body: ListInventoryFlowRequest) {
  return request<BaseResponse<PageResult<InventoryFlowVO>>>('/mes/inventory-flow/list', {
    method: 'POST',
    headers: jsonHeaders,
    data: withLimit(body),
  })
}

export async function createEngineeringOrder(body: CreateEngineeringOrderRequest) {
  return request<BaseResponse<number>>('/mes/engineering-order/draft/create', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function updateEngineeringOrder(body: UpdateEngineeringOrderRequest) {
  return request<BaseResponse<boolean>>('/mes/engineering-order/draft/update', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function deleteEngineeringOrderDraft(body: { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/engineering-order/draft/delete', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function submitEngineeringOrder(id: number | { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/engineering-order/submit', {
    method: 'POST',
    headers: jsonHeaders,
    data: typeof id === 'number' ? { id } : id,
  })
}

export async function getEngineeringOrder(params: { id: number }) {
  return request<BaseResponse<EngineeringOrderVO>>('/mes/engineering-order/get', {
    method: 'GET',
    params,
  })
}

export async function listEngineeringOrder(body: ListEngineeringOrderRequest) {
  return request<BaseResponse<PageResult<EngineeringOrderVO>>>('/mes/engineering-order/list', {
    method: 'POST',
    headers: jsonHeaders,
    data: withLimit(body),
  })
}

export async function createProcessDraft(body: CreateProcessDraftRequest) {
  return request<BaseResponse<number>>('/mes/process/draft/create', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function updateProcessDraft(body: UpdateProcessDraftRequest) {
  return request<BaseResponse<boolean>>('/mes/process/draft/update', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function deleteProcessDraft(body: { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/process/draft/delete', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function submitProcess(id: number | { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/process/submit', {
    method: 'POST',
    headers: jsonHeaders,
    data: typeof id === 'number' ? { id } : id,
  })
}

export async function getProcess(params: { id: number }) {
  return request<BaseResponse<ProcessVO>>('/mes/process/get', {
    method: 'GET',
    params,
  })
}

export async function listProcess(body: ListProcessRequest) {
  return request<BaseResponse<PageResult<ProcessVO>>>('/mes/process/list', {
    method: 'POST',
    headers: jsonHeaders,
    data: withLimit(body),
  })
}

export async function createWorkOrderDraft(body: CreateWorkOrderDraftRequest) {
  return request<BaseResponse<number>>('/mes/work-order/draft/create', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function updateWorkOrderDraft(body: UpdateWorkOrderDraftRequest) {
  return request<BaseResponse<boolean>>('/mes/work-order/draft/update', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function deleteWorkOrderDraft(body: { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/work-order/draft/delete', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function submitWorkOrder(id: number | { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/work-order/submit', {
    method: 'POST',
    headers: jsonHeaders,
    data: typeof id === 'number' ? { id } : id,
  })
}

export async function getWorkOrder(params: { id: number }) {
  return request<BaseResponse<WorkOrderVO>>('/mes/work-order/get', {
    method: 'GET',
    params,
  })
}

export async function listWorkOrder(body: ListWorkOrderRequest) {
  return request<BaseResponse<PageResult<WorkOrderVO>>>('/mes/work-order/list', {
    method: 'POST',
    headers: jsonHeaders,
    data: withLimit(body),
  })
}

export async function markWorkOrderRead(body: { id?: number }) {
  return request<BaseResponse<boolean>>('/mes/work-order/read', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function addItem(body: AddItemRequest) {
  return request<BaseResponse<number>>('/mes/item/add', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function updateItem(body: UpdateItemRequest) {
  return request<BaseResponse<boolean>>('/mes/item/update', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function getItem(params: { id: number }) {
  return request<BaseResponse<ItemVO>>('/mes/item/get', {
    method: 'GET',
    params,
  })
}

export async function searchItems(params: SearchItemsRequest) {
  const { limit, ...rest } = params
  return request<BaseResponse<PageResult<ItemVO>>>('/mes/item/search', {
    method: 'GET',
    params: {
      pageNum: rest.pageNum ?? 1,
      pageSize: rest.pageSize ?? limit ?? 20,
      ...rest,
    },
  })
}

export async function listItems(body: ListItemRequest) {
  return request<BaseResponse<PageResult<ItemVO>>>('/mes/item/list', {
    method: 'POST',
    headers: jsonHeaders,
    data: withLimit(body),
  })
}

export const listItem = listItems

export async function listItemUnit(body: ListItemUnitRequest) {
  return request<BaseResponse<PageResult<ItemUnitVO>>>('/mes/item-unit/list', {
    method: 'POST',
    headers: jsonHeaders,
    data: withLimit(body),
  })
}

export async function addItemUnit(body: AddItemUnitRequest) {
  return request<BaseResponse<number>>('/mes/item-unit/add', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}

export async function getItemUnit(params: { id: number }) {
  return request<BaseResponse<ItemUnitVO>>('/mes/item-unit/get', {
    method: 'GET',
    params,
  })
}

export async function updateItemUnitStatus(body: UpdateItemUnitStatusRequest) {
  return request<BaseResponse<boolean>>('/mes/item-unit/status/update', {
    method: 'POST',
    headers: jsonHeaders,
    data: body,
  })
}
