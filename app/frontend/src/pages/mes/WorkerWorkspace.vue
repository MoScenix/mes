<template>
  <div class="workspace-page">
    <div class="workspace-header">
      <MesListSearchPicker
        v-if="selectedType !== 'receive' && selectedType !== 'inspect' && selectedType !== 'scan'"
        v-model="searchText"
        placeholder="搜索 MES 码或物品名"
        class="search-input"
        @search="onSearch"
        @select-item="selectSearchItem"
        @clear="clearSearch"
      />
      <a-button
        v-if="selectedType === 'itemUnits' || selectedType === 'engineering'"
        type="primary"
        @click="openCreate"
      >
        <PlusOutlined /> 新建
      </a-button>
    </div>

    <ReceiveWorkflow v-if="selectedType === 'receive'" :flow-id="receiveFlowId" />
    <InspectWorkflow v-else-if="selectedType === 'inspect'" :order-id="inspectOrderId" />
    <PurchaseScanPanel
      v-else-if="selectedType === 'scan'"
      v-model:flow-code="scanFlowCode"
      v-model:scan-value="scanValue"
      :scan-flow="scanFlow"
      :operation-key="scanOperationKey"
      @load-flow="loadScanFlow"
      @add-scan-input="addScanInput"
      @back="backToInboundScan"
    />
    <WorkerRecordList
      v-else
      :selected-type="selectedType"
      :data-list="dataList"
      :loading="loading"
      :loading-more="loadingMore"
      :has-more="listPage.hasMore"
      @load-more="loadMore"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { PlusOutlined } from '@ant-design/icons-vue'
import MesListSearchPicker from '@/components/mes/MesListSearchPicker.vue'
import ReceiveWorkflow from './worker/ReceiveWorkflow.vue'
import InspectWorkflow from './worker/InspectWorkflow.vue'
import PurchaseScanPanel from './purchase/PurchaseScanPanel.vue'
import WorkerRecordList from './worker/WorkerRecordList.vue'
import { useWorkerList } from './worker/useWorkerList'
import type { WorkerPanelType } from './worker/types'
import { usePurchaseScan } from './purchase/usePurchaseScan'

const router = useRouter()
const route = useRoute()

const panelFromRoute = () => {
  const panel = String(route.query.panel || 'itemUnits')
  return ['itemUnits', 'engineering', 'receive', 'inspect', 'scan'].includes(panel)
    ? (panel as WorkerPanelType)
    : 'itemUnits'
}

const selectedType = ref<WorkerPanelType>(panelFromRoute())

const receiveFlowId = computed(() => Number(route.query.flowId || 0) || undefined)
const inspectOrderId = computed(() => Number(route.query.orderId || 0) || undefined)
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

const {
  searchText,
  dataList,
  loading,
  loadingMore,
  listPage,
  fetchData,
  loadMore,
  onSearch,
  selectSearchItem,
  clearSearch,
} = useWorkerList(selectedType)

const openCreate = () => {
  router.push({
    path: '/mes/create',
    query: { type: selectedType.value === 'engineering' ? 'engineering' : 'itemUnit' },
  })
}

const syncPanel = async () => {
  selectedType.value = panelFromRoute()
  if (selectedType.value === 'scan') {
    const flowId = Number(route.query.flowId || 0)
    if (flowId > 0) await loadScanFlowById(flowId)
    return
  }
  await fetchData()
}

watch(() => [route.query.panel, route.query.engineeringOrderId, route.query.flowId], syncPanel)

onMounted(async () => {
  await syncPanel()
})
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
</style>
