<template>
  <main class="create-page">
    <header class="create-head">
      <button type="button" class="back-button" @click="handleLeave">返回</button>
      <h1>{{ title }}</h1>
    </header>

    <section class="form-surface">
      <a-form layout="vertical" @finish="submit">
        <template v-if="type === 'item'">
          <a-form-item label="物品名称" required>
            <a-input v-model:value="itemForm.name" placeholder="如 M8 螺栓" />
          </a-form-item>
          <a-form-item label="计量单位">
            <a-input v-model:value="itemForm.unit" placeholder="个 / 件 / kg" />
          </a-form-item>
          <a-form-item label="说明">
            <a-textarea v-model:value="itemForm.description" :rows="4" />
          </a-form-item>
        </template>

        <template v-else-if="type === 'itemUnit'">
          <a-form-item label="物品" required>
            <MesItemPicker v-model="unitForm.itemId" placeholder="输入物品名或 ID" />
          </a-form-item>
          <a-form-item label="绑定工程单">
            <MesEngineeringOrderPicker
              v-model="unitForm.engineeringOrderId"
              :item-id="unitForm.itemId"
              placeholder="可选，输入工程单名或 ID"
              @select-order="selectUnitEngineeringOrder"
            />
          </a-form-item>
          <div class="form-row">
            <a-form-item label="库存状态">
              <a-select v-model:value="unitForm.stockStatus" :options="stockOptions" />
            </a-form-item>
            <a-form-item label="质量状态">
              <a-select v-model:value="unitForm.qualityStatus" :options="qualityOptions" />
            </a-form-item>
          </div>
          <a-form-item label="说明">
            <a-input v-model:value="unitForm.description" />
          </a-form-item>
        </template>

        <template v-else-if="type === 'flow'">
          <a-form-item label="流转单名称" required>
            <a-input v-model:value="flowForm.name" placeholder="请输入流转单名称" />
          </a-form-item>
          <div class="form-row">
            <a-form-item label="流转方向">
              <a-segmented v-model:value="flowForm.flowType" :options="flowTypeOptions" />
            </a-form-item>
            <a-form-item label="接收人">
              <MesUserPicker v-model="flowForm.toUserId" placeholder="输入人名、账号或 ID" />
            </a-form-item>
          </div>
          <a-form-item label="说明">
            <a-textarea v-model:value="flowForm.description" :rows="4" />
          </a-form-item>
          <a-form-item label="物品明细" required>
            <div class="flow-lines">
              <div v-for="(line, index) in flowForm.items" :key="index" class="flow-line">
                <MesItemPicker
                  v-model="line.itemId"
                  class="flow-line-item"
                  placeholder="输入物品名或 ID"
                />
                <a-input-number
                  v-model:value="line.applyQuantity"
                  :min="1"
                  class="flow-line-quantity"
                  placeholder="数量"
                />
                <a-button
                  type="text"
                  danger
                  :disabled="flowForm.items.length === 1"
                  @click="removeFlowLine(index)"
                >
                  删除
                </a-button>
              </div>
              <a-button type="dashed" block @click="addFlowLine">添加物品</a-button>
            </div>
          </a-form-item>
        </template>

        <template v-else-if="type === 'process'">
          <a-form-item label="工艺名称" required>
            <a-input v-model:value="processForm.name" placeholder="请输入工艺名称" />
          </a-form-item>
          <a-form-item label="产出物品" required>
            <MesItemPicker v-model="processForm.itemId" placeholder="输入物品名或 ID" />
          </a-form-item>
          <a-form-item label="消耗物品" required>
            <div class="flow-lines">
              <div v-for="(line, index) in processForm.items" :key="index" class="flow-line">
                <MesItemPicker
                  v-model="line.consumeItemId"
                  class="flow-line-item"
                  placeholder="输入消耗物品名或 ID"
                />
                <a-input-number
                  v-model:value="line.quantity"
                  :min="1"
                  class="flow-line-quantity"
                  placeholder="数量"
                />
                <a-button
                  type="text"
                  danger
                  :disabled="processForm.items.length === 1"
                  @click="removeProcessLine(index)"
                >
                  删除
                </a-button>
              </div>
              <a-button type="dashed" block @click="addProcessLine">添加消耗物品</a-button>
            </div>
          </a-form-item>
          <a-form-item label="说明">
            <a-textarea v-model:value="processForm.description" :rows="4" />
          </a-form-item>
        </template>

        <template v-else-if="type === 'engineering'">
          <a-steps class="create-steps" :current="engineeringStep">
            <a-step title="选物品" />
            <a-step title="选工艺" />
            <a-step title="填信息" />
          </a-steps>

          <a-form-item v-if="!engineeringForm.itemId" label="生产物品" required>
            <MesItemPicker
              v-model="engineeringForm.itemId"
              placeholder="输入物品名或 ID"
              @select-item="selectEngineeringItem"
            />
          </a-form-item>

          <template v-if="engineeringForm.itemId">
            <div class="selected-strip">
              <span>生产物品</span>
              <strong class="selected-item-name">
                <MesItemName :id="engineeringForm.itemId" :item="selectedEngineeringItem" />
              </strong>
              <a-button size="small" type="text" @click="resetEngineeringItem">重新选择</a-button>
            </div>
            <a-form-item label="生产工艺" required>
              <a-select
                v-model:value="engineeringForm.processId"
                :loading="processLoading"
                :options="processOptions"
                placeholder="选择该物品的已提交工艺"
                style="width: 100%"
              />
              <div v-if="!processLoading && !processOptions.length" class="form-help">
                该物品还没有可用工艺，请先由工艺工程师维护并提交工艺单。
              </div>
            </a-form-item>
          </template>

          <template v-if="engineeringForm.itemId && engineeringForm.processId">
            <a-form-item label="工程单名称" required>
              <a-input v-model:value="engineeringForm.name" placeholder="请输入工程单名称" />
            </a-form-item>
            <div class="form-row">
              <a-form-item label="预计产量">
                <a-input-number v-model:value="engineeringForm.expectedQuantity" :min="1" class="wide" />
              </a-form-item>
              <a-form-item label="合格目标">
                <a-input-number v-model:value="engineeringForm.qualifiedQuantity" :min="0" class="wide" />
              </a-form-item>
            </div>
            <a-form-item label="说明">
              <a-textarea v-model:value="engineeringForm.description" :rows="4" />
            </a-form-item>
          </template>
        </template>

        <template v-else-if="type === 'workOrder'">
          <a-form-item label="工单名称" required>
            <a-input v-model:value="workOrderForm.name" placeholder="请输入工单名称" />
          </a-form-item>
          <a-form-item label="接收人" required>
            <MesUserPicker v-model="workOrderForm.toUserId" placeholder="输入人名、账号或 ID" />
          </a-form-item>
          <a-form-item label="描述">
            <a-textarea v-model:value="workOrderForm.description" :rows="5" />
          </a-form-item>
        </template>

        <a-empty v-else description="未知的新建类型" />

        <footer v-if="type" class="form-actions">
          <a-button
            type="primary"
            :html-type="canSaveDraft ? 'button' : 'submit'"
            :loading="saving || savingDraft"
            @click="canSaveDraft ? submitDraft() : undefined"
          >
            提交
          </a-button>
        </footer>
      </a-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import {
  FLOW_TYPE_IN,
  FLOW_TYPE_OUT,
  QUALITY_STATUS_PENDING,
  QUALITY_STATUS_QUALIFIED,
  STOCK_STATUS_OUT_STOCK,
  DraftStatus,
  MesListScope,
  addItem,
  addItemUnit,
  createEngineeringOrder,
  createInventoryFlowDraft,
  createProcessDraft,
  createWorkOrderDraft,
  getEngineeringOrder,
  getInventoryFlow,
  getProcess,
  getWorkOrder,
  listProcess,
  submitEngineeringOrder,
  submitInventoryFlow,
  submitProcess,
  submitWorkOrder,
  updateEngineeringOrder,
  updateInventoryFlowDraft,
  updateProcessDraft,
  updateWorkOrderDraft,
  type EngineeringOrderVO,
  type ItemVO,
} from '@/api/mesController'
import MesEngineeringOrderPicker from '@/components/mes/MesEngineeringOrderPicker.vue'
import MesItemName from '@/components/mes/MesItemName.vue'
import MesItemPicker from '@/components/mes/MesItemPicker.vue'
import MesUserPicker from '@/components/mes/MesUserPicker.vue'

type CreateType = 'item' | 'itemUnit' | 'flow' | 'process' | 'engineering' | 'workOrder'

const route = useRoute()
const router = useRouter()
const saving = ref(false)
const savingDraft = ref(false)
const leaving = ref(false)
const draftId = ref<number | undefined>(Number(route.query.id || 0) || undefined)
const type = computed(() => String(route.query.type || '') as CreateType)
const canSaveDraft = computed(() => ['flow', 'process', 'engineering', 'workOrder'].includes(type.value))

const titleMap: Record<CreateType, string> = {
  item: '新建物品类型',
  itemUnit: '新建库存单体',
  flow: '新建流转单',
  process: '新建工艺',
  engineering: '新建工程单',
  workOrder: '新建工单',
}

const title = computed(() => titleMap[type.value] || '新建')
const initialStockStatus = Number(route.query.stockStatus || 0) || STOCK_STATUS_OUT_STOCK
const initialQualityStatus = Number(route.query.qualityStatus || 0) || QUALITY_STATUS_PENDING

const itemForm = reactive({ name: '', unit: '个', description: '' })
const unitForm = reactive({
  itemId: Number(route.query.itemId || 0) || undefined,
  stockStatus: initialStockStatus,
  qualityStatus: initialQualityStatus,
  description: '',
  engineeringOrderId: Number(route.query.engineeringOrderId || 0) || undefined as number | undefined,
})
const flowForm = reactive({
  name: '',
  flowType: FLOW_TYPE_IN,
  toUserId: undefined as number | undefined,
  description: '',
  items: [{ itemId: undefined as number | undefined, applyQuantity: 1 }],
})
const processForm = reactive({
  name: '',
  itemId: undefined as number | undefined,
  description: '',
  items: [{ consumeItemId: undefined as number | undefined, quantity: 1 }],
})
const engineeringForm = reactive({
  name: '',
  itemId: Number(route.query.itemId || 0) || undefined as number | undefined,
  processId: Number(route.query.processId || 0) || undefined as number | undefined,
  expectedQuantity: 1,
  qualifiedQuantity: 1,
  description: '',
})
const workOrderForm = reactive({
  name: '',
  toUserId: undefined as number | undefined,
  description: '',
})

const stockOptions = [{ label: '不在库', value: STOCK_STATUS_OUT_STOCK }]
const qualityOptions = computed(() => [
  initialQualityStatus === QUALITY_STATUS_QUALIFIED
    ? { label: '合格', value: QUALITY_STATUS_QUALIFIED }
    : { label: '待检测', value: QUALITY_STATUS_PENDING },
])
const flowTypeOptions = [
  { label: '入库', value: FLOW_TYPE_IN },
  { label: '出库', value: FLOW_TYPE_OUT },
]
const processLoading = ref(false)
const processOptions = ref<{ label: string; value: number }[]>([])
const selectedEngineeringItem = ref<ItemVO>()
const selectedUnitEngineeringOrder = ref<EngineeringOrderVO>()
const engineeringStep = computed(() => {
  if (!engineeringForm.itemId) return 0
  if (!engineeringForm.processId) return 1
  return 2
})

const selectEngineeringItem = (item: ItemVO) => {
  selectedEngineeringItem.value = item
}

const selectUnitEngineeringOrder = (order: EngineeringOrderVO) => {
  selectedUnitEngineeringOrder.value = order
  if (!order.itemId) return
  if (!unitForm.itemId) {
    unitForm.itemId = order.itemId
    return
  }
  if (unitForm.itemId !== order.itemId) {
    message.error('单体物品与工程单生产物品不一致')
    unitForm.engineeringOrderId = undefined
    selectedUnitEngineeringOrder.value = undefined
  }
}

const validateUnitEngineeringOrder = async () => {
  if (!unitForm.engineeringOrderId) return true
  let order = selectedUnitEngineeringOrder.value
  if (!order || order.id !== unitForm.engineeringOrderId) {
    const res = await getEngineeringOrder({ id: unitForm.engineeringOrderId })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '工程单不存在')
      return false
    }
    order = res.data.data
    selectedUnitEngineeringOrder.value = order
  }
  if (order.status !== DraftStatus.Submitted) {
    message.error('只能绑定已提交工程单')
    return false
  }
  if (unitForm.itemId && order.itemId && unitForm.itemId !== order.itemId) {
    message.error('单体物品与工程单生产物品不一致')
    return false
  }
  if (!unitForm.itemId && order.itemId) {
    unitForm.itemId = order.itemId
  }
  return true
}

const loadEngineeringProcesses = async () => {
  if (!engineeringForm.itemId) {
    processOptions.value = []
    engineeringForm.processId = undefined
    return
  }
  processLoading.value = true
  try {
    const res = await listProcess({
      itemId: engineeringForm.itemId,
      status: DraftStatus.Submitted,
      scope: MesListScope.All,
      pageSize: 50,
    })
    if (res.data.code !== 0) {
      message.error(res.data.message || '读取工艺失败')
      processOptions.value = []
      return
    }
    processOptions.value = (res.data.data?.records || [])
      .filter((item) => item.id)
      .map((item) => ({
        value: item.id!,
        label: item.name ? `${item.name} #${item.id}` : `工艺 #${item.id}`,
      }))
    if (engineeringForm.processId && !processOptions.value.some((item) => item.value === engineeringForm.processId)) {
      engineeringForm.processId = undefined
    }
  } finally {
    processLoading.value = false
  }
}

const resetEngineeringItem = () => {
  engineeringForm.itemId = undefined
  engineeringForm.processId = undefined
  selectedEngineeringItem.value = undefined
  processOptions.value = []
}

const addFlowLine = () => {
  flowForm.items.push({ itemId: undefined, applyQuantity: 1 })
}

const removeFlowLine = (index: number) => {
  if (flowForm.items.length <= 1) return
  flowForm.items.splice(index, 1)
}

const addProcessLine = () => {
  processForm.items.push({ consumeItemId: undefined, quantity: 1 })
}

const removeProcessLine = (index: number) => {
  if (processForm.items.length <= 1) return
  processForm.items.splice(index, 1)
}

const normalizedFlowItems = () =>
  flowForm.items
    .filter((item) => item.itemId && item.applyQuantity)
    .map((item) => ({ itemId: item.itemId, applyQuantity: item.applyQuantity }))

const normalizedProcessItems = () =>
  processForm.items
    .filter((item) => item.consumeItemId && item.quantity)
    .map((item) => ({ consumeItemId: item.consumeItemId, quantity: item.quantity }))

const hasContent = computed(() => {
  if (type.value === 'item') {
    return Boolean(itemForm.name.trim() || itemForm.description.trim() || itemForm.unit !== '个')
  }
  if (type.value === 'itemUnit') {
    return Boolean(unitForm.itemId || unitForm.description.trim())
  }
  if (type.value === 'flow') {
    return Boolean(flowForm.name.trim() || flowForm.toUserId || flowForm.description.trim() || normalizedFlowItems().length)
  }
  if (type.value === 'process') {
    return Boolean(processForm.name.trim() || processForm.itemId || processForm.description.trim() || normalizedProcessItems().length)
  }
  if (type.value === 'engineering') {
    return Boolean(
      engineeringForm.name.trim() ||
      engineeringForm.itemId ||
      engineeringForm.processId ||
      engineeringForm.description.trim() ||
      engineeringForm.expectedQuantity !== 1 ||
      engineeringForm.qualifiedQuantity !== 1,
    )
  }
  if (type.value === 'workOrder') {
    return Boolean(workOrderForm.name.trim() || workOrderForm.toUserId || workOrderForm.description.trim())
  }
  return false
})

const draftPayload = () => {
  if (type.value === 'flow') {
    return { ...flowForm, items: normalizedFlowItems() }
  }
  if (type.value === 'process') {
    return { ...processForm, items: normalizedProcessItems() }
  }
  if (type.value === 'engineering') {
    return { ...engineeringForm }
  }
  if (type.value === 'workOrder') {
    return { ...workOrderForm }
  }
  return {}
}

const validateDraft = () => {
  if (type.value === 'flow' && !flowForm.name.trim()) {
    message.warning('请输入流转单名称')
    return false
  }
  if (type.value === 'flow' && !normalizedFlowItems().length) {
    message.warning('请添加物品和数量')
    return false
  }
  if (type.value === 'process' && !processForm.name.trim()) {
    message.warning('请输入工艺名称')
    return false
  }
  if (type.value === 'process' && !processForm.itemId) {
    message.warning('请选择产出物品')
    return false
  }
  if (type.value === 'process' && !normalizedProcessItems().length) {
    message.warning('请添加消耗物品和数量')
    return false
  }
  if (type.value === 'engineering' && !engineeringForm.name.trim()) {
    message.warning('请输入工程单名称')
    return false
  }
  if (type.value === 'engineering' && !engineeringForm.itemId) {
    message.warning('请选择生产物品')
    return false
  }
  if (type.value === 'engineering' && !engineeringForm.processId) {
    message.warning('请选择生产工艺')
    return false
  }
  if (type.value === 'workOrder' && !workOrderForm.name.trim()) {
    message.warning('请输入工单名称')
    return false
  }
  if (type.value === 'workOrder' && !workOrderForm.toUserId) {
    message.warning('请输入接收人')
    return false
  }
  return true
}

const saveDraft = async (silent = false, syncRoute = true) => {
  if (!canSaveDraft.value) return false
  if (!validateDraft()) return false
  savingDraft.value = true
  try {
    let res: any
    if (type.value === 'flow') {
      res = draftId.value
        ? await updateInventoryFlowDraft({ id: draftId.value, ...draftPayload() })
        : await createInventoryFlowDraft(draftPayload())
    } else if (type.value === 'engineering') {
      res = draftId.value
        ? await updateEngineeringOrder({ id: draftId.value, ...draftPayload() })
        : await createEngineeringOrder(draftPayload())
    } else if (type.value === 'process') {
      res = draftId.value
        ? await updateProcessDraft({ id: draftId.value, ...draftPayload() })
        : await createProcessDraft(draftPayload())
    } else if (type.value === 'workOrder') {
      res = draftId.value
        ? await updateWorkOrderDraft({ id: draftId.value, ...draftPayload() })
        : await createWorkOrderDraft(draftPayload())
    }

    if (res?.data?.code !== 0) {
      message.error(res?.data?.message || '草稿保存失败')
      return false
    }
    const wasUpdate = Boolean(draftId.value)
    if (!draftId.value && typeof res?.data?.data === 'number') {
      draftId.value = res.data.data
      if (syncRoute) {
        await router.replace({ query: { ...route.query, id: String(draftId.value) } })
      }
    }
    if (!silent) {
      message.success(wasUpdate ? '草稿已更新' : '草稿已保存')
    }
    return true
  } finally {
    savingDraft.value = false
  }
}

const loadDraft = async () => {
  if (!draftId.value) return
  if (type.value === 'workOrder') {
    const res = await getWorkOrder({ id: draftId.value })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '读取草稿失败')
      return
    }
    workOrderForm.name = res.data.data.name || ''
    workOrderForm.toUserId = res.data.data.toUserId
    workOrderForm.description = res.data.data.description || ''
    return
  }
  if (type.value === 'process') {
    const res = await getProcess({ id: draftId.value })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '读取草稿失败')
      return
    }
    processForm.name = res.data.data.name || ''
    processForm.itemId = res.data.data.itemId
    processForm.description = res.data.data.description || ''
    processForm.items.splice(
      0,
      processForm.items.length,
      ...((res.data.data.items || []).map((item) => ({
        consumeItemId: item.consumeItemId,
        quantity: item.quantity || 1,
      })) || [{ consumeItemId: undefined, quantity: 1 }]),
    )
    if (!processForm.items.length) {
      processForm.items.push({ consumeItemId: undefined, quantity: 1 })
    }
    return
  }
  if (type.value === 'engineering') {
    const res = await getEngineeringOrder({ id: draftId.value })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '读取草稿失败')
      return
    }
    engineeringForm.name = res.data.data.name || ''
    engineeringForm.itemId = res.data.data.itemId
    selectedEngineeringItem.value = res.data.data.item
    engineeringForm.processId = res.data.data.processId
    engineeringForm.expectedQuantity = res.data.data.expectedQuantity || 1
    engineeringForm.qualifiedQuantity = res.data.data.qualifiedQuantity || 1
    engineeringForm.description = res.data.data.description || ''
    return
  }
  if (type.value !== 'flow') return
  const res = await getInventoryFlow({ id: draftId.value })
  if (res.data.code !== 0 || !res.data.data) {
    message.error(res.data.message || '读取草稿失败')
    return
  }
  const flow = res.data.data
  flowForm.name = flow.name || ''
  flowForm.flowType = flow.flowType || FLOW_TYPE_IN
  flowForm.toUserId = flow.toUserId
  flowForm.description = flow.description || ''
  flowForm.items.splice(
    0,
    flowForm.items.length,
    ...((flow.items || []).map((item) => ({
      itemId: item.itemId,
      applyQuantity: item.applyQuantity || 1,
    })) || [{ itemId: undefined, applyQuantity: 1 }]),
  )
  if (!flowForm.items.length) {
    flowForm.items.push({ itemId: undefined, applyQuantity: 1 })
  }
}

const leaveBack = async () => {
  leaving.value = true
  await exitAfterLeave()
}

const exitAfterLeave = async () => {
  if (type.value === 'engineering') {
    await router.push({ path: '/mes/leader', query: { panel: 'engineering', view: 'engineering' } })
    return
  }
  if (type.value === 'process') {
    await router.push({ path: '/mes/processes', query: { panel: 'processes', view: 'processes' } })
    return
  }
  if (type.value === 'workOrder') {
    await router.push({ path: '/mes/workorders', query: { mode: 'sent' } })
    return
  }
  if (type.value === 'item') {
    await router.push({ path: '/mes/purchase', query: { panel: 'items', view: 'catalog' } })
    return
  }
  if (type.value === 'itemUnit') {
    await router.push({ path: '/mes/purchase', query: { panel: 'itemUnits', view: 'units' } })
    return
  }
  await router.push({ path: '/mes/purchase', query: { panel: 'flows', view: 'flows' } })
}

const confirmLeave = () =>
  new Promise<boolean>((resolve) => {
    if (!hasContent.value) {
      resolve(true)
      return
    }

    Modal.confirm({
      title: canSaveDraft.value ? '离开前保存草稿？' : '放弃未提交内容？',
      content: canSaveDraft.value
        ? '当前填写内容还没有提交，可以先保存为草稿。'
        : '当前内容还没有提交，离开后不会保存。',
      okText: canSaveDraft.value ? (draftId.value ? '更新草稿并退出' : '保存草稿并退出') : '放弃离开',
      cancelText: canSaveDraft.value ? '不保存退出' : '继续编辑',
      okButtonProps: canSaveDraft.value ? {} : { danger: true },
      async onOk() {
        if (!canSaveDraft.value) {
          resolve(true)
          return
        }
        const ok = await saveDraft(false, false)
        resolve(ok)
      },
      onCancel() {
        resolve(canSaveDraft.value)
      },
    })
  })

const handleLeave = async () => {
  if (leaving.value) return
  const ok = await confirmLeave()
  if (ok) await leaveBack()
}

const exitAfterSubmit = async () => {
  leaving.value = true
  if (type.value === 'engineering') {
    await router.push({ path: '/mes/leader', query: { panel: 'engineering', view: 'engineering' } })
    return
  }
  if (type.value === 'process') {
    await router.push({ path: '/mes/processes', query: { panel: 'processes', view: 'processes' } })
    return
  }
  if (type.value === 'workOrder') {
    await router.push({ path: '/mes/workorders', query: { view: 'workOrders' } })
    return
  }
  if (type.value === 'item') {
    await router.push({ path: '/mes/purchase', query: { panel: 'items', view: 'catalog' } })
    return
  }
  if (type.value === 'itemUnit') {
    await router.push({ path: '/mes/purchase', query: { panel: 'itemUnits', view: 'units' } })
    return
  }
  await router.push({ path: '/mes/purchase', query: { panel: 'flows', view: 'purchase' } })
}

const submitDraft = async () => {
  if (!canSaveDraft.value) return
  saving.value = true
  try {
    const ok = await saveDraft(true)
    if (!ok || !draftId.value) return
    const res =
      type.value === 'flow'
        ? await submitInventoryFlow({ id: draftId.value })
        : type.value === 'process'
          ? await submitProcess({ id: draftId.value })
        : type.value === 'engineering'
          ? await submitEngineeringOrder({ id: draftId.value })
          : await submitWorkOrder({ id: draftId.value })
    if (res.data.code !== 0) {
      message.error(res.data.message || '提交失败')
      return
    }
    message.success('已提交')
    await exitAfterSubmit()
  } finally {
    saving.value = false
  }
}

const submit = async () => {
  if (canSaveDraft.value) {
    await submitDraft()
    return
  }

  saving.value = true
  try {
    let id: number | undefined
    if (type.value === 'item') {
      if (!itemForm.name.trim()) {
        message.warning('请输入物品名称')
        return
      }
      const res = await addItem(itemForm)
      id = res.data.data
    } else if (type.value === 'itemUnit') {
      if (!unitForm.itemId) {
        message.warning('请选择物品')
        return
      }
      if (!(await validateUnitEngineeringOrder())) {
        return
      }
      const res = await addItemUnit(unitForm)
      id = res.data.data
    } else if (type.value === 'engineering') {
      if (!engineeringForm.itemId) {
        message.warning('请选择生产物品')
        return
      }
      if (!engineeringForm.processId) {
        message.warning('请选择生产工艺')
        return
      }
      const res = await createEngineeringOrder(engineeringForm)
      id = res.data.data
      if (id) await submitEngineeringOrder({ id })
    } else if (type.value === 'workOrder') {
      if (!workOrderForm.name.trim()) {
        message.warning('请输入工单名称')
        return
      }
      if (!workOrderForm.toUserId) {
        message.warning('请输入接收人')
        return
      }
      const res = await createWorkOrderDraft(workOrderForm)
      id = res.data.data
      if (id) await submitWorkOrder({ id })
    }
    if (id) {
      message.success('已提交')
      await exitAfterSubmit()
    }
  } finally {
    saving.value = false
  }
}

onBeforeRouteLeave(async () => {
  if (leaving.value || saving.value) return true
  return await confirmLeave()
})

watch(
  () => engineeringForm.itemId,
  async (itemId, oldItemId) => {
    if (itemId !== oldItemId) {
      engineeringForm.processId = undefined
      if (selectedEngineeringItem.value?.id !== itemId) {
        selectedEngineeringItem.value = undefined
      }
    }
    if (type.value === 'engineering') {
      await loadEngineeringProcesses()
    }
  },
)

watch(
  () => unitForm.itemId,
  (itemId) => {
    if (!selectedUnitEngineeringOrder.value || !itemId) return
    if (selectedUnitEngineeringOrder.value.itemId && selectedUnitEngineeringOrder.value.itemId !== itemId) {
      unitForm.engineeringOrderId = undefined
      selectedUnitEngineeringOrder.value = undefined
    }
  },
)

onMounted(async () => {
  await loadDraft()
  if (type.value === 'engineering') {
    await loadEngineeringProcesses()
  }
})
</script>

<style scoped>
.create-page {
  max-width: 760px;
}

.create-head {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 18px;
}

.create-head h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.back-button {
  border: 0;
  padding: 0;
  background: transparent;
  color: var(--primary);
  cursor: pointer;
}

.form-surface {
  padding-top: 18px;
  border-top: 1px solid var(--border);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.full-input {
  width: 100%;
}

.create-steps {
  margin-bottom: 22px;
}

.selected-strip {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
  padding: 12px 14px;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  background: #eff6ff;
}

.selected-strip > span:first-child {
  flex: 0 0 auto;
  color: #2563eb;
  font-size: 12px;
  font-weight: 700;
}

.selected-strip > span:first-child::after {
  content: ':';
}

.selected-item-name {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  color: #111827;
  font-size: 14px;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.form-help {
  margin-top: 8px;
  color: var(--muted-foreground);
  font-size: 13px;
}

.wide {
  width: 100%;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.flow-lines {
  display: grid;
  gap: 10px;
}

.flow-line {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 120px 64px;
  gap: 10px;
  align-items: center;
}

.flow-line-item,
.flow-line-quantity {
  width: 100%;
}

@media (max-width: 720px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .flow-line {
    grid-template-columns: 1fr;
  }
}
</style>
