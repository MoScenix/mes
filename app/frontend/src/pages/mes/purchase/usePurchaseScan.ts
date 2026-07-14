import { ref } from 'vue'
import { Modal, message } from 'ant-design-vue'
import type { RouteLocationNormalizedLoadedGeneric, Router } from 'vue-router'
import {
  completeInventoryFlow,
  FLOW_TYPE_IN,
  FlowStatus,
  getInventoryFlow,
  getItemUnit,
  QUALITY_STATUS_QUALIFIED,
  STOCK_STATUS_OUT_STOCK,
  type InventoryFlowVO,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'

export function usePurchaseScan(route: RouteLocationNormalizedLoadedGeneric, router: Router) {
  const scanFlowCode = ref('')
  const scanFlow = ref<InventoryFlowVO>()
  const scanValue = ref('')
  const scanOperationKey = ref(0)

  const clearScannedUnits = () => {
    scanValue.value = ''
  }

  const loadScanFlowById = async (id: number) => {
    const res = await getInventoryFlow({ id })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '读取流转单失败')
      return
    }
    const flow = res.data.data
    if (flow.flowType !== FLOW_TYPE_IN) {
      message.error('只能进入入库流转单')
      return
    }
    if (flow.flowStatus !== FlowStatus.Approved) {
      message.error('只能录入已审批通过的入库流转单')
      return
    }
    scanFlow.value = flow
    clearScannedUnits()
  }

  const loadScanFlow = async (value: string) => {
    const parsed = parseMesCode(value, 'FLOW')
    if (parsed.kind !== 'FLOW' || !parsed.id) {
      message.warning('请输入有效的流转单码')
      return
    }
    await loadScanFlowById(parsed.id)
    if (scanFlow.value?.id) {
      await router.replace({
        query: { ...route.query, panel: 'scan', flowId: String(scanFlow.value.id) },
      })
    }
  }

  const addScanInput = async (value: string) => {
    const reopenScanner = () => {
      scanValue.value = ''
      scanOperationKey.value += 1
    }
    if (!scanFlow.value) {
      message.warning('请先扫描入库流转单')
      reopenScanner()
      return
    }
    const parsed = parseMesCode(value, 'ITEM_UNIT')
    if (parsed.kind !== 'ITEM_UNIT' || !parsed.id) {
      message.warning('请输入有效的库存单体码')
      reopenScanner()
      return
    }
    const unitRes = await getItemUnit({ id: parsed.id })
    const unit = unitRes.data.data
    if (unitRes.data.code !== 0 || !unit?.itemId) {
      message.error(unitRes.data.message || '读取库存单体失败')
      reopenScanner()
      return
    }
    if (unit.stockStatus !== STOCK_STATUS_OUT_STOCK) {
      message.error('只能录入不在库的单体')
      reopenScanner()
      return
    }
    if (unit.qualityStatus !== QUALITY_STATUS_QUALIFIED) {
      message.error('只能录入合格单体')
      reopenScanner()
      return
    }
    const flowItem = (scanFlow.value.items || []).find((item) => item.itemId === unit.itemId)
    if (!flowItem) {
      message.error('该单体物品不在当前流转单明细中')
      reopenScanner()
      return
    }
    Modal.confirm({
      title: '确认入库',
      content: `确认单体 #${parsed.id} 入库？`,
      okText: '确认',
      cancelText: '取消',
      async onOk() {
        scanValue.value = ''
        const res = await completeInventoryFlow({
          id: scanFlow.value?.id,
          itemUnitIds: [parsed.id],
        })
        if (res.data.code !== 0) {
          throw new Error(res.data.message || '入库失败')
        }
        message.success('入库已确认')
        scanOperationKey.value += 1
      },
      onCancel: reopenScanner,
    })
  }

  const clearScan = async () => {
    scanFlow.value = undefined
    scanFlowCode.value = ''
    clearScannedUnits()
    const query = { ...route.query }
    delete query.flowId
    await router.replace({ query: { ...query, panel: 'scan' } })
  }

  const backToInboundScan = async () => {
    await router.push({ path: '/mes/scan', query: { mode: 'inbound' } })
  }

  return {
    scanFlowCode,
    scanFlow,
    scanValue,
    scanOperationKey,
    loadScanFlowById,
    loadScanFlow,
    addScanInput,
    clearScan,
    backToInboundScan,
  }
}
