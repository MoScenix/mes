<template>
  <div class="assistant-root" :class="{ open, 'page-mode': pageMode }">
    <!-- Floating mode: ball trigger -->
    <button v-if="!pageMode && !open" class="assistant-ball" type="button" @click="toggle">
      <svg class="assistant-gpt-logo" viewBox="0 0 48 48" aria-hidden="true">
        <path
          d="M13.5 14.5h21A5.5 5.5 0 0 1 40 20v9a5.5 5.5 0 0 1-5.5 5.5H26l-7.5 5v-5h-5A5.5 5.5 0 0 1 8 29v-9a5.5 5.5 0 0 1 5.5-5.5Z"
          fill="none"
          stroke="currentColor"
          stroke-width="3"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
        <circle cx="18" cy="24.5" r="1.8" fill="currentColor" />
        <circle cx="24" cy="24.5" r="1.8" fill="currentColor" />
        <circle cx="30" cy="24.5" r="1.8" fill="currentColor" />
      </svg>
      <span class="sr-only">打开 MES 助手</span>
    </button>

    <!-- Chat window -->
    <section v-else-if="open" class="assistant-window">
      <div class="chat-pane">
        <header class="assistant-head">
          <div>
            <strong>MES 助手</strong>
            <span>{{ activeTitle }}</span>
          </div>
          <div class="head-actions">
            <a-tooltip title="新建对话">
              <a-button type="text" shape="circle" size="small" :loading="creating" @click="createNewChat">
                <EditOutlined />
              </a-button>
            </a-tooltip>
            <a-tooltip title="历史记录">
              <a-button type="text" shape="circle" size="small" @click="toggleHistory">
                <HistoryOutlined />
              </a-button>
            </a-tooltip>
            <a-button v-if="!pageMode" type="text" shape="circle" size="small" @click="open = false">
              <CloseOutlined />
            </a-button>
          </div>
        </header>

        <aside v-if="historyOpen" class="history-panel">
          <div class="history-title">
            <span>历史记录</span>
            <a-button size="small" type="text" :loading="loadingSessions" @click="loadSessions">刷新</a-button>
          </div>
          <div class="history-cards custom-scrollbar">
            <template v-if="sessions.length">
              <div
                v-for="session in sessions"
                :key="session.id"
                class="history-card-item"
                :class="{ active: Number(activeAppId) === session.id }"
              >
                <button class="history-card-main" type="button" @click="selectSession(session)">
                  <span class="history-card-title">{{ session.appName || '未命名对话' }}</span>
                  <time class="history-card-time">{{ session.updateTime || session.createTime || '' }}</time>
                </button>
                <a-popconfirm title="确定要删除这条对话吗？" @confirm="deleteSession(session)">
                  <a-tooltip title="删除">
                    <a-button
                      class="history-delete"
                      danger
                      type="text"
                      shape="circle"
                      size="small"
                      :loading="deletingSessionId === session.id"
                      @click.stop
                    >
                      <DeleteOutlined />
                    </a-button>
                  </a-tooltip>
                </a-popconfirm>
              </div>
            </template>
            <a-empty v-else-if="!loadingSessions" description="暂无历史" />
            <a-spin v-else class="loading-spin" />
          </div>
        </aside>

        <main ref="messagesContainer" class="chat-list custom-scrollbar">
          <div v-if="!messages.length" class="empty-chat">
            <div class="empty-title">今天要处理什么？</div>
            <div class="empty-subtitle">询问工单、库存、工程单或流转单，我会按你的角色调用工具。</div>
          </div>

        <template v-for="(item, index) in messages" :key="item.id || index">
          <div v-if="item.type === 'system'" class="system-row">
            <span>{{ item.content }}</span>
          </div>

          <div v-else-if="item.type === 'user'" class="message-row user-row">
            <div class="user-message">
              <div v-if="item.isPush" class="push-label">push</div>
              <div
                v-if="item.isFile && item.fileMeta"
                class="file-message"
                @dblclick="openFileMessage(item.fileMeta)"
              >
                <FileTextOutlined class="file-icon" />
                <div class="file-main">
                  <div class="file-name">{{ item.fileMeta.filename }}</div>
                  <div class="file-meta">
                    <span>{{ formatFileSize(item.fileMeta.size) }}</span>
                    <span>{{ item.fileMeta.isBig ? '已分块' : '已解析' }}</span>
                    <span v-if="item.fileMeta.parentCount">{{ item.fileMeta.parentCount }} 个父块</span>
                  </div>
                </div>
              </div>
              <MarkdownRenderer v-else class="message-markdown user-markdown" :content="item.content.trimEnd()" />
            </div>
          </div>

          <div v-else class="message-row assistant-row">
            <div class="assistant-message">
              <div v-if="item.agent" class="agent-label">{{ item.agent }}</div>
              <template v-if="item.parts && item.parts.length">
                <template v-for="part in item.parts" :key="part.id">
                  <MarkdownRenderer
                    v-if="part.type === 'text' && part.content"
                    class="message-markdown assistant-markdown"
                    :content="part.content.trimEnd()"
                  />
                  <div v-else-if="part.type === 'tool'" class="tool-list">
                    <details class="tool-item">
                      <summary>
                        <span class="tool-corner"></span>
                        <span class="tool-status">{{ part.tool.status === 'running' ? 'running' : part.tool.status === 'error' ? 'failed' : 'ran' }}</span>
                        <span class="tool-name">{{ part.tool.name }}</span>
                      </summary>
                      <pre v-if="part.tool.args" class="tool-code">{{ formatJSON(part.tool.args) }}</pre>
                      <pre v-if="part.tool.result" class="tool-code">{{ part.tool.result }}</pre>
                    </details>
                  </div>
                </template>
              </template>
              <template v-else>
                <MarkdownRenderer v-if="item.content" class="message-markdown assistant-markdown" :content="item.content.trimEnd()" />
                <div v-if="item.toolCalls && item.toolCalls.length" class="tool-list">
                  <details v-for="tool in item.toolCalls" :key="tool.id" class="tool-item">
                    <summary>
                      <span class="tool-corner"></span>
                      <span class="tool-status">{{ tool.status === 'running' ? 'running' : tool.status === 'error' ? 'failed' : 'ran' }}</span>
                      <span class="tool-name">{{ tool.name }}</span>
                    </summary>
                    <pre v-if="tool.args" class="tool-code">{{ formatJSON(tool.args) }}</pre>
                    <pre v-if="tool.result" class="tool-code">{{ formatJSON(tool.result) }}</pre>
                  </details>
                </div>
              </template>
              <div v-if="item.loading" class="assistant-loading">
                <span class="loading-dot"></span>
                <span>{{ aiState?.status || 'unknown' }}</span>
              </div>
            </div>
          </div>
        </template>
        </main>

        <footer class="composer-wrap">
          <div v-if="currentQuestion && currentQuestionItem" class="question-panel">
          <div class="question-title">
            <span>{{ currentQuestion.agent || 'AI' }} 需要确认</span>
            <span v-if="currentQuestionItems.length > 1" class="question-count">
              {{ currentQuestionIndex + 1 }} / {{ currentQuestionItems.length }}
            </span>
          </div>
          <div class="question-content">{{ currentQuestionItem.question }}</div>
          <div v-if="currentQuestionItem.options.length" class="question-options">
            <button
              v-for="option in currentQuestionItem.options"
              :key="option"
              type="button"
              class="question-option"
              :class="{ 'question-option-active': currentAnswerSelection === option }"
              @click="selectAnswerOption(option)"
            >
              {{ option }}
            </button>
          </div>
          <a-textarea
            v-model:value="answerInput"
            :rows="2"
            placeholder="输入其他回答，Enter 继续"
            class="answer-input"
            @keydown.enter.prevent="submitQuestionStep"
          />
          <div class="question-actions">
            <a-button size="small" @click="currentQuestion = null">稍后</a-button>
            <a-button
              type="primary"
              size="small"
              :disabled="!canSubmitAnswer || answeringQuestion"
              :loading="answeringQuestion"
              @click="submitQuestionStep"
            >
              {{ isLastQuestion ? '发送' : '继续' }}
            </a-button>
          </div>
          </div>

          <PromptInputBox
            v-if="!currentQuestion"
            v-model="userInput"
            :is-loading="isGenerating"
            :is-submitting="sendingMessage || creating"
            :placeholder="activeAppId ? '继续和 MES 助手对话' : '发送后创建一条新的助手对话'"
            @send="(msg: string, files?: File[]) => handleSendMessage(msg, files)"
            @cancel="cancelCurrentTask"
          />
        </footer>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  CloseOutlined,
  DeleteOutlined,
  EditOutlined,
  FileTextOutlined,
  HistoryOutlined,
} from '@ant-design/icons-vue'
import { addApp, deleteApp, listMyAppVoByPage } from '@/api/appController'
import { listAppChatHistory } from '@/api/chatHistoryController'
import { cancelAI } from '@/api/aiController'
import request from '@/request'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import PromptInputBox from '@/components/PromptInputBox.vue'
import { useAIEvents, type AIFileMeta, type AIMessage } from '@/composables/useAIEvents'

const props = withDefaults(defineProps<{ pageMode?: boolean }>(), {
  pageMode: false,
})

const router = useRouter()
const open = ref(props.pageMode)
const creating = ref(false)
const sendingMessage = ref(false)
const answeringQuestion = ref(false)
const userInput = ref('')
const activeAppId = ref<number>()
const activeTitle = ref('选择历史记录或直接提问')
const sessions = ref<API.AppVO[]>([])
const historyOpen = ref(false)
const loadingSessions = ref(false)
const deletingSessionId = ref<number>()
const messagesContainer = ref<HTMLElement>()
const currentQuestionIndex = ref(0)
const answerInput = ref('')
const answerSelections = ref<Record<number, string>>({})


const {
  messages,
  aiState,
  isGenerating,
  currentQuestion,
  sendMessage,
  pushMessage,
  answerQuestion,
  loadInitialState,
  stop,
} = useAIEvents(activeAppId, { onDone: loadSessions })

const currentQuestionItems = computed(() => {
  if (!currentQuestion.value) return []
  return currentQuestion.value.questions?.length
    ? currentQuestion.value.questions
    : [{ question: safeText(currentQuestion.value.content), options: [] }]
})

const currentQuestionItem = computed(() => currentQuestionItems.value[currentQuestionIndex.value])

const currentAnswerSelection = computed(() => answerSelections.value[currentQuestionIndex.value] || '')

const isLastQuestion = computed(() => currentQuestionIndex.value >= currentQuestionItems.value.length - 1)

const canSubmitAnswer = computed(() => Boolean(answerInput.value.trim() || currentAnswerSelection.value.trim()))

const toggle = async () => {
  if (isMobileViewport()) {
    await router.push('/mes/assistant')
    return
  }
  open.value = !open.value
  if (open.value) {
    await loadSessions()
    await nextTick(scrollToBottom)
  }
}

const isMobileViewport = () => window.matchMedia?.('(max-width: 768px)').matches || window.innerWidth <= 768

async function loadSessions() {
  loadingSessions.value = true
  try {
    const res = await listMyAppVoByPage({ pageNum: 1, pageSize: 30 })
    if (res.data.code === 0 && res.data.data?.records) {
      sessions.value = [...res.data.data.records].sort((a, b) => {
        const at = new Date(a.updateTime || a.createTime || '').getTime()
        const bt = new Date(b.updateTime || b.createTime || '').getTime()
        return bt - at
      })
    }
  } finally {
    loadingSessions.value = false
  }
}

const toggleHistory = async () => {
  historyOpen.value = !historyOpen.value
  if (historyOpen.value) {
    await loadSessions()
  }
}

const selectSession = async (idOrSession: number | API.AppVO) => {
  const sessionId = typeof idOrSession === 'number' ? idOrSession : idOrSession.id
  if (!sessionId) return
  activeAppId.value = sessionId
  activeTitle.value = typeof idOrSession === 'number' ? `对话 #${sessionId}` : idOrSession.appName || `对话 #${sessionId}`
  historyOpen.value = false
  await loadHistory(sessionId)
  await loadInitialState()
  await nextTick(scrollToBottom)
}

const deleteSession = async (session: API.AppVO) => {
  if (!session.id || deletingSessionId.value) return
  deletingSessionId.value = session.id
  try {
    const res = await deleteApp({ id: session.id })
    if (res.data.code !== 0 || !res.data.data) {
      message.error(res.data.message || '删除失败')
      return
    }
    message.success('删除成功')
    sessions.value = sessions.value.filter((item) => item.id !== session.id)
    if (activeAppId.value === session.id) {
      activeAppId.value = undefined
      activeTitle.value = '选择历史记录或直接提问'
      messages.value = []
      currentQuestion.value = null
      stop()
    }
    await loadSessions()
  } catch (error) {
    console.error('删除对话失败：', error)
    message.error('删除失败')
  } finally {
    deletingSessionId.value = undefined
  }
}

const createNewChat = async () => {
  activeAppId.value = undefined
  activeTitle.value = '新对话'
  messages.value = []
  currentQuestion.value = null
  historyOpen.value = false
  stop()
}

const chatTitleFromPrompt = (prompt: string) => {
  const title = prompt.trim().replace(/\s+/g, ' ')
  return title.length > 24 ? `${title.slice(0, 24)}...` : title || '新对话'
}

const ensureAssistantApp = async (initialPrompt = '') => {
  if (activeAppId.value) return true
  creating.value = true
  try {
    const prompt = initialPrompt.trim()
    const res = await addApp({ initPrompt: prompt })
    if (res.data.code === 0 && res.data.data) {
      activeAppId.value = res.data.data
      activeTitle.value = chatTitleFromPrompt(prompt)
      await loadSessions()
      return true
    }
    message.error(res.data.message || 'AI 助手初始化失败')
    return false
  } finally {
    creating.value = false
  }
}

const parseFileMeta = (content?: string): AIFileMeta | undefined => {
  if (!content) return undefined
  try {
    const meta = JSON.parse(content) as AIFileMeta
    return meta && typeof meta === 'object' ? meta : undefined
  } catch {
    return undefined
  }
}

const buildHistoryMessage = (chat: API.ChatHistory): AIMessage => {
  const fileMeta = chat.isFile ? parseFileMeta(chat.message) : undefined
  return {
    id: chat.id?.toString() || `${chat.messageType}-${chat.createTime || Math.random()}`,
    type: (chat.messageType === 'user' ? 'user' : 'ai') as 'user' | 'ai',
    content: fileMeta ? fileMeta.filename || '' : chat.message || '',
    createTime: chat.createTime,
    isFile: Boolean(chat.isFile && fileMeta),
    fileMeta,
  }
}

const loadHistory = async (appId: number) => {
  const res = await listAppChatHistory({ appId, pageSize: 30 })
  if (res.data.code !== 0) return
  const records = res.data.data?.records || []
  messages.value = records.map(buildHistoryMessage).reverse()
}

const handleSendMessage = async (rawMessage: string, rawFiles?: File[]) => {
  const content = rawMessage.trim()
  if ((!content && !rawFiles?.length) || sendingMessage.value) return
  const ready = await ensureAssistantApp(content)
  if (!ready) return

  sendingMessage.value = true
  try {
    if (rawFiles?.length && activeAppId.value) {
      const appId = String(activeAppId.value)
      for (const file of rawFiles) {
        const formData = new FormData()
        formData.append('appId', appId)
        formData.append('file', file)
        try {
          const res = await request.post('/app/file/add', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            timeout: 10 * 60 * 1000,
          })
          const data = res.data as any
          if (data?.code !== 0) {
            message.error(`文件 ${file.name} 上传失败: ${data?.message || '未知错误'}`)
          } else {
            await loadHistory(Number(appId))
          }
        } catch (err) {
          console.error('Upload failed:', err)
          message.error(`文件 ${file.name} 上传失败`)
        }
      }
    }

    if (!content) {
      sendingMessage.value = false
      await loadSessions()
      return
    }

    const ok = isGenerating.value ? await pushMessage(content) : await sendMessage(content)
    if (!ok) {
      message.error(isGenerating.value ? '追加失败' : '提交失败')
      return
    }
    userInput.value = ''
    await loadSessions()
  } finally {
    sendingMessage.value = false
    await nextTick(scrollToBottom)
  }
}

const cancelCurrentTask = async () => {
  if (!activeAppId.value) return
  await cancelAI({ appId: activeAppId.value, reason: 'cancelled by user' })
  stop()
}

const selectAnswerOption = (option: string) => {
  answerSelections.value = {
    ...answerSelections.value,
    [currentQuestionIndex.value]: option,
  }
}

const currentStepAnswer = () => (answerInput.value.trim() || currentAnswerSelection.value.trim()).trim()

const buildAnswerContent = () => {
  return currentQuestionItems.value
    .map((question, index) => {
      const value = answerSelections.value[index]?.trim()
      if (!value) return ''
      if (currentQuestionItems.value.length === 1) return value
      return `问题：${question.question}\n回答：${value}`
    })
    .filter(Boolean)
    .join('\n\n')
    .trim()
}

const submitQuestionStep = async () => {
  const stepAnswer = currentStepAnswer()
  if (!stepAnswer) return
  answerSelections.value = {
    ...answerSelections.value,
    [currentQuestionIndex.value]: stepAnswer,
  }
  answerInput.value = ''

  if (!isLastQuestion.value) {
    currentQuestionIndex.value += 1
    return
  }

  const answer = buildAnswerContent()
  if (!answer) return
  answeringQuestion.value = true
  try {
    const ok = await answerQuestion(answer)
    if (!ok) {
      message.error('回答提交失败')
      return
    }
    currentQuestion.value = null
    answerSelections.value = {}
  } finally {
    answeringQuestion.value = false
  }
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const formatFileSize = (size?: number) => {
  if (!size || size <= 0) return '未知大小'
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

const openFileMessage = (fileMeta: AIFileMeta) => {
  if (!activeAppId.value || !fileMeta.fileId || !fileMeta.filename) {
    message.warning('文件地址不存在')
    return
  }
  const base = (window as any).STATIC_BASE_URL || ''
  const filename = fileMeta.filename.split('/').map(encodeURIComponent).join('/')
  window.open(`${base}/document/${activeAppId.value}/${fileMeta.fileId}/${filename}`, '_blank')
}

const formatJSON = (value?: string) => {
  if (!value) return ''
  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch {
    return value
  }
}

const safeText = (value?: unknown) => (typeof value === 'string' ? value : value == null ? '' : String(value))

watch(
  () => currentQuestion.value?.id,
  () => {
    currentQuestionIndex.value = 0
    answerInput.value = ''
    answerSelections.value = {}
  },
)

watch(
  () => {
    const last = messages.value[messages.value.length - 1] as AIMessage | undefined
    return `${messages.value.length}:${last?.id || ''}:${last?.content?.length || 0}:${last?.toolCalls?.length || 0}`
  },
  () => nextTick(scrollToBottom),
)

onMounted(async () => {
  if (props.pageMode) {
    await loadSessions()
    await nextTick(scrollToBottom)
  }
})
</script>
<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 10px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}

.assistant-root {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 56px;
  height: 56px;
  z-index: 40;
}

.assistant-root.page-mode {
  position: static;
  width: 100%;
  height: calc(100vh - 92px);
  z-index: auto;
}

.assistant-root.page-mode .assistant-window {
  position: static;
  width: min(760px, 100%);
  height: 100%;
  margin: 0 auto;
  border-radius: 16px;
  background: #ffffff;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
}

.assistant-window {
  position: absolute;
  right: 0;
  bottom: 72px;
  width: min(520px, calc(100vw - 28px));
  height: min(680px, calc(100vh - 112px));
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  background: #f8fafc;
  box-shadow: 0 24px 80px rgba(15, 23, 42, 0.22);
}

.chat-pane {
  position: relative;
  min-height: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.assistant-head {
  height: 58px;
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  border-bottom: 1px solid #e5e7eb;
  background: rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(18px);
}

.assistant-head strong {
  display: block;
  color: #0f172a;
  font-size: 15px;
  font-weight: 700;
}

.assistant-head span {
  display: block;
  margin-top: 2px;
  max-width: 360px;
  color: #64748b;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.head-actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.head-actions :deep(.ant-btn) {
  border: 0;
  box-shadow: none;
}

.history-panel {
  position: absolute;
  top: 64px;
  right: 12px;
  z-index: 3;
  width: min(320px, calc(100% - 24px));
  max-height: min(420px, calc(100% - 96px));
  display: flex;
  flex-direction: column;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  background: #ffffff;
  box-shadow: 0 18px 48px rgba(15, 23, 42, 0.18);
  overflow: hidden;
}

.history-title {
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 0 12px;
  border-bottom: 1px solid #e5e7eb;
  color: #0f172a;
  font-size: 13px;
  font-weight: 700;
}

.history-cards {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.history-card-item {
  width: 100%;
  min-height: 46px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 4px 6px 4px 10px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  text-align: left;
  transition: background 0.15s ease, color 0.15s ease;
}

.history-card-item:hover {
  background: #f3f4f6;
}

.history-card-item.active {
  background: #ececf1;
}

.history-card-main {
  min-width: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
  padding: 3px 0;
  border: 0;
  background: transparent;
  text-align: left;
  cursor: pointer;
}

.history-delete {
  flex: 0 0 auto;
  opacity: 0;
  transition: opacity 0.15s ease, background 0.15s ease;
}

.history-card-item:hover .history-delete,
.history-delete:focus-visible {
  opacity: 1;
}

.history-card-title {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #111827;
  font-size: 14px;
  font-weight: 500;
  line-height: 1.35;
}

.history-card-time {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #9ca3af;
  font-size: 11px;
  line-height: 1.35;
}

.chat-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 22px 22px 18px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.empty-chat {
  min-height: 320px;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  user-select: none;
}

.empty-title {
  color: #111827;
  font-size: 16px;
  font-weight: 600;
}

.empty-subtitle {
  margin-top: 6px;
  color: #9ca3af;
  font-size: 13px;
}

.system-row {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #9ca3af;
  font-size: 12px;
}

.system-row::before,
.system-row::after {
  content: '';
  height: 1px;
  flex: 1;
  background: #f1f5f9;
}

@keyframes systemFadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}

.message-row {
  display: flex;
}

.user-row {
  justify-content: flex-end;
}

.assistant-row {
  justify-content: flex-start;
}

.user-message,
.assistant-message {
  max-width: min(88%, 430px);
  animation: messageSlideIn 0.25s ease;
}

.assistant-message {
  animation-delay: 0.1s;
}

@keyframes messageSlideIn {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}

.user-message {
  padding: 10px 14px;
  border-radius: 18px;
  background: #f3f4f6;
  color: #111827;
}

.assistant-message {
  padding: 0;
  color: #111827;
}

.agent-label,
.push-label {
  margin-bottom: 4px;
  color: #9ca3af;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
}

.message-markdown {
  font-size: 14px;
}

.user-markdown {
  line-height: 1.55;
}

.file-message {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 240px;
  max-width: 360px;
  cursor: pointer;
  user-select: none;
}

.file-message:hover .file-name {
  color: #2563eb;
}

.file-icon {
  flex: 0 0 auto;
  color: #2563eb;
  font-size: 22px;
}

.file-main {
  min-width: 0;
}

.file-name {
  overflow: hidden;
  color: #111827;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.35;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 3px;
  color: #64748b;
  font-size: 11px;
  line-height: 1.4;
}

.assistant-markdown {
  line-height: 1.7;
}

.tool-list {
  margin-top: 10px;
  display: grid;
  gap: 6px;
}

.tool-item {
  color: #6b7280;
  font-size: 12px;
}

.tool-item summary {
  display: flex;
  align-items: center;
  gap: 7px;
  cursor: pointer;
  list-style: none;
}

.tool-item summary::-webkit-details-marker {
  display: none;
}

.tool-corner {
  width: 12px;
  height: 13px;
  border-left: 1px solid #d1d5db;
  border-bottom: 1px solid #d1d5db;
  border-bottom-left-radius: 5px;
}

.tool-status {
  color: #9ca3af;
}

.tool-name {
  color: #374151;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.tool-code {
  margin: 7px 0 2px 20px;
  max-height: 180px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f9fafb;
  padding: 8px 10px;
  color: #374151;
  font-size: 11px;
  line-height: 1.5;
}

.assistant-loading {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  color: #9ca3af;
  font-size: 12px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.loading-dot {
  width: 7px;
  height: 7px;
  border-radius: 999px;
  background: #9ca3af;
  animation: pulse-dot 1.2s ease-in-out infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

.composer-wrap {
  flex: 0 0 auto;
  padding: 14px 16px 12px;
  border-top: 1px solid #f1f5f9;
  background: white;
}

.question-panel {
  margin-bottom: 12px;
  padding: 12px;
  border: 1px solid #dbeafe;
  border-radius: 14px;
  background: #eff6ff;
}

.question-title,
.question-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.question-title {
  margin-bottom: 8px;
  color: #1e3a8a;
  font-size: 13px;
  font-weight: 700;
}

.question-count {
  color: #64748b;
  font-size: 12px;
}

.question-content {
  margin-bottom: 8px;
  color: #0f172a;
  font-size: 14px;
}

.question-options {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.question-option {
  border: 1px solid #bfdbfe;
  border-radius: 999px;
  padding: 5px 10px;
  background: white;
  color: #1d4ed8;
  cursor: pointer;
}

.question-option-active {
  background: #2563eb;
  color: white;
}

.answer-input {
  margin-bottom: 8px;
  border-radius: 12px;
}

.assistant-ball {
  position: absolute;
  inset: 0;
  width: 56px;
  height: 56px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  border-radius: 50%;
  background: #101010;
  color: #ffffff;
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.22);
  cursor: pointer;
}

.assistant-gpt-logo {
  width: 34px;
  height: 34px;
  display: block;
  color: #fff;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

@media (max-width: 768px) {
  .assistant-root {
    right: 12px;
    bottom: 76px;
  }

  .assistant-window {
    width: min(380px, calc(100vw - 24px));
    height: min(620px, calc(100vh - 128px));
  }

  .assistant-root.page-mode {
    height: calc(100vh - 68px);
  }

  .assistant-root.page-mode .assistant-window {
    width: 100%;
    height: 100%;
    border-radius: 0;
  }

  .chat-list {
    padding: 14px;
  }

  .empty-title {
    font-size: 22px;
  }
}
</style>
