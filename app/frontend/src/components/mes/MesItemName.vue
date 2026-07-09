<template>
  <span>{{ label }}</span>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { getItem, type ItemVO } from '@/api/mesController'

const props = defineProps<{
  id?: number
  item?: ItemVO
}>()

const loaded = ref<ItemVO>()

const label = computed(() => {
  const id = props.id || props.item?.id || loaded.value?.id
  const name = props.item?.name || loaded.value?.name || '物品'
  return id ? `${name} #${id}` : '-'
})

watch(
  () => [props.id, props.item?.id, props.item?.name] as const,
  async () => {
    loaded.value = undefined
    const id = props.id || props.item?.id
    if (!id || props.item?.name) return
    const res = await getItem({ id })
    if (res.data.code === 0 && res.data.data) {
      loaded.value = res.data.data
    }
  },
  { immediate: true },
)
</script>
