<template>
  <div class="qr-code" :style="{ '--qr-size': `${size}px` }">
    <img v-if="qrUrl" :src="qrUrl" :alt="value" />
    <span v-else>生成中</span>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import QRCode from 'qrcode'

const props = withDefaults(
  defineProps<{
    value: string
    size?: number
  }>(),
  {
    size: 136,
  },
)

const qrUrl = ref('')

watch(
  () => [props.value, props.size] as const,
  async ([value, size]) => {
    qrUrl.value = value
      ? await QRCode.toDataURL(value, {
          width: size,
          margin: 1,
          errorCorrectionLevel: 'M',
          color: {
            dark: '#1d1d1f',
            light: '#ffffff',
          },
        })
      : ''
  },
  { immediate: true },
)
</script>

<style scoped>
.qr-code {
  width: fit-content;
  padding: 8px;
  border: 1px solid #d2d2d7;
  border-radius: 8px;
  background: #ffffff;
}

.qr-code img {
  display: block;
  width: var(--qr-size);
  height: var(--qr-size);
}

.qr-code span {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: var(--qr-size);
  height: var(--qr-size);
  color: #7a7a7a;
  font-size: 13px;
}
</style>
