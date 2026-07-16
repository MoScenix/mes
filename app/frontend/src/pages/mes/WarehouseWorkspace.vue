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
      <a-select
        v-if="selectedType === 'audit' || selectedType === 'flows'"
        v-model:value="businessTypeFilter"
        allow-clear
        placeholder="业务类型"
        :options="businessTypeOptions"
        @change="fetchData()"
      />
      <a-date-picker
        v-if="selectedType === 'audit' || selectedType === 'flows'"
        v-model:value="createdDate"
        value-format="YYYY-MM-DD"
        placeholder="创建日期"
        allow-clear
        @change="fetchData()"
      />
      <a-button v-if="selectedType === 'workOrders'" type="primary" @click="createWorkOrder">
        新建工单
      </a-button>
    </div>

    <!-- 流转单 -->
    <template v-if="selectedType === 'audit' || selectedType === 'flows'">
      <a-table
        row-key="id"
        :columns="auditColumns"
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
          <template v-else-if="column.dataIndex === 'businessType'">
            <a-tag>{{ businessTypeLabel(record.businessType) }}</a-tag>
          </template>
          <template v-else-if="column.dataIndex === 'description'">
            <span>{{ record.description || '-' }}</span>
          </template>
          <template v-else-if="column.dataIndex === 'createTime'">
            {{ formatTime(record.createTime) }}
          </template>
          <template v-else-if="column.dataIndex === 'fromUserId'">
            <MesUserName :id="record.fromUserId" />
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="viewDetail(record)">详情</a-button>
              <template v-if="selectedType === 'audit'">
                <a-button type="primary" size="small" @click="approveFlow(record)">批准</a-button>
                <a-button danger size="small" @click="rejectFlow(record)">拒绝</a-button>
              </template>
            </a-space>
          </template>
        </template>
      </a-table>
    </template>

    <!-- 物资情况 -->
    <template v-else-if="selectedType === 'inventory'">
      <a-table
        row-key="id"
        :columns="inventoryColumns"
        :data-source="dataList"
        :pagination="false"
        :loading="loading"
        :scroll="{ x: 'max-content' }"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'id'">
            <span>#{{ record.id }}</span>
          </template>
        </template>
      </a-table>
    </template>

    <!-- 工单 -->
    <template v-else>
      <WorkOrderMailList
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
    </template>
    <a-pagination
      v-if="selectedType !== 'workOrders' && dataList.length"
      class="page-pagination"
      :current="currentPage"
      :page-size="listPage.pageSize"
      :total="(currentPage - 1) * listPage.pageSize + dataList.length + (listPage.hasMore ? 1 : 0)"
      :show-size-changer="false"
      @change="changePage"
    />

    <!-- 审批 Modal -->
    <a-modal
      v-model:open="auditOpen"
      :title="auditTitle"
      :confirm-loading="auditSaving"
      @ok="handleAudit"
    >
      <a-form layout="vertical">
        <a-form-item label="审批意见">
          <a-textarea
            v-model:value="auditReason"
            :rows="3"
            :placeholder="isApprove ? '输入批准说明（可选）' : '输入拒绝原因'"
          />
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
  FLOW_STATUS_SUBMITTED,
  FlowBusinessType,
  MesListScope,
  listInventoryFlow,
  auditInventoryFlow,
  listWorkOrder,
  listItems,
  deleteWorkOrderDraft,
  submitWorkOrder,
  type ItemVO,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'
import MesListSearchPicker from '@/components/mes/MesListSearchPicker.vue'
import MesUserName from '@/components/mes/MesUserName.vue'
import WorkOrderMailList from '@/components/mes/WorkOrderMailList.vue'

const router = useRouter()
const route = useRoute()

type DataType = 'audit' | 'flows' | 'inventory' | 'workOrders'
const panelFromRoute = () => {
  const panel = String(route.query.panel || 'audit')
  return ['audit', 'flows', 'inventory', 'workOrders'].includes(panel)
    ? (panel as DataType)
    : 'audit'
}
const selectedType = ref<DataType>(panelFromRoute())
const businessTypeFilter = ref<FlowBusinessType>()
const createdDate = ref<string>()
const businessTypeOptions = [
  { label: '采购入库', value: FlowBusinessType.PurchaseInbound },
  { label: '申请货物', value: FlowBusinessType.MaterialRequest },
  { label: '生产入库', value: FlowBusinessType.ProductionInbound },
]

const searchText = ref('')
const searchItemId = ref<number>()
const onTypeChange = () => {
  changePage(1)
}
const onSearch = (value: string) => {
  const parsed = parseMesCode(value)
  if (parsed.kind && parsed.id) {
    router.push({ path: '/mes/detail', query: { kind: parsed.kind, id: String(parsed.id) } })
    return
  }
  searchItemId.value = undefined
  changePage(1)
}
const selectSearchItem = (item: ItemVO) => {
  searchText.value = item.name || ''
  searchItemId.value = item.id
  changePage(1)
}
const clearSearch = () => {
  searchText.value = ''
  searchItemId.value = undefined
  changePage(1)
}

const auditColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '业务类型', dataIndex: 'businessType', width: 110 },
  { title: '申请人', dataIndex: 'fromUserId', width: 150 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  {
    title: '进度',
    key: 'flowProgress',
    width: 120,
    customRender: ({ record }: any) => flowProgressText(record),
  },
  { title: '创建时间', dataIndex: 'createTime', width: 160 },
  { title: '操作', key: 'action', width: 250 },
]

const flowProgressText = (record: any) => {
  const items = record.items || []
  if (!items.length) return '-'
  const finished = items.reduce((sum: number, item: any) => sum + (item.finishedQuantity || 0), 0)
  const applied = items.reduce((sum: number, item: any) => sum + (item.applyQuantity || 0), 0)
  return `${finished}/${applied}`
}

const inventoryColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 160 },
  { title: '单位', dataIndex: 'unit', width: 60 },
  { title: '总库存', dataIndex: 'totalCount', width: 80 },
  { title: '在库', dataIndex: 'inStockCount', width: 80 },
  { title: '可用', dataIndex: 'availableCount', width: 80 },
  { title: '预留', dataIndex: 'reservedCount', width: 80 },
  { title: '待检', dataIndex: 'pendingCount', width: 80 },
  { title: '合格', dataIndex: 'qualifiedCount', width: 80 },
  { title: '不合格', dataIndex: 'unqualifiedCount', width: 80 },
]

const woColumns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '发起人', dataIndex: 'fromUserId', width: 150 },
  { title: '接收人', dataIndex: 'toUserId', width: 150 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 150 },
]

const dataList = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const currentPage = ref(1)
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
    const p = {
      pageNum: next ? currentPage.value + 1 : currentPage.value,
      pageSize: listPage.pageSize,
    }
    let res: any
    if (selectedType.value === 'audit') {
      res = await listInventoryFlow({
        ...p,
        flowStatus: FLOW_STATUS_SUBMITTED,
        keyword: searchText.value.trim() || undefined,
        createdDate: createdDate.value,
        scope: MesListScope.Audit,
        businessType: businessTypeFilter.value,
      })
    } else if (selectedType.value === 'flows') {
      res = await listInventoryFlow({
        ...p,
        keyword: searchText.value.trim() || undefined,
        createdDate: createdDate.value,
        scope: MesListScope.All,
        businessType: businessTypeFilter.value,
      })
    } else if (selectedType.value === 'inventory') {
      res = await listItems({
        ...p,
        namePrefix: searchText.value.trim() || undefined,
        cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
        cursorId: next ? listPage.nextCursorId : undefined,
      })
    } else {
      res = await listWorkOrder({
        ...p,
        namePrefix: searchText.value.trim() || undefined,
        cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
        cursorId: next ? listPage.nextCursorId : undefined,
      })
    }
    if (res.data.code === 0 && res.data.data) {
      dataList.value = res.data.data.records ?? []
      currentPage.value = p.pageNum
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
const changePage = (page: number) => {
  currentPage.value = page
  fetchData()
}

const viewDetail = (record: any) => {
  const kind = selectedType.value === 'workOrders' ? 'WORK_ORDER' : 'FLOW'
  router.push({ path: '/mes/detail', query: { kind, id: String(record.id) } })
}

const createWorkOrder = () => {
  router.push({ path: '/mes/create', query: { type: 'workOrder' } })
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

// --- 审批 ---
const auditOpen = ref(false)
const auditSaving = ref(false)
const auditTarget = ref<any>(null)
const isApprove = ref(true)
const auditReason = ref('')
const auditTitle = computed(() => (isApprove.value ? '批准流转单' : '拒绝流转单'))

const approveFlow = (record: any) => {
  isApprove.value = true
  auditTarget.value = record
  auditReason.value = ''
  auditOpen.value = true
}
const rejectFlow = (record: any) => {
  isApprove.value = false
  auditTarget.value = record
  auditReason.value = ''
  auditOpen.value = true
}

const handleAudit = async () => {
  if (!auditTarget.value) return
  auditSaving.value = true
  try {
    const res = await auditInventoryFlow({ id: auditTarget.value.id, approved: isApprove.value })
    if (res.data.code === 0) {
      message.success(isApprove.value ? '已批准' : '已拒绝')
      auditOpen.value = false
      fetchData()
    } else {
      message.error(res.data.message || '操作失败')
    }
  } finally {
    auditSaving.value = false
  }
}

const formatTime = (t?: string) => (t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-')
const businessTypeLabel = (type?: number) =>
  type === 1 ? '采购入库' : type === 2 ? '申请货物' : type === 3 ? '生产入库' : '未知'

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
.page-pagination { margin-top: 18px; text-align: right; }
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
