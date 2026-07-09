<template>
  <div class="workspace-page">
    <!-- 顶部操作栏 -->
    <div class="workspace-header">
      <div class="header-left">
        <MesListSearchPicker
          v-if="selectedType !== 'scan'"
          v-model="searchText"
          placeholder="输入码或物品名搜索"
          class="search-input"
          @search="onSearch"
          @select-item="selectSearchItem"
          @clear="clearSearch"
        />
      </div>
      <div v-if="selectedType !== 'scan'" class="header-right">
        <div class="header-filters">
          <a-select
            v-if="selectedType === 'flows'"
            v-model:value="flowStatusFilter"
            allow-clear
            placeholder="流转状态"
            :options="flowStatusOptions"
            @change="fetchData()"
          />
          <a-select
            v-if="selectedType === 'itemUnits'"
            v-model:value="stockStatusFilter"
            allow-clear
            placeholder="库存状态"
            :options="stockFilterOptions"
            @change="fetchData()"
          />
          <a-select
            v-if="selectedType === 'itemUnits'"
            v-model:value="qualityStatusFilter"
            allow-clear
            placeholder="质量状态"
            :options="qualityFilterOptions"
            @change="fetchData()"
          />
        </div>
        <a-button v-if="selectedType !== 'items'" type="primary" @click="openCreate">
          <PlusOutlined /> {{ createButtonText }}
        </a-button>
      </div>
    </div>

    <section v-if="selectedType === 'scan'" class="scan-panel">
      <a-button v-if="scanFlow" class="scan-back" type="text" @click="backToInboundScan">返回</a-button>
      <div v-if="!scanFlow" class="mobile-scan-card">
        <strong>扫描入库</strong>
      </div>
      <CodeTool
        v-if="!scanFlow"
        v-model="scanFlowCode"
        class="mobile-code-tool"
        kind="FLOW"
        auto-open
        scanner-only
        scanner-variant="bare"
        @submit="loadScanFlow"
      />
      <template v-else>
        <div class="mobile-scan-card">
          <strong>扫描单体入库</strong>
        </div>
        <CodeTool
          :key="scanOperationKey"
          v-model="scanValue"
          class="mobile-code-tool"
          kind="ITEM_UNIT"
          auto-open
          scanner-only
          scanner-variant="bare"
          scanner-display="inline"
          @submit="addScanInput"
        />
      </template>
    </section>

    <!-- 数据表格 -->
    <a-table
      v-else
      row-key="id"
      :columns="currentColumns"
      :data-source="dataList"
      :pagination="false"
      :loading="loading"
      :scroll="{ x: 'max-content' }"
      size="middle"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'id'">
          <a class="id-link" @click="viewDetail(record)">#{{ record.id }}</a>
        </template>
        <template v-else-if="column.dataIndex === 'flowType'">
          <a-tag>{{ record.flowType === FLOW_TYPE_IN ? '入库' : record.flowType === FLOW_TYPE_OUT ? '出库' : '未知' }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'flowStatus'">
          <a-tag :color="flowStatusColor(record.flowStatus)">{{ flowStatusLabel(record.flowStatus) }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'stockStatus'">
          <a-tag>{{ stockLabel(record.stockStatus) }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'qualityStatus'">
          <a-tag>{{ qualityLabel(record.qualityStatus) }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'itemId'">
          <MesItemName :id="record.itemId" />
        </template>
        <template v-else-if="column.dataIndex === 'createTime' || column.dataIndex === 'updateTime'">
          {{ formatTime(record[column.dataIndex]) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space class="row-actions" size="small">
            <a-button type="link" size="small" @click="viewDetail(record)">详情</a-button>
            <a-button
              v-if="selectedType === 'flows' && record.flowStatus === 1"
              type="link"
              size="small"
              @click="editFlowDraft(record)"
            >
              编辑
            </a-button>
            <a-button
              v-if="selectedType === 'flows' && record.flowStatus === 1"
              type="link"
              size="small"
              @click="submitFlowDraft(record)"
            >
              提交
            </a-button>
            <a-popconfirm
              v-if="selectedType === 'flows' && record.flowStatus === 1"
              title="删除这个草稿？"
              @confirm="deleteFlowDraft(record)"
            >
              <a-button type="link" danger size="small">删除</a-button>
            </a-popconfirm>
            <a-button v-if="selectedType === 'itemUnits'" type="link" size="small" @click="editUnit(record)">编辑</a-button>
            <a-button v-if="selectedType === 'items' && record.id" type="link" size="small" @click="addUnitForItem(record)">添加单体</a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    <div v-if="selectedType !== 'scan' && dataList.length" class="list-more">
      <a-button v-if="listPage.hasMore" :loading="loadingMore" @click="loadMore">加载更多</a-button>
      <span v-else class="muted-text">没有更多了</span>
    </div>

    <!-- 旧弹窗仅保留编辑；新建统一走 /mes/create -->
    <a-modal
      v-if="false"
      v-model:open="createOpen"
      :title="createTitle"
      :confirm-loading="saving"
      width="520px"
      @ok="handleCreate"
      @cancel="createOpen = false"
    >
      <!-- 新建物品类型 -->
      <template v-if="selectedType === 'items'">
        <a-form layout="vertical" :model="itemForm">
          <a-form-item label="物品名称" required>
            <a-input v-model:value="itemForm.name" placeholder="如 M8 螺栓" />
          </a-form-item>
          <a-form-item label="计量单位">
            <a-input v-model:value="itemForm.unit" placeholder="个 / 件 / kg" />
          </a-form-item>
          <a-form-item label="说明">
            <a-textarea v-model:value="itemForm.description" :rows="3" placeholder="规格或备注" />
          </a-form-item>
        </a-form>
      </template>

      <!-- 新建库存单体 -->
      <template v-else-if="selectedType === 'itemUnits'">
        <a-form layout="vertical" :model="unitForm">
          <a-form-item label="物品" required>
            <a-select
              v-model:value="unitForm.itemId"
              show-search
              :filter-option="false"
              placeholder="搜索物品名或输入 ID"
              style="width: 100%"
              @search="searchItemsForSelect"
            >
              <a-select-option v-for="opt in itemSelectOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </a-select-option>
            </a-select>
          </a-form-item>
          <div class="form-row">
            <a-form-item label="库存状态">
              <a-select v-model:value="unitForm.stockStatus" :options="stockOptions" style="width: 100%" />
            </a-form-item>
            <a-form-item label="质量状态">
              <a-select v-model:value="unitForm.qualityStatus" :options="qualityOptions" style="width: 100%" />
            </a-form-item>
          </div>
          <a-form-item label="说明">
            <a-input v-model:value="unitForm.description" placeholder="批次、位置或备注" />
          </a-form-item>
        </a-form>
      </template>

      <!-- 新建流转单 -->
      <template v-else-if="selectedType === 'flows'">
        <a-form layout="vertical" :model="flowForm">
          <a-form-item label="流转方向">
            <a-segmented v-model:value="flowForm.flowType" :options="[{ label: '入库', value: FLOW_TYPE_IN }, { label: '出库', value: FLOW_TYPE_OUT }]" block />
          </a-form-item>
          <a-form-item label="接收人 ID">
            <a-input-number v-model:value="flowForm.toUserId" :min="1" style="width: 100%" placeholder="输入接收人用户 ID" />
          </a-form-item>
          <a-form-item label="说明">
            <a-textarea v-model:value="flowForm.description" :rows="3" placeholder="流转说明" />
          </a-form-item>
          <a-form-item label="物品列表">
            <div class="flow-items">
              <div v-for="(line, idx) in flowForm.items" :key="idx" class="flow-item-row">
                <a-input-number v-model:value="line.itemId" :min="1" placeholder="物品 ID" style="flex:1" />
                <a-input-number v-model:value="line.applyQuantity" :min="1" placeholder="数量" style="width:80px" />
                <a-button type="text" danger @click="flowForm.items.splice(idx, 1)">
                  <DeleteOutlined />
                </a-button>
              </div>
              <a-button type="dashed" block @click="flowForm.items.push({ itemId: undefined, applyQuantity: 1 })">
                <PlusOutlined /> 添加物品
              </a-button>
            </div>
          </a-form-item>
        </a-form>
      </template>
    </a-modal>

    <!-- 编辑库存单体 Modal -->
    <a-modal
      v-model:open="editOpen"
      title="编辑库存单体"
      :confirm-loading="editSaving"
      width="480px"
      @ok="handleEditUnit"
    >
      <a-form layout="vertical" :model="editUnitForm">
        <a-form-item label="库存状态">
          <a-select v-model:value="editUnitForm.stockStatus" :options="stockOptions" />
        </a-form-item>
        <a-form-item label="质量状态">
          <a-select v-model:value="editUnitForm.qualityStatus" :options="qualityOptions" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import {
  FLOW_TYPE_IN,
  FLOW_TYPE_OUT,
  FlowStatus,
  MesListScope,
  STOCK_STATUS_IN_STOCK,
  STOCK_STATUS_RESERVED,
  STOCK_STATUS_OUT_STOCK,
  QUALITY_STATUS_PENDING,
  QUALITY_STATUS_QUALIFIED,
  QUALITY_STATUS_UNQUALIFIED,
  listInventoryFlow,
  listItems,
  listItemUnit,
  addItem,
  addItemUnit,
  deleteInventoryFlowDraft,
  getInventoryFlow,
  getItemUnit,
  completeInventoryFlow,
  submitInventoryFlow,
  updateInventoryFlowDraft,
  updateItemUnitStatus,
  searchItems,
  type InventoryFlowVO,
  type ItemUnitVO,
  type ItemVO,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'
import CodeTool from '@/components/mes/CodeTool.vue'
import MesListSearchPicker from '@/components/mes/MesListSearchPicker.vue'
import MesItemName from '@/components/mes/MesItemName.vue'

const route = useRoute()
const router = useRouter()

// --- 数据类型 ---
type DataType = 'flows' | 'items' | 'itemUnits' | 'scan'

const panelFromRoute = () => {
  const panel = String(route.query.panel || 'items')
  return ['flows', 'items', 'itemUnits', 'scan'].includes(panel) ? (panel as DataType) : 'items'
}

const selectedType = ref<DataType>(panelFromRoute())

// --- 搜索 ---
const searchText = ref('')
const searchItemId = ref<number>()
const flowStatusFilter = ref<number>()
const stockStatusFilter = ref<number>()
const qualityStatusFilter = ref<number>()

const onTypeChange = () => {
  searchText.value = ''
  searchItemId.value = undefined
  flowStatusFilter.value = undefined
  stockStatusFilter.value = undefined
  qualityStatusFilter.value = undefined
  fetchData()
}

const syncPanel = () => {
  selectedType.value = panelFromRoute()
  onTypeChange()
}

const onSearch = (value: string) => {
  const parsed = parseMesCode(value)
  if (parsed.kind && parsed.id) {
    router.push({ path: '/mes/detail', query: { kind: parsed.kind, id: String(parsed.id) } })
    return
  }
  searchItemId.value = undefined
  fetchData()
}

const selectSearchItem = (item: ItemVO) => {
  searchText.value = item.name || ''
  searchItemId.value = item.id
  fetchData()
}

const clearSearch = () => {
  searchText.value = ''
  searchItemId.value = undefined
  fetchData()
}

// --- 列定义 ---
const flowColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '类型', dataIndex: 'flowType', width: 80 },
  { title: '状态', dataIndex: 'flowStatus', width: 80 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '单体', key: 'unitCount', width: 60, customRender: ({ record }: any) => record.itemUnits?.length || 0 },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 150 },
]

const flowStatusOptions = [
  { label: '草稿', value: FlowStatus.Draft },
  { label: '待处理', value: FlowStatus.Submitted },
  { label: '已通过', value: FlowStatus.Approved },
  { label: '已拒绝', value: FlowStatus.Rejected },
]

const itemColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 160 },
  { title: '单位', dataIndex: 'unit', width: 80 },
  { title: '库存', key: 'totalCount', width: 60, customRender: ({ record }: any) => record.totalCount ?? 0 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 140 },
]

const unitColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '物品', dataIndex: 'itemId', width: 160 },
  { title: '库存', dataIndex: 'stockStatus', width: 80 },
  { title: '质量', dataIndex: 'qualityStatus', width: 80 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '工程单', key: 'engineeringOrderId', width: 80, customRender: ({ record }: any) => record.engineeringOrderId ? `#${record.engineeringOrderId}` : '-' },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 120 },
]

const currentColumns = computed(() => {
  if (selectedType.value === 'items') return itemColumns
  if (selectedType.value === 'itemUnits') return unitColumns
  return flowColumns
})

// --- 数据 ---
const dataList = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const listPage = reactive({
  pageSize: 20,
  hasMore: false,
  nextCursorUpdatedAt: '',
  nextCursorName: '',
  nextCursorId: 0,
})

const syncCursor = (data?: { hasMore?: boolean; nextCursorUpdatedAt?: string; nextCursorName?: string; nextCursorId?: number }) => {
  listPage.hasMore = Boolean(data?.hasMore)
  listPage.nextCursorUpdatedAt = data?.nextCursorUpdatedAt || ''
  listPage.nextCursorName = data?.nextCursorName || ''
  listPage.nextCursorId = data?.nextCursorId || 0
}

const fetchData = async (next = false) => {
  if (selectedType.value === 'scan') return
  if (next) loadingMore.value = true
  else loading.value = true
  try {
    const type = selectedType.value
    const params = { pageSize: listPage.pageSize }
    let res: any

    if (type === 'flows') {
      res = await listInventoryFlow({
        ...params,
        itemNamePrefix: searchText.value.trim() || undefined,
        scope: MesListScope.Mine,
        flowStatus: flowStatusFilter.value,
        recentSeconds: 30 * 24 * 60 * 60,
        cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
        cursorId: next ? listPage.nextCursorId : undefined,
      })
    } else if (type === 'items') {
      const namePrefix = searchText.value.trim() || undefined
      res = await listItems({
        ...params,
        namePrefix,
        cursorName: next ? listPage.nextCursorName : undefined,
        cursorId: next ? listPage.nextCursorId : undefined,
      })
    } else {
      res = await listItemUnit({
        ...params,
        itemId: searchItemId.value,
        itemNamePrefix: searchItemId.value ? undefined : searchText.value.trim() || undefined,
        stockStatus: stockStatusFilter.value,
        qualityStatus: qualityStatusFilter.value,
        cursorId: next ? listPage.nextCursorId : undefined,
      })
    }

    if (res.data.code === 0 && res.data.data) {
      dataList.value = next ? [...dataList.value, ...(res.data.data.records ?? [])] : (res.data.data.records ?? [])
      syncCursor(res.data.data)
    } else {
      message.error(res.data.message || '获取数据失败')
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const loadMore = () => {
  if (!listPage.hasMore) return
  fetchData(true)
}

const scanFlowCode = ref('')
const scanFlow = ref<InventoryFlowVO>()
const scanValue = ref('')
const scannedUnitIds = ref<number[]>([])
const scannedUnits = ref<ItemUnitVO[]>([])
const scanSubmitting = ref(false)
const scanOperationKey = ref(0)

const scanExpectedQuantity = computed(() =>
  (scanFlow.value?.items || []).reduce((sum, item) => sum + (item.applyQuantity || 0), 0),
)

const scanQuantityByItem = computed(() => {
  const result = new Map<number, number>()
  for (const unit of scannedUnits.value) {
    if (!unit.itemId) continue
    result.set(unit.itemId, (result.get(unit.itemId) || 0) + 1)
  }
  return result
})

const loadScanFlowById = async (id: number) => {
  const res = await getInventoryFlow({ id })
  if (res.data.code !== 0 || !res.data.data) {
    message.error(res.data.message || '读取流转单失败')
    return
  }
  const flow = res.data.data
  if (flow.flowType !== FLOW_TYPE_IN) {
    message.error('只能进入入库流转单')
    return
  }
  if (flow.flowStatus !== FlowStatus.Approved) {
    message.error('只能录入已审批通过的入库流转单')
    return
  }
  scanFlow.value = flow
  clearScannedUnits()
}

const loadScanFlow = async (value: string) => {
  const parsed = parseMesCode(value, 'FLOW')
  if (!parsed.id) {
    message.warning('请输入有效的流转单码')
    return
  }
  await loadScanFlowById(parsed.id)
  if (scanFlow.value?.id) {
    await router.replace({ query: { ...route.query, panel: 'scan', flowId: String(scanFlow.value.id) } })
  }
}

const addScanInput = async (value: string) => {
  const reopenScanner = () => {
    scanValue.value = ''
    scanOperationKey.value += 1
  }
  if (!scanFlow.value) {
    message.warning('请先扫描入库流转单')
    reopenScanner()
    return
  }
  const parsed = parseMesCode(value, 'ITEM_UNIT')
  if (!parsed.id) {
    message.warning('请输入有效的库存单体码')
    reopenScanner()
    return
  }
  const unitRes = await getItemUnit({ id: parsed.id })
  const unit = unitRes.data.data
  if (unitRes.data.code !== 0 || !unit?.itemId) {
    message.error(unitRes.data.message || '读取库存单体失败')
    reopenScanner()
    return
  }
  if (unit.stockStatus !== STOCK_STATUS_OUT_STOCK) {
    message.error('只能录入不在库的单体')
    reopenScanner()
    return
  }
  if (unit.qualityStatus !== QUALITY_STATUS_QUALIFIED) {
    message.error('只能录入合格单体')
    reopenScanner()
    return
  }
  const flowItem = (scanFlow.value.items || []).find((item) => item.itemId === unit.itemId)
  if (!flowItem) {
    message.error('该单体物品不在当前流转单明细中')
    reopenScanner()
    return
  }
  Modal.confirm({
    title: '确认入库',
    content: `确认单体 #${parsed.id} 入库？`,
    okText: '确认',
    cancelText: '取消',
    async onOk() {
      scanValue.value = ''
      const res = await completeInventoryFlow({
        id: scanFlow.value?.id,
        itemUnitIds: [parsed.id],
      })
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '入库失败')
      }
      message.success('入库已确认')
      scanOperationKey.value += 1
    },
    onCancel: reopenScanner,
  })
}

const removeScanUnit = (id: number) => {
  scannedUnitIds.value = scannedUnitIds.value.filter((item) => item !== id)
  scannedUnits.value = scannedUnits.value.filter((item) => item.id !== id)
}

const clearScannedUnits = () => {
  scannedUnitIds.value = []
  scannedUnits.value = []
  scanValue.value = ''
}

const clearScan = async () => {
  scanFlow.value = undefined
  scanFlowCode.value = ''
  clearScannedUnits()
  const query = { ...route.query }
  delete query.flowId
  await router.replace({ query: { ...query, panel: 'scan' } })
}

const backToInboundScan = async () => {
  await router.push({ path: '/mes/scan', query: { mode: 'inbound' } })
}

const submitScanInbound = async () => {
  if (!scanFlow.value?.id) {
    message.warning('请先扫描入库流转单')
    return
  }
  if (!scannedUnitIds.value.length) {
    message.warning('请先扫描库存单体')
    return
  }
  if (scannedUnitIds.value.length !== scanExpectedQuantity.value) {
    message.warning('已扫数量需要和流转单申请数量一致')
    return
  }
  scanSubmitting.value = true
  try {
    for (const item of scanFlow.value.items || []) {
      if ((scanQuantityByItem.value.get(item.itemId || 0) || 0) !== (item.applyQuantity || 0)) {
        throw new Error(`物品 #${item.itemId} 的扫码数量与申请数量不一致`)
      }
    }
    const res = await completeInventoryFlow({
      id: scanFlow.value.id,
      itemUnitIds: scannedUnitIds.value,
    })
    if (res.data.code !== 0) {
      throw new Error(res.data.message || '提交入库失败')
    }
    message.success('入库已完成')
    await backToInboundScan()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '提交失败')
  } finally {
    scanSubmitting.value = false
  }
}

// --- 新建 ---
const createOpen = ref(false)
const saving = ref(false)

const createTitle = computed(() => {
  if (selectedType.value === 'items') return '新建物品类型'
  if (selectedType.value === 'itemUnits') return '新建库存单体'
  return '新建流转单'
})

const createButtonText = computed(() => {
  if (selectedType.value === 'items') return '新增物料'
  if (selectedType.value === 'itemUnits') return '新增单体'
  return '新建流转单'
})

const itemForm = reactive({ name: '', unit: '', description: '' })
const unitForm = reactive({ itemId: undefined as number | undefined, stockStatus: STOCK_STATUS_OUT_STOCK, qualityStatus: QUALITY_STATUS_QUALIFIED, description: '' })
const flowForm = reactive({ flowType: FLOW_TYPE_IN, toUserId: undefined as number | undefined, description: '', items: [] as { itemId?: number; applyQuantity?: number }[] })

const stockOptions = [
  { label: '在库', value: STOCK_STATUS_IN_STOCK },
  { label: '预留', value: STOCK_STATUS_RESERVED },
  { label: '不在库', value: STOCK_STATUS_OUT_STOCK },
]
const qualityOptions = [
  { label: '待检测', value: QUALITY_STATUS_PENDING },
  { label: '合格', value: QUALITY_STATUS_QUALIFIED },
  { label: '不合格', value: QUALITY_STATUS_UNQUALIFIED },
]
const stockFilterOptions = [...stockOptions]
const qualityFilterOptions = [...qualityOptions]

const itemSelectOptions = ref<{ label: string; value: number }[]>([])

const searchItemsForSelect = async (query: string) => {
  if (!query) { itemSelectOptions.value = []; return }
  const res = await searchItems({ namePrefix: query, limit: 10 })
  if (res.data.code === 0) {
    itemSelectOptions.value = (res.data.data?.records || []).map((i: ItemVO) => ({
      label: `${i.name} (#${i.id})`,
      value: i.id!,
    }))
  }
}

const openCreate = () => {
  const typeMap: Record<DataType, string> = {
    items: 'item',
    itemUnits: 'itemUnit',
    flows: 'flow',
    scan: 'flow',
  }
  const query: Record<string, string> = { type: typeMap[selectedType.value] }
  if (selectedType.value === 'itemUnits') {
    query.stockStatus = String(STOCK_STATUS_OUT_STOCK)
    query.qualityStatus = String(QUALITY_STATUS_QUALIFIED)
  }
  router.push({ path: '/mes/create', query })
}

const handleCreate = async () => {
  saving.value = true
  try {
    if (selectedType.value === 'items') {
      if (!itemForm.name.trim()) { message.warning('请输入物品名称'); return }
      const res = await addItem({ ...itemForm })
      if (res.data.code === 0) {
        message.success('物品类型已创建')
        createOpen.value = false
        fetchData()
      } else { message.error(res.data.message || '创建失败') }
    } else if (selectedType.value === 'itemUnits') {
      if (!unitForm.itemId) { message.warning('请选择物品'); return }
      const res = await addItemUnit({ ...unitForm })
      if (res.data.code === 0) {
        message.success('库存单体已创建')
        createOpen.value = false
        fetchData()
      } else { message.error(res.data.message || '创建失败') }
    }
  } finally {
    saving.value = false
  }
}

// --- 编辑库存单体 ---
const editOpen = ref(false)
const editSaving = ref(false)
const editUnitForm = reactive({ id: undefined as number | undefined, stockStatus: STOCK_STATUS_IN_STOCK, qualityStatus: QUALITY_STATUS_PENDING })

const editUnit = (record: any) => {
  editUnitForm.id = record.id
  editUnitForm.stockStatus = record.stockStatus ?? STOCK_STATUS_IN_STOCK
  editUnitForm.qualityStatus = record.qualityStatus ?? QUALITY_STATUS_PENDING
  editOpen.value = true
}

const handleEditUnit = async () => {
  if (!editUnitForm.id) return
  editSaving.value = true
  try {
    const res = await updateItemUnitStatus({ id: editUnitForm.id, stockStatus: editUnitForm.stockStatus, qualityStatus: editUnitForm.qualityStatus })
    if (res.data.code === 0) {
      message.success('已更新')
      editOpen.value = false
      fetchData()
    } else { message.error(res.data.message || '更新失败') }
  } finally {
    editSaving.value = false
  }
}

// --- 从物品类型快速跳转到添加单体 ---
const addUnitForItem = (item: ItemVO) => {
  router.push({
    path: '/mes/create',
    query: {
      type: 'itemUnit',
      itemId: String(item.id || ''),
      stockStatus: String(STOCK_STATUS_OUT_STOCK),
      qualityStatus: String(QUALITY_STATUS_QUALIFIED),
    },
  })
}

const editFlowDraft = (record: any) => {
  router.push({ path: '/mes/create', query: { type: 'flow', id: String(record.id) } })
}

const submitFlowDraft = async (record: any) => {
  if (!record.id) return
  const res = await submitInventoryFlow({ id: record.id })
  if (res.data.code === 0) {
    message.success('流转单已提交')
    fetchData()
  } else {
    message.error(res.data.message || '提交失败')
  }
}

const deleteFlowDraft = async (record: any) => {
  if (!record.id) return
  const res = await deleteInventoryFlowDraft({ id: record.id })
  if (res.data.code === 0) {
    message.success('草稿已删除')
    fetchData()
  } else {
    message.error(res.data.message || '删除失败')
  }
}

// --- 详情 ---
const viewDetail = (record: any) => {
  const kind = selectedType.value === 'flows' ? 'FLOW' : selectedType.value === 'items' ? 'ITEM' : 'ITEM_UNIT'
  router.push({ path: '/mes/detail', query: { kind, id: String(record.id) } })
}

// --- 工具函数 ---
const formatTime = (t?: string) => t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-'
const flowStatusColor = (s?: number) => {
  if (s === 1) return 'default'
  if (s === 2) return 'blue'
  if (s === 3) return 'green'
  if (s === 4) return 'red'
  return 'default'
}
const flowStatusLabel = (s?: number) => {
  if (s === 1) return '草稿'
  if (s === 2) return '待处理'
  if (s === 3) return '已通过'
  if (s === 4) return '已拒绝'
  return '未知'
}
const stockLabel = (s?: number) => s === STOCK_STATUS_IN_STOCK ? '在库' : s === STOCK_STATUS_RESERVED ? '预留' : s === STOCK_STATUS_OUT_STOCK ? '出库' : '未知'
const qualityLabel = (s?: number) => s === QUALITY_STATUS_PENDING ? '待检测' : s === QUALITY_STATUS_QUALIFIED ? '合格' : s === QUALITY_STATUS_UNQUALIFIED ? '不合格' : '未知'

// --- 初始化 ---
watch(() => route.query.panel, syncPanel)
onMounted(async () => {
  if (selectedType.value === 'scan') {
    const flowId = Number(route.query.flowId || 0)
    if (flowId > 0) {
      await loadScanFlowById(flowId)
      return
    }
  }
  await fetchData()
})
</script>

<style scoped>
.workspace-page {
  position: relative;
}

.workspace-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}

.header-filters {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-filters :deep(.ant-select) {
  min-width: 118px;
}

.scan-panel {
  position: relative;
  min-height: min(620px, calc(100vh - 180px));
  display: grid;
  align-content: center;
  justify-items: center;
  gap: 0;
  max-width: 520px;
  margin: 0 auto;
}

.scan-back {
  position: fixed;
  top: 72px;
  left: 96px;
  z-index: 20;
}

.mobile-scan-card {
  display: grid;
  justify-items: center;
  gap: 0;
  padding: 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.mobile-scan-card > span,
.scan-context span {
  color: var(--muted-foreground);
  font-size: 13px;
}

.mobile-scan-card strong {
  color: var(--foreground);
  font-size: 42px;
  line-height: 1.25;
  font-weight: 650;
}

.mobile-scan-card p {
  margin: 0;
  color: var(--muted-foreground);
  line-height: 1.6;
}

.mobile-code-tool {
  width: 100%;
}

.mobile-code-tool :deep(.ant-btn-primary),
.mobile-code-tool :deep(button[type='submit']) {
  min-height: 48px;
  font-size: 16px;
  font-weight: 600;
}

.scan-list {
  min-height: 34px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.scan-actions {
  display: grid;
  gap: 12px;
}

.scan-actions :deep(.ant-space) {
  width: 100%;
}

.scan-actions :deep(.ant-space-item) {
  flex: 1;
}

.scan-actions :deep(.ant-btn) {
  width: 100%;
  min-height: 44px;
}

.muted {
  color: var(--muted-foreground);
  font-size: 13px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.flow-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.flow-item-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.id-link {
  color: var(--primary);
  cursor: pointer;
  font-weight: 500;
}

.id-link:hover {
  text-decoration: underline;
}

:deep(.ant-table-wrapper) {
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.list-more {
  display: flex;
  justify-content: center;
  padding-top: 14px;
}

.muted-text { color: var(--muted-foreground, #94a3b8); font-size: 13px; }

.row-actions {
  white-space: nowrap;
}

.row-actions :deep(.ant-btn) {
  padding-inline: 2px;
}
</style>
