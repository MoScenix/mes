declare namespace API {
  type AIEventType =
    | 'accepted'
    | 'answer'
    | 'agent_start'
    | 'message'
    | 'push'
    | 'tool_call'
    | 'tool_result'
    | 'question'
    | 'done'
    | 'cancelled'
    | 'error'

  type AIStatusType =
    | 'queued'
    | 'running'
    | 'waiting_answer'
    | 'interrupted'
    | 'done'
    | 'cancelled'
    | 'error'

  type AIEvent = {
    id?: string
    historyId?: string
    type?: AIEventType
    agent?: string
    content?: string
    targetId?: string
    name?: string
    status?: string
    payloadJson?: string
    createdAt?: number
    questions?: AIQuestion[]
  }

  type AIPendingInterrupt = {
    id?: string
    agent?: string
    content?: string
    payloadJson?: string
  }

  type AIQuestion = {
    question?: string
    options?: string[]
  }

  type AIState = {
    exists?: boolean
    status?: AIStatusType
    agent?: string
    lastEventId?: string
    pendingInterrupts?: AIPendingInterrupt[]
    message?: string
  }

  type AIEvents = {
    events?: AIEvent[]
    lastId?: string
  }

  type AISubmitRequest = {
    historyId: number
    message?: string
  }

  type AIControlRequest = {
    historyId: number
    content?: string
    reason?: string
    answers?: Record<string, AIAnswer>
  }

  type AIAnswer = {
    content?: string
    payload?: Record<string, any>
  }

  type AIStateRequest = {
    historyId: number
  }

  type AIEventsRequest = {
    historyId: number
    lastId?: string
    blockMs?: number
    count?: number
  }

  type BaseResponseAIState = {
    code?: number
    data?: AIState
    message?: string
  }

  type BaseResponseAIEvents = {
    code?: number
    data?: AIEvents
    message?: string
  }

  type BaseResponseBoolean = {
    code?: number
    data?: boolean
    message?: string
  }

  type BaseResponseLoginUserVO = {
    code?: number
    data?: LoginUserVO
    message?: string
  }

  type BaseResponseLong = {
    code?: number
    data?: number
    message?: string
  }

  type BaseResponsePageHistoryMessage = {
    code?: number
    data?: PageHistoryMessage
    message?: string
  }

  type BaseResponsePageHistoryVO = {
    code?: number
    data?: PageHistoryVO
    message?: string
  }

  type BaseResponsePageUserVO = {
    code?: number
    data?: PageUserVO
    message?: string
  }

  type BaseResponseString = {
    code?: number
    data?: string
    message?: string
  }

  type BaseResponseUser = {
    code?: number
    data?: User
    message?: string
  }

  type BaseResponseUserVO = {
    code?: number
    data?: UserVO
    message?: string
  }

  type HistoryMessage = {
    id?: number
    message?: string
    messageType?: string
    historyId?: number
    userId?: number
    createTime?: string
    updateTime?: string
    isDelete?: number
    isFile?: boolean
  }

  type HistoryMessageQueryRequest = {
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    id?: number
    message?: string
    messageType?: string
    historyId?: number
    userId?: number
    lastCreateTime?: string
  }

  type DeleteRequest = {
    id?: number
  }

  type getUserByIdParams = {
    id: number
  }

  type getUserVOByIdParams = {
    id: number
  }

  type listHistoryMessagesParams = {
    historyId: number
    pageSize?: number
    lastCreateTime?: string
  }

  type HistoryVO = {
    id?: number
    historyName?: string
    userId?: number
    createTime?: string
    updateTime?: string
  }

  type PageHistoryVO = {
    records?: HistoryVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type LoginUserVO = {
    id?: number
    userAccount?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    createTime?: string
    updateTime?: string
  }

  type PageHistoryMessage = {
    records?: HistoryMessage[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type PageUserVO = {
    records?: UserVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
  }

  type ServerSentEventString = true

  type serveStaticResourceParams = {
    deployKey: string
  }

  type User = {
    id?: number
    userAccount?: string
    userPassword?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    editTime?: string
    createTime?: string
    updateTime?: string
    isDelete?: number
  }

  type UserAddRequest = {
    userName?: string
    userAccount?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
  }

  type UserLoginRequest = {
    userAccount?: string
    userPassword?: string
  }

  type UserQueryRequest = {
    pageNum?: number
    pageSize?: number
    sortField?: string
    sortOrder?: string
    id?: number
    userName?: string
    userAccount?: string
    userProfile?: string
    userRole?: string
  }

  type UserRegisterRequest = {
    userAccount?: string
    userPassword?: string
    checkPassword?: string
  }

  type UserUpdateRequest = {
    id?: number
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
  }

  type UserVO = {
    id?: number
    userAccount?: string
    userName?: string
    userAvatar?: string
    userProfile?: string
    userRole?: string
    createTime?: string
  }

  type WorkOrderStatus = 0 | 1 | 2
  type WorkOrderReadStatus = 0 | 1 | 2
  type FlowType = 0 | 1 | 2
  type FlowStatus = 0 | 1 | 2 | 3 | 4
  type StockStatus = 0 | 1 | 2 | 3
  type QualityStatus = 0 | 1 | 2 | 3
  type DraftStatus = 0 | 1 | 2 | 3
  type MesListScope = 0 | 1 | 2 | 3 | 4

  type GetByIdRequest = {
    id?: number
  }

  type MESListRequest = {
    pageNum?: number
    pageSize?: number
    limit?: number
    sinceTime?: string
    recentSeconds?: number
    cursorUpdatedAt?: string
    cursorId?: number
  }

  type CreateWorkOrderDraftRequest = {
    name?: string
    toUserId?: number
    description?: string
  }

  type UpdateWorkOrderDraftRequest = CreateWorkOrderDraftRequest & {
    id?: number
  }

  type ListWorkOrderRequest = MESListRequest & {
    id?: number
    namePrefix?: string
    scope?: MesListScope
    isTo?: boolean
    isUnread?: boolean
  }

  type WorkOrderVO = {
    id?: number
    name?: string
    fromUserId?: number
    toUserId?: number
    description?: string
    status?: WorkOrderStatus
    createTime?: string
    updateTime?: string
    readStatus?: WorkOrderReadStatus
  }

  type PageWorkOrderVO = {
    records?: WorkOrderVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
    hasMore?: boolean
    nextCursorUpdatedAt?: string
    nextCursorId?: number
  }

  type AddItemRequest = {
    name?: string
    unit?: string
    description?: string
  }

  type UpdateItemRequest = AddItemRequest & {
    id?: number
  }

  type ListItemRequest = MESListRequest & {
    namePrefix?: string
    cursorUpdatedAt?: string
  }

  type SearchItemsRequest = MESListRequest & {
    namePrefix?: string
    cursorUpdatedAt?: string
  }

  type ItemVO = {
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

  type PageItemVO = {
    records?: ItemVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
    hasMore?: boolean
    nextCursorUpdatedAt?: string
    nextCursorId?: number
  }

  type AddItemUnitRequest = {
    itemId?: number
    stockStatus?: StockStatus
    qualityStatus?: QualityStatus
    description?: string
    engineeringOrderId?: number
  }

  type UpdateItemUnitStatusRequest = {
    id?: number
    stockStatus?: StockStatus
    qualityStatus?: QualityStatus
  }

  type ListItemUnitRequest = MESListRequest & {
    itemId?: number
    itemNamePrefix?: string
    scope?: MesListScope
    stockStatus?: StockStatus
    qualityStatus?: QualityStatus
    engineeringOrderId?: number
    inventoryFlowId?: number
    cursorUpdatedAt?: string
  }

  type ItemUnitVO = {
    id?: number
    itemId?: number
    stockStatus?: StockStatus
    qualityStatus?: QualityStatus
    description?: string
    createTime?: string
    updateTime?: string
    engineeringOrderId?: number
  }

  type PageItemUnitVO = {
    records?: ItemUnitVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
    hasMore?: boolean
    nextCursorUpdatedAt?: string
    nextCursorId?: number
  }

  type ProcessItemRequest = {
    consumeItemId?: number
    quantity?: number
  }

  type ProcessItemVO = {
    id?: number
    processId?: number
    consumeItemId?: number
    quantity?: number
    consumeItem?: ItemVO
  }

  type ProcessVO = {
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

  type CreateProcessDraftRequest = {
    ownerUserId?: number
    itemId?: number
    name?: string
    description?: string
    items?: ProcessItemRequest[]
  }

  type UpdateProcessDraftRequest = CreateProcessDraftRequest & {
    id?: number
  }

  type ListProcessRequest = MESListRequest & {
    ownerUserId?: number
    itemId?: number
    namePrefix?: string
    itemNamePrefix?: string
    scope?: MesListScope
    status?: DraftStatus
  }

  type PageProcessVO = {
    records?: ProcessVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
    hasMore?: boolean
    nextCursorUpdatedAt?: string
    nextCursorId?: number
  }

  type CreateEngineeringOrderRequest = {
    name?: string
    leaderUserId?: number
    itemId?: number
    expectedQuantity?: number
    qualifiedQuantity?: number
    description?: string
    processId?: number
  }

  type UpdateEngineeringOrderRequest = CreateEngineeringOrderRequest & {
    id?: number
  }

  type ListEngineeringOrderRequest = MESListRequest & {
    leaderUserId?: number
    itemId?: number
    namePrefix?: string
    itemNamePrefix?: string
    scope?: MesListScope
    processId?: number
    status?: DraftStatus
  }

  type EngineeringOrderVO = {
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

  type PageEngineeringOrderVO = {
    records?: EngineeringOrderVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
    hasMore?: boolean
    nextCursorUpdatedAt?: string
    nextCursorId?: number
  }

  type InventoryFlowItemRequest = {
    itemId?: number
    applyQuantity?: number
  }

  type CreateInventoryFlowDraftRequest = {
    name?: string
    toUserId?: number
    flowType?: FlowType
    description?: string
    items?: InventoryFlowItemRequest[]
    itemUnitIds?: number[]
  }

  type UpdateInventoryFlowDraftRequest = CreateInventoryFlowDraftRequest & {
    id?: number
  }

  type AuditInventoryFlowRequest = {
    id?: number
    approved?: boolean
  }

  type ListInventoryFlowRequest = MESListRequest & {
    userId?: number
    isTo?: boolean
    flowStatus?: FlowStatus
    namePrefix?: string
    itemNamePrefix?: string
    scope?: MesListScope
    itemUnitId?: number
  }

  type InventoryFlowItemVO = {
    id?: number
    inventoryFlowId?: number
    itemId?: number
    applyQuantity?: number
    finishedQuantity?: number
    item?: ItemVO
  }

  type InventoryFlowVO = {
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

  type PageInventoryFlowVO = {
    records?: InventoryFlowVO[]
    pageNumber?: number
    pageSize?: number
    totalPage?: number
    totalRow?: number
    optimizeCountQuery?: boolean
    hasMore?: boolean
    nextCursorUpdatedAt?: string
    nextCursorId?: number
  }

  type BaseResponseWorkOrderVO = {
    code?: number
    data?: WorkOrderVO
    message?: string
  }

  type BaseResponsePageWorkOrderVO = {
    code?: number
    data?: PageWorkOrderVO
    message?: string
  }

  type BaseResponseItemVO = {
    code?: number
    data?: ItemVO
    message?: string
  }

  type BaseResponsePageItemVO = {
    code?: number
    data?: PageItemVO
    message?: string
  }

  type BaseResponseItemUnitVO = {
    code?: number
    data?: ItemUnitVO
    message?: string
  }

  type BaseResponsePageItemUnitVO = {
    code?: number
    data?: PageItemUnitVO
    message?: string
  }

  type BaseResponseProcessVO = {
    code?: number
    data?: ProcessVO
    message?: string
  }

  type BaseResponsePageProcessVO = {
    code?: number
    data?: PageProcessVO
    message?: string
  }

  type BaseResponseEngineeringOrderVO = {
    code?: number
    data?: EngineeringOrderVO
    message?: string
  }

  type BaseResponsePageEngineeringOrderVO = {
    code?: number
    data?: PageEngineeringOrderVO
    message?: string
  }

  type BaseResponseInventoryFlowVO = {
    code?: number
    data?: InventoryFlowVO
    message?: string
  }

  type BaseResponsePageInventoryFlowVO = {
    code?: number
    data?: PageInventoryFlowVO
    message?: string
  }
}
