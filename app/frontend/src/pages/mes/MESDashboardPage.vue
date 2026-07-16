<template>
  <main class="dashboard-page">
    <header class="dashboard-header">
      <div>
        <p class="eyebrow">PRODUCTION CONTROL</p>
        <h1>智能生产看板</h1>
        <p class="subtitle">实时汇总生产、质检与计划执行情况</p>
      </div>
      <a-button :loading="loading" @click="loadDashboard">
        <template #icon><ReloadOutlined /></template>
        刷新数据
      </a-button>
    </header>

    <a-alert v-if="error" type="error" :message="error" show-icon closable @close="error = ''" />

    <a-spin :spinning="loading && !overview">
      <section class="metrics-grid">
        <article v-for="metric in metrics" :key="metric.label" class="metric-panel">
          <div class="metric-icon" :class="metric.tone"><component :is="metric.icon" /></div>
          <div>
            <p>{{ metric.label }}</p>
            <strong>{{ metric.value }}</strong>
            <span>{{ metric.unit }}</span>
          </div>
          <small>{{ metric.note }}</small>
        </article>
      </section>

      <section class="dashboard-grid">
        <article class="panel production-panel">
          <div class="panel-heading">
            <div>
              <p class="panel-kicker">OUTPUT TREND</p>
              <h2>近 7 日产量</h2>
            </div>
            <span class="live-badge"><i></i> 数据已同步</span>
          </div>
          <div class="trend-chart">
            <svg viewBox="0 0 700 260" role="img" aria-label="近七日产量折线图">
              <line v-for="y in [30, 85, 140, 195]" :key="y" x1="42" :y1="y" x2="680" :y2="y" class="grid-line" />
              <polyline :points="trendPoints" class="trend-line" />
              <g v-for="(day, index) in trend" :key="day.date">
                <line :x1="pointX(index)" y1="30" :x2="pointX(index)" y2="205" class="vertical-guide" />
                <circle :cx="pointX(index)" :cy="pointY(day.quantity)" r="5" class="trend-dot" />
                <text :x="pointX(index)" :y="pointY(day.quantity) - 13" class="value-label">{{ day.quantity }}</text>
                <text :x="pointX(index)" y="235" class="date-label">{{ formatDay(day.date) }}</text>
              </g>
            </svg>
          </div>
        </article>

        <article class="panel completion-panel">
          <div class="panel-heading">
            <div>
              <p class="panel-kicker">PLAN DELIVERY</p>
              <h2>计划完成度</h2>
            </div>
          </div>
          <div class="completion-content">
            <div class="donut" :style="completionStyle">
              <div><strong>{{ completionRate }}</strong><span>%</span><small>合格完成率</small></div>
            </div>
            <div class="completion-detail">
              <div><span>已完成合格量</span><strong>{{ overview?.planCompletedQuantity ?? 0 }}</strong></div>
              <div><span>计划目标总量</span><strong>{{ overview?.planExpectedQuantity ?? 0 }}</strong></div>
            </div>
          </div>
        </article>

        <article class="panel status-panel">
          <div class="panel-heading">
            <div>
              <p class="panel-kicker">PLAN STATUS</p>
              <h2>生产计划分布</h2>
            </div>
            <strong class="plan-total">{{ planTotal }}<small> 项计划</small></strong>
          </div>
          <div class="status-content">
            <div class="status-donut" :style="statusStyle"><div>{{ planTotal }}</div></div>
            <div class="legend">
              <div v-for="item in statusItems" :key="item.label">
                <i :style="{ background: item.color }"></i><span>{{ item.label }}</span><strong>{{ item.value }}</strong>
              </div>
            </div>
          </div>
        </article>

        <article class="panel insight-panel">
          <div class="insight-mark"><BulbOutlined /></div>
          <div>
            <p class="panel-kicker">SMART INSIGHT</p>
            <h2>今日生产洞察</h2>
            <p>{{ insight }}</p>
          </div>
        </article>
      </section>
    </a-spin>

    <footer>数据生成于 {{ generatedTime }}</footer>
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { BulbOutlined, CheckCircleOutlined, ClockCircleOutlined, ReloadOutlined, RiseOutlined } from '@ant-design/icons-vue'
import { getDashboardOverview, type DashboardOverview } from '@/api/mesController'

const overview = ref<DashboardOverview>()
const loading = ref(false)
const error = ref('')

const loadDashboard = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getDashboardOverview()
    if (response.data.code !== 0 || !response.data.data) throw new Error(response.data.message || '看板数据加载失败')
    overview.value = response.data.data
  } catch (cause: any) {
    error.value = cause?.response?.data?.message || cause?.message || '看板数据加载失败'
  } finally {
    loading.value = false
  }
}

const metrics = computed(() => [
  { label: '今日产量', value: overview.value?.todayProduction ?? 0, unit: '件', note: '今日生产单体', icon: RiseOutlined, tone: 'green' },
  { label: '本周产量', value: overview.value?.weekProduction ?? 0, unit: '件', note: '本周累计生产', icon: CheckCircleOutlined, tone: 'blue' },
  { label: '待检数量', value: overview.value?.pendingInspection ?? 0, unit: '件', note: '等待质量检验', icon: ClockCircleOutlined, tone: 'amber' },
  { label: '计划完成率', value: overview.value?.planCompletionRate ?? 0, unit: '%', note: '按合格产出统计', icon: RiseOutlined, tone: 'red' },
])
const trend = computed(() => overview.value?.dailyProduction ?? Array.from({ length: 7 }, (_, index) => ({ date: '', quantity: index === 6 ? 0 : 0 })))
const maxQuantity = computed(() => Math.max(1, ...trend.value.map((item) => item.quantity)))
const pointX = (index: number) => 55 + index * 102
const pointY = (quantity: number) => 205 - (quantity / maxQuantity.value) * 160
const trendPoints = computed(() => trend.value.map((item, index) => `${pointX(index)},${pointY(item.quantity)}`).join(' '))
const formatDay = (date: string) => date ? `${Number(date.slice(5, 7))}/${Number(date.slice(8, 10))}` : '--'
const completionRate = computed(() => Math.min(100, overview.value?.planCompletionRate ?? 0))
const completionStyle = computed(() => ({ background: `conic-gradient(#1677ff 0 ${completionRate.value}%, #edf1f5 ${completionRate.value}% 100%)` }))
const statusItems = computed(() => [
  { label: '未开始', value: overview.value?.planStatus.notStarted ?? 0, color: '#a8b1bd' },
  { label: '进行中', value: overview.value?.planStatus.inProgress ?? 0, color: '#1677ff' },
  { label: '已完成', value: overview.value?.planStatus.completed ?? 0, color: '#20a66a' },
])
const planTotal = computed(() => statusItems.value.reduce((sum, item) => sum + item.value, 0))
const statusStyle = computed(() => {
  if (!planTotal.value) return { background: '#edf1f5' }
  const first = statusItems.value[0].value / planTotal.value * 100
  const second = first + statusItems.value[1].value / planTotal.value * 100
  return { background: `conic-gradient(#a8b1bd 0 ${first}%, #1677ff ${first}% ${second}%, #20a66a ${second}% 100%)` }
})
const insight = computed(() => {
  if (!overview.value) return '正在汇总生产数据。'
  if (overview.value.pendingInspection > overview.value.todayProduction) return `当前有 ${overview.value.pendingInspection} 件待检，质检队列高于今日产出，建议优先安排质量检验。`
  if (completionRate.value >= 80) return `生产计划总体完成率达到 ${completionRate.value}%，执行状态良好，可关注剩余进行中计划。`
  return `今日已产出 ${overview.value.todayProduction} 件，当前计划完成率 ${completionRate.value}%，建议跟进进行中的生产计划。`
})
const generatedTime = computed(() => overview.value?.generatedAt ? new Date(overview.value.generatedAt).toLocaleString('zh-CN', { hour12: false }) : '--')

onMounted(loadDashboard)
</script>

<style scoped>
.dashboard-page { min-height: 100%; padding: 28px; background: #f3f6f8; color: #18212b; }
.dashboard-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 24px; }
.eyebrow,.panel-kicker { margin: 0 0 6px; color: #1677ff; font-size: 11px; font-weight: 700; letter-spacing: 0; }
h1 { margin: 0; font-size: 28px; line-height: 1.25; } .subtitle { margin: 7px 0 0; color: #6b7785; }
.metrics-grid { display: grid; grid-template-columns: repeat(4,minmax(0,1fr)); gap: 16px; margin: 18px 0; }
.metric-panel,.panel { background: #fff; border: 1px solid #e5eaef; border-radius: 8px; box-shadow: 0 5px 18px rgba(24,33,43,.045); }
.metric-panel { position: relative; display: flex; align-items: center; gap: 14px; min-height: 112px; padding: 20px; overflow: hidden; }
.metric-panel::after { content: ''; position: absolute; top: 0; right: 0; width: 4px; height: 100%; background: #dbe2e8; }
.metric-icon { width: 44px; height: 44px; display: grid; place-items: center; border-radius: 8px; font-size: 20px; }
.metric-icon.green { color:#168a58;background:#e7f7ef }.metric-icon.blue{color:#1677ff;background:#eaf3ff}.metric-icon.amber{color:#bd7410;background:#fff3dd}.metric-icon.red{color:#d94841;background:#fff0ef}
.metric-panel p { margin: 0 0 3px; color: #697684; font-size: 13px; }.metric-panel strong { font-size: 29px; line-height: 1; }.metric-panel span { margin-left: 5px; color:#697684 }.metric-panel small { position:absolute;left:78px;bottom:14px;color:#98a2ad }
.dashboard-grid { display:grid;grid-template-columns:minmax(0,1.8fr) minmax(300px,1fr);gap:16px }.panel { padding:22px;min-width:0 }.panel-heading { display:flex;align-items:flex-start;justify-content:space-between }.panel-heading h2,.insight-panel h2 { margin:0;font-size:17px }.live-badge{color:#687480;font-size:12px}.live-badge i{display:inline-block;width:7px;height:7px;margin-right:6px;border-radius:50%;background:#20a66a}
.trend-chart { height:270px;margin-top:6px }.trend-chart svg { width:100%;height:100%;overflow:visible }.grid-line{stroke:#e9edf1;stroke-width:1}.vertical-guide{stroke:transparent}.trend-line{fill:none;stroke:#1677ff;stroke-width:3;stroke-linecap:round;stroke-linejoin:round}.trend-dot{fill:#fff;stroke:#1677ff;stroke-width:3}.value-label,.date-label{text-anchor:middle;fill:#687480;font-size:12px}.value-label{fill:#26313c;font-weight:600}
.completion-content,.status-content{display:flex;align-items:center;gap:28px;min-height:220px}.donut,.status-donut{width:150px;aspect-ratio:1;border-radius:50%;display:grid;place-items:center;flex:none}.donut>div,.status-donut>div{width:108px;aspect-ratio:1;border-radius:50%;background:#fff;display:grid;place-content:center;text-align:center}.donut strong{font-size:34px}.donut span{font-size:16px}.donut small{display:block;color:#7d8995}.completion-detail{flex:1}.completion-detail div,.legend div{display:flex;align-items:center;justify-content:space-between;padding:12px 0;border-bottom:1px solid #edf0f3}.completion-detail span,.legend span{color:#6b7785}.completion-detail strong{font-size:18px}.status-donut{width:140px}.status-donut>div{width:98px;font-size:28px;font-weight:700}.legend{flex:1}.legend i{width:9px;height:9px;border-radius:2px;margin-right:9px}.legend span{margin-right:auto}.plan-total{font-size:22px}.plan-total small{font-size:12px;color:#76828e;font-weight:400}
.insight-panel{display:flex;gap:16px;align-items:flex-start;border-left:4px solid #20a66a}.insight-mark{width:42px;height:42px;display:grid;place-items:center;flex:none;border-radius:8px;background:#e8f7ef;color:#168a58;font-size:20px}.insight-panel p:last-child{margin:12px 0 0;color:#5c6975;line-height:1.7}.dashboard-page footer{text-align:right;margin-top:14px;color:#8b96a1;font-size:12px}
@media(max-width:1100px){.metrics-grid{grid-template-columns:repeat(2,1fr)}.dashboard-grid{grid-template-columns:1fr}}@media(max-width:640px){.dashboard-page{padding:16px}.dashboard-header{gap:16px}.metrics-grid{grid-template-columns:1fr}.completion-content,.status-content{align-items:flex-start;gap:18px}.donut,.status-donut{width:120px}.donut>div,.status-donut>div{width:84px}.trend-chart{overflow-x:auto}.trend-chart svg{min-width:600px}}
</style>
