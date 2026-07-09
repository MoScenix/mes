import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import {
  FLOW_STATUS_APPROVED,
  FLOW_TYPE_OUT,
  QUALITY_STATUS_QUALIFIED,
  STOCK_STATUS_IN_STOCK,
  completeInventoryFlow,
  getInventoryFlow,
  getItemUnit,
  type InventoryFlowVO,
  type ItemUnitVO,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'

export const useReceiveWorkflow = () => {
  const router = useRouter()

  const receiveFlowCode = ref('')
  const receiveFlow = ref<InventoryFlowVO>()
  const receiveUnitCode = ref('')
  const receiveUnitIds = ref<number[]>([])
  const receiveUnits = ref<ItemUnitVO[]>([])
  const receiveSubmitting = ref(false)
  const receiveOperationKey = ref(0)

  const receiveExpectedQuantity = computed(() => {
    const units = receiveFlow.value?.itemUnits || []
    if (units.length) return units.length
    return (receiveFlow.value?.items || []).reduce((sum, item) => sum + (item.applyQuantity || 0), 0)
  })

  const receiveQuantityByItem = computed(() => {
    const result = new Map<number, number>()
    for (const unit of receiveUnits.value) {
      if (!unit.itemId) continue
      result.set(unit.itemId, (result.get(unit.itemId) || 0) + 1)
    }
    return result
  })

  const clearReceiveUnits = () => {
    receiveUnitIds.value = []
    receiveUnits.value = []
    receiveUnitCode.value = ''
  }

  const loadReceiveFlow = async (value: string) => {
    const parsed = parseMesCode(value, 'FLOW')
    if (!parsed.id) return
    const res = await getInventoryFlow({ id: parsed.id })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '读取流转单失败')
      return
    }
    if (res.data.data.flowType !== FLOW_TYPE_OUT) {
      message.error('只能领取出库流转单')
      return
    }
    if (res.data.data.flowStatus !== FLOW_STATUS_APPROVED) {
      message.error('流转单尚未审批通过，不能领取')
      return
    }
    receiveFlow.value = res.data.data
    clearReceiveUnits()
  }

  const resetReceive = () => {
    receiveFlowCode.value = ''
    receiveUnitCode.value = ''
    receiveFlow.value = undefined
    clearReceiveUnits()
  }

  const backToReceiveScan = async () => {
    await router.push({ path: '/mes/scan', query: { mode: 'receive' } })
  }

  const addReceiveUnit = async (value: string) => {
    const reopenScanner = () => {
      receiveUnitCode.value = ''
      receiveOperationKey.value += 1
    }
    const parsed = parseMesCode(value, 'ITEM_UNIT')
    if (!parsed.id) {
      reopenScanner()
      return
    }
    const flowUnits = receiveFlow.value?.itemUnits || []
    if (flowUnits.length && !flowUnits.some((unit) => unit.id === parsed.id)) {
      message.error('单体不属于当前流转单')
      reopenScanner()
      return
    }
    const unit = await getItemUnit({ id: parsed.id })
    if (unit.data.code !== 0 || !unit.data.data) {
      message.error(unit.data.message || '读取单体失败')
      reopenScanner()
      return
    }
    if (unit.data.data.stockStatus !== STOCK_STATUS_IN_STOCK) {
      message.error('只能领取在库单体')
      reopenScanner()
      return
    }
    if (unit.data.data.qualityStatus !== QUALITY_STATUS_QUALIFIED) {
      message.error('只能领取合格单体')
      reopenScanner()
      return
    }
    const flowItems = receiveFlow.value?.items || []
    const flowItem = flowItems.find((item) => item.itemId === unit.data.data?.itemId)
    if (!flowUnits.length && !flowItem) {
      message.error('该单体物品不在当前流转单明细中')
      reopenScanner()
      return
    }
    Modal.confirm({
      title: '确认领取',
      content: `确认领取单体 #${parsed.id}？`,
      okText: '确认',
      cancelText: '取消',
      async onOk() {
        receiveUnitCode.value = ''
        const res = await completeInventoryFlow({
          id: receiveFlow.value?.id,
          itemUnitIds: [parsed.id],
        })
        if (res.data.code !== 0) {
          throw new Error(res.data.message || '领取失败')
        }
        message.success('领取已确认')
        receiveOperationKey.value += 1
      },
      onCancel: reopenScanner,
    })
  }

  const removeReceiveUnit = (id: number) => {
    receiveUnitIds.value = receiveUnitIds.value.filter((item) => item !== id)
    receiveUnits.value = receiveUnits.value.filter((item) => item.id !== id)
  }

  const submitReceive = async () => {
    if (!receiveUnitIds.value.length) {
      message.warning('请先扫描单体')
      return
    }
    if (receiveUnitIds.value.length !== receiveExpectedQuantity.value) {
      message.warning('已扫数量需要和流转单数量一致')
      return
    }
    receiveSubmitting.value = true
    try {
      for (const item of receiveFlow.value?.items || []) {
        if ((receiveQuantityByItem.value.get(item.itemId || 0) || 0) !== (item.applyQuantity || 0)) {
          throw new Error(`物品 #${item.itemId} 的扫码数量与申请数量不一致`)
        }
      }
      message.success('领取已确认')
      resetReceive()
    } catch (error) {
      message.error(error instanceof Error ? error.message : '提交失败')
    } finally {
      receiveSubmitting.value = false
    }
  }

  return {
    receiveFlowCode,
    receiveFlow,
    receiveUnitCode,
    receiveUnitIds,
    receiveUnits,
    receiveSubmitting,
    receiveOperationKey,
    receiveExpectedQuantity,
    receiveQuantityByItem,
    loadReceiveFlow,
    resetReceive,
    backToReceiveScan,
    addReceiveUnit,
    removeReceiveUnit,
    clearReceiveUnits,
    submitReceive,
  }
}
