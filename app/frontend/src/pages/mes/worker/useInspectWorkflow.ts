import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  QUALITY_STATUS_PENDING,
  QUALITY_STATUS_QUALIFIED,
  STOCK_STATUS_OUT_STOCK,
  getEngineeringOrder,
  getItemUnit,
  updateItemUnitStatus,
  type EngineeringOrderVO,
  type ItemUnitVO,
} from '@/api/mesController'
import { parseMesCode } from '@/utils/mesCode'

export const useInspectWorkflow = () => {
  const router = useRouter()

  const inspectOrderCode = ref('')
  const inspectOrder = ref<EngineeringOrderVO>()
  const inspectUnitCode = ref('')
  const inspectUnit = ref<ItemUnitVO>()
  const inspectSubmitting = ref(false)
  const inspectOperationKey = ref(0)

  const loadInspectOrder = async (value: string) => {
    const parsed = parseMesCode(value, 'ENGINEERING_ORDER')
    if (!parsed.id) return
    const res = await getEngineeringOrder({ id: parsed.id })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '读取工程单失败')
      return
    }
    inspectOrder.value = res.data.data
  }

  const resetInspect = () => {
    inspectOrderCode.value = ''
    inspectUnitCode.value = ''
    inspectOrder.value = undefined
    inspectUnit.value = undefined
  }

  const backToInspectScan = async () => {
    await router.push({ path: '/mes/scan', query: { mode: 'inspect' } })
  }

  const loadInspectUnit = async (value: string) => {
    const reopenScanner = () => {
      inspectUnitCode.value = ''
      inspectOperationKey.value += 1
    }
    const parsed = parseMesCode(value, 'ITEM_UNIT')
    if (!parsed.id) {
      reopenScanner()
      return
    }
    const res = await getItemUnit({ id: parsed.id })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '读取单体失败')
      reopenScanner()
      return
    }
    if (inspectOrder.value?.id && res.data.data.engineeringOrderId !== inspectOrder.value.id) {
      message.error('单体不属于当前工程单')
      reopenScanner()
      return
    }
    if (res.data.data.qualityStatus !== QUALITY_STATUS_PENDING) {
      message.error('只能检测待检测单体')
      reopenScanner()
      return
    }
    inspectUnit.value = res.data.data
  }

  const submitInspect = async (qualityStatus: number) => {
    if (!inspectUnit.value?.id) return
    inspectSubmitting.value = true
    try {
      const res = await updateItemUnitStatus({
        id: inspectUnit.value.id,
        stockStatus: inspectUnit.value.stockStatus || STOCK_STATUS_OUT_STOCK,
        qualityStatus,
      })
      if (res.data.code !== 0) throw new Error(res.data.message || '更新失败')
      message.success(qualityStatus === QUALITY_STATUS_QUALIFIED ? '已标记合格' : '已标记不合格')
      inspectUnit.value = undefined
      inspectUnitCode.value = ''
      inspectOperationKey.value += 1
    } catch (error) {
      message.error(error instanceof Error ? error.message : '提交失败')
    } finally {
      inspectSubmitting.value = false
    }
  }

  return {
    inspectOrderCode,
    inspectOrder,
    inspectUnitCode,
    inspectUnit,
    inspectSubmitting,
    inspectOperationKey,
    loadInspectOrder,
    resetInspect,
    backToInspectScan,
    loadInspectUnit,
    submitInspect,
  }
}
