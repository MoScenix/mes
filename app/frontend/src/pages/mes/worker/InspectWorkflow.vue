<template>
  <section class="scan-workflow">
    <a-button v-if="inspectOrder" class="scan-back" type="text" @click="backToInspectScan">返回</a-button>
    <div v-if="!inspectOrder" class="workflow-card">
      <strong>检测单体</strong>
    </div>
    <CodeTool
      v-if="!inspectOrder"
      v-model="inspectOrderCode"
      kind="ENGINEERING_ORDER"
      auto-open
      scanner-only
      scanner-variant="bare"
      @submit="loadInspectOrder"
    />
    <template v-else>
      <div class="workflow-context">
        <strong>工程单 #{{ inspectOrder.id }}</strong>
        <MesItemName :id="inspectOrder.itemId" :item="inspectOrder.item" />
      </div>
      <CodeTool
        :key="inspectOperationKey"
        v-model="inspectUnitCode"
        kind="ITEM_UNIT"
        auto-open
        scanner-only
        scanner-variant="bare"
        scanner-display="inline"
        @submit="loadInspectUnit"
      />
      <div v-if="inspectUnit" class="inspect-result">
        <span>单体 #{{ inspectUnit.id }}</span>
        <a-space>
          <a-button :loading="inspectSubmitting" @click="submitInspect(QUALITY_STATUS_QUALIFIED)">合格</a-button>
          <a-button danger :loading="inspectSubmitting" @click="submitInspect(QUALITY_STATUS_UNQUALIFIED)">不合格</a-button>
        </a-space>
      </div>
    </template>
  </section>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import {
  QUALITY_STATUS_QUALIFIED,
  QUALITY_STATUS_UNQUALIFIED,
} from '@/api/mesController'
import CodeTool from '@/components/mes/CodeTool.vue'
import MesItemName from '@/components/mes/MesItemName.vue'
import { useInspectWorkflow } from './useInspectWorkflow'

const props = defineProps<{ orderId?: number }>()

const {
  inspectOrderCode,
  inspectOrder,
  inspectUnitCode,
  inspectUnit,
  inspectSubmitting,
  inspectOperationKey,
  loadInspectOrder,
  backToInspectScan,
  loadInspectUnit,
  submitInspect,
} = useInspectWorkflow()

watch(
  () => props.orderId,
  async (orderId) => {
    if (orderId && orderId > 0) {
      await loadInspectOrder(`MES:ENGINEERING_ORDER:${orderId}`)
    }
  },
  { immediate: true },
)
</script>

<style scoped>
.scan-workflow { position: relative; max-width: 560px; margin: 0 auto; display: grid; align-content: center; justify-items: center; gap: 0; min-height: min(620px, calc(100vh - 180px)); }
.scan-back { position: fixed; top: 72px; left: 96px; z-index: 20; }
.workflow-card { display: grid; justify-items: center; gap: 0; padding: 0; border: 0; border-radius: 0; background: transparent; box-shadow: none; }
.workflow-card > span,
.workflow-context span { color: var(--muted-foreground); font-size: 13px; }
.workflow-card strong { color: var(--foreground); font-size: 42px; line-height: 1.25; font-weight: 650; }
.workflow-card p { margin: 0; color: var(--muted-foreground); line-height: 1.6; }
.inspect-result { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
</style>
