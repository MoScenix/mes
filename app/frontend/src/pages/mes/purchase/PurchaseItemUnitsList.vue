<template>
  <a-table
    row-key="id"
    :columns="columns"
    :data-source="units"
    :pagination="false"
    :loading="loading"
    :scroll="{ x: 'max-content' }"
    size="middle"
  >
    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'id'">
        <a class="id-link" @click="$emit('view-detail', record)">#{{ record.id }}</a>
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
      <template v-else-if="column.dataIndex === 'updateTime'">
        {{ formatTime(record.updateTime) }}
      </template>
      <template v-else-if="column.key === 'action'">
        <a-space class="row-actions" size="small">
          <a-button type="link" size="small" @click="$emit('view-detail', record)">详情</a-button>
          <a-button type="link" size="small" @click="$emit('edit-unit', record)">编辑</a-button>
        </a-space>
      </template>
    </template>
  </a-table>
  <div v-if="units.length" class="list-more">
    <a-button v-if="hasMore" :loading="loadingMore" @click="$emit('load-more')">加载更多</a-button>
    <span v-else class="muted-text">没有更多了</span>
  </div>
</template>

<script setup lang="ts">
import type { ItemUnitVO } from '@/api/mesController'
import MesItemName from '@/components/mes/MesItemName.vue'
import { formatTime, qualityLabel, stockLabel } from './display'

defineProps<{
  units: ItemUnitVO[]
  loading: boolean
  loadingMore: boolean
  hasMore: boolean
}>()

defineEmits<{
  (e: 'view-detail', record: ItemUnitVO): void
  (e: 'edit-unit', record: ItemUnitVO): void
  (e: 'load-more'): void
}>()

const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '物品', dataIndex: 'itemId', width: 160 },
  { title: '库存', dataIndex: 'stockStatus', width: 80 },
  { title: '质量', dataIndex: 'qualityStatus', width: 80 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '工程单', key: 'engineeringOrderId', width: 80, customRender: ({ record }: any) => record.engineeringOrderId ? `#${record.engineeringOrderId}` : '-' },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 120 },
]
</script>

<style scoped>
:deep(.ant-table-wrapper) {
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.id-link {
  color: var(--primary);
  cursor: pointer;
  font-weight: 500;
}

.id-link:hover {
  text-decoration: underline;
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

.row-actions {
  white-space: nowrap;
}

.row-actions :deep(.ant-btn) {
  padding-inline: 2px;
}
</style>
