<template>
  <a-tooltip title="扫码">
    <a-button type="text" class="scan-btn" aria-label="扫码" @click="openScanner">
      <ScanOutlined />
    </a-button>
  </a-tooltip>
  <a-modal
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
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { ScanOutlined } from '@ant-design/icons-vue'
import { parseMesCode } from '@/utils/mesCode'

type BarcodeDetectorLike = {
  detect: (source: HTMLVideoElement) => Promise<Array<{ rawValue?: string }>>
}

declare global {
  interface Window {
    BarcodeDetector?: new (options?: { formats?: string[] }) => BarcodeDetectorLike
  }
}

const router = useRouter()
const scannerOpen = ref(false)
const scannerMessage = ref('打开摄像头后对准二维码')
const videoRef = ref<HTMLVideoElement>()
const stream = ref<MediaStream>()
const scanning = ref(false)

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
      scannerOpen.value = false
      stopScanner()
      message.success('已扫码')
      const parsed = parseMesCode(rawValue)
      if (parsed.kind && parsed.id) {
        await router.push({ path: '/mes/detail', query: { kind: parsed.kind, id: String(parsed.id) } })
      }
      return
    }
  } catch (error) {
    scannerMessage.value = error instanceof Error ? error.message : '扫码失败'
  }
  requestAnimationFrame(() => scanLoop(detector))
}

const openScanner = async () => {
  if (!window.BarcodeDetector) {
    message.warning('当前浏览器不支持摄像头扫码，请直接输入码')
    return
  }
  scannerOpen.value = true
  scannerMessage.value = '打开摄像头后对准二维码'
  await nextTick()
  try {
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: { ideal: 'environment' } },
      audio: false,
    })
    if (!videoRef.value) return
    videoRef.value.srcObject = stream.value
    await videoRef.value.play()
    scanning.value = true
    scanLoop(new window.BarcodeDetector({ formats: ['qr_code'] }))
  } catch (error) {
    scannerMessage.value = '无法打开摄像头，请检查浏览器权限'
    message.error(error instanceof Error ? error.message : '无法打开摄像头')
    stopScanner()
  }
}

onBeforeUnmount(stopScanner)
</script>

<style scoped>
.scan-btn {
  width: 42px;
  height: 42px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #1d1d1f;
  border: 0;
  background: transparent;
}

.scan-btn :deep(svg) {
  width: 24px;
  height: 24px;
}

.scan-btn:hover {
  color: #0066cc;
  background: transparent;
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
  text-align: center;
}
</style>
