<template>
  <a-table
    row-key="id"
    :columns="columns"
    :data-source="flows"
    :pagination="false"
    :loading="loading"
    :scroll="{ x: 'max-content' }"
    size="middle"
  >
    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'id'">
        <a class="id-link" @click="$emit('view-detail', record)">#{{ record.id }}</a>
      </template>
      <template v-else-if="column.dataIndex === 'flowType'">
        <a-tag>{{ flowTypeLabel(record.flowType) }}</a-tag>
      </template>
      <template v-else-if="column.dataIndex === 'flowStatus'">
        <a-tag :color="flowStatusColor(record.flowStatus)">{{ flowStatusLabel(record.flowStatus) }}</a-tag>
      </template>
      <template v-else-if="column.dataIndex === 'updateTime'">
        {{ formatTime(record.updateTime) }}
      </template>
      <template v-else-if="column.key === 'action'">
        <a-space class="row-actions" size="small">
          <a-button type="link" size="small" @click="$emit('view-detail', record)">详情</a-button>
          <a-button
            v-if="record.flowStatus === 1"
            type="link"
            size="small"
            @click="$emit('edit-draft', record)"
          >
            编辑
          </a-button>
          <a-button
            v-if="record.flowStatus === 1"
            type="link"
            size="small"
            @click="$emit('submit-draft', record)"
          >
            提交
          </a-button>
          <a-popconfirm
            v-if="record.flowStatus === 1"
            title="删除这个草稿？"
            @confirm="$emit('delete-draft', record)"
          >
            <a-button type="link" danger size="small">删除</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </template>
  </a-table>
  <div v-if="flows.length" class="list-more">
    <a-button v-if="hasMore" :loading="loadingMore" @click="$emit('load-more')">加载更多</a-button>
    <span v-else class="muted-text">没有更多了</span>
  </div>
</template>

<script setup lang="ts">
import type { InventoryFlowVO } from '@/api/mesController'
import { flowStatusColor, flowStatusLabel, flowTypeLabel, formatTime } from './display'

defineProps<{
  flows: InventoryFlowVO[]
  loading: boolean
  loadingMore: boolean
  hasMore: boolean
}>()

defineEmits<{
  (e: 'view-detail', record: InventoryFlowVO): void
  (e: 'edit-draft', record: InventoryFlowVO): void
  (e: 'submit-draft', record: InventoryFlowVO): void
  (e: 'delete-draft', record: InventoryFlowVO): void
  (e: 'load-more'): void
}>()

const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '类型', dataIndex: 'flowType', width: 80 },
  { title: '状态', dataIndex: 'flowStatus', width: 80 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '进度', key: 'flowProgress', width: 120, customRender: ({ record }: any) => flowProgressText(record) },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 150 },
]

const flowProgressText = (record: InventoryFlowVO) => {
  const items = record.items || []
  if (!items.length) return '-'
  const finished = items.reduce((sum, item) => sum + (item.finishedQuantity || 0), 0)
  const applied = items.reduce((sum, item) => sum + (item.applyQuantity || 0), 0)
  return `${finished}/${applied}`
}
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
