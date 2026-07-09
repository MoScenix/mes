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
    @change="emit('update:modelValue', Number($event) || undefined)"
  />
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { getUserVoById, listUserVoByPage } from '@/api/userController'

const props = defineProps<{
  modelValue?: number
  placeholder?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value?: number]
}>()

const loading = ref(false)
const options = ref<{ label: string; value: number }[]>([])
let timer: ReturnType<typeof setTimeout> | undefined

const toOption = (user: API.UserVO) => ({
  label: `${user.userName || user.userAccount || '用户'} #${user.id}`,
  value: user.id || 0,
})

const ensureSelectedLabel = async (id?: number) => {
  if (!id || options.value.some((item) => item.value === id)) return
  const res = await getUserVoById({ id })
  if (res.data.code === 0 && res.data.data) {
    options.value = [toOption(res.data.data), ...options.value]
  } else {
    options.value = [{ label: `用户 #${id}`, value: id }, ...options.value]
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
        const res = await getUserVoById({ id: numericId })
        if (res.data.code === 0 && res.data.data) {
          options.value = [toOption(res.data.data)]
          return
        }
        options.value = [{ label: `用户 #${numericId}`, value: numericId }]
        return
      }

      const [nameRes, accountRes] = await Promise.all([
        listUserVoByPage({ userName: text, pageNum: 1, pageSize: 8 }),
        listUserVoByPage({ userAccount: text, pageNum: 1, pageSize: 8 }),
      ])
      const seen = new Set<number>()
      const users = [...(nameRes.data.data?.records || []), ...(accountRes.data.data?.records || [])]
      options.value = users.filter((user) => {
        const id = user.id || 0
        if (!id || seen.has(id)) return false
        seen.add(id)
        return true
      }).map(toOption)
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
