<template>
  <a-select
    :value="modelValue"
    :label-in-value="false"
    show-search
    :filter-option="false"
    :placeholder="placeholder"
    :options="options"
    :loading="loading"
    @search="handleSearch"
    @change="handleChange"
  />
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { getItem, searchItems, type ItemVO } from '@/api/mesController'

const props = defineProps<{
  modelValue?: number
  placeholder?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value?: number]
  'select-item': [item: ItemVO]
}>()

const loading = ref(false)
const options = ref<{ label: string; value: number; item?: ItemVO }[]>([])
let timer: ReturnType<typeof setTimeout> | undefined

const toOption = (item: ItemVO) => ({
  label: `${item.name || '物品'} #${item.id}`,
  value: item.id || 0,
  item,
})

const handleChange = (value?: number | string) => {
  const id = Number(value) || undefined
  emit('update:modelValue', id)
  const selected = options.value.find((item) => item.value === id)?.item
  if (selected) {
    emit('select-item', selected)
  }
}

const ensureSelectedLabel = async (id?: number) => {
  if (!id || options.value.some((item) => item.value === id)) return
  const res = await getItem({ id })
  if (res.data.code === 0 && res.data.data) {
    options.value = [toOption(res.data.data), ...options.value]
  } else {
    options.value = [{ label: `物品 #${id}`, value: id }, ...options.value]
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
        const res = await getItem({ id: numericId })
        if (res.data.code === 0 && res.data.data) {
          options.value = [toOption(res.data.data)]
          return
        }
        options.value = [{ label: `物品 #${numericId}`, value: numericId }]
        return
      }

      const res = await searchItems({ namePrefix: text, pageNum: 1, pageSize: 10 })
      if (res.data.code === 0) {
        options.value = (res.data.data?.records || []).map(toOption)
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
</script>
