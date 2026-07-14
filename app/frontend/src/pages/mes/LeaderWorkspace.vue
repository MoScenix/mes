<template>
  <div class="workspace-page">
    <div class="workspace-header">
      <MesListSearchPicker
        v-model="searchText"
        placeholder="搜索 MES 码或物品名"
        class="search-input"
        @search="onSearch"
        @select-item="selectSearchItem"
        @clear="clearSearch"
      />
      <a-space>
        <a-button v-if="selectedType === 'flows'" @click="openCreateFlow">新建流转单</a-button>
        <a-button v-else-if="selectedType === 'engineering'" type="primary" @click="openCreateEng"
          >新建工程单</a-button
        >
        <a-button v-else type="primary" @click="openCreateWorkOrder">新建工单</a-button>
      </a-space>
    </div>

    <WorkOrderMailList
      v-if="selectedType === 'workOrders'"
      mode="sent"
      :orders="dataList"
      :loading="loading"
      :has-more="listPage.hasMore"
      :loading-more="loadingMore"
      @view-detail="viewDetail"
      @view-draft="editWorkOrderDraft"
      @submit-draft="submitWorkOrderDraft"
      @delete-draft="deleteWorkOrderDraftRow"
      @load-more="loadMore"
    />
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
          <a-tag>{{ record.flowType === FLOW_TYPE_IN ? '入库' : '出库' }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'flowStatus' || column.dataIndex === 'status'">
          <a-tag :color="statusColor(record[column.dataIndex])">{{
            statusLabel(record[column.dataIndex])
          }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'toUserId'">
          <MesUserName :id="record.toUserId" />
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
          </a-space>
        </template>
      </template>
    </a-table>
    <div v-if="selectedType !== 'workOrders' && dataList.length" class="list-more">
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

    <a-modal
      v-if="false"
      v-model:open="engOpen"
      title="新建工程单"
      :confirm-loading="engSaving"
      width="520px"
      @ok="handleCreateEng"
    >
      <a-form layout="vertical" :model="engForm">
        <a-form-item label="生产物品" required>
          <a-select
            v-model:value="engForm.itemId"
            show-search
            :filter-option="false"
            placeholder="搜索物品"
            @search="onSearchItem"
          >
            <a-select-option v-for="opt in itemOpts" :key="opt.value" :value="opt.value">{{
              opt.label
            }}</a-select-option>
          </a-select>
        </a-form-item>
        <div class="form-row">
          <a-form-item label="预计产量">
            <a-input-number v-model:value="engForm.expectedQuantity" :min="1" style="width: 100%" />
          </a-form-item>
          <a-form-item label="合格目标">
            <a-input-number
              v-model:value="engForm.qualifiedQuantity"
              :min="1"
              style="width: 100%"
            />
          </a-form-item>
        </div>
        <a-form-item label="说明">
          <a-textarea v-model:value="engForm.description" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-if="false"
      v-model:open="woOpen"
      title="新建工单"
      :confirm-loading="woSaving"
      @ok="handleCreateWO"
    >
      <a-form layout="vertical" :model="woForm">
        <a-form-item label="工单名称" required>
          <a-input v-model:value="woForm.name" placeholder="请输入工单名称" />
        </a-form-item>
        <a-form-item label="接收人 ID" required>
          <a-input-number v-model:value="woForm.toUserId" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="woForm.description" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import {
  FLOW_TYPE_IN,
  FLOW_TYPE_OUT,
  MesListScope,
  listInventoryFlow,
  listEngineeringOrder,
  listWorkOrder,
  createInventoryFlowDraft,
  deleteInventoryFlowDraft,
  submitInventoryFlow,
  createEngineeringOrder,
  submitEngineeringOrder,
  createWorkOrderDraft,
  deleteWorkOrderDraft,
  submitWorkOrder,
  searchItems,
  type ItemVO,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'
import { useLoginUserStore } from '@/stores/loginUser'
import MesListSearchPicker from '@/components/mes/MesListSearchPicker.vue'
import MesUserName from '@/components/mes/MesUserName.vue'
import WorkOrderMailList from '@/components/mes/WorkOrderMailList.vue'

const router = useRouter()
const route = useRoute()
const loginUserStore = useLoginUserStore()

type DataType = 'flows' | 'engineering' | 'workOrders'
const panelFromRoute = () => {
  const panel = String(route.query.panel || 'flows')
  return ['flows', 'engineering', 'workOrders'].includes(panel) ? (panel as DataType) : 'flows'
}
const selectedType = ref<DataType>(panelFromRoute())

const searchText = ref('')
const searchItemId = ref<number>()
const onTypeChange = () => {
  fetchData()
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

const engColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  {
    title: '物品',
    key: 'itemName',
    width: 170,
    customRender: ({ record }: any) => `${record.item?.name || '物品'} #${record.itemId}`,
  },
  { title: '预计', dataIndex: 'expectedQuantity', width: 80 },
  { title: '已产出', dataIndex: 'producedQuantity', width: 80 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 150 },
]
const woColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '接收人', dataIndex: 'toUserId', width: 150 },
  { title: '内容', dataIndex: 'description', ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]
const currentColumns = computed(() => {
  if (selectedType.value === 'engineering') return engColumns
  if (selectedType.value === 'workOrders') return woColumns
  return flowColumns
})

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
    const baseParams = {
      pageSize: listPage.pageSize,
      cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
      cursorId: next ? listPage.nextCursorId : undefined,
    }
    const entitySearchParams = {
      ...baseParams,
      namePrefix: searchText.value.trim() || undefined,
    }
    let res: any
    if (selectedType.value === 'engineering') {
      res = await listEngineeringOrder({
        ...baseParams,
        itemId: searchItemId.value,
        itemNamePrefix: searchItemId.value ? undefined : searchText.value.trim() || undefined,
        scope: MesListScope.Mine,
      })
    } else if (selectedType.value === 'workOrders') {
      res = await listWorkOrder({ ...entitySearchParams })
    } else {
      res = await listInventoryFlow({
        ...baseParams,
        itemNamePrefix: searchText.value.trim() || undefined,
        scope: MesListScope.Mine,
      })
    }
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
  const kind =
    selectedType.value === 'engineering'
      ? 'ENGINEERING_ORDER'
      : selectedType.value === 'workOrders'
        ? 'WORK_ORDER'
        : 'FLOW'
  router.push({ path: '/mes/detail', query: { kind, id: String(record.id) } })
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

const editWorkOrderDraft = (record: any) => {
  router.push({ path: '/mes/create', query: { type: 'workOrder', id: String(record.id) } })
}

const submitWorkOrderDraft = async (record: any) => {
  if (!record.id) return
  const res = await submitWorkOrder({ id: record.id })
  if (res.data.code === 0) {
    message.success('工单已提交')
    fetchData()
  } else {
    message.error(res.data.message || '提交失败')
  }
}

const deleteWorkOrderDraftRow = async (record: any) => {
  if (!record.id) return
  const res = await deleteWorkOrderDraft({ id: record.id })
  if (res.data.code === 0) {
    message.success('草稿已删除')
    fetchData()
  } else {
    message.error(res.data.message || '删除失败')
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
    const res = await createInventoryFlowDraft({ ...flowForm })
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

// --- 新建工程单 ---
const engOpen = ref(false)
const engSaving = ref(false)
const itemOpts = ref<{ label: string; value: number }[]>([])
const engForm = reactive({
  itemId: undefined as number | undefined,
  expectedQuantity: 1,
  qualifiedQuantity: 1,
  description: '',
})
const onSearchItem = async (q: string) => {
  if (!q) {
    itemOpts.value = []
    return
  }
  const res = await searchItems({ namePrefix: q, limit: 10 })
  if (res.data.code === 0) {
    itemOpts.value = (res.data.data?.records || []).map((i: ItemVO) => ({
      label: `${i.name} (#${i.id})`,
      value: i.id!,
    }))
  }
}
const openCreateEng = () => {
  router.push({ path: '/mes/create', query: { type: 'engineering' } })
}
const handleCreateEng = async () => {
  if (!engForm.itemId) {
    message.warning('请选择生产物品')
    return
  }
  engSaving.value = true
  try {
    const res = await createEngineeringOrder({ ...engForm })
    if (res.data.code === 0 && res.data.data) {
      await submitEngineeringOrder({ id: res.data.data })
      message.success('工程单已提交')
      engOpen.value = false
      fetchData()
    } else {
      message.error(res.data.message || '创建失败')
    }
  } finally {
    engSaving.value = false
  }
}

// --- 新建工单 ---
const woOpen = ref(false)
const woSaving = ref(false)
const woForm = reactive({ name: '', toUserId: undefined as number | undefined, description: '' })
const openCreateWorkOrder = () => {
  router.push({ path: '/mes/create', query: { type: 'workOrder' } })
}
const handleCreateWO = async () => {
  if (!woForm.name.trim()) {
    message.warning('请输入工单名称')
    return
  }
  if (!woForm.toUserId) {
    message.warning('请输入接收人')
    return
  }
  woSaving.value = true
  try {
    const res = await createWorkOrderDraft({ ...woForm })
    if (res.data.code === 0 && res.data.data) {
      await submitWorkOrder({ id: res.data.data })
      message.success('工单已提交')
      woOpen.value = false
      fetchData()
    } else {
      message.error(res.data.message || '创建失败')
    }
  } finally {
    woSaving.value = false
  }
}

const formatTime = (t?: string) => (t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-')
const statusColor = (s?: number) =>
  s === 1 ? 'default' : s === 2 ? 'blue' : s === 3 ? 'green' : 'red'
const statusLabel = (s?: number) =>
  s === 1 ? '草稿' : s === 2 ? '待处理' : s === 3 ? '已通过' : '已拒绝'

watch(
  () => route.query.panel,
  () => {
    selectedType.value = panelFromRoute()
    onTypeChange()
  },
)
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
.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
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
