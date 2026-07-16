import { reactive, ref, type Ref } from 'vue'
import { message } from 'ant-design-vue'
import { listInventoryFlow, listItems, listItemUnit, MesListScope } from '@/api/mesController'
import type { PurchasePanel } from './types'
import { useLoginUserStore } from '@/stores/loginUser'

type UsePurchaseListOptions = {
  selectedType: Ref<PurchasePanel>
  searchText: Ref<string>
  searchItemId: Ref<number | undefined>
  flowStatusFilter: Ref<number | undefined>
  onlyDraft: Ref<boolean>
  flowBusinessType: Ref<number | undefined>
  createdDate: Ref<string | undefined>
  stockStatusFilter: Ref<number | undefined>
  qualityStatusFilter: Ref<number | undefined>
  engineeringOrderId: Ref<number | undefined>
  flowId: Ref<number | undefined>
}

export function usePurchaseList(options: UsePurchaseListOptions) {
  const loginUserStore = useLoginUserStore()
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
    if (options.selectedType.value === 'scan') return
    if (next) loadingMore.value = true
    else loading.value = true
    try {
      const type = options.selectedType.value
      const pageNum = next ? currentPage.value + 1 : currentPage.value
      const params = { pageNum, pageSize: listPage.pageSize }
      let res: any

      if (type === 'flows') {
        res = await listInventoryFlow({
          ...params,
          keyword: options.searchText.value.trim() || undefined,
          createdDate: options.createdDate.value,
          scope: ['admin', 'administrator', '管理员'].includes(
            String(loginUserStore.loginUser.userRole || '').toLowerCase(),
          )
            ? MesListScope.All
            : MesListScope.Mine,
          flowStatus: options.onlyDraft.value ? 1 : options.flowStatusFilter.value,
          onlyDraft: options.onlyDraft.value || undefined,
          businessType: options.flowBusinessType.value,
        })
      } else if (type === 'items') {
        const namePrefix = options.searchText.value.trim() || undefined
        res = await listItems({
          ...params,
          namePrefix,
        })
      } else {
        res = await listItemUnit({
          ...params,
          itemId: options.searchItemId.value,
          itemNamePrefix: options.searchItemId.value
            ? undefined
            : options.searchText.value.trim() || undefined,
          stockStatus: options.stockStatusFilter.value,
          qualityStatus: options.qualityStatusFilter.value,
          engineeringOrderId: options.engineeringOrderId.value,
          inventoryFlowId: options.flowId.value,
          createdDate: options.createdDate.value,
        })
      }

      if (res.data.code === 0 && res.data.data) {
        dataList.value = res.data.data.records ?? []
        currentPage.value = pageNum
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

  const changePage = (page: number) => {
    currentPage.value = page
    fetchData()
  }

  return {
    dataList,
    loading,
    loadingMore,
    listPage,
    currentPage,
    changePage,
    fetchData,
    loadMore,
  }
}
