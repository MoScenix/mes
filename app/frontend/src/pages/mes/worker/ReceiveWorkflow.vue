<template>
  <section class="scan-workflow">
    <a-button v-if="receiveFlow" class="scan-back" type="text" @click="backToReceiveScan"
      >返回</a-button
    >
    <div v-if="!receiveFlow" class="workflow-card">
      <strong>领取货物</strong>
    </div>
    <CodeTool
      v-if="!receiveFlow"
      v-model="receiveFlowCode"
      class="workflow-code-tool"
      kind="FLOW"
      auto-open
      scanner-only
      scanner-variant="bare"
      @submit="loadReceiveFlow"
    />
    <template v-else>
      <div class="workflow-card">
        <strong>扫描单体出库</strong>
      </div>
      <CodeTool
        :key="receiveOperationKey"
        v-model="receiveUnitCode"
        class="workflow-code-tool"
        kind="ITEM_UNIT"
        auto-open
        scanner-only
        scanner-variant="bare"
        scanner-display="inline"
        @submit="addReceiveUnit"
      />
    </template>
  </section>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import CodeTool from '@/components/mes/CodeTool.vue'
import { useReceiveWorkflow } from './useReceiveWorkflow'

const props = defineProps<{ flowId?: number }>()

const {
  receiveFlowCode,
  receiveFlow,
  receiveUnitCode,
  receiveOperationKey,
  loadReceiveFlow,
  backToReceiveScan,
  addReceiveUnit,
} = useReceiveWorkflow()

watch(
  () => props.flowId,
  async (flowId) => {
    if (flowId && flowId > 0) {
      await loadReceiveFlow(`MES:FLOW:${flowId}`)
    }
  },
  { immediate: true },
)
</script>

<style scoped>
.scan-workflow {
  position: relative;
  max-width: 560px;
  margin: 0 auto;
  display: grid;
  align-content: center;
  justify-items: center;
  gap: 0;
  min-height: min(620px, calc(100vh - 180px));
}
.scan-back {
  position: fixed;
  top: 72px;
  left: 96px;
  z-index: 20;
}
.workflow-card {
  display: grid;
  justify-items: center;
  gap: 0;
  padding: 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}
.workflow-card > span {
  color: var(--muted-foreground);
  font-size: 13px;
}
.workflow-card strong {
  color: var(--foreground);
  font-size: 42px;
  line-height: 1.25;
  font-weight: 650;
}
.workflow-card p {
  margin: 0;
  color: var(--muted-foreground);
  line-height: 1.6;
}
.workflow-code-tool {
  width: 100%;
}
.workflow-code-tool :deep(.ant-btn-primary),
.workflow-code-tool :deep(button[type='submit']) {
  min-height: 48px;
  font-size: 16px;
  font-weight: 600;
}
</style>
