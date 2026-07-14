<template>
  <div
    class="code-tool"
    :class="{
      'surface-mode': surface,
      'scanner-compact-mode': scannerOnly && scannerVariant === 'compact',
      'scanner-bare-mode': scannerOnly && scannerVariant === 'bare',
    }"
  >
    <div v-if="scannerOnly && scannerDisplay === 'inline'" class="inline-scanner-box">
      <video ref="videoRef" muted playsinline></video>
      <p>{{ scannerMessage }}</p>
    </div>

    <button
      v-else-if="scannerOnly"
      type="button"
      class="scanner-only-trigger"
      aria-label="打开扫码"
      @click="openScanner"
    >
      <span v-if="scannerVariant !== 'bare'" class="scanner-only-fill"></span>
      <span class="scanner-only-icon">
        <ScanOutlined />
      </span>
    </button>

    <div v-else class="code-input-shell">
      <a-input-group v-if="buttonPlacement === 'inline'" compact class="code-input-row">
        <a-input-search
          :id="inputId"
          :name="inputId"
          :value="modelValue"
          :placeholder="placeholder || defaultPlaceholder"
          :enter-button="buttonText"
          :size="size"
          :loading="loading"
          @focus="focused = true"
          @blur="focused = false"
          @update:value="updateValue"
          @search="handleSearch"
        />
        <a-tooltip title="扫码">
          <a-button :size="size" class="scan-button" @click="openScanner" aria-label="扫码">
            <ScanOutlined />
          </a-button>
        </a-tooltip>
      </a-input-group>
      <div v-else class="code-command">
        <div class="code-command-input">
          <a-input
            :id="inputId"
            :name="inputId"
            :value="modelValue"
            :placeholder="placeholder || defaultPlaceholder"
            :size="size"
            :disabled="loading"
            @focus="focused = true"
            @blur="focused = false"
            @update:value="updateValue"
            @pressEnter="handleSearch(modelValue)"
          />
          <a-tooltip title="扫码">
            <a-button
              :size="size"
              class="scan-button command-scan-button"
              @click="openScanner"
              aria-label="扫码"
            >
              <ScanOutlined />
            </a-button>
          </a-tooltip>
        </div>
        <a-button
          type="primary"
          :size="size"
          block
          :loading="loading"
          class="code-command-submit"
          @click="handleSearch(modelValue)"
        >
          {{ buttonText }}
        </a-button>
      </div>

      <div v-if="showSuggestions" class="code-suggestions">
        <button
          v-for="item in suggestions"
          :key="item.value"
          type="button"
          @mousedown.prevent="chooseSuggestion(item.value)"
        >
          <span>{{ item.label }}</span>
          <code>{{ item.value }}</code>
        </button>
      </div>
    </div>

    <div v-if="showGenerated && generatedCode" class="generated-code">
      <MesQrCode :value="generatedCode" />
      <button type="button" @click="copyCode">{{ generatedCode }}</button>
    </div>
    <a-modal
      v-if="scannerDisplay === 'modal'"
      v-model:open="scannerOpen"
      title="扫码"
      :footer="null"
      width="420px"
      @cancel="stopScanner"
    >
      <div class="scanner-box">
        <video ref="videoRef" muted playsinline></video>
        <p>{{ scannerMessage }}</p>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import { ScanOutlined } from '@ant-design/icons-vue'
import { makeMesCode, parseMesCode, type MesCodeKind, type ParsedMesCode } from '@/utils/mesCode'
import MesQrCode from '@/components/mes/MesQrCode.vue'

type BarcodeDetectorLike = {
  detect: (source: HTMLVideoElement) => Promise<Array<{ rawValue?: string }>>
}

declare global {
  interface Window {
    BarcodeDetector?: new (options?: { formats?: string[] }) => BarcodeDetectorLike
  }
}

const props = withDefaults(
  defineProps<{
    modelValue: string
    kind?: MesCodeKind
    generatedId?: number
    placeholder?: string
    buttonText?: string
    size?: 'small' | 'middle' | 'large'
    showGenerated?: boolean
    loading?: boolean
    buttonPlacement?: 'inline' | 'below'
    surface?: boolean
    autoOpen?: boolean
    scannerOnly?: boolean
    scannerVariant?: 'spotlight' | 'compact' | 'bare'
    scannerDisplay?: 'modal' | 'inline'
  }>(),
  {
    buttonText: '加入',
    size: 'large',
    showGenerated: false,
    buttonPlacement: 'inline',
    autoOpen: false,
    scannerOnly: false,
    scannerVariant: 'spotlight',
    scannerDisplay: 'modal',
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
  submit: [value: string]
  parsed: [value: ParsedMesCode]
}>()

const defaultPlaceholder = computed(() => {
  if (!props.kind) return '输入 id 或 MES 码'
  return `输入 id，或 MES:${props.kind}:123`
})

const generatedCode = computed(() => (props.kind ? makeMesCode(props.kind, props.generatedId) : ''))
const inputId = `mes-code-input-${Math.random().toString(36).slice(2)}`
const scannerOpen = ref(false)
const scannerMessage = ref('打开摄像头后对准二维码')
const videoRef = ref<HTMLVideoElement>()
const stream = ref<MediaStream>()
const scanning = ref(false)
const focused = ref(false)
const buttonPlacement = computed(() => props.buttonPlacement)

const kindLabel: Record<MesCodeKind, string> = {
  FLOW: '流转单',
  ITEM_UNIT: '库存单体',
  ENGINEERING_ORDER: '工程单',
}

const suggestions = computed(() => {
  const value = props.modelValue.trim()
  if (!value) return []
  const parsed = parseMesCode(value, props.kind)
  if (parsed.kind && parsed.id && (!props.kind || parsed.kind === props.kind)) {
    return [
      {
        label: `打开${kindLabel[parsed.kind]} #${parsed.id}`,
        value: makeMesCode(parsed.kind, parsed.id),
      },
    ]
  }
  return []
})

const showSuggestions = computed(() => focused.value && suggestions.value.length > 0)

const updateValue = (value: string) => {
  emit('update:modelValue', value)
}

const handleSearch = (value: string) => {
  emit('parsed', parseMesCode(value, props.kind))
  emit('submit', value)
}

const chooseSuggestion = (value: string) => {
  emit('update:modelValue', value)
  handleSearch(value)
}

const copyCode = async () => {
  if (!generatedCode.value) return
  await navigator.clipboard?.writeText(generatedCode.value)
  message.success('已复制')
}

const stopScanner = () => {
  scanning.value = false
  stream.value?.getTracks().forEach((track) => track.stop())
  stream.value = undefined
}

const scanLoop = async (detector: BarcodeDetectorLike) => {
  if (!scanning.value || !videoRef.value) return
  try {
    const codes = await detector.detect(videoRef.value)
    const rawValue = codes.find((code) => code.rawValue)?.rawValue
    if (rawValue) {
      emit('update:modelValue', rawValue)
      handleSearch(rawValue)
      if (props.scannerDisplay === 'modal') {
        scannerOpen.value = false
      }
      stopScanner()
      message.success('已扫码')
      return
    }
  } catch (error) {
    scannerMessage.value = error instanceof Error ? error.message : '扫码失败'
  }
  requestAnimationFrame(() => scanLoop(detector))
}

const openScanner = async () => {
  const supportsDetector = Boolean(window.BarcodeDetector)
  scannerOpen.value = props.scannerDisplay === 'modal'
  scannerMessage.value = supportsDetector
    ? '打开摄像头后对准二维码'
    : '当前浏览器不支持自动识别，请检查浏览器支持'
  await nextTick()
  try {
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: { ideal: 'environment' } },
      audio: false,
    })
    if (!videoRef.value) return
    videoRef.value.srcObject = stream.value
    await videoRef.value.play()
    if (supportsDetector && window.BarcodeDetector) {
      scanning.value = true
      scanLoop(new window.BarcodeDetector({ formats: ['qr_code'] }))
    }
  } catch (error) {
    scannerMessage.value = '无法打开摄像头，请检查浏览器权限'
    message.error(error instanceof Error ? error.message : '无法打开摄像头')
    stopScanner()
  }
}

onBeforeUnmount(stopScanner)

onMounted(() => {
  if (props.autoOpen) {
    openScanner()
  }
})
</script>

<style scoped>
.code-tool {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.surface-mode {
  gap: 0;
  padding: 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.code-input-shell {
  position: relative;
}

.scanner-only-trigger {
  width: 100%;
  min-height: 96px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 56px;
  align-items: center;
  gap: 14px;
  border: 1px solid #dfe3ea;
  border-radius: 34px;
  padding: 18px 20px 18px 28px;
  background: #ffffff;
  color: #1d1d1f;
  box-shadow: 0 24px 70px rgba(15, 23, 42, 0.12);
  cursor: pointer;
  transition:
    border-color 0.16s ease,
    box-shadow 0.16s ease,
    transform 0.16s ease;
}

.scanner-compact-mode .scanner-only-trigger {
  min-height: 64px;
  grid-template-columns: minmax(0, 1fr) 48px;
  border-radius: 22px;
  padding: 8px 10px 8px 20px;
  box-shadow: none;
}

.scanner-compact-mode .scanner-only-fill {
  height: 32px;
}

.scanner-compact-mode .scanner-only-icon {
  width: 48px;
  height: 48px;
}

.scanner-bare-mode {
  align-items: center;
}

.scanner-bare-mode .scanner-only-trigger {
  display: none;
}

.inline-scanner-box {
  width: min(100%, 520px);
  min-height: 460px;
  display: grid;
  gap: 12px;
  align-content: center;
  justify-items: center;
}

.inline-scanner-box video {
  width: min(78vw, 420px);
  aspect-ratio: 1;
  border-radius: 18px;
  background: #111827;
  object-fit: cover;
}

.inline-scanner-box p {
  margin: 0;
  color: #6b7280;
  font-size: 13px;
}

.scanner-only-trigger:hover {
  border-color: #94a3b8;
  box-shadow: 0 28px 84px rgba(15, 23, 42, 0.16);
  transform: translateY(-1px);
}

.scanner-only-fill {
  min-width: 0;
  height: 44px;
}

.scanner-only-icon {
  width: 56px;
  height: 56px;
  display: grid;
  place-items: center;
  border: 1px solid #e5e7eb;
  border-radius: 999px;
  background: #f8fafc;
  color: #334155;
}

.scanner-only-icon :deep(svg) {
  width: 22px;
  height: 22px;
}

.code-input-row {
  display: flex;
}

.code-input-row :deep(.ant-input-search) {
  flex: 1;
}

.scan-button {
  width: 44px;
}

.code-command {
  display: grid;
  gap: 14px;
}

.code-command-input {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 52px;
  align-items: stretch;
  overflow: hidden;
  border: 1px solid #dfe3ea;
  border-radius: 999px;
  background: #ffffff;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
  transition:
    border-color 0.16s ease,
    box-shadow 0.16s ease;
}

.surface-mode .code-command-input {
  border: 1px solid #dfe1e5;
  border-radius: 24px;
  background: #ffffff;
  box-shadow: none;
}

.surface-mode .code-command-input :deep(.ant-input) {
  height: 46px;
  padding-inline: 18px 8px;
  font-size: 16px;
}

.surface-mode .command-scan-button {
  width: 48px;
  height: 46px;
  color: #5f6368;
}

.surface-mode .command-scan-button :deep(svg) {
  width: 18px;
  height: 18px;
}

.surface-mode .code-command-submit {
  display: none;
}

.surface-mode .code-command-input:hover,
.surface-mode .code-command-input:focus-within {
  border-color: #dfe1e5;
  box-shadow: 0 1px 6px rgba(32, 33, 36, 0.18);
}

.code-command-input:focus-within {
  border-color: #94a3b8;
  box-shadow: 0 12px 34px rgba(15, 23, 42, 0.12);
}

.code-command-input :deep(.ant-input) {
  height: 56px;
  border: 0;
  border-radius: 0;
  padding-inline: 22px 10px;
  background: transparent;
  font-size: 16px;
  box-shadow: none;
}

.code-command-input :deep(.ant-input:focus) {
  box-shadow: none;
}

.command-scan-button {
  width: 52px;
  height: 56px;
  border: 0;
  border-radius: 0;
  background: transparent;
  color: #475569;
}

.command-scan-button:hover {
  background: #f8fafc;
  color: #1677ff;
}

.surface-mode .command-scan-button:hover {
  background: rgba(15, 23, 42, 0.04);
}

.code-command-submit {
  height: 44px;
  border-radius: 999px;
  font-weight: 600;
}

.code-suggestions {
  position: absolute;
  z-index: 20;
  top: calc(100% + 6px);
  left: 0;
  right: 44px;
  display: grid;
  overflow: hidden;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 14px 30px rgba(0, 0, 0, 0.08);
}

.code-suggestions button {
  min-width: 0;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  border: 0;
  border-bottom: 1px solid #f0f0f0;
  padding: 10px 12px;
  background: transparent;
  color: #1d1d1f;
  text-align: left;
  cursor: pointer;
}

.code-suggestions button:last-child {
  border-bottom: 0;
}

.code-suggestions button:hover {
  background: #f5f5f7;
}

.code-suggestions span {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.code-suggestions code {
  color: #7a7a7a;
  font-size: 12px;
}

.generated-code {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
}

.generated-code button {
  max-width: 100%;
  border: 1px solid #d2d2d7;
  border-radius: 8px;
  padding: 6px 10px;
  background: #ffffff;
  color: #0066cc;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  overflow-wrap: anywhere;
  cursor: pointer;
}

.scanner-box {
  display: grid;
  gap: 12px;
}

.scanner-box video {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 8px;
  background: #1d1d1f;
  object-fit: cover;
}

.scanner-box p {
  margin: 0;
  color: #7a7a7a;
  font-size: 13px;
}
</style>
