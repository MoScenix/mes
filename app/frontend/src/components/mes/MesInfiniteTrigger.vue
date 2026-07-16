<template>
  <div ref="trigger" class="infinite-trigger" aria-live="polite">
    <a-spin v-if="loading" size="small" />
    <span>{{ loading ? '正在加载' : hasMore ? '继续滚动加载' : '没有更多了' }}</span>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps<{ hasMore: boolean; loading: boolean }>()
const emit = defineEmits<{ (event: 'load'): void }>()
const trigger = ref<HTMLElement>()
let observer: IntersectionObserver | undefined
let visible = false

const requestLoad = () => {
  if (visible && props.hasMore && !props.loading) emit('load')
}

onMounted(() => {
  observer = new IntersectionObserver(
    ([entry]) => {
      visible = Boolean(entry?.isIntersecting)
      requestLoad()
    },
    { rootMargin: '240px 0px' },
  )
  if (trigger.value) observer.observe(trigger.value)
})
watch(() => [props.hasMore, props.loading], requestLoad)
onBeforeUnmount(() => observer?.disconnect())
</script>

<style scoped>
.infinite-trigger {
  min-height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: var(--muted-foreground);
  font-size: 13px;
}
</style>
