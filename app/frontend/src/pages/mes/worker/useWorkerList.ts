import { reactive, ref, type Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { listEngineeringOrder, listItemUnit, type ItemVO } from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'
import type { WorkerPanelType } from './types'

export const useWorkerList = (selectedType: Ref<WorkerPanelType>) => {
  const router = useRouter()
  const route = useRoute()

  const searchText = ref('')
  const searchItemId = ref<number>()
  const dataList = ref<any[]>([])
  const loading = ref(false)
  const loadingMore = ref(false)
  const listPage = reactive({ pageSize: 30, hasMore: false, nextCursorUpdatedAt: '', nextCursorId: 0 })

  const syncCursor = (data?: { hasMore?: boolean; nextCursorUpdatedAt?: string; nextCursorId?: number }) => {
    listPage.hasMore = Boolean(data?.hasMore)
    listPage.nextCursorUpdatedAt = data?.nextCursorUpdatedAt || ''
    listPage.nextCursorId = data?.nextCursorId || 0
  }

  const fetchData = async (next = false) => {
    if (selectedType.value === 'receive' || selectedType.value === 'inspect') return
    if (next) loadingMore.value = true
    else loading.value = true
    try {
      const params = { pageSize: listPage.pageSize }
      const res = selectedType.value === 'engineering'
        ? await listEngineeringOrder({
            ...params,
            itemId: searchItemId.value,
            itemNamePrefix: searchItemId.value ? undefined : searchText.value.trim() || undefined,
            cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
            cursorId: next ? listPage.nextCursorId : undefined,
          })
        : await listItemUnit({
            ...params,
            itemId: searchItemId.value,
            itemNamePrefix: searchItemId.value ? undefined : searchText.value.trim() || undefined,
            engineeringOrderId: Number(route.query.engineeringOrderId || 0) || undefined,
            cursorId: next ? listPage.nextCursorId : undefined,
          })
      if (res.data.code === 0 && res.data.data) {
        dataList.value = next ? [...dataList.value, ...(res.data.data.records ?? [])] : (res.data.data.records ?? [])
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

  return {
    searchText,
    searchItemId,
    dataList,
    loading,
    loadingMore,
    listPage,
    fetchData,
    loadMore,
    onSearch,
    selectSearchItem,
    clearSearch,
  }
}
