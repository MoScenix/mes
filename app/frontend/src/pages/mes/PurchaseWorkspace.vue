<template>
  <div class="workspace-page">
    <div class="workspace-header">
      <div class="header-left">
        <MesListSearchPicker
          v-if="selectedType !== 'scan'"
          v-model="searchText"
          placeholder="输入码或物品名搜索"
          class="search-input"
          @search="onSearch"
          @select-item="selectSearchItem"
          @clear="clearSearch"
        />
      </div>
      <div v-if="selectedType !== 'scan'" class="header-right">
        <div class="header-filters">
          <a-select
            v-if="selectedType === 'flows'"
            v-model:value="flowStatusFilter"
            allow-clear
            placeholder="流转状态"
            :options="flowStatusOptions"
            @change="fetchData()"
          />
          <a-select
            v-if="selectedType === 'itemUnits'"
            v-model:value="stockStatusFilter"
            allow-clear
            placeholder="库存状态"
            :options="stockFilterOptions"
            @change="fetchData()"
          />
          <a-select
            v-if="selectedType === 'itemUnits'"
            v-model:value="qualityStatusFilter"
            allow-clear
            placeholder="质量状态"
            :options="qualityFilterOptions"
            @change="fetchData()"
          />
        </div>
        <a-button v-if="selectedType !== 'items'" type="primary" @click="openCreate">
          <PlusOutlined /> {{ createButtonText }}
        </a-button>
      </div>
    </div>

    <PurchaseScanPanel
      v-if="selectedType === 'scan'"
      v-model:flow-code="scanFlowCode"
      v-model:scan-value="scanValue"
      :scan-flow="scanFlow"
      :operation-key="scanOperationKey"
      @load-flow="loadScanFlow"
      @add-scan-input="addScanInput"
      @back="backToInboundScan"
    />
    <PurchaseItemsList
      v-else-if="selectedType === 'items'"
      :items="dataList"
      :loading="loading"
      :loading-more="loadingMore"
      :has-more="listPage.hasMore"
      show-add-unit
      @view-detail="viewDetail"
      @add-unit="addUnitForItem"
      @load-more="loadMore"
    />
    <PurchaseItemUnitsList
      v-else-if="selectedType === 'itemUnits'"
      :units="dataList"
      :loading="loading"
      :loading-more="loadingMore"
      :has-more="listPage.hasMore"
      @view-detail="viewDetail"
      @edit-unit="editUnit"
      @load-more="loadMore"
    />
    <PurchaseFlowsList
      v-else
      :flows="dataList"
      :loading="loading"
      :loading-more="loadingMore"
      :has-more="listPage.hasMore"
      @view-detail="viewDetail"
      @edit-draft="editFlowDraft"
      @submit-draft="submitFlowDraft"
      @delete-draft="deleteFlowDraft"
      @load-more="loadMore"
    />

    <a-modal
      v-model:open="editOpen"
      title="编辑库存单体"
      :confirm-loading="editSaving"
      width="480px"
      @ok="handleEditUnit"
    >
      <a-form layout="vertical" :model="editUnitForm">
        <a-form-item label="库存状态">
          <a-select v-model:value="editUnitForm.stockStatus" :options="stockOptions" />
        </a-form-item>
        <a-form-item label="质量状态">
          <a-select v-model:value="editUnitForm.qualityStatus" :options="qualityOptions" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import {
  deleteInventoryFlowDraft,
  QUALITY_STATUS_PENDING,
  QUALITY_STATUS_QUALIFIED,
  STOCK_STATUS_IN_STOCK,
  STOCK_STATUS_OUT_STOCK,
  submitInventoryFlow,
  updateItemUnitStatus,
  type InventoryFlowVO,
  type ItemUnitVO,
  type ItemVO,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'
import MesListSearchPicker from '@/components/mes/MesListSearchPicker.vue'
import PurchaseFlowsList from './purchase/PurchaseFlowsList.vue'
import PurchaseItemsList from './purchase/PurchaseItemsList.vue'
import PurchaseItemUnitsList from './purchase/PurchaseItemUnitsList.vue'
import PurchaseScanPanel from './purchase/PurchaseScanPanel.vue'
import { flowStatusOptions, qualityFilterOptions, qualityOptions, stockFilterOptions, stockOptions } from './purchase/options'
import type { PurchasePanel } from './purchase/types'
import { usePurchaseList } from './purchase/usePurchaseList'
import { usePurchaseScan } from './purchase/usePurchaseScan'

const route = useRoute()
const router = useRouter()

const panelFromRoute = () => {
  const panel = String(route.query.panel || 'items')
  return ['flows', 'items', 'itemUnits', 'scan'].includes(panel) ? (panel as PurchasePanel) : 'items'
}

const selectedType = ref<PurchasePanel>(panelFromRoute())
const searchText = ref('')
const searchItemId = ref<number>()
const flowStatusFilter = ref<number>()
const stockStatusFilter = ref<number>()
const qualityStatusFilter = ref<number>()
const engineeringOrderIdFilter = computed(() => Number(route.query.engineeringOrderId || 0) || undefined)
const flowIdFilter = computed(() => Number(route.query.flowId || 0) || undefined)

const {
  dataList,
  loading,
  loadingMore,
  listPage,
  fetchData,
  loadMore,
} = usePurchaseList({
  selectedType,
  searchText,
  searchItemId,
  flowStatusFilter,
  stockStatusFilter,
  qualityStatusFilter,
  engineeringOrderId: engineeringOrderIdFilter,
  flowId: flowIdFilter,
})

const {
  scanFlowCode,
  scanFlow,
  scanValue,
  scanOperationKey,
  loadScanFlowById,
  loadScanFlow,
  addScanInput,
  backToInboundScan,
} = usePurchaseScan(route, router)

const onTypeChange = () => {
  searchText.value = ''
  searchItemId.value = undefined
  flowStatusFilter.value = undefined
  stockStatusFilter.value = undefined
  qualityStatusFilter.value = undefined
  fetchData()
}

const syncPanel = () => {
  selectedType.value = panelFromRoute()
  onTypeChange()
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

const createButtonText = computed(() => {
  if (selectedType.value === 'items') return '新增物料'
  if (selectedType.value === 'itemUnits') return '新增单体'
  return '新建流转单'
})

const openCreate = () => {
  const typeMap: Record<PurchasePanel, string> = {
    items: 'item',
    itemUnits: 'itemUnit',
    flows: 'flow',
    scan: 'flow',
  }
  const query: Record<string, string> = { type: typeMap[selectedType.value] }
  if (selectedType.value === 'itemUnits') {
    query.stockStatus = String(STOCK_STATUS_OUT_STOCK)
    query.qualityStatus = String(QUALITY_STATUS_QUALIFIED)
  }
  router.push({ path: '/mes/create', query })
}

const editOpen = ref(false)
const editSaving = ref(false)
const editUnitForm = reactive({
  id: undefined as number | undefined,
  stockStatus: STOCK_STATUS_IN_STOCK,
  qualityStatus: QUALITY_STATUS_PENDING,
})

const editUnit = (record: ItemUnitVO) => {
  editUnitForm.id = record.id
  editUnitForm.stockStatus = record.stockStatus ?? STOCK_STATUS_IN_STOCK
  editUnitForm.qualityStatus = record.qualityStatus ?? QUALITY_STATUS_PENDING
  editOpen.value = true
}

const handleEditUnit = async () => {
  if (!editUnitForm.id) return
  editSaving.value = true
  try {
    const res = await updateItemUnitStatus({
      id: editUnitForm.id,
      stockStatus: editUnitForm.stockStatus,
      qualityStatus: editUnitForm.qualityStatus,
    })
    if (res.data.code === 0) {
      message.success('已更新')
      editOpen.value = false
      fetchData()
    } else {
      message.error(res.data.message || '更新失败')
    }
  } finally {
    editSaving.value = false
  }
}

const addUnitForItem = (item: ItemVO) => {
  router.push({
    path: '/mes/create',
    query: {
      type: 'itemUnit',
      itemId: String(item.id || ''),
      stockStatus: String(STOCK_STATUS_OUT_STOCK),
      qualityStatus: String(QUALITY_STATUS_QUALIFIED),
    },
  })
}

const editFlowDraft = (record: InventoryFlowVO) => {
  router.push({ path: '/mes/create', query: { type: 'flow', id: String(record.id) } })
}

const submitFlowDraft = async (record: InventoryFlowVO) => {
  if (!record.id) return
  const res = await submitInventoryFlow({ id: record.id })
  if (res.data.code === 0) {
    message.success('流转单已提交')
    fetchData()
  } else {
    message.error(res.data.message || '提交失败')
  }
}

const deleteFlowDraft = async (record: InventoryFlowVO) => {
  if (!record.id) return
  const res = await deleteInventoryFlowDraft({ id: record.id })
  if (res.data.code === 0) {
    message.success('草稿已删除')
    fetchData()
  } else {
    message.error(res.data.message || '删除失败')
  }
}

const viewDetail = (record: InventoryFlowVO | ItemVO | ItemUnitVO) => {
  const kind = selectedType.value === 'flows' ? 'FLOW' : selectedType.value === 'items' ? 'ITEM' : 'ITEM_UNIT'
  router.push({ path: '/mes/detail', query: { kind, id: String(record.id) } })
}

watch(() => [route.query.panel, route.query.engineeringOrderId, route.query.flowId], syncPanel)
onMounted(async () => {
  if (selectedType.value === 'scan') {
    const flowId = Number(route.query.flowId || 0)
    if (flowId > 0) {
      await loadScanFlowById(flowId)
      return
    }
  }
  await fetchData()
})
</script>

<style scoped>
.workspace-page {
  position: relative;
}

.workspace-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}

.header-filters {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-filters :deep(.ant-select) {
  min-width: 118px;
}
</style>
