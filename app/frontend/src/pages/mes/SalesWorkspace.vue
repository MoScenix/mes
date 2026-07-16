<template>
  <div class="workspace-page">
    <div class="workspace-header">
      <MesListSearchPicker
        v-model="searchText"
        placeholder="搜索 MES 码或物品名"
        class="search-input"
        :item-search="false"
        @search="onSearch"
        @clear="clearSearch"
      />
      <a-button type="primary" @click="openCreateFlow">新建流转单</a-button>
    </div>

    <a-table
      row-key="id"
      :columns="flowColumns"
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
          <a-tag>{{ record.flowType === FLOW_TYPE_IN ? '入库' : '出库' }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'flowStatus'">
          <a-tag :color="statusColor(record.flowStatus)">{{
            statusLabel(record.flowStatus)
          }}</a-tag>
        </template>
        <template
          v-else-if="column.dataIndex === 'createTime' || column.dataIndex === 'updateTime'"
        >
          {{ formatTime(record[column.dataIndex]) }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space class="row-actions" size="small">
            <a-button type="link" size="small" @click="viewDetail(record)">详情</a-button>
            <a-button
              v-if="record.flowStatus === 1"
              type="link"
              size="small"
              @click="editFlowDraft(record)"
            >
              编辑
            </a-button>
            <a-button
              v-if="record.flowStatus === 1"
              type="link"
              size="small"
              @click="submitFlowDraft(record)"
            >
              提交
            </a-button>
            <a-popconfirm
              v-if="record.flowStatus === 1"
              title="删除这个草稿？"
              @confirm="deleteFlowDraft(record)"
            >
              <a-button type="link" danger size="small">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
    <div v-if="dataList.length" class="list-more">
      <a-button v-if="listPage.hasMore" :loading="loadingMore" @click="loadMore">加载更多</a-button>
      <span v-else class="muted-text">没有更多了</span>
    </div>

    <a-modal
      v-if="false"
      v-model:open="flowOpen"
      title="新建流转单"
      :confirm-loading="flowSaving"
      @ok="handleCreateFlow"
    >
      <a-form layout="vertical" :model="flowForm">
        <a-form-item label="流转方向">
          <a-segmented
            v-model:value="flowForm.flowType"
            :options="[
              { label: '入库', value: FLOW_TYPE_IN },
              { label: '出库', value: FLOW_TYPE_OUT },
            ]"
            block
          />
        </a-form-item>
        <a-form-item label="说明">
          <a-textarea v-model:value="flowForm.description" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import {
  FLOW_TYPE_IN,
  FLOW_TYPE_OUT,
  FLOW_BUSINESS_MATERIAL_REQUEST,
  FLOW_BUSINESS_PURCHASE_INBOUND,
  MesListScope,
  listInventoryFlow,
  createInventoryFlowDraft,
  deleteInventoryFlowDraft,
  submitInventoryFlow,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'
import { useLoginUserStore } from '@/stores/loginUser'
import MesListSearchPicker from '@/components/mes/MesListSearchPicker.vue'

const router = useRouter()
const loginUserStore = useLoginUserStore()

type DataType = 'flows'
const selectedType = ref<DataType>('flows')

const searchText = ref('')
const onTypeChange = () => {
  fetchData()
}
const onSearch = (value: string) => {
  const parsed = parseMesCode(value)
  if (parsed.kind && parsed.id) {
    router.push({ path: '/mes/detail', query: { kind: parsed.kind, id: String(parsed.id) } })
    return
  }
  fetchData()
}
const clearSearch = () => {
  searchText.value = ''
  fetchData()
}

const flowColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '类型', dataIndex: 'flowType', width: 80 },
  { title: '状态', dataIndex: 'flowStatus', width: 80 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  {
    title: '进度',
    key: 'flowProgress',
    width: 120,
    customRender: ({ record }: any) => flowProgressText(record),
  },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 150 },
]

const flowProgressText = (record: any) => {
  const items = record.items || []
  if (!items.length) return '-'
  const finished = items.reduce((sum: number, item: any) => sum + (item.finishedQuantity || 0), 0)
  const applied = items.reduce((sum: number, item: any) => sum + (item.applyQuantity || 0), 0)
  return `${finished}/${applied}`
}

const dataList = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const listPage = reactive({
  pageSize: 30,
  hasMore: false,
  nextCursorUpdatedAt: '',
  nextCursorId: 0,
})

const syncCursor = (data?: {
  hasMore?: boolean
  nextCursorUpdatedAt?: string
  nextCursorId?: number
}) => {
  listPage.hasMore = Boolean(data?.hasMore)
  listPage.nextCursorUpdatedAt = data?.nextCursorUpdatedAt || ''
  listPage.nextCursorId = data?.nextCursorId || 0
}

const fetchData = async (next = false) => {
  if (next) loadingMore.value = true
  else loading.value = true
  try {
    const res = await listInventoryFlow({
      pageSize: listPage.pageSize,
      itemNamePrefix: searchText.value.trim() || undefined,
      scope: MesListScope.Mine,
      cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
      cursorId: next ? listPage.nextCursorId : undefined,
    })
    if (res.data.code === 0 && res.data.data) {
      dataList.value = next
        ? [...dataList.value, ...(res.data.data.records ?? [])]
        : (res.data.data.records ?? [])
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
  router.push({ path: '/mes/detail', query: { kind: 'FLOW', id: String(record.id) } })
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

// --- 新建流转单 ---
const flowOpen = ref(false)
const flowSaving = ref(false)
const flowForm = reactive({ flowType: FLOW_TYPE_IN, description: '' })
const openCreateFlow = () => {
  router.push({ path: '/mes/create', query: { type: 'flow' } })
}
const handleCreateFlow = async () => {
  flowSaving.value = true
  try {
    const res = await createInventoryFlowDraft({
      ...flowForm,
      businessType:
        flowForm.flowType === FLOW_TYPE_OUT
          ? FLOW_BUSINESS_MATERIAL_REQUEST
          : FLOW_BUSINESS_PURCHASE_INBOUND,
    })
    if (res.data.code === 0 && res.data.data) {
      await submitInventoryFlow({ id: res.data.data })
      message.success('流转单已提交')
      flowOpen.value = false
      fetchData()
    } else {
      message.error(res.data.message || '创建失败')
    }
  } finally {
    flowSaving.value = false
  }
}

const formatTime = (t?: string) => (t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-')
const statusColor = (s?: number) =>
  s === 1 ? 'default' : s === 2 ? 'blue' : s === 3 ? 'green' : 'red'
const statusLabel = (s?: number) =>
  s === 1 ? '草稿' : s === 2 ? '待处理' : s === 3 ? '已通过' : '已拒绝'

onMounted(fetchData)
</script>

<style scoped>
.workspace-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}
.search-input {
  width: 280px;
  max-width: 100%;
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
