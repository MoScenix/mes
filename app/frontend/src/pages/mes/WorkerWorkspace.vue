<template>
  <div class="workspace-page">
    <div class="workspace-header">
      <MesListSearchPicker
        v-if="selectedType !== 'receive' && selectedType !== 'inspect'"
        v-model="searchText"
        placeholder="搜索码或物品名"
        class="search-input"
        @search="onSearch"
        @select-item="selectSearchItem"
        @clear="clearSearch"
      />
      <a-button v-if="selectedType === 'itemUnits' || selectedType === 'engineering'" type="primary" @click="openCreate">
        <PlusOutlined /> 新建
      </a-button>
    </div>

    <section v-if="selectedType === 'receive'" class="scan-workflow">
      <a-button v-if="receiveFlow" class="scan-back" type="text" @click="backToReceiveScan">返回</a-button>
      <div v-if="!receiveFlow" class="workflow-card">
        <strong>领取货物</strong>
      </div>
      <CodeTool
        v-if="!receiveFlow"
        v-model="receiveFlowCode"
        class="workflow-code-tool"
        kind="FLOW"
        auto-open
        scanner-only
        scanner-variant="bare"
        @submit="loadReceiveFlow"
      />
      <template v-else>
        <div class="workflow-card">
          <strong>扫描单体出库</strong>
        </div>
        <CodeTool
          :key="receiveOperationKey"
          v-model="receiveUnitCode"
          class="workflow-code-tool"
          kind="ITEM_UNIT"
          auto-open
          scanner-only
          scanner-variant="bare"
          scanner-display="inline"
          @submit="addReceiveUnit"
        />
      </template>
    </section>

    <section v-else-if="selectedType === 'inspect'" class="scan-workflow">
      <a-button v-if="inspectOrder" class="scan-back" type="text" @click="backToInspectScan">返回</a-button>
      <div v-if="!inspectOrder" class="workflow-card">
        <strong>检测单体</strong>
      </div>
      <CodeTool
        v-if="!inspectOrder"
        v-model="inspectOrderCode"
        kind="ENGINEERING_ORDER"
        auto-open
        scanner-only
        scanner-variant="bare"
        @submit="loadInspectOrder"
      />
      <template v-else>
        <div class="workflow-context">
          <strong>工程单 #{{ inspectOrder.id }}</strong>
          <MesItemName :id="inspectOrder.itemId" :item="inspectOrder.item" />
        </div>
        <CodeTool
          :key="inspectOperationKey"
          v-model="inspectUnitCode"
          kind="ITEM_UNIT"
          auto-open
          scanner-only
          scanner-variant="bare"
          scanner-display="inline"
          @submit="loadInspectUnit"
        />
        <div v-if="inspectUnit" class="inspect-result">
          <span>单体 #{{ inspectUnit.id }}</span>
          <a-space>
            <a-button :loading="inspectSubmitting" @click="submitInspect(QUALITY_STATUS_QUALIFIED)">合格</a-button>
            <a-button danger :loading="inspectSubmitting" @click="submitInspect(QUALITY_STATUS_UNQUALIFIED)">不合格</a-button>
          </a-space>
        </div>
      </template>
    </section>

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
        <template v-else-if="column.dataIndex === 'stockStatus'">
          <a-tag>{{ stockLabel(record.stockStatus) }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'qualityStatus'">
          <a-tag>{{ qualityLabel(record.qualityStatus) }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'flowStatus'">
          <a-tag :color="flowStatusColor(record.flowStatus)">{{ flowStatusLabel(record.flowStatus) }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'itemId'">
          <MesItemName :id="record.itemId" />
        </template>
        <template v-else-if="column.dataIndex === 'createTime' || column.dataIndex === 'updateTime'">
          {{ formatTime(record[column.dataIndex]) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space size="small">
            <a-button type="link" size="small" @click="viewDetail(record)">详情</a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    <div v-if="selectedType !== 'receive' && selectedType !== 'inspect' && dataList.length" class="list-more">
      <a-button v-if="listPage.hasMore" :loading="loadingMore" @click="loadMore">加载更多</a-button>
      <span v-else class="muted-text">没有更多了</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { Modal, message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import {
  FLOW_TYPE_OUT, FLOW_STATUS_APPROVED,
  STOCK_STATUS_IN_STOCK, STOCK_STATUS_RESERVED, STOCK_STATUS_OUT_STOCK,
  QUALITY_STATUS_PENDING, QUALITY_STATUS_QUALIFIED, QUALITY_STATUS_UNQUALIFIED,
  completeInventoryFlow, getEngineeringOrder, getInventoryFlow, getItemUnit,
  listItemUnit, listEngineeringOrder, updateItemUnitStatus,
  type EngineeringOrderVO,
  type InventoryFlowVO,
  type ItemUnitVO,
  type ItemVO,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'
import MesListSearchPicker from '@/components/mes/MesListSearchPicker.vue'
import MesItemName from '@/components/mes/MesItemName.vue'
import CodeTool from '@/components/mes/CodeTool.vue'

const router = useRouter()
const route = useRoute()

type DataType = 'itemUnits' | 'engineering' | 'receive' | 'inspect'
const panelFromRoute = () => {
  const panel = String(route.query.panel || 'itemUnits')
  return ['itemUnits', 'engineering', 'receive', 'inspect'].includes(panel) ? (panel as DataType) : 'itemUnits'
}
const selectedType = ref<DataType>(panelFromRoute())

const searchText = ref('')
const searchItemId = ref<number>()
const onTypeChange = () => { fetchData() }
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

const unitColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '物品', dataIndex: 'itemId', width: 160 },
  { title: '库存', dataIndex: 'stockStatus', width: 80 },
  { title: '质量', dataIndex: 'qualityStatus', width: 80 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '工程单', key: 'engineeringOrderId', width: 80, customRender: ({ record }: any) => record.engineeringOrderId ? `#${record.engineeringOrderId}` : '-' },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

const engColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '生产物品', key: 'itemName', width: 170, customRender: ({ record }: any) => `${record.item?.name || '物品'} #${record.itemId}` },
  { title: '预计', dataIndex: 'expectedQuantity', width: 80 },
  { title: '已产出', dataIndex: 'producedQuantity', width: 80 },
  { title: '合格', dataIndex: 'qualifiedQuantity', width: 80 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

const currentColumns = computed(() => selectedType.value === 'engineering' ? engColumns : unitColumns)

const dataList = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const listPage = reactive({ pageSize: 20, hasMore: false, nextCursorUpdatedAt: '', nextCursorId: 0 })

const syncCursor = (data?: { hasMore?: boolean; nextCursorUpdatedAt?: string; nextCursorId?: number }) => {
  listPage.hasMore = Boolean(data?.hasMore)
  listPage.nextCursorUpdatedAt = data?.nextCursorUpdatedAt || ''
  listPage.nextCursorId = data?.nextCursorId || 0
}

const fetchData = async (next = false) => {
  if (selectedType.value === 'receive' || selectedType.value === 'inspect') return
  if (next) loadingMore.value = true
  else loading.value = true
  try {
    const params = { pageSize: listPage.pageSize }
    const res = selectedType.value === 'engineering'
      ? await listEngineeringOrder({
          ...params,
          namePrefix: searchText.value.trim() || undefined,
          itemId: searchItemId.value,
          recentSeconds: 30 * 24 * 60 * 60,
          cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
          cursorId: next ? listPage.nextCursorId : undefined,
        })
      : await listItemUnit({
          ...params,
          itemId: searchItemId.value,
          engineeringOrderId: Number(route.query.engineeringOrderId || 0) || undefined,
          cursorId: next ? listPage.nextCursorId : undefined,
        })
    if (res.data.code === 0 && res.data.data) {
      dataList.value = next ? [...dataList.value, ...(res.data.data.records ?? [])] : (res.data.data.records ?? [])
      syncCursor(res.data.data)
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

const viewDetail = (record: any) => {
  const kind = selectedType.value === 'engineering' ? 'ENGINEERING_ORDER' : 'ITEM_UNIT'
  router.push({ path: '/mes/detail', query: { kind, id: String(record.id) } })
}

const formatTime = (t?: string) => t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-'
const stockLabel = (s?: number) => s === STOCK_STATUS_IN_STOCK ? '在库' : s === STOCK_STATUS_RESERVED ? '预留' : s === STOCK_STATUS_OUT_STOCK ? '出库' : '未知'
const qualityLabel = (s?: number) => s === QUALITY_STATUS_PENDING ? '待检测' : s === QUALITY_STATUS_QUALIFIED ? '合格' : s === QUALITY_STATUS_UNQUALIFIED ? '不合格' : '未知'
const flowStatusColor = (s?: number) => s === 1 ? 'default' : s === 2 ? 'blue' : s === 3 ? 'green' : 'red'
const flowStatusLabel = (s?: number) => s === 1 ? '草稿' : s === 2 ? '待处理' : s === 3 ? '已通过' : '已拒绝'

const openCreate = () => {
  router.push({ path: '/mes/create', query: { type: selectedType.value === 'engineering' ? 'engineering' : 'itemUnit' } })
}

const receiveFlowCode = ref('')
const receiveFlow = ref<InventoryFlowVO>()
const receiveUnitCode = ref('')
const receiveUnitIds = ref<number[]>([])
const receiveUnits = ref<ItemUnitVO[]>([])
const receiveSubmitting = ref(false)
const receiveOperationKey = ref(0)

const receiveExpectedQuantity = computed(() => {
  const units = receiveFlow.value?.itemUnits || []
  if (units.length) return units.length
  return (receiveFlow.value?.items || []).reduce((sum, item) => sum + (item.applyQuantity || 0), 0)
})

const receiveQuantityByItem = computed(() => {
  const result = new Map<number, number>()
  for (const unit of receiveUnits.value) {
    if (!unit.itemId) continue
    result.set(unit.itemId, (result.get(unit.itemId) || 0) + 1)
  }
  return result
})

const loadReceiveFlow = async (value: string) => {
  const parsed = parseMesCode(value, 'FLOW')
  if (!parsed.id) return
  const res = await getInventoryFlow({ id: parsed.id })
  if (res.data.code !== 0 || !res.data.data) {
    message.error(res.data.message || '读取流转单失败')
    return
  }
  if (res.data.data.flowType !== FLOW_TYPE_OUT) {
    message.error('只能领取出库流转单')
    return
  }
  if (res.data.data.flowStatus !== FLOW_STATUS_APPROVED) {
    message.error('流转单尚未审批通过，不能领取')
    return
  }
  receiveFlow.value = res.data.data
  clearReceiveUnits()
}

const resetReceive = () => {
  receiveFlowCode.value = ''
  receiveUnitCode.value = ''
  receiveFlow.value = undefined
  clearReceiveUnits()
}

const backToReceiveScan = async () => {
  await router.push({ path: '/mes/scan', query: { mode: 'receive' } })
}

const addReceiveUnit = async (value: string) => {
  const reopenScanner = () => {
    receiveUnitCode.value = ''
    receiveOperationKey.value += 1
  }
  const parsed = parseMesCode(value, 'ITEM_UNIT')
  if (!parsed.id) {
    reopenScanner()
    return
  }
  const flowUnits = receiveFlow.value?.itemUnits || []
  if (flowUnits.length && !flowUnits.some((unit) => unit.id === parsed.id)) {
    message.error('单体不属于当前流转单')
    reopenScanner()
    return
  }
  const unit = await getItemUnit({ id: parsed.id })
  if (unit.data.code !== 0 || !unit.data.data) {
    message.error(unit.data.message || '读取单体失败')
    reopenScanner()
    return
  }
  if (unit.data.data.stockStatus !== STOCK_STATUS_IN_STOCK) {
    message.error('只能领取在库单体')
    reopenScanner()
    return
  }
  if (unit.data.data.qualityStatus !== QUALITY_STATUS_QUALIFIED) {
    message.error('只能领取合格单体')
    reopenScanner()
    return
  }
  const flowItems = receiveFlow.value?.items || []
  const flowItem = flowItems.find((item) => item.itemId === unit.data.data?.itemId)
  if (!flowUnits.length && !flowItem) {
    message.error('该单体物品不在当前流转单明细中')
    reopenScanner()
    return
  }
  Modal.confirm({
    title: '确认领取',
    content: `确认领取单体 #${parsed.id}？`,
    okText: '确认',
    cancelText: '取消',
    async onOk() {
      receiveUnitCode.value = ''
      const res = await completeInventoryFlow({
        id: receiveFlow.value?.id,
        itemUnitIds: [parsed.id],
      })
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '领取失败')
      }
      message.success('领取已确认')
      receiveOperationKey.value += 1
    },
    onCancel: reopenScanner,
  })
}

const removeReceiveUnit = (id: number) => {
  receiveUnitIds.value = receiveUnitIds.value.filter((item) => item !== id)
  receiveUnits.value = receiveUnits.value.filter((item) => item.id !== id)
}

const clearReceiveUnits = () => {
  receiveUnitIds.value = []
  receiveUnits.value = []
  receiveUnitCode.value = ''
}

const submitReceive = async () => {
  if (!receiveUnitIds.value.length) {
    message.warning('请先扫描单体')
    return
  }
  if (receiveUnitIds.value.length !== receiveExpectedQuantity.value) {
    message.warning('已扫数量需要和流转单数量一致')
    return
  }
  receiveSubmitting.value = true
  try {
    for (const item of receiveFlow.value?.items || []) {
      if ((receiveQuantityByItem.value.get(item.itemId || 0) || 0) !== (item.applyQuantity || 0)) {
        throw new Error(`物品 #${item.itemId} 的扫码数量与申请数量不一致`)
      }
    }
    message.success('领取已确认')
    resetReceive()
  } catch (error) {
    message.error(error instanceof Error ? error.message : '提交失败')
  } finally {
    receiveSubmitting.value = false
  }
}

const inspectOrderCode = ref('')
const inspectOrder = ref<EngineeringOrderVO>()
const inspectUnitCode = ref('')
const inspectUnit = ref<ItemUnitVO>()
const inspectSubmitting = ref(false)
const inspectOperationKey = ref(0)

const loadInspectOrder = async (value: string) => {
  const parsed = parseMesCode(value, 'ENGINEERING_ORDER')
  if (!parsed.id) return
  const res = await getEngineeringOrder({ id: parsed.id })
  if (res.data.code !== 0 || !res.data.data) {
    message.error(res.data.message || '读取工程单失败')
    return
  }
  inspectOrder.value = res.data.data
}

const resetInspect = () => {
  inspectOrderCode.value = ''
  inspectUnitCode.value = ''
  inspectOrder.value = undefined
  inspectUnit.value = undefined
}

const backToInspectScan = async () => {
  await router.push({ path: '/mes/scan', query: { mode: 'inspect' } })
}

const loadInspectUnit = async (value: string) => {
  const reopenScanner = () => {
    inspectUnitCode.value = ''
    inspectOperationKey.value += 1
  }
  const parsed = parseMesCode(value, 'ITEM_UNIT')
  if (!parsed.id) {
    reopenScanner()
    return
  }
  const res = await getItemUnit({ id: parsed.id })
  if (res.data.code !== 0 || !res.data.data) {
    message.error(res.data.message || '读取单体失败')
    reopenScanner()
    return
  }
  if (inspectOrder.value?.id && res.data.data.engineeringOrderId !== inspectOrder.value.id) {
    message.error('单体不属于当前工程单')
    reopenScanner()
    return
  }
  if (res.data.data.qualityStatus !== QUALITY_STATUS_PENDING) {
    message.error('只能检测待检测单体')
    reopenScanner()
    return
  }
  inspectUnit.value = res.data.data
}

const submitInspect = async (qualityStatus: number) => {
  if (!inspectUnit.value?.id) return
  inspectSubmitting.value = true
  try {
    const res = await updateItemUnitStatus({
      id: inspectUnit.value.id,
      stockStatus: inspectUnit.value.stockStatus || STOCK_STATUS_OUT_STOCK,
      qualityStatus,
    })
    if (res.data.code !== 0) throw new Error(res.data.message || '更新失败')
    message.success(qualityStatus === QUALITY_STATUS_QUALIFIED ? '已标记合格' : '已标记不合格')
    inspectUnit.value = undefined
    inspectUnitCode.value = ''
    inspectOperationKey.value += 1
  } catch (error) {
    message.error(error instanceof Error ? error.message : '提交失败')
  } finally {
    inspectSubmitting.value = false
  }
}
watch(() => [route.query.panel, route.query.engineeringOrderId], async () => {
  selectedType.value = panelFromRoute()
  await hydrateWorkflowFromRoute()
  onTypeChange()
})

const hydrateWorkflowFromRoute = async () => {
  if (selectedType.value === 'receive') {
    const flowId = Number(route.query.flowId || 0)
    if (flowId > 0) {
      await loadReceiveFlow(`MES:FLOW:${flowId}`)
    }
  }
  if (selectedType.value === 'inspect') {
    const orderId = Number(route.query.orderId || 0)
    if (orderId > 0) {
      await loadInspectOrder(`MES:ENGINEERING_ORDER:${orderId}`)
    }
  }
}

onMounted(async () => {
  await hydrateWorkflowFromRoute()
  await fetchData()
})
</script>

<style scoped>
.workspace-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 16px; }
.search-input { width: 280px; max-width: 100%; }
.id-link { color: var(--primary); cursor: pointer; font-weight: 500; }
.id-link:hover { text-decoration: underline; }
:deep(.ant-table-wrapper) { border: 1px solid var(--border); border-radius: var(--radius); }
.list-more { display: flex; justify-content: center; padding-top: 14px; }
.muted-text { color: var(--muted-foreground, #94a3b8); font-size: 13px; }
.scan-workflow { position: relative; max-width: 560px; margin: 0 auto; display: grid; align-content: center; justify-items: center; gap: 0; min-height: min(620px, calc(100vh - 180px)); }
.scan-back { position: fixed; top: 72px; left: 96px; z-index: 20; }
.workflow-card { display: grid; justify-items: center; gap: 0; padding: 0; border: 0; border-radius: 0; background: transparent; box-shadow: none; }
.workflow-card > span,
.workflow-context span { color: var(--muted-foreground); font-size: 13px; }
.workflow-card strong { color: var(--foreground); font-size: 42px; line-height: 1.25; font-weight: 650; }
.workflow-card p { margin: 0; color: var(--muted-foreground); line-height: 1.6; }
.workflow-code-tool { width: 100%; }
.workflow-code-tool :deep(.ant-btn-primary),
.workflow-code-tool :deep(button[type='submit']) { min-height: 48px; font-size: 16px; font-weight: 600; }
.scan-list { min-height: 34px; display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }
.workflow-actions { display: grid; gap: 12px; padding-top: 6px; }
.workflow-actions :deep(.ant-space) { width: 100%; justify-content: flex-end; }
.inspect-result { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.muted { color: var(--muted-foreground); font-size: 13px; }
</style>
