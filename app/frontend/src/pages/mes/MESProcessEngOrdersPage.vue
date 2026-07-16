<template>
  <main class="eng-list-page">
    <section class="detail-head">
      <div>
        <p>工艺单</p>
        <h1>关联工程单</h1>
      </div>
      <a-button @click="goBack">返回</a-button>
    </section>

    <section class="detail-surface">
      <a-spin :spinning="loading">
        <a-table
          :data-source="engRows"
          :columns="columns"
          :pagination="false"
          :scroll="{ x: 'max-content' }"
          row-key="id"
          size="middle"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'id'">
              <a class="id-link" @click="viewEng(record)">#{{ record.id }}</a>
            </template>
            <template v-else-if="column.key === 'itemName'">
              {{ record.item?.name || '物品' }} #{{ record.itemId }}
            </template>
            <template v-else-if="column.dataIndex === 'updateTime'">
              {{ formatTime(record.updateTime) }}
            </template>
            <template v-else-if="column.key === 'action'">
              <a-button type="link" size="small" @click="viewEng(record)">详情</a-button>
            </template>
          </template>
        </a-table>
        <a-empty
          v-if="!engRows.length && !loading"
          class="empty-state"
          description="暂无关联工程单"
        />
        <div v-if="engRows.length" class="mail-hint">
          <span>单击编号或详情进入工程单正文。</span>
        </div>
        <div v-if="engRows.length" class="list-more">
          <MesInfiniteTrigger :has-more="engHasMore" :loading="loading" @load="loadMore" />
        </div>
      </a-spin>
    </section>
  </main>
</template>

<script setup lang="ts">
import MesInfiniteTrigger from '@/components/mes/MesInfiniteTrigger.vue'
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import { DraftStatus, listEngineeringOrder, type EngineeringOrderVO } from '@/api/mesController'

const route = useRoute()
const router = useRouter()
const engRows = ref<EngineeringOrderVO[]>([])
const loading = ref(false)
const engHasMore = ref(false)
let cursorUpdatedAt = ''
let cursorId = 0

const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '生产物品', key: 'itemName', width: 180 },
  { title: '预计', dataIndex: 'expectedQuantity', width: 90 },
  { title: '已产出', dataIndex: 'producedQuantity', width: 90 },
  { title: '合格', dataIndex: 'qualifiedQuantity', width: 90 },
  { title: '正文', dataIndex: 'description', ellipsis: true },
  { title: '更新时间', dataIndex: 'updateTime', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

const processId = () => Number(route.query.processId || 0)

const loadEngOrders = async (next = false) => {
  const pid = processId()
  if (!pid) return
  loading.value = true
  try {
    const res = await listEngineeringOrder({
      processId: pid,
      status: DraftStatus.Submitted,
      pageSize: 30,
      cursorUpdatedAt: next ? cursorUpdatedAt : undefined,
      cursorId: next ? cursorId : undefined,
    })
    if (res.data.code === 0 && res.data.data) {
      const records = next
        ? [...engRows.value, ...(res.data.data.records || [])]
        : res.data.data.records || []
      engRows.value = records as EngineeringOrderVO[]
      engHasMore.value = !!res.data.data.hasMore
      cursorUpdatedAt = res.data.data.nextCursorUpdatedAt || ''
      cursorId = res.data.data.nextCursorId || 0
    }
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  if (!engHasMore.value) return
  loadEngOrders(true)
}

const viewEng = (eng: EngineeringOrderVO) => {
  if (eng.id) {
    router.push({ path: '/mes/detail', query: { kind: 'ENGINEERING_ORDER', id: String(eng.id) } })
  }
}

const goBack = () => {
  const pid = processId()
  if (pid) {
    router.push({ path: '/mes/detail', query: { kind: 'PROCESS', id: String(pid) } })
  } else {
    router.back()
  }
}

const formatTime = (t?: string) => (t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-')

onMounted(() => loadEngOrders())
</script>

<style scoped>
.eng-list-page {
  display: grid;
  gap: 18px;
}

.detail-head {
  min-height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.detail-head p {
  margin: 0 0 4px;
  color: #7a7a7a;
  font-size: 12px;
}

.detail-head h1 {
  margin: 0;
  color: #1d1d1f;
  font-size: 24px;
  line-height: 1.18;
  font-weight: 600;
}

.detail-surface {
  min-height: 240px;
  padding: 20px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background: #fff;
}

.id-link {
  color: var(--primary, #2563eb);
  cursor: pointer;
  font-weight: 600;
}
.id-link:hover {
  text-decoration: underline;
}

.empty-state {
  padding: 36px 0;
}

.mail-hint {
  margin-top: 10px;
  color: var(--muted-foreground, #94a3b8);
  font-size: 12px;
}

.list-more {
  display: flex;
  justify-content: center;
  padding-top: 14px;
}

.muted-text {
  color: #94a3b8;
  font-size: 13px;
}
</style>
