<template>
  <a-table
    row-key="id"
    :columns="columns"
    :data-source="items"
    :pagination="false"
    :loading="loading"
    :scroll="{ x: 'max-content' }"
    size="middle"
  >
    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'id'">
        <a class="id-link" @click="$emit('view-detail', record)">#{{ record.id }}</a>
      </template>
      <template v-else-if="column.dataIndex === 'updateTime'">
        {{ formatTime(record.updateTime) }}
      </template>
      <template v-else-if="column.key === 'action'">
        <a-space class="row-actions" size="small">
          <a-button type="link" size="small" @click="$emit('view-detail', record)">详情</a-button>
          <a-button v-if="record.id" type="link" size="small" @click="$emit('add-unit', record)">添加单体</a-button>
        </a-space>
      </template>
    </template>
  </a-table>
  <div v-if="items.length" class="list-more">
    <a-button v-if="hasMore" :loading="loadingMore" @click="$emit('load-more')">加载更多</a-button>
    <span v-else class="muted-text">没有更多了</span>
  </div>
</template>

<script setup lang="ts">
import type { ItemVO } from '@/api/mesController'
import { formatTime } from './display'

defineProps<{
  items: ItemVO[]
  loading: boolean
  loadingMore: boolean
  hasMore: boolean
}>()

defineEmits<{
  (e: 'view-detail', record: ItemVO): void
  (e: 'add-unit', record: ItemVO): void
  (e: 'load-more'): void
}>()

const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 160 },
  { title: '单位', dataIndex: 'unit', width: 80 },
  { title: '库存', key: 'totalCount', width: 60, customRender: ({ record }: any) => record.totalCount ?? 0 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 140 },
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
