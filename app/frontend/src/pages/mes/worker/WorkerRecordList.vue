<template>
  <a-table
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
        <a-tag :color="flowStatusColor(record.flowStatus)">{{
          flowStatusLabel(record.flowStatus)
        }}</a-tag>
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
  <div v-if="dataList.length" class="list-more">
    <MesInfiniteTrigger :has-more="hasMore" :loading="loadingMore" @load="$emit('load-more')" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import {
  QUALITY_STATUS_PENDING,
  QUALITY_STATUS_QUALIFIED,
  QUALITY_STATUS_UNQUALIFIED,
  STOCK_STATUS_IN_STOCK,
  STOCK_STATUS_OUT_STOCK,
} from '@/api/mesController'
import MesItemName from '@/components/mes/MesItemName.vue'
import MesInfiniteTrigger from '@/components/mes/MesInfiniteTrigger.vue'
import type { WorkerPanelType } from './types'

const props = defineProps<{
  selectedType: WorkerPanelType
  dataList: any[]
  loading: boolean
  loadingMore: boolean
  hasMore: boolean
}>()

defineEmits<{ (event: 'load-more'): void }>()

const router = useRouter()

const unitColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '物品', dataIndex: 'itemId', width: 160 },
  { title: '库存', dataIndex: 'stockStatus', width: 80 },
  { title: '质量', dataIndex: 'qualityStatus', width: 80 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  {
    title: '工程单',
    key: 'engineeringOrderId',
    width: 80,
    customRender: ({ record }: any) =>
      record.engineeringOrderId ? `#${record.engineeringOrderId}` : '-',
  },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

const engColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  {
    title: '生产物品',
    key: 'itemName',
    width: 170,
    customRender: ({ record }: any) => `${record.itemName} #${record.itemId}`,
  },
  { title: '预计', dataIndex: 'expectedQuantity', width: 80 },
  { title: '已产出', dataIndex: 'producedQuantity', width: 80 },
  { title: '合格', dataIndex: 'qualifiedQuantity', width: 80 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

const currentColumns = computed(() =>
  props.selectedType === 'engineering' ? engColumns : unitColumns,
)

const viewDetail = (record: any) => {
  const kind = props.selectedType === 'engineering' ? 'ENGINEERING_ORDER' : 'ITEM_UNIT'
  router.push({ path: '/mes/detail', query: { kind, id: String(record.id) } })
}

const formatTime = (t?: string) => (t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-')
const stockLabel = (s?: number) =>
  s === STOCK_STATUS_IN_STOCK ? '在库' : s === STOCK_STATUS_OUT_STOCK ? '不在库' : '未知'
const qualityLabel = (s?: number) =>
  s === QUALITY_STATUS_PENDING
    ? '待检测'
    : s === QUALITY_STATUS_QUALIFIED
      ? '合格'
      : s === QUALITY_STATUS_UNQUALIFIED
        ? '不合格'
        : '未知'
const flowStatusColor = (s?: number) =>
  s === 1 ? 'default' : s === 2 ? 'blue' : s === 3 ? 'green' : 'red'
const flowStatusLabel = (s?: number) =>
  s === 1 ? '草稿' : s === 2 ? '待处理' : s === 3 ? '已通过' : '已拒绝'
</script>

<style scoped>
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
.muted-text {
  color: var(--muted-foreground, #94a3b8);
  font-size: 13px;
}
</style>
