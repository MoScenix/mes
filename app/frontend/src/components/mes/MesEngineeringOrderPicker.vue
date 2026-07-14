<template>
  <a-select
    :value="modelValue"
    :label-in-value="false"
    show-search
    :filter-option="false"
    :placeholder="placeholder"
    :options="options"
    :loading="loading"
    allow-clear
    @search="handleSearch"
    @change="handleChange"
  />
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  DraftStatus,
  MesListScope,
  getEngineeringOrder,
  listEngineeringOrder,
  type EngineeringOrderVO,
} from '@/api/mesController'

const props = defineProps<{
  modelValue?: number
  placeholder?: string
  itemId?: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value?: number]
  'select-order': [order: EngineeringOrderVO]
}>()

const loading = ref(false)
const options = ref<{ label: string; value: number; order?: EngineeringOrderVO }[]>([])
let timer: ReturnType<typeof setTimeout> | undefined

const toOption = (order: EngineeringOrderVO) => {
  const produced = order.producedQuantity ?? 0
  const expected = order.expectedQuantity ?? 0
  const itemText = order.item?.name
    ? `${order.item.name} #${order.itemId}`
    : `物品 #${order.itemId || '-'}`
  return {
    label: `${order.name || '工程单'} #${order.id} / ${itemText} / ${produced}/${expected}`,
    value: order.id || 0,
    order,
  }
}

const acceptOrder = (order: EngineeringOrderVO) => {
  if (props.itemId && order.itemId && props.itemId !== order.itemId) return false
  return true
}

const handleChange = (value?: number | string) => {
  const id = Number(value) || undefined
  emit('update:modelValue', id)
  const selected = options.value.find((item) => item.value === id)?.order
  if (selected) {
    emit('select-order', selected)
  }
}

const ensureSelectedLabel = async (id?: number) => {
  if (!id || options.value.some((item) => item.value === id)) return
  const res = await getEngineeringOrder({ id })
  if (res.data.code === 0 && res.data.data) {
    options.value = [toOption(res.data.data), ...options.value]
  } else {
    options.value = [{ label: `工程单 #${id}`, value: id }, ...options.value]
  }
}

const handleSearch = (query: string) => {
  if (timer) clearTimeout(timer)
  timer = setTimeout(async () => {
    const text = query.trim()
    if (!text) {
      options.value = []
      return
    }

    loading.value = true
    try {
      const numericId = Number(text)
      if (Number.isFinite(numericId) && numericId > 0) {
        const res = await getEngineeringOrder({ id: numericId })
        if (res.data.code === 0 && res.data.data) {
          options.value = [res.data.data].filter(acceptOrder).map(toOption)
          return
        }
        options.value = [{ label: `工程单 #${numericId}`, value: numericId }]
        return
      }

      const res = await listEngineeringOrder({
        namePrefix: text,
        itemId: props.itemId,
        status: DraftStatus.Submitted,
        scope: MesListScope.All,
        pageNum: 1,
        pageSize: 10,
      })
      if (res.data.code === 0) {
        options.value = (res.data.data?.records || []).filter(acceptOrder).map(toOption)
      }
    } finally {
      loading.value = false
    }
  }, 220)
}

watch(
  () => props.modelValue,
  (id) => {
    void ensureSelectedLabel(id)
  },
  { immediate: true },
)

watch(
  () => props.itemId,
  () => {
    if (!props.modelValue) return
    const selected = options.value.find((item) => item.value === props.modelValue)?.order
    if (selected && !acceptOrder(selected)) {
      emit('update:modelValue', undefined)
    }
  },
)
</script>
