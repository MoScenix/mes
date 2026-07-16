<template>
  <div class="workspace-page">
    <div class="workspace-header">
      <MesListSearchPicker
        v-model="searchText"
        :placeholder="isItemsPanel ? '搜索物料或输入 MES 码' : '搜索工艺、物品名或输入 MES 码'"
        class="search-input"
        @search="onSearch"
        @select-item="selectSearchItem"
        @clear="clearSearch"
      />
      <a-button type="primary" @click="openCreate">
        <PlusOutlined /> {{ isItemsPanel ? '新增物料' : '新建工艺' }}
      </a-button>
    </div>

    <PurchaseItemsList
      v-if="isItemsPanel"
      :items="itemRows"
      :loading="loading"
      :loading-more="loadingMore"
      :has-more="itemPage.hasMore"
      @view-detail="viewItemDetail"
      @add-unit="() => undefined"
      @load-more="loadMoreItems"
    />

    <a-table
      v-else
      row-key="rowKey"
      :columns="columns"
      :data-source="dataList"
      :pagination="false"
      :loading="loading"
      :scroll="{ x: 'max-content' }"
      size="middle"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'id'">#{{ record.id }}</template>
        <template v-else-if="column.dataIndex === 'status'">
          <a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'itemId'">
          <MesItemName :id="record.itemId" :item="record.item" />
        </template>
        <template v-else-if="column.key === 'consume'">
          <span>{{ consumeSummary(record) }}</span>
        </template>
        <template v-else-if="column.dataIndex === 'updateTime'">
          {{ formatTime(record.updateTime) }}
        </template>
        <template v-else-if="column.key === 'engOrders'">
          <a class="eng-jump-link" @click="jumpToEngOrders(record)">跳转</a>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space class="row-actions" size="small">
            <a-button type="link" size="small" @click="viewDetail(record)">详情</a-button>
            <a-button
              v-if="record.status === DraftStatus.Draft || record.status === DraftStatus.Submitted"
              type="link"
              size="small"
              @click="editDraft(record)"
            >
              编辑
            </a-button>
            <template v-if="record.status === DraftStatus.Draft">
              <a-button type="link" size="small" @click="submitDraft(record)">提交</a-button>
              <a-popconfirm title="删除这个草稿？" @confirm="deleteDraft(record)">
                <a-button type="link" danger size="small">删除</a-button>
              </a-popconfirm>
            </template>
          </a-space>
        </template>
      </template>
    </a-table>

    <div v-if="!isItemsPanel && dataList.length" class="list-more">
      <MesInfiniteTrigger :has-more="hasMore" :loading="loadingMore" @load="loadMore" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import {
  DraftStatus,
  MesListScope,
  deleteProcessDraft,
  listItems,
  listProcess,
  submitProcess,
  type ItemVO,
  type PageResult,
  type ProcessVO,
} from '@/api/mesController'
import MesItemName from '@/components/mes/MesItemName.vue'
import MesListSearchPicker from '@/components/mes/MesListSearchPicker.vue'
import MesInfiniteTrigger from '@/components/mes/MesInfiniteTrigger.vue'
import { parseMesCode } from '@/utils/mesCode'
import PurchaseItemsList from './purchase/PurchaseItemsList.vue'

type ProcessRow = ProcessVO & { rowKey: string }

const router = useRouter()
const route = useRoute()
const searchText = ref('')
const searchItemId = ref<number>()
const dataList = ref<ProcessRow[]>([])
const itemRows = ref<ItemVO[]>([])
const loading = ref(false)
const loadingMore = ref(false)

const draftPage = reactive({ hasMore: false, nextCursorUpdatedAt: '', nextCursorId: 0 })
const submittedPage = reactive({ hasMore: false, nextCursorUpdatedAt: '', nextCursorId: 0 })
const itemPage = reactive({ hasMore: false, nextCursorUpdatedAt: '', nextCursorId: 0 })
const pageSize = 30
const isItemsPanel = computed(() => String(route.query.panel || 'processes') === 'items')

const columns = [
  { title: 'ID', key: 'id', width: 72 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 90 },
  { title: '产出物品', dataIndex: 'itemId', width: 180 },
  { title: '消耗明细', key: 'consume', ellipsis: true },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '关联工程单', key: 'engOrders', width: 100 },
  { title: '操作', key: 'action', width: 180 },
]

const hasMore = computed(() => draftPage.hasMore || submittedPage.hasMore)

const syncPage = (
  page: typeof draftPage,
  data?: { hasMore?: boolean; nextCursorUpdatedAt?: string; nextCursorId?: number },
) => {
  page.hasMore = Boolean(data?.hasMore)
  page.nextCursorUpdatedAt = data?.nextCursorUpdatedAt || ''
  page.nextCursorId = data?.nextCursorId || 0
}

const queryStatus = async (status: DraftStatus, next: boolean) => {
  const page = status === DraftStatus.Draft ? draftPage : submittedPage
  if (next && !page.hasMore) {
    return []
  }
  const res = await listProcess({
    status,
    itemId: searchItemId.value,
    itemNamePrefix: searchItemId.value ? undefined : searchText.value.trim() || undefined,
    scope: MesListScope.Mine,
    pageSize,
    cursorUpdatedAt: next ? page.nextCursorUpdatedAt : undefined,
    cursorId: next ? page.nextCursorId : undefined,
  })
  if (res.data.code !== 0) {
    throw new Error(res.data.message || '读取工艺失败')
  }
  syncPage(page, res.data.data)
  return toRows(res.data.data, status)
}

const toRows = (data: PageResult<ProcessVO> | undefined, status: DraftStatus) =>
  (data?.records || []).map((item) => ({
    ...item,
    status: item.status || status,
    rowKey: `${status}-${item.id}`,
  }))

const sortRows = (rows: ProcessRow[]) =>
  rows.sort((a, b) => {
    const timeDiff = dayjs(b.updateTime || 0).valueOf() - dayjs(a.updateTime || 0).valueOf()
    if (timeDiff !== 0) return timeDiff
    return (b.id || 0) - (a.id || 0)
  })

const fetchData = async (next = false) => {
  if (isItemsPanel.value) {
    await fetchItems(next)
    return
  }
  if (next) loadingMore.value = true
  else loading.value = true
  try {
    const [drafts, submitted] = await Promise.all([
      queryStatus(DraftStatus.Draft, next),
      queryStatus(DraftStatus.Submitted, next),
    ])
    const merged = sortRows([...drafts, ...submitted])
    dataList.value = next ? sortRows([...dataList.value, ...merged]) : merged
  } catch (error) {
    message.error(error instanceof Error ? error.message : '读取工艺失败')
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const fetchItems = async (next = false) => {
  if (next) loadingMore.value = true
  else loading.value = true
  try {
    const res = await listItems({
      pageSize,
      namePrefix: searchText.value.trim() || undefined,
      cursorUpdatedAt: next ? itemPage.nextCursorUpdatedAt : undefined,
      cursorId: next ? itemPage.nextCursorId : undefined,
    })
    if (res.data.code !== 0) {
      throw new Error(res.data.message || '读取物料失败')
    }
    const records = res.data.data?.records || []
    itemRows.value = next ? [...itemRows.value, ...records] : records
    itemPage.hasMore = Boolean(res.data.data?.hasMore)
    itemPage.nextCursorUpdatedAt = res.data.data?.nextCursorUpdatedAt || ''
    itemPage.nextCursorId = res.data.data?.nextCursorId || 0
  } catch (error) {
    message.error(error instanceof Error ? error.message : '读取物料失败')
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const loadMore = () => {
  if (!hasMore.value) return
  fetchData(true)
}

const loadMoreItems = () => {
  if (!itemPage.hasMore) return
  fetchItems(true)
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

const viewItemDetail = (record: ItemVO) => {
  router.push({ path: '/mes/detail', query: { kind: 'ITEM', id: String(record.id) } })
}

const viewDetail = (record: ProcessVO) => {
  router.push({ path: '/mes/detail', query: { kind: 'PROCESS', id: String(record.id) } })
}

const jumpToEngOrders = (record: ProcessVO) => {
  router.push({ path: '/mes/process-eng-orders', query: { processId: String(record.id) } })
}

const openCreate = () => {
  router.push({ path: '/mes/create', query: { type: isItemsPanel.value ? 'item' : 'process' } })
}

const editDraft = (record: ProcessVO) => {
  router.push({ path: '/mes/create', query: { type: 'process', id: String(record.id) } })
}

const submitDraft = async (record: ProcessVO) => {
  if (!record.id) return
  const res = await submitProcess({ id: record.id })
  if (res.data.code !== 0) {
    message.error(res.data.message || '提交失败')
    return
  }
  message.success('工艺已提交')
  await fetchData()
}

const deleteDraft = async (record: ProcessVO) => {
  if (!record.id) return
  const res = await deleteProcessDraft({ id: record.id })
  if (res.data.code !== 0) {
    message.error(res.data.message || '删除失败')
    return
  }
  message.success('草稿已删除')
  await fetchData()
}

const consumeSummary = (record: ProcessVO) => {
  const items = record.items || []
  if (!items.length) return record.description || '-'
  return items
    .slice(0, 3)
    .map(
      (item) =>
        `${item.consumeItem?.name || '物品'} #${item.consumeItemId} x ${item.quantity || 0}`,
    )
    .join(' / ')
}

const statusLabel = (status?: DraftStatus) => {
  if (status === DraftStatus.Draft) return '草稿'
  if (status === DraftStatus.Submitted) return '已提交'
  if (status === DraftStatus.Done) return '已完成'
  return '未知'
}

const statusColor = (status?: DraftStatus) => {
  if (status === DraftStatus.Draft) return 'default'
  if (status === DraftStatus.Submitted) return 'blue'
  if (status === DraftStatus.Done) return 'green'
  return 'default'
}

const formatTime = (t?: string) => (t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-')

watch(
  () => route.query.panel,
  () => {
    searchText.value = ''
    searchItemId.value = undefined
    fetchData()
  },
)

onMounted(() => fetchData())
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
  width: 320px;
  max-width: 100%;
}

.row-actions {
  white-space: nowrap;
}

.muted {
  color: var(--muted-foreground);
  font-size: 13px;
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
.eng-jump-link {
  color: var(--primary, #2563eb);
  cursor: pointer;
  font-weight: 500;
}
.eng-jump-link:hover {
  text-decoration: underline;
}
</style>
