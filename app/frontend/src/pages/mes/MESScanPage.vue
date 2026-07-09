<template>
  <main class="scan-entry">
    <section class="scan-panel" aria-label="扫码入口">
      <div class="scan-heading">
        <h1>{{ title }}</h1>
      </div>

      <CodeTool
        v-model="scanValue"
        :kind="expectedKind"
        :button-text="buttonText"
        :placeholder="placeholder"
        button-placement="below"
        surface
        @submit="openCode"
      />
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import CodeTool from '@/components/mes/CodeTool.vue'
import { parseMesCode, type MesCodeKind } from '@/utils/mesCode'

const route = useRoute()
const router = useRouter()
const scanValue = ref('')
const codeKinds: MesCodeKind[] = ['FLOW', 'ITEM_UNIT', 'ENGINEERING_ORDER']
const mode = computed(() => String(route.query.mode || 'detail'))
const expectedKind = computed<MesCodeKind | undefined>(() => {
  if (mode.value === 'inbound' || mode.value === 'receive') return 'FLOW'
  if (mode.value === 'inspect') return 'ENGINEERING_ORDER'
  return undefined
})
const title = computed(() => {
  if (mode.value === 'inbound') return '扫描入库'
  if (mode.value === 'receive') return '领取货物'
  if (mode.value === 'inspect') return '检测单体'
  return '扫描 MES 码'
})
const buttonText = computed(() => mode.value === 'detail' ? '打开' : '进入')
const placeholder = computed(() => {
  if (mode.value === 'inbound' || mode.value === 'receive') return '扫描或输入流转单码'
  if (mode.value === 'inspect') return '扫描或输入工程单码'
  return '扫描或输入 MES 码'
})
const openDetail = async (kind: MesCodeKind, id: number) => {
  await router.push({
    path: '/mes/detail',
    query: { kind, id: String(id) },
  })
}

const openCode = async (value: string) => {
  const parsed = parseMesCode(value, expectedKind.value)
  if (!parsed.kind || !parsed.id || (expectedKind.value && parsed.kind !== expectedKind.value)) {
    message.warning('请输入有效的 MES 对象码')
    return
  }
  if (mode.value === 'inbound') {
    await router.push({ path: '/mes/purchase', query: { panel: 'scan', flowId: String(parsed.id) } })
    return
  }
  if (mode.value === 'receive') {
    await router.push({ path: '/mes/worker', query: { panel: 'receive', flowId: String(parsed.id) } })
    return
  }
  if (mode.value === 'inspect') {
    await router.push({ path: '/mes/worker', query: { panel: 'inspect', orderId: String(parsed.id) } })
    return
  }
  await openDetail(parsed.kind, parsed.id)
}

onMounted(async () => {
  const queryKind = String(route.query.kind || '').toUpperCase() as MesCodeKind
  const kind = codeKinds.includes(queryKind) ? queryKind : undefined
  const id = Number(route.query.id || 0)
  if (kind && id > 0) {
    await openDetail(kind, id)
  }
})
</script>

<style scoped>
.scan-entry {
  min-height: calc(100vh - 64px - 72px);
  display: grid;
  align-items: start;
  justify-items: center;
  padding: 44px 16px;
  background: #ffffff;
}

.scan-panel {
  width: min(100%, 680px);
  display: grid;
  gap: 28px;
  padding-top: min(10vh, 92px);
}

.scan-heading {
  display: grid;
  gap: 6px;
  text-align: center;
}

.scan-heading h1 {
  margin: 0;
  color: #1d1d1f;
  font-size: 42px;
  line-height: 1.12;
  font-weight: 600;
  letter-spacing: 0;
}

@media (max-width: 768px) {
  .scan-entry {
    padding-top: 40px;
  }

  .scan-heading h1 {
    font-size: 34px;
  }
}
</style>
