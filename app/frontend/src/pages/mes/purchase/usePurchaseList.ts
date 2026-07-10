import { reactive, ref, type Ref } from 'vue'
import { message } from 'ant-design-vue'
import {
  listInventoryFlow,
  listItems,
  listItemUnit,
  MesListScope,
} from '@/api/mesController'
import type { PurchasePanel } from './types'

type UsePurchaseListOptions = {
  selectedType: Ref<PurchasePanel>
  searchText: Ref<string>
  searchItemId: Ref<number | undefined>
  flowStatusFilter: Ref<number | undefined>
  stockStatusFilter: Ref<number | undefined>
  qualityStatusFilter: Ref<number | undefined>
  engineeringOrderId: Ref<number | undefined>
  flowId: Ref<number | undefined>
}

export function usePurchaseList(options: UsePurchaseListOptions) {
  const dataList = ref<any[]>([])
  const loading = ref(false)
  const loadingMore = ref(false)
  const listPage = reactive({
    pageSize: 30,
    hasMore: false,
    nextCursorUpdatedAt: '',
    nextCursorId: 0,
  })

  const syncCursor = (data?: { hasMore?: boolean; nextCursorUpdatedAt?: string; nextCursorId?: number }) => {
    listPage.hasMore = Boolean(data?.hasMore)
    listPage.nextCursorUpdatedAt = data?.nextCursorUpdatedAt || ''
    listPage.nextCursorId = data?.nextCursorId || 0
  }

  const fetchData = async (next = false) => {
    if (options.selectedType.value === 'scan') return
    if (next) loadingMore.value = true
    else loading.value = true
    try {
      const type = options.selectedType.value
      const params = { pageSize: listPage.pageSize }
      let res: any

      if (type === 'flows') {
        res = await listInventoryFlow({
          ...params,
          itemNamePrefix: options.searchText.value.trim() || undefined,
          scope: MesListScope.Mine,
          flowStatus: options.flowStatusFilter.value,
          cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
          cursorId: next ? listPage.nextCursorId : undefined,
        })
      } else if (type === 'items') {
        const namePrefix = options.searchText.value.trim() || undefined
        res = await listItems({
          ...params,
          namePrefix,
          cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
          cursorId: next ? listPage.nextCursorId : undefined,
        })
      } else {
        res = await listItemUnit({
          ...params,
          itemId: options.searchItemId.value,
          itemNamePrefix: options.searchItemId.value ? undefined : options.searchText.value.trim() || undefined,
          stockStatus: options.stockStatusFilter.value,
          qualityStatus: options.qualityStatusFilter.value,
          engineeringOrderId: options.engineeringOrderId.value,
          inventoryFlowId: options.flowId.value,
          cursorUpdatedAt: next ? listPage.nextCursorUpdatedAt : undefined,
          cursorId: next ? listPage.nextCursorId : undefined,
        })
      }

      if (res.data.code === 0 && res.data.data) {
        dataList.value = next ? [...dataList.value, ...(res.data.data.records ?? [])] : (res.data.data.records ?? [])
        syncCursor(res.data.data)
      } else {
        message.error(res.data.message || '获取数据失败')
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

  return {
    dataList,
    loading,
    loadingMore,
    listPage,
    fetchData,
    loadMore,
  }
}
