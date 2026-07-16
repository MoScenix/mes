import { onUnmounted, ref, watch, type Ref } from 'vue'
import { answerAI, cancelAI, getAIState, listAIEvents, pushAI, submitAI } from '@/api/aiController'

export interface AIMessage {
  id: string
  type: 'user' | 'ai' | 'system'
  content: string
  loading?: boolean
  agent?: string
  isPush?: boolean
  isFile?: boolean
  fileMeta?: AIFileMeta
  createTime?: string
  toolCalls?: AIToolCall[]
  parts?: AIMessagePart[]
}

export interface AIFileMeta {
  fileId?: number
  filename?: string
  contentType?: string
  size?: number
  textFilename?: string
  textSize?: number
  isBig?: boolean
  chunkCount?: number
  parentCount?: number
  text?: string
}

export interface AIToolCall {
  id: string
  name: string
  status: 'running' | 'success' | 'error'
  args?: string
  result?: string
}

export type AIMessagePart =
  | { type: 'text'; id: string; content: string }
  | { type: 'tool'; id: string; tool: AIToolCall }

export interface AIQuestion {
  id: string
  agent?: string
  content: string
  questions: AIQuestionItem[]
  payload?: any
}

export interface AIQuestionItem {
  question: string
  options: string[]
}

export function useAIEvents(historyId: Ref<any>, options?: { onDone?: () => void }) {
  const messages = ref<AIMessage[]>([])
  const aiState = ref<API.AIState | null>(null)
  const isGenerating = ref(false)
  const currentQuestion = ref<AIQuestion | null>(null)

  let lastEventId = '0'
  let polling = false
  let abortController: AbortController | null = null

  function stopEventPolling() {
    polling = false
    abortController?.abort()
    abortController = null
  }

  function stop() {
    stopEventPolling()
  }

  onUnmounted(stop)

  watch(
    () => historyId.value,
    () => {
      stopEventPolling()
      lastEventId = '0'
      aiState.value = null
      currentQuestion.value = null
      isGenerating.value = false
    },
  )

  function appendSystem(content: string) {
    if (!content) return
    const last = messages.value[messages.value.length - 1]
    if (last?.type === 'system' && last.content === content) return
    messages.value.push({ id: `system-${Date.now()}`, type: 'system', content })
  }

  function hasMessage(id?: string) {
    if (!id) return false
    return messages.value.some((msg) => msg.id === id)
  }

  function ensureAIMessage(event?: API.AIEvent) {
    const last = messages.value[messages.value.length - 1]
    if (last?.type === 'ai') return last
    const msg: AIMessage = {
      id: event?.id || `ai-${Date.now()}`,
      type: 'ai',
      content: '',
      agent: event?.agent,
      loading: true,
      toolCalls: [],
    }
    messages.value.push(msg)
    return msg
  }

  function isActiveStatus(status?: string) {
    return (
      status === 'queued' ||
      status === 'running' ||
      status === 'waiting_answer' ||
      status === 'interrupted'
    )
  }

  function setLocalState(status: API.AIStatusType, patch: Partial<API.AIState> = {}) {
    aiState.value = {
      exists: true,
      status,
      agent: patch.agent ?? aiState.value?.agent,
      lastEventId: patch.lastEventId ?? lastEventId,
      pendingInterrupts: patch.pendingInterrupts ?? aiState.value?.pendingInterrupts,
      message: patch.message ?? aiState.value?.message,
    }
  }

  function restoreRunningMessage(state: API.AIState) {
    if (
      !isActiveStatus(state.status) ||
      state.status === 'waiting_answer' ||
      state.status === 'interrupted'
    )
      return
    const last = messages.value[messages.value.length - 1]
    if (last?.type === 'ai') {
      last.loading = true
      if (state.agent) last.agent = state.agent
      return
    }
    messages.value.push({
      id: `ai-running-${Date.now()}`,
      type: 'ai',
      content: ((state as any).buffer as string) || '',
      agent: state.agent,
      loading: true,
      toolCalls: [],
    })
  }

  function finishAIMessage() {
    const last = messages.value[messages.value.length - 1]
    if (last?.type === 'ai') last.loading = false
  }

  function appendToolError(event: API.AIEvent) {
    const content = event.content || 'tool failed'
    const match = content.match(/tool\[name:([^\]\s]+)|toolName=([^,\s]+)/)
    const name = event.name || match?.[1] || match?.[2] || 'tool'
    const msg = ensureAIMessage(event)
    if (!msg.toolCalls) msg.toolCalls = []
    const tool = {
      id: event.id || `tool-error-${Date.now()}`,
      name,
      status: 'error',
      result: content,
    } satisfies AIToolCall
    msg.toolCalls.push(tool)
    appendToolPart(msg, tool)
    msg.loading = false
  }

  function appendTextPart(msg: AIMessage, content: string, event?: API.AIEvent) {
    if (!content) return
    if (!msg.parts) msg.parts = []
    const last = msg.parts[msg.parts.length - 1]
    if (last?.type === 'text') {
      last.content += content
      return
    }
    msg.parts.push({
      type: 'text',
      id: event?.id || `text-${Date.now()}-${msg.parts.length}`,
      content,
    })
  }

  function appendToolPart(msg: AIMessage, tool: AIToolCall) {
    if (!msg.parts) msg.parts = []
    if (msg.parts.some((part) => part.type === 'tool' && part.tool.id === tool.id)) return
    msg.parts.push({
      type: 'tool',
      id: tool.id,
      tool,
    })
  }

  function processEvent(event: API.AIEvent) {
    if (!event.type) return
    if (event.id) lastEventId = event.id

    switch (event.type) {
      case 'accepted':
        setLocalState('queued', { lastEventId: event.id, message: event.content })
        appendSystem(event.content || '任务已接收')
        break
      case 'agent_start':
        setLocalState('running', { agent: event.agent, lastEventId: event.id })
        isGenerating.value = true
        break
      case 'message': {
        const content = event.content || ''
        if (isQuestionPayload(content)) break
        const msg = ensureAIMessage(event)
        msg.content += content
        appendTextPart(msg, content, event)
        msg.loading = false
        break
      }
      case 'tool_call': {
        if (event.name === 'ask_user') {
          setLocalState('waiting_answer', { agent: event.agent, lastEventId: event.id })
          currentQuestion.value = normalizeQuestion(event)
          isGenerating.value = false
          break
        }
        const msg = ensureAIMessage(event)
        if (!msg.toolCalls) msg.toolCalls = []
        const toolID = event.targetId || event.id || `tool-${Date.now()}`
        let tool = msg.toolCalls.find((item) => item.id === toolID)
        if (tool) {
          tool.name = event.name || event.content || tool.name
          tool.status = 'running'
          tool.args = event.payloadJson || tool.args
        } else {
          tool = {
            id: toolID,
            name: event.name || event.content || 'tool',
            status: 'running',
            args: event.payloadJson || '',
          }
          msg.toolCalls.push(tool)
        }
        appendToolPart(msg, tool)
        msg.loading = false
        break
      }
      case 'tool_result': {
        const msg = ensureAIMessage(event)
        if (!msg.toolCalls) msg.toolCalls = []
        const tool = msg.toolCalls.find(
          (item) => item.name === event.name || item.id === event.targetId,
        )
        if (tool) {
          tool.status = event.status === 'error' ? 'error' : 'success'
          tool.result = event.content || ''
        } else {
          const fallbackTool = {
            id: event.targetId || event.id || `tool-result-${Date.now()}`,
            name: event.name || 'tool',
            status: event.status === 'error' ? 'error' : 'success',
            result: event.content || '',
          } satisfies AIToolCall
          msg.toolCalls.push(fallbackTool)
          appendToolPart(msg, fallbackTool)
        }
        msg.loading = false
        break
      }
      case 'push':
        if (hasMessage(event.id)) break
        finishAIMessage()
        messages.value.push({
          id: event.id || `push-${Date.now()}`,
          type: 'user',
          content: event.content || '',
          isPush: true,
        })
        ensureAIMessage()
        break
      case 'answer':
        if (!currentQuestion.value || currentQuestion.value.id === event.targetId) {
          currentQuestion.value = null
          isGenerating.value = true
          setLocalState('running', { lastEventId: event.id })
          ensureAIMessage()
        }
        break
      case 'question':
        setLocalState('waiting_answer', { agent: event.agent, lastEventId: event.id })
        currentQuestion.value = normalizeQuestion(event)
        isGenerating.value = false
        break
      case 'done':
        setLocalState('done', { lastEventId: event.id, message: event.content })
        finishAIMessage()
        isGenerating.value = false
        polling = false
        options?.onDone?.()
        break
      case 'cancelled':
        setLocalState('cancelled', { lastEventId: event.id, message: event.content })
        finishAIMessage()
        isGenerating.value = false
        polling = false
        options?.onDone?.()
        break
      case 'error': {
        setLocalState('error', { lastEventId: event.id, message: event.content })
        if ((event.content || '').includes('failed to invoke tool')) {
          appendToolError(event)
        } else {
          const msg = ensureAIMessage(event)
          msg.content = event.content || 'AI 处理失败'
          msg.loading = false
        }
        isGenerating.value = false
        polling = false
        options?.onDone?.()
        break
      }
    }
  }

  async function pollEvents(initialLastId?: string) {
    if (polling) return
    const pollingHistoryId = Number(historyId.value)
    if (!pollingHistoryId) return
    polling = true
    lastEventId = initialLastId || lastEventId || '0'

    while (polling && Number(historyId.value) === pollingHistoryId) {
      try {
        abortController = new AbortController()
        const res = await listAIEvents(
          {
            historyId: pollingHistoryId,
            lastId: lastEventId,
            blockMs: 30000,
            count: 50,
          },
          { signal: abortController.signal },
        )
        if (Number(historyId.value) !== pollingHistoryId) break
        if (res.data.code === 0) {
          const events = res.data.data?.events || []
          for (const event of events) processEvent(event)
        }
      } catch (e: any) {
        if (e?.name === 'CanceledError' || e?.code === 'ERR_CANCELED') break
        await refreshState()
        if (!isActiveStatus(aiState.value?.status)) {
          polling = false
          break
        }
        await new Promise((resolve) => setTimeout(resolve, 2000))
      }
    }
  }

  async function refreshState() {
    if (!historyId.value) return
    const res = await getAIState({ historyId: Number(historyId.value) })
    if (res.data.code !== 0 || !res.data.data) return
    aiState.value = res.data.data
    if (!res.data.data.exists) {
      currentQuestion.value = null
      return
    }
    syncQuestionFromState(res.data.data)
  }

  function syncQuestionFromState(state: API.AIState) {
    const interrupt = state.pendingInterrupts?.[0]
    if (interrupt) {
      if (!currentQuestion.value || currentQuestion.value.id !== interrupt.id) {
        currentQuestion.value = normalizeQuestion(interrupt)
      }
      return
    }
    if (state.status !== 'waiting_answer' && state.status !== 'interrupted') {
      currentQuestion.value = null
    }
  }

  async function loadInitialState() {
    if (!historyId.value) return
    await refreshState()
    if (!aiState.value?.exists) return
    lastEventId = aiState.value.lastEventId || '0'
    const status = aiState.value.status
    if (isActiveStatus(status)) {
      isGenerating.value = status !== 'waiting_answer'
      restoreRunningMessage(aiState.value)
      void pollEvents(lastEventId)
    }
  }

  async function sendMessage(content: string) {
    if (!historyId.value || !content.trim() || isGenerating.value) return false
    messages.value.push({ id: `user-${Date.now()}`, type: 'user', content: content.trim() })
    messages.value.push({
      id: `ai-${Date.now()}`,
      type: 'ai',
      content: '',
      loading: true,
      toolCalls: [],
    })
    currentQuestion.value = null
    isGenerating.value = true
    lastEventId = '0'
    setLocalState('queued')
    const res = await submitAI({ historyId: Number(historyId.value), message: content.trim() })
    if (res.data.code !== 0) {
      finishAIMessage()
      isGenerating.value = false
      setLocalState('error', { message: res.data.message || '提交失败' })
      return false
    }
    void pollEvents('0')
    return true
  }

  async function pushMessage(content: string) {
    if (!historyId.value || !content.trim() || !isGenerating.value) return false
    const res = await pushAI({ historyId: Number(historyId.value), content: content.trim() })
    if (res.data.code !== 0) return false
    const id = res.data.data || `push-${Date.now()}`
    finishAIMessage()
    messages.value.push({ id, type: 'user', content: content.trim(), isPush: true })
    messages.value.push({
      id: `ai-${Date.now()}`,
      type: 'ai',
      content: '',
      loading: true,
      toolCalls: [],
    })
    setLocalState('running')
    return true
  }

  async function answerQuestion(content: string) {
    if (!historyId.value || !currentQuestion.value || !content.trim()) return false
    const res = await answerAI({
      historyId: Number(historyId.value),
      answers: {
        [currentQuestion.value.id]: {
          content: content.trim(),
        },
      },
    })
    if (res.data.code !== 0) return false
    await refreshState()
    if (!aiState.value?.exists) return true
    isGenerating.value = true
    setLocalState('running')
    ensureAIMessage()
    void pollEvents(lastEventId)
    return true
  }

  async function cancelCurrentTask() {
    if (!historyId.value) return
    await cancelAI({ historyId: Number(historyId.value), reason: '用户取消' })
    finishAIMessage()
    appendSystem('已取消')
    isGenerating.value = false
    setLocalState('cancelled', { message: '用户取消' })
    stopEventPolling()
  }

  return {
    messages,
    aiState,
    isGenerating,
    currentQuestion,
    sendMessage,
    pushMessage,
    answerQuestion,
    cancelCurrentTask,
    loadInitialState,
    stop,
  }
}

function normalizeQuestionItem(value: any): AIQuestionItem | null {
  if (typeof value === 'string') {
    const question = value.trim()
    return question ? { question, options: [] } : null
  }
  const question = String(value?.question || value?.content || '').trim()
  if (!question) return null
  const options = Array.isArray(value?.options)
    ? value.options.map((option: any) => String(option).trim()).filter(Boolean)
    : []
  return { question, options }
}

function normalizeQuestion(source: API.AIEvent | API.AIPendingInterrupt): AIQuestion {
  const payload = parseMaybeJSON(source.payloadJson)
  const contentValue = parseMaybeJSON(source.content)
  const data = typeof contentValue === 'object' && contentValue ? contentValue : payload
  const eventQuestions = 'questions' in source ? source.questions : undefined
  const questions = normalizeQuestionItems(eventQuestions || data?.questions)
  const content = questions.length
    ? questions.map((item) => item.question).join('\n')
    : String(data?.question || data?.content || data?.message || source.content || '')
  return {
    id: ('targetId' in source ? source.targetId : '') || source.id || '',
    agent: source.agent,
    content,
    questions,
    payload: typeof data === 'object' ? data : payload,
  }
}

function isQuestionPayload(content?: string) {
  const data = parseMaybeJSON(content)
  return Boolean(
    data &&
      typeof data === 'object' &&
      Array.isArray(data.questions) &&
      normalizeQuestionItems(data.questions).length,
  )
}

function normalizeQuestionItems(value: any): AIQuestionItem[] {
  if (!Array.isArray(value)) return []
  return value.map(normalizeQuestionItem).filter((item): item is AIQuestionItem => Boolean(item))
}

function parseMaybeJSON(value?: string): any {
  if (!value) return undefined
  let current: any = value
  for (let i = 0; i < 2; i++) {
    if (typeof current !== 'string') return current
    const trimmed = current.trim()
    if (!trimmed.startsWith('{') && !trimmed.startsWith('[') && !trimmed.startsWith('"'))
      return current
    try {
      current = JSON.parse(trimmed)
    } catch {
      return current
    }
  }
  return current
}
