<template>
  <main class="detail-page">
    <section class="detail-head">
      <div>
        <p>详情</p>
        <h1>{{ title }}</h1>
      </div>
      <a-space>
        <a-button v-if="kind === 'PROCESS'" @click="jumpToEngList">关联工程单</a-button>
        <a-button v-if="hasUnitRows" @click="openUnitList">绑定单体</a-button>
        <ScanButton />
        <MesCodeMenu v-if="codeKind && id" :kind="codeKind" :id="id" />
      </a-space>
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
                <span v-else>{{ row.value }}</span>
              </dd>
            </div>
          </dl>

          <a-table
            v-if="unitRows.length"
            ref="unitTableRef"
            :data-source="unitRows"
            :columns="unitColumns"
            :pagination="false"
            size="small"
            row-key="id"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'action'">
                <a-button type="link" size="small" @click="viewUnit(record)">详情</a-button>
              </template>
            </template>
          </a-table>
        </template>
      </a-spin>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import ScanButton from '@/components/ScanButton.vue'
import MesCodeMenu from '@/components/mes/MesCodeMenu.vue'
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
  getProcess,
  getWorkOrder,
  type EngineeringOrderVO,
  type InventoryFlowVO,
  type ItemVO,
  type ItemUnitVO,
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
  fullWidth?: boolean
}

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const error = ref('')
const unitTableRef = ref()
const detail = ref<InventoryFlowVO | ItemUnitVO | EngineeringOrderVO | WorkOrderVO | ItemVO | ProcessVO>()
const detailKinds: MesDetailKind[] = ['FLOW', 'ITEM_UNIT', 'ENGINEERING_ORDER', 'WORK_ORDER', 'ITEM', 'PROCESS']
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
  const current = detail.value as (InventoryFlowVO | EngineeringOrderVO | WorkOrderVO | ItemVO | ProcessVO) | undefined
  return current && 'name' in current ? current.name : ''
})
const title = computed(() => {
  if (!kind.value || !id.value) return 'MES 对象'
  return titleName.value ? `${titleName.value} #${id.value}` : `${titleByKind[kind.value]} #${id.value}`
})
const codeValue = computed(() => (codeKind.value && id.value ? makeMesCode(codeKind.value, id.value) : ''))

const copyCode = async () => {
  if (!codeValue.value) return
  await navigator.clipboard?.writeText(codeValue.value)
  message.success('已复制')
}

const unitColumns = [
  { title: '单体 ID', dataIndex: 'id', key: 'id', width: 120 },
  { title: '物品 ID', dataIndex: 'itemId', key: 'itemId', width: 120 },
  { title: '库存', dataIndex: 'stockStatusText', key: 'stockStatusText' },
  { title: '质量', dataIndex: 'qualityStatusText', key: 'qualityStatusText' },
  { title: '操作', key: 'action', width: 80 },
]

const stockText = (status?: StockStatus) => {
  if (status === StockStatus.InStock) return '在库'
  if (status === StockStatus.Reserved) return '预留'
  if (status === StockStatus.OutStock) return '出库'
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

const unitRows = computed(() => {
  const source = detail.value as InventoryFlowVO | EngineeringOrderVO | undefined
  return (source?.itemUnits || []).map((unit) => ({
    ...unit,
    stockStatusText: stockText(unit.stockStatus),
    qualityStatusText: qualityText(unit.qualityStatus),
  }))
})

const hasUnitRows = computed(() => {
  return (kind.value === 'ENGINEERING_ORDER' || kind.value === 'FLOW') && unitRows.value.length > 0
})

const updateTimeText = computed(() => {
  const current = detail.value as { updateTime?: string } | undefined
  return current?.updateTime || '-'
})

const bodyText = computed(() => {
  const current = detail.value as { description?: string } | undefined
  return current?.description?.trim() || '暂无正文内容。'
})

const viewUnit = (unit: any) => {
  if (unit && unit.id) {
    router.push({ path: '/mes/detail', query: { kind: 'ITEM_UNIT', id: String(unit.id) } })
  }
}

const jumpToEngList = () => {
  router.push({ path: '/mes/process-eng-orders', query: { processId: String(id.value) } })
}

const openUnitList = () => {
  if (kind.value === 'ENGINEERING_ORDER') {
    router.push({
      path: '/mes/worker',
      query: {
        panel: 'itemUnits',
        view: 'units',
        engineeringOrderId: String(id.value),
      },
    })
    return
  }
  const el = unitTableRef.value?.$el as HTMLElement | undefined
  el?.scrollIntoView({ behavior: 'smooth', block: 'start' })
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
      { label: '生产物品', itemId: order.itemId, item: order.item },
      { label: '预计数量', value: String(order.expectedQuantity ?? 0) },
      { label: '合格数量', value: String(order.qualifiedQuantity ?? 0) },
      { label: '已产出', value: String(order.producedQuantity ?? 0) },
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

const loadDetail = async () => {
  error.value = ''
  detail.value = undefined
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

dt {
  color: #7a7a7a;
}

dd {
  min-width: 0;
  margin: 0;
  color: #1d1d1f;
  overflow-wrap: anywhere;
}

@media (max-width: 768px) {
  .detail-head {
    align-items: flex-start;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .detail-surface {
    padding: 14px;
  }
}

</style>
