<template>
  <main class="detail-page">
    <section class="detail-head">
      <div>
        <p>详情</p>
        <h1>{{ title }}</h1>
      </div>
    </section>

    <section class="detail-surface">
      <a-spin :spinning="loading">
        <a-result v-if="error" status="warning" title="没有读取到对象" :sub-title="error">
          <template #extra>
            <a-button type="primary" @click="loadDetail">重试</a-button>
          </template>
        </a-result>

        <template v-else>
          <div v-if="codeValue" class="code-strip">
            <MesQrCode :value="codeValue" :size="132" />
            <div class="code-text">
              <span>对象码</span>
              <button type="button" @click="copyCode">{{ codeValue }}</button>
            </div>
          </div>

          <section class="mail-body">
            <div class="mail-meta">
              <span>{{ titleByKind[kind!] || '对象' }}</span>
              <span>#{{ id }}</span>
              <span>{{ updateTimeText }}</span>
            </div>
            <p>{{ bodyText }}</p>
          </section>

          <dl class="detail-grid">
            <div v-for="row in rows" :key="row.label" :class="{ 'full-width': row.fullWidth }">
              <dt>{{ row.label }}</dt>
              <dd>
                <MesUserName v-if="row.userId" :id="row.userId" />
                <MesItemName v-else-if="row.itemId || row.item" :id="row.itemId" :item="row.item" />
                <button
                  v-else-if="row.action"
                  type="button"
                  class="detail-link"
                  @click="runRowAction(row.action)"
                >
                  跳转
                </button>
                <span v-else>{{ row.value }}</span>
              </dd>
            </div>
          </dl>

          <section v-if="kind === 'FLOW'" class="flow-progress-panel">
            <div class="trace-head">
              <h2>物品进度</h2>
              <span>{{ flowProgressRows.length }} 个物品</span>
            </div>
            <a-table
              row-key="key"
              :columns="flowProgressColumns"
              :data-source="flowProgressRows"
              :pagination="false"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'item'">
                  <MesItemName :id="record.itemId" :item="record.item" />
                </template>
                <template v-else-if="column.key === 'progress'">
                  {{ record.finishedQuantity }}/{{ record.applyQuantity }}
                </template>
              </template>
            </a-table>
          </section>

          <section v-if="kind === 'ITEM_UNIT'" class="trace-panel">
            <div class="trace-head">
              <h2>单体追踪</h2>
              <span v-if="traceLoading">读取中...</span>
              <span v-else>{{ traceItems.length }} 条记录</span>
            </div>
            <a-empty v-if="!traceLoading && !traceItems.length" description="暂无追踪记录" />
            <div v-else class="trace-list">
              <button
                v-for="item in traceItems"
                :key="item.key"
                class="trace-item"
                :class="{ clickable: item.kind && item.id }"
                type="button"
                @click="openTraceItem(item)"
              >
                <span class="trace-dot"></span>
                <span class="trace-time">{{ item.time || '-' }}</span>
                <span class="trace-main">
                  <strong>{{ item.title }}</strong>
                  <small>{{ item.description }}</small>
                </span>
              </button>
            </div>
          </section>
        </template>
      </a-spin>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import MesQrCode from '@/components/mes/MesQrCode.vue'
import {
  DraftStatus,
  FlowStatus,
  FlowType,
  QualityStatus,
  StockStatus,
  WorkOrderStatus,
  getEngineeringOrder,
  getInventoryFlow,
  getItem,
  getItemUnit,
  listInventoryFlow,
  getProcess,
  getWorkOrder,
  type EngineeringOrderVO,
  type InventoryFlowVO,
  type ItemVO,
  type ItemUnitVO,
  MesListScope,
  type ProcessVO,
  type WorkOrderVO,
} from '@/api/mesController'
import MesItemName from '@/components/mes/MesItemName.vue'
import MesUserName from '@/components/mes/MesUserName.vue'
import { makeMesCode, type MesCodeKind, type MesDetailKind } from '@/utils/mesCode'

type DetailRow = {
  label: string
  value?: string
  userId?: number
  itemId?: number
  item?: ItemVO
  action?: 'processOrders' | 'relatedUnits'
  fullWidth?: boolean
}

type TraceItem = {
  key: string
  title: string
  description: string
  time?: string
  kind?: MesDetailKind
  id?: number
}

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const error = ref('')
const detail = ref<
  InventoryFlowVO | ItemUnitVO | EngineeringOrderVO | WorkOrderVO | ItemVO | ProcessVO
>()
const traceLoading = ref(false)
const traceOrder = ref<EngineeringOrderVO>()
const traceFlows = ref<InventoryFlowVO[]>([])
const detailKinds: MesDetailKind[] = [
  'FLOW',
  'ITEM_UNIT',
  'ENGINEERING_ORDER',
  'WORK_ORDER',
  'ITEM',
  'PROCESS',
]
const scannableKinds: MesCodeKind[] = ['FLOW', 'ITEM_UNIT', 'ENGINEERING_ORDER']

const kind = computed(() => {
  const value = String(route.query.kind || '').toUpperCase() as MesDetailKind
  return detailKinds.includes(value) ? value : undefined
})
const id = computed(() => Number(route.query.id || 0))
const codeKind = computed(() => {
  const value = kind.value as MesCodeKind
  return scannableKinds.includes(value) ? value : undefined
})

const titleByKind: Record<MesDetailKind, string> = {
  FLOW: '流转单',
  ITEM_UNIT: '库存单体',
  ENGINEERING_ORDER: '工程单',
  WORK_ORDER: '工单',
  ITEM: '物料',
  PROCESS: '工艺单',
}

const titleName = computed(() => {
  const current = detail.value as
    | (InventoryFlowVO | EngineeringOrderVO | WorkOrderVO | ItemVO | ProcessVO)
    | undefined
  return current && 'name' in current ? current.name : ''
})
const title = computed(() => {
  if (!kind.value || !id.value) return 'MES 对象'
  return titleName.value
    ? `${titleName.value} #${id.value}`
    : `${titleByKind[kind.value]} #${id.value}`
})
const codeValue = computed(() =>
  codeKind.value && id.value ? makeMesCode(codeKind.value, id.value) : '',
)

const copyCode = async () => {
  if (!codeValue.value) return
  await navigator.clipboard?.writeText(codeValue.value)
  message.success('已复制')
}

const stockText = (status?: StockStatus) => {
  if (status === StockStatus.InStock) return '在库'
  if (status === StockStatus.OutStock) return '不在库'
  return '未知'
}

const qualityText = (status?: QualityStatus) => {
  if (status === QualityStatus.Pending) return '待检测'
  if (status === QualityStatus.Qualified) return '合格'
  if (status === QualityStatus.Unqualified) return '不合格'
  return '未知'
}

const flowText = (type?: FlowType) => {
  if (type === FlowType.In) return '入库'
  if (type === FlowType.Out) return '出库'
  return '未知'
}

const flowStatusText = (status?: FlowStatus) => {
  if (status === FlowStatus.Draft) return '草稿'
  if (status === FlowStatus.Submitted) return '已提交'
  if (status === FlowStatus.Approved) return '已通过'
  if (status === FlowStatus.Rejected) return '已拒绝'
  return '未知'
}

const workOrderStatusText = (status?: WorkOrderStatus) => {
  if (status === WorkOrderStatus.Draft) return '草稿'
  if (status === WorkOrderStatus.Submitted) return '已提交'
  return '未知'
}

const processStatusText = (status?: DraftStatus) => {
  if (status === DraftStatus.Draft) return '草稿'
  if (status === DraftStatus.Submitted) return '已提交'
  if (status === DraftStatus.Done) return '已完成'
  return '未知'
}

const updateTimeText = computed(() => {
  const current = detail.value as { updateTime?: string } | undefined
  return current?.updateTime || '-'
})

const bodyText = computed(() => {
  const current = detail.value as { description?: string } | undefined
  return current?.description?.trim() || '暂无正文内容。'
})

const openTraceItem = (item: TraceItem) => {
  if (!item.kind || !item.id) return
  router.push({ path: '/mes/detail', query: { kind: item.kind, id: String(item.id) } })
}

const traceItems = computed<TraceItem[]>(() => {
  if (kind.value !== 'ITEM_UNIT') return []
  const unit = detail.value as ItemUnitVO | undefined
  if (!unit) return []
  const items: TraceItem[] = []
  if (unit.createTime) {
    items.push({
      key: 'unit-created',
      title: '单体创建',
      description: `库存状态：${stockText(unit.stockStatus)}，质量状态：${qualityText(unit.qualityStatus)}`,
      time: unit.createTime,
    })
  }
  if (traceOrder.value?.id) {
    const order = traceOrder.value
    items.push({
      key: `engineering-${order.id}`,
      title: `关联工程单 #${order.id}`,
      description: order.name || order.description || '工程单',
      time: order.updateTime || order.createTime,
      kind: 'ENGINEERING_ORDER',
      id: order.id,
    })
  }
  for (const flow of traceFlows.value) {
    if (!flow.id) continue
    items.push({
      key: `flow-${flow.id}`,
      title: `${flowText(flow.flowType)}流转单 #${flow.id}`,
      description: `${flowStatusText(flow.flowStatus)} · ${flow.name || flow.description || '流转单'}`,
      time: flow.approvedAt || flow.updateTime || flow.createTime,
      kind: 'FLOW',
      id: flow.id,
    })
  }
  return items.sort((a, b) => (b.time || '').localeCompare(a.time || ''))
})

const openProcessOrderList = () => {
  router.push({ path: '/mes/process-eng-orders', query: { processId: String(id.value) } })
}

const openUnitList = () => {
  if (kind.value === 'ENGINEERING_ORDER') {
    router.push({
      path: '/mes/purchase',
      query: { panel: 'itemUnits', view: 'units', engineeringOrderId: String(id.value) },
    })
    return
  }
  if (kind.value === 'FLOW') {
    router.push({
      path: '/mes/purchase',
      query: { panel: 'itemUnits', view: 'units', flowId: String(id.value) },
    })
  }
}

const runRowAction = (action: DetailRow['action']) => {
  if (action === 'processOrders') {
    openProcessOrderList()
    return
  }
  if (action === 'relatedUnits') {
    openUnitList()
  }
}

const loadItemUnitTrace = async (unit: ItemUnitVO) => {
  traceOrder.value = undefined
  traceFlows.value = []
  if (!unit.id) return
  traceLoading.value = true
  try {
    const [orderRes, flowRes] = await Promise.all([
      unit.engineeringOrderId
        ? getEngineeringOrder({ id: unit.engineeringOrderId })
        : Promise.resolve(undefined),
      listInventoryFlow({
        itemUnitId: unit.id,
        scope: MesListScope.All,
        pageSize: 50,
      }),
    ])
    if (orderRes?.data.code === 0) {
      traceOrder.value = orderRes.data.data
    }
    if (flowRes.data.code === 0) {
      traceFlows.value = flowRes.data.data?.records || []
    }
  } catch (err) {
    message.error(err instanceof Error ? err.message : '读取单体追踪失败')
  } finally {
    traceLoading.value = false
  }
}

const rows = computed<DetailRow[]>(() => {
  const current = detail.value
  if (!current) return []

  if (kind.value === 'FLOW') {
    const flow = current as InventoryFlowVO
    return [
      { label: '编号', value: String(flow.id ?? '-') },
      { label: '名称', value: flow.name || '-' },
      { label: '方向', value: flowText(flow.flowType) },
      { label: '状态', value: flowStatusText(flow.flowStatus) },
      { label: '提交人', userId: flow.fromUserId },
      { label: '接收人', userId: flow.toUserId },
      { label: '关联单体', action: 'relatedUnits' },
      { label: '更新时间', value: flow.updateTime || '-' },
    ]
  }

  if (kind.value === 'ITEM_UNIT') {
    const unit = current as ItemUnitVO
    return [
      { label: '编号', value: String(unit.id ?? '-') },
      { label: '物品', itemId: unit.itemId },
      { label: '库存状态', value: stockText(unit.stockStatus) },
      { label: '质量状态', value: qualityText(unit.qualityStatus) },
      { label: '工程单', value: unit.engineeringOrderId ? `#${unit.engineeringOrderId}` : '-' },
      { label: '更新时间', value: unit.updateTime || '-' },
    ]
  }

  if (kind.value === 'ENGINEERING_ORDER') {
    const order = current as EngineeringOrderVO
    return [
      { label: '编号', value: String(order.id ?? '-') },
      { label: '名称', value: order.name || '-' },
      { label: '状态', value: processStatusText(order.status) },
      { label: '生产物品', itemId: order.itemId, item: order.item },
      { label: '预计数量', value: String(order.expectedQuantity ?? 0) },
      { label: '合格数量', value: String(order.qualifiedQuantity ?? 0) },
      { label: '已产出', value: String(order.producedQuantity ?? 0) },
      { label: '关联单体', action: 'relatedUnits' },
      { label: '更新时间', value: order.updateTime || '-' },
    ]
  }

  if (kind.value === 'ITEM') {
    const item = current as ItemVO
    return [
      { label: '编号', value: String(item.id ?? '-') },
      { label: '物品', itemId: item.id, item },
      { label: '单位', value: item.unit || '-' },
      { label: '总库存', value: String(item.totalCount ?? 0) },
      { label: '在库', value: String(item.inStockCount ?? 0) },
      { label: '可用', value: String(item.availableCount ?? 0) },
      { label: '更新时间', value: item.updateTime || '-' },
    ]
  }

  if (kind.value === 'PROCESS') {
    const process = current as ProcessVO
    return [
      { label: '编号', value: String(process.id ?? '-') },
      { label: '名称', value: process.name || '-' },
      { label: '关联物品', itemId: process.itemId, item: process.item },
      { label: '产出物品', value: process.item?.name || '-' },
      { label: '状态', value: processStatusText(process.status) },
      { label: '关联工程单', action: 'processOrders' },
      { label: '更新时间', value: process.updateTime || '-' },
    ]
  }

  const workOrder = current as WorkOrderVO
  return [
    { label: '编号', value: String(workOrder.id ?? '-') },
    { label: '名称', value: workOrder.name || '-' },
    { label: '发起人', userId: workOrder.fromUserId },
    { label: '接收人', userId: workOrder.toUserId },
    { label: '状态', value: workOrderStatusText(workOrder.status) },
    { label: '更新时间', value: workOrder.updateTime || '-' },
  ]
})

const flowProgressColumns = [
  { title: '物品', key: 'item' },
  { title: '已操作', dataIndex: 'finishedQuantity', width: 100 },
  { title: '申请数量', dataIndex: 'applyQuantity', width: 100 },
  { title: '进度', key: 'progress', width: 100 },
]

const flowProgressRows = computed(() => {
  if (kind.value !== 'FLOW') return []
  const flow = detail.value as InventoryFlowVO | undefined
  return (flow?.items || []).map((item, index) => ({
    key: `${item.itemId || 'item'}-${index}`,
    itemId: item.itemId,
    item: item.item,
    finishedQuantity: item.finishedQuantity || 0,
    applyQuantity: item.applyQuantity || 0,
  }))
})

const loadDetail = async () => {
  error.value = ''
  detail.value = undefined
  traceOrder.value = undefined
  traceFlows.value = []
  if (!kind.value || !id.value) {
    error.value = '缺少对象类型或 ID'
    return
  }

  loading.value = true
  try {
    const res =
      kind.value === 'FLOW'
        ? await getInventoryFlow({ id: id.value })
        : kind.value === 'ITEM_UNIT'
          ? await getItemUnit({ id: id.value })
          : kind.value === 'ENGINEERING_ORDER'
            ? await getEngineeringOrder({ id: id.value })
            : kind.value === 'ITEM'
              ? await getItem({ id: id.value })
              : kind.value === 'PROCESS'
                ? await getProcess({ id: id.value })
                : await getWorkOrder({ id: id.value })

    if (res.data.code !== 0 || !res.data.data) {
      error.value = res.data.message || '对象不存在或无权查看'
      return
    }
    detail.value = res.data.data
    if (kind.value === 'ITEM_UNIT') {
      await loadItemUnitTrace(res.data.data as ItemUnitVO)
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : '读取失败'
  } finally {
    loading.value = false
  }
}

watch(() => route.query, loadDetail)
onMounted(loadDetail)
</script>

<style scoped>
.detail-page {
  display: grid;
  gap: 18px;
}

.detail-head {
  min-height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.detail-head p {
  margin: 0 0 4px;
  color: #7a7a7a;
  font-size: 12px;
}

.detail-head h1 {
  margin: 0;
  color: #1d1d1f;
  font-size: 24px;
  line-height: 1.18;
  font-weight: 600;
}

.detail-surface {
  min-height: 240px;
  padding: 20px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background: #fff;
}

.mail-body {
  display: grid;
  gap: 10px;
  margin-bottom: 18px;
  padding: 18px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: #fff;
}

.mail-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  color: var(--muted-foreground);
  font-size: 12px;
}

.mail-meta span + span::before {
  content: '·';
  margin-right: 10px;
  color: #d1d5db;
}

.mail-body p {
  margin: 0;
  color: #1d1d1f;
  font-size: 14px;
  line-height: 1.8;
  white-space: pre-wrap;
}

.code-strip {
  display: flex;
  align-items: center;
  gap: 18px;
  margin-bottom: 20px;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--border);
}

.code-text {
  display: grid;
  gap: 8px;
}

.code-text span {
  color: #7a7a7a;
  font-size: 12px;
}

.code-text button {
  max-width: 100%;
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 8px 10px;
  background: transparent;
  color: var(--primary);
  font-family: var(--font-mono);
  overflow-wrap: anywhere;
  cursor: pointer;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1px;
  margin: 0 0 20px;
  overflow: hidden;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  background: #f0f0f0;
}

.detail-grid div {
  min-width: 0;
  display: grid;
  grid-template-columns: 96px minmax(0, 1fr);
  gap: 12px;
  padding: 12px;
}

.detail-grid div.full-width {
  grid-column: 1 / -1;
}

.flow-progress-panel,
.trace-panel {
  margin: 0 0 20px;
  padding: 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: #fff;
}

.flow-progress-panel :deep(.ant-table-wrapper) {
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  overflow: hidden;
}

.trace-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.trace-head h2 {
  margin: 0;
  color: #1d1d1f;
  font-size: 16px;
  font-weight: 600;
}

.trace-head span {
  color: var(--muted-foreground);
  font-size: 12px;
}

.trace-list {
  display: grid;
  gap: 0;
}

.trace-item {
  position: relative;
  min-width: 0;
  display: grid;
  grid-template-columns: 18px minmax(0, 1fr);
  gap: 3px 10px;
  border: 0;
  margin: 0;
  padding: 0 0 18px;
  background: transparent;
  color: inherit;
  text-align: left;
}

.trace-item::before {
  content: '';
  position: absolute;
  left: 4px;
  top: 13px;
  bottom: 0;
  width: 1px;
  background: #e5e7eb;
}

.trace-item:last-child::before {
  display: none;
}

.trace-item:last-child {
  padding-bottom: 0;
}

.trace-item.clickable {
  cursor: pointer;
}

.trace-item.clickable:hover strong {
  color: var(--primary);
}

.trace-dot {
  grid-column: 1;
  grid-row: 1 / span 2;
  position: relative;
  left: 0;
  top: 3px;
  width: 9px;
  height: 9px;
  border-radius: 999px;
  background: var(--primary);
}

.trace-time {
  grid-column: 2;
  min-width: 0;
  color: var(--muted-foreground);
  font-size: 12px;
  line-height: 1.45;
  white-space: normal;
  overflow-wrap: anywhere;
}

.trace-main {
  grid-column: 2;
  min-width: 0;
  display: grid;
  gap: 4px;
}

.trace-main strong {
  color: #1d1d1f;
  font-size: 13px;
  font-weight: 600;
}

.trace-main small {
  color: var(--muted-foreground);
  font-size: 12px;
  line-height: 1.45;
  overflow-wrap: anywhere;
}

dt {
  color: #7a7a7a;
}

dd {
  min-width: 0;
  margin: 0;
  color: #1d1d1f;
  overflow-wrap: anywhere;
}

.detail-link {
  border: 0;
  padding: 0;
  background: transparent;
  color: var(--primary);
  font: inherit;
  font-weight: 500;
  cursor: pointer;
}

.detail-link:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .detail-head {
    align-items: flex-start;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .trace-item {
    grid-template-columns: 18px minmax(0, 1fr);
  }

  .detail-surface {
    padding: 14px;
  }
}
</style>
