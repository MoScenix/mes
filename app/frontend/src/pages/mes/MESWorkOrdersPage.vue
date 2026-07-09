<template>
  <main class="workorders-page">
    <section class="toolbar">
      <a-segmented v-model:value="mode" :options="modeOptions" @change="loadOrders" />
      <a-space>
        <a-button type="primary" @click="createOrder">新建工单</a-button>
        <a-input-search
          v-model:value="namePrefix"
          class="name-filter"
          allow-clear
          placeholder="按工单名称搜索"
          @search="loadOrders"
          @change="handleNameFilterChange"
        />
        <MesUserPicker
          v-model="filterUserId"
          class="user-filter"
          :placeholder="mode === 'inbox' ? '筛选发起人' : '筛选接收人'"
        />
        <a-switch v-model:checked="unreadOnly" size="small" @change="loadOrders" />
        <span class="muted">未读</span>
        <a-button :loading="loading" @click="loadOrders">刷新</a-button>
      </a-space>
    </section>

    <WorkOrderMailList
      :mode="mode"
      :orders="filteredOrders"
      :loading="loading"
      :has-more="page.hasMore"
      :loading-more="loadingMore"
      @view-detail="openDetail"
      @view-draft="editDraft"
      @mark-read="markRead"
      @submit-draft="submitDraft"
      @delete-draft="deleteDraft"
      @load-more="loadMore"
    />
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  deleteWorkOrderDraft,
  listWorkOrder,
  markWorkOrderRead,
  submitWorkOrder,
  type WorkOrderVO,
} from '@/api/mesController'
import MesUserPicker from '@/components/mes/MesUserPicker.vue'
import WorkOrderMailList from '@/components/mes/WorkOrderMailList.vue'

const router = useRouter()
const route = useRoute()
const routeMode = () => (route.query.mode === 'sent' ? 'sent' : 'inbox')
const mode = ref<'inbox' | 'sent'>(routeMode())
const unreadOnly = ref(false)
const filterUserId = ref<number>()
const namePrefix = ref('')
const loading = ref(false)
const loadingMore = ref(false)
const orders = ref<WorkOrderVO[]>([])
const page = reactive({
  hasMore: false,
  nextCursorUpdatedAt: '',
  nextCursorId: 0,
})

const modeOptions = [
  { label: '收件', value: 'inbox' },
  { label: '发件', value: 'sent' },
]

const filteredOrders = computed(() => {
  if (!filterUserId.value) return orders.value
  return orders.value.filter((order) =>
    mode.value === 'inbox' ? order.fromUserId === filterUserId.value : order.toUserId === filterUserId.value,
  )
})

const requestParams = (next = false) => ({
  isTo: mode.value === 'inbox',
  isUnread: mode.value === 'inbox' ? unreadOnly.value : false,
  namePrefix: namePrefix.value.trim() || undefined,
  pageSize: 20,
  cursorUpdatedAt: next ? page.nextCursorUpdatedAt : undefined,
  cursorId: next ? page.nextCursorId : undefined,
})

const handleNameFilterChange = () => {
  if (!namePrefix.value.trim()) {
    void loadOrders()
  }
}

const syncPage = (data?: { hasMore?: boolean; nextCursorUpdatedAt?: string; nextCursorId?: number }) => {
  page.hasMore = Boolean(data?.hasMore)
  page.nextCursorUpdatedAt = data?.nextCursorUpdatedAt || ''
  page.nextCursorId = data?.nextCursorId || 0
}

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await listWorkOrder(requestParams())
    if (res.data.code !== 0) {
      throw new Error(res.data.message || '工单加载失败')
    }
    orders.value = res.data.data?.records || []
    syncPage(res.data.data)
  } catch (error) {
    message.error(error instanceof Error ? error.message : '工单加载失败')
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  if (!page.hasMore) return
  loadingMore.value = true
  try {
    const res = await listWorkOrder(requestParams(true))
    if (res.data.code !== 0) {
      throw new Error(res.data.message || '加载失败')
    }
    orders.value = [...orders.value, ...(res.data.data?.records || [])]
    syncPage(res.data.data)
  } catch (error) {
    message.error(error instanceof Error ? error.message : '加载失败')
  } finally {
    loadingMore.value = false
  }
}

const markRead = async (order: WorkOrderVO) => {
  if (!order.id) return
  const res = await markWorkOrderRead({ id: order.id })
  if (res.data.code === 0) {
    order.readStatus = 2
  } else {
    message.error(res.data.message || '标记已读失败')
  }
}

const createOrder = async () => {
  await router.push({ path: '/mes/create', query: { type: 'workOrder' } })
}

const editDraft = async (order: WorkOrderVO) => {
  if (!order.id) return
  await router.push({ path: '/mes/create', query: { type: 'workOrder', id: String(order.id) } })
}

const submitDraft = async (order: WorkOrderVO) => {
  if (!order.id) return
  const res = await submitWorkOrder({ id: order.id })
  if (res.data.code === 0) {
    message.success('工单已提交')
    await loadOrders()
  } else {
    message.error(res.data.message || '提交失败')
  }
}

const deleteDraft = async (order: WorkOrderVO) => {
  if (!order.id) return
  const res = await deleteWorkOrderDraft({ id: order.id })
  if (res.data.code === 0) {
    message.success('草稿已删除')
    await loadOrders()
  } else {
    message.error(res.data.message || '删除失败')
  }
}

const openDetail = async (order: WorkOrderVO) => {
  if (!order.id) return
  await router.push({ path: '/mes/detail', query: { kind: 'WORK_ORDER', id: String(order.id) } })
}

watch(mode, async () => {
  unreadOnly.value = false
  filterUserId.value = undefined
  namePrefix.value = ''
  if (route.query.mode !== mode.value) {
    await router.replace({ query: { ...route.query, mode: mode.value } })
  }
  void loadOrders()
})

watch(
  () => route.query.mode,
  () => {
    const nextMode = routeMode()
    if (mode.value !== nextMode) {
      mode.value = nextMode
    }
  },
)

watch(filterUserId, () => {
  if (!orders.value.length) {
    void loadOrders()
  }
})

onMounted(loadOrders)
</script>

<style scoped>
.workorders-page {
  display: grid;
  gap: 16px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.user-filter {
  width: 220px;
}

.name-filter {
  width: 220px;
}

.muted {
  color: #7a7a7a;
}

@media (max-width: 768px) {
  .toolbar {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
