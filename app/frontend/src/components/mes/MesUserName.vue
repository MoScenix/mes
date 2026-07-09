<template>
  <span>{{ label }}</span>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { getUserVoById } from '@/api/userController'

const props = defineProps<{
  id?: number
}>()

const user = ref<API.UserVO>()

const label = computed(() => {
  const id = props.id || user.value?.id
  const name = user.value?.userName || user.value?.userAccount || '用户'
  return id ? `${name} #${id}` : '-'
})

watch(
  () => props.id,
  async (id) => {
    user.value = undefined
    if (!id) return
    const res = await getUserVoById({ id })
    if (res.data.code === 0 && res.data.data) {
      user.value = res.data.data
    }
  },
  { immediate: true },
)
</script>
