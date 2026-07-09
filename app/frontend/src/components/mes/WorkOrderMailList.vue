<template>
  <a-spin :spinning="loading">
    <section v-if="orders.length" class="mail-list" aria-label="工单列表">
      <article
        v-for="record in orders"
        :key="record.id"
        class="mail-row"
        :class="{ unread: mode === 'inbox' && record.readStatus === 1 }"
        role="button"
        tabindex="0"
        @click="openOrder(record)"
        @keydown.enter.prevent="openOrder(record)"
      >
        <div class="row-main">
          <div class="row-title">
            <span class="row-id">#{{ record.id }}</span>
            <strong>{{ record.name || '未命名工单' }}</strong>
            <a-tag v-if="mode === 'inbox' && record.readStatus === 1" color="blue">未读</a-tag>
          </div>
          <div class="row-meta">
            <span>{{ mode === 'inbox' ? '发件人' : '收件人' }}</span>
            <span class="party-id">#{{ peerUserId(record) || '-' }}</span>
            <MesUserName :id="peerUserId(record)" />
            <span class="meta-dot"></span>
            <span>{{ statusText(record.status) }}</span>
            <span class="meta-dot"></span>
            <time>{{ displayTime(record.updateTime) }}</time>
          </div>
          <p v-if="record.description" class="row-desc">{{ record.description }}</p>
        </div>

        <div class="row-actions" @click.stop>
          <a-button size="small" @click="openOrder(record)">
            {{ mode === 'sent' && record.status === WorkOrderStatus.Draft ? '编辑' : '详情' }}
          </a-button>
          <a-button
            v-if="mode === 'sent' && record.status === WorkOrderStatus.Draft"
            size="small"
            type="primary"
            @click="emitSubmitDraft(record)"
          >
            提交
          </a-button>
          <a-popconfirm
            v-if="mode === 'sent' && record.status === WorkOrderStatus.Draft"
            title="删除这个草稿？"
            @confirm="emitDeleteDraft(record)"
          >
            <a-button size="small" danger>删除</a-button>
          </a-popconfirm>
          <a-button
            v-if="mode === 'inbox' && record.readStatus === 1"
            size="small"
            type="primary"
            @click="emitMarkRead(record)"
          >
            标记已读
          </a-button>
        </div>
      </article>
    </section>
    <a-empty v-else-if="!loading" description="暂无工单" />
  </a-spin>

  <div v-if="orders.length" class="pager">
    <a-button :disabled="!hasMore" :loading="loadingMore" @click="emitLoadMore">
      加载更多
    </a-button>
  </div>
</template>

<script setup lang="ts">
import type { WorkOrderVO } from '@/api/mesController'
import { WorkOrderStatus } from '@/api/mesController'
import MesUserName from '@/components/mes/MesUserName.vue'

const props = withDefaults(defineProps<{
  mode?: 'inbox' | 'sent'
  orders: WorkOrderVO[]
  loading: boolean
  hasMore: boolean
  loadingMore: boolean
}>(), {
  mode: 'inbox',
})

const emit = defineEmits<{
  (e: 'view-detail', order: WorkOrderVO): void
  (e: 'view-draft', order: WorkOrderVO): void
  (e: 'mark-read', order: WorkOrderVO): void
  (e: 'submit-draft', order: WorkOrderVO): void
  (e: 'delete-draft', order: WorkOrderVO): void
  (e: 'load-more'): void
}>()

const peerUserId = (order: WorkOrderVO) =>
  props.mode === 'inbox' ? order.fromUserId : order.toUserId

const displayTime = (time?: string) => time || '-'

const statusText = (status?: WorkOrderStatus) => {
  if (status === WorkOrderStatus.Draft) return '草稿'
  if (status === WorkOrderStatus.Submitted) return '已提交'
  return '未知'
}

const openOrder = (order: WorkOrderVO) => {
  if (props.mode === 'sent' && order.status === WorkOrderStatus.Draft) {
    emit('view-draft', order)
    return
  }
  emit('view-detail', order)
}

const emitSubmitDraft = (order: WorkOrderVO) => emit('submit-draft', order)
const emitDeleteDraft = (order: WorkOrderVO) => emit('delete-draft', order)
const emitMarkRead = (order: WorkOrderVO) => emit('mark-read', order)
const emitLoadMore = () => emit('load-more')
</script>

<style scoped>
.mail-list {
  display: grid;
  overflow: hidden;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: #ffffff;
}

.mail-row {
  min-width: 0;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background 0.16s ease, box-shadow 0.16s ease;
}

.mail-row:last-child {
  border-bottom: 0;
}

.mail-row:hover,
.mail-row:focus-visible {
  background: #f8fafc;
  outline: none;
}

.mail-row.unread {
  box-shadow: inset 3px 0 0 var(--primary);
}

.row-main {
  min-width: 0;
  display: grid;
  gap: 6px;
}

.row-title {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.row-title strong {
  min-width: 0;
  overflow: hidden;
  color: var(--foreground);
  font-size: 15px;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.row-id,
.party-id {
  color: var(--muted-foreground);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
}

.row-meta {
  min-width: 0;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  color: var(--muted-foreground);
  font-size: 13px;
}

.row-desc {
  margin: 0;
  overflow: hidden;
  color: #4b5563;
  font-size: 13px;
  line-height: 1.5;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta-dot {
  width: 3px;
  height: 3px;
  border-radius: 999px;
  background: #cbd5e1;
}

.row-actions {
  max-width: 220px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 6px;
}

:deep(.ant-btn) {
  border-radius: 6px;
}

.pager {
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .mail-row {
    grid-template-columns: 1fr;
  }

  .row-actions {
    max-width: none;
    justify-content: flex-start;
  }
}
</style>
