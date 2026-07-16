<template>
  <label class="mes-list-search">
    <input
      :value="innerValue"
      :placeholder="placeholder"
      type="search"
      @input="updateValue(($event.target as HTMLInputElement).value)"
      @keydown.enter.prevent="submitCurrent()"
    />
  </label>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import type { ItemVO } from '@/api/mesController'

const props = defineProps<{
  modelValue?: string
  placeholder?: string
  itemSearch?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  search: [value: string]
  clear: []
  'select-item': [item: ItemVO]
}>()

const innerValue = ref(props.modelValue || '')
let searchTimer: ReturnType<typeof setTimeout> | undefined

const updateValue = (value: string) => {
  innerValue.value = value
  emit('update:modelValue', value)
  if (searchTimer) clearTimeout(searchTimer)
  if (!value.trim()) {
    emit('clear')
    return
  }
  searchTimer = setTimeout(() => submitCurrent(value), 300)
}

const submitCurrent = (value?: string) => {
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = undefined
  }
  const text = String(value ?? innerValue.value).trim()
  innerValue.value = text
  emit('update:modelValue', text)
  if (!text) {
    emit('clear')
    return
  }
  emit('search', text)
}

onBeforeUnmount(() => {
  if (searchTimer) clearTimeout(searchTimer)
})

watch(
  () => props.modelValue,
  (value) => {
    if (value !== innerValue.value) {
      innerValue.value = value || ''
    }
  },
)
</script>

<style scoped>
.mes-list-search {
  width: 280px;
  max-width: 100%;
  height: 38px;
  display: block;
  padding: 0 16px;
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 16px;
  background: rgba(248, 250, 252, 0.92);
  color: #64748b;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.7),
    0 8px 24px rgba(15, 23, 42, 0.06);
  transition:
    border-color 0.16s ease,
    box-shadow 0.16s ease,
    background 0.16s ease;
}

.mes-list-search:focus-within {
  border-color: rgba(59, 130, 246, 0.42);
  background: #ffffff;
  box-shadow:
    0 0 0 3px rgba(59, 130, 246, 0.1),
    0 10px 28px rgba(15, 23, 42, 0.08);
}

.mes-list-search input {
  width: 100%;
  height: 100%;
  min-width: 0;
  border: 0;
  outline: none;
  background: transparent;
  color: #0f172a;
  font-size: 14px;
  line-height: 38px;
}

.mes-list-search input::placeholder {
  color: #94a3b8;
}
</style>
