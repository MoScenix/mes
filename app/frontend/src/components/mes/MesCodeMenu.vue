<template>
  <a-popover trigger="click" placement="bottomRight">
    <template #content>
      <div class="code-menu">
        <div class="code-tabs">
          <button
            v-for="option in options"
            :key="option.value"
            type="button"
            :class="{ active: activeValue === option.value }"
            @click="activeValue = option.value"
          >
            {{ option.label }}
          </button>
        </div>
        <MesQrCode :value="activeValue" :size="148" />
        <code>{{ activeValue }}</code>
        <a-space>
          <a-button size="small" @click="copy(activeValue)">复制</a-button>
          <a-button size="small" type="primary" @click="openDetail(activeValue)">打开详情</a-button>
        </a-space>
      </div>
    </template>
    <a-button size="small">
      <QrcodeOutlined />
      码
    </a-button>
  </a-popover>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { QrcodeOutlined } from '@ant-design/icons-vue'
import MesQrCode from '@/components/mes/MesQrCode.vue'
import { makeMesCode, parseMesCode, type MesCodeKind } from '@/utils/mesCode'

const props = defineProps<{
  kind: MesCodeKind
  id?: number
}>()

const router = useRouter()

const options = computed(() => {
  return props.id ? [{ label: '对象码', value: makeMesCode(props.kind, props.id) }] : []
})

const activeValue = ref('')

watch(
  options,
  (items) => {
    activeValue.value = items[0]?.value || ''
  },
  { immediate: true },
)

const copy = async (value: string) => {
  if (!value) return
  await navigator.clipboard?.writeText(value)
  message.success('已复制')
}

const openDetail = async (value: string) => {
  const parsed = parseMesCode(value)
  if (parsed.kind && parsed.id) {
    await router.push({
      path: '/mes/detail',
      query: {
        kind: parsed.kind,
        id: String(parsed.id),
      },
    })
    return
  }
  await router.push('/mes/scan')
}
</script>

<style scoped>
.code-menu {
  display: grid;
  justify-items: center;
  gap: 10px;
  max-width: 240px;
}

.code-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
}

.code-tabs button {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  padding: 4px 8px;
  background: #fff;
  color: #1d1d1f;
  font-size: 12px;
  cursor: pointer;
}

.code-tabs button.active {
  border-color: #1677ff;
  color: #1677ff;
  background: #eef5ff;
}

code {
  max-width: 100%;
  overflow-wrap: anywhere;
  color: #1677ff;
  font-size: 12px;
}
</style>
