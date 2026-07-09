<template>
  <section class="scan-panel">
    <a-button v-if="scanFlow" class="scan-back" type="text" @click="$emit('back')">返回</a-button>
    <div v-if="!scanFlow" class="mobile-scan-card">
      <strong>扫描入库</strong>
    </div>
    <CodeTool
      v-if="!scanFlow"
      :model-value="flowCode"
      class="mobile-code-tool"
      kind="FLOW"
      auto-open
      scanner-only
      scanner-variant="bare"
      @update:model-value="$emit('update:flowCode', $event)"
      @submit="$emit('load-flow', $event)"
    />
    <template v-else>
      <div class="mobile-scan-card">
        <strong>扫描单体入库</strong>
      </div>
      <CodeTool
        :key="operationKey"
        :model-value="scanValue"
        class="mobile-code-tool"
        kind="ITEM_UNIT"
        auto-open
        scanner-only
        scanner-variant="bare"
        scanner-display="inline"
        @update:model-value="$emit('update:scanValue', $event)"
        @submit="$emit('add-scan-input', $event)"
      />
    </template>
  </section>
</template>

<script setup lang="ts">
import type { InventoryFlowVO } from '@/api/mesController'
import CodeTool from '@/components/mes/CodeTool.vue'

defineProps<{
  flowCode: string
  scanFlow?: InventoryFlowVO
  scanValue: string
  operationKey: number
}>()

defineEmits<{
  (e: 'update:flowCode', value: string): void
  (e: 'update:scanValue', value: string): void
  (e: 'load-flow', value: string): void
  (e: 'add-scan-input', value: string): void
  (e: 'back'): void
}>()
</script>

<style scoped>
.scan-panel {
  position: relative;
  min-height: min(620px, calc(100vh - 180px));
  display: grid;
  align-content: center;
  justify-items: center;
  gap: 0;
  max-width: 520px;
  margin: 0 auto;
}

.scan-back {
  position: fixed;
  top: 72px;
  left: 96px;
  z-index: 20;
}

.mobile-scan-card {
  display: grid;
  justify-items: center;
  gap: 0;
  padding: 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.mobile-scan-card > span,
.scan-context span {
  color: var(--muted-foreground);
  font-size: 13px;
}

.mobile-scan-card strong {
  color: var(--foreground);
  font-size: 42px;
  line-height: 1.25;
  font-weight: 650;
}

.mobile-scan-card p {
  margin: 0;
  color: var(--muted-foreground);
  line-height: 1.6;
}

.mobile-code-tool {
  width: 100%;
}

.mobile-code-tool :deep(.ant-btn-primary),
.mobile-code-tool :deep(button[type='submit']) {
  min-height: 48px;
  font-size: 16px;
  font-weight: 600;
}
</style>
