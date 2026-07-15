const { get, post } = require('../../utils/request')
const { scan } = require('../../utils/mes-code')

const FLOW_IN = 1
const FLOW_OUT = 2
const FLOW_APPROVED = 3
const STOCK_IN = 1
const STOCK_OUT = 3
const QUALITY_PENDING = 1
const QUALITY_QUALIFIED = 2

Page({
  data: {
    mode: 'detail', title: '扫描 MES 码', subtitle: '物料、流转单和工程单均可识别',
    buttonText: '开始扫码', showScanButton: true, scanning: false, submitting: false,
    context: null, detail: null, unit: null, lastMessage: ''
  },
  onLoad(options) {
    const mode = options.mode || 'detail'
    const copy = {
      inbound: ['扫描入库', '先扫描入库流转单，再逐个扫描合格单体'],
      receive: ['领取货物', '先扫描出库流转单，再逐个扫描在库单体'],
      inspect: ['检测单体', '先扫描工程单，再扫描其中的待检测单体'],
      detail: ['扫描 MES 码', '物料、流转单和工程单均可识别']
    }[mode]
    this.setData({ mode, title: copy[0], subtitle: copy[1] })
    if (options.kind && options.id) this.loadDetail(options.kind, Number(options.id))
  },
  expectedKind() {
    if (this.data.context) return 'ITEM_UNIT'
    if (this.data.mode === 'inbound' || this.data.mode === 'receive') return 'FLOW'
    if (this.data.mode === 'inspect') return 'ENGINEERING_ORDER'
    return undefined
  },
  async startScan() {
    this.setData({ scanning: true, lastMessage: '' })
    try {
      const parsed = await scan(this.expectedKind())
      if (!this.data.context) await this.loadRoot(parsed)
      else if (this.data.mode === 'inspect') await this.loadInspectUnit(parsed.id)
      else await this.completeFlowUnit(parsed.id)
    } catch (error) {
      if (error) wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ scanning: false })
    }
  },
  async loadRoot(parsed) {
    if (this.data.mode === 'detail') {
      await this.loadDetail(parsed.kind, parsed.id)
      return
    }
    if (this.data.mode === 'inspect') {
      const order = await get('/mes/engineering-order/get', { id: parsed.id })
      this.setData({
        context: { id: order.id, name: order.name || '工程单', kindLabel: '工程单', progress: `${order.producedQuantity || 0} / ${order.expectedQuantity || 0}` },
        buttonText: '扫描待检测单体'
      })
      return
    }
    const flow = await get('/mes/inventory-flow/get', { id: parsed.id })
    const requiredType = this.data.mode === 'inbound' ? FLOW_IN : FLOW_OUT
    if (flow.flowType !== requiredType) throw new Error(this.data.mode === 'inbound' ? '只能扫描入库流转单' : '只能扫描出库流转单')
    if (flow.flowStatus !== FLOW_APPROVED) throw new Error('流转单尚未审批通过')
    const expected = (flow.items || []).reduce((sum, item) => sum + (item.applyQuantity || 0), 0)
    const finished = (flow.items || []).reduce((sum, item) => sum + (item.finishedQuantity || 0), 0)
    this.flow = flow
    this.setData({
      context: { id: flow.id, name: flow.name || '流转单', kindLabel: '流转单', progress: `${finished} / ${expected}` },
      buttonText: this.data.mode === 'inbound' ? '扫描合格单体入库' : '扫描在库单体领取'
    })
  },
  async completeFlowUnit(id) {
    const unit = await get('/mes/item-unit/get', { id })
    const inbound = this.data.mode === 'inbound'
    if (unit.qualityStatus !== QUALITY_QUALIFIED) throw new Error(inbound ? '扫码入库必须为合格品' : '只能领取合格单体')
    if (unit.stockStatus !== (inbound ? STOCK_OUT : STOCK_IN)) throw new Error(inbound ? '只能录入不在库单体' : '只能领取在库单体')
    const flowItem = (this.flow.items || []).find(item => item.itemId === unit.itemId)
    if (!flowItem) throw new Error('该单体物品不在流转单明细中')
    await this.confirm(inbound ? `确认单体 #${id} 入库？` : `确认领取单体 #${id}？`)
    await post('/mes/inventory-flow/complete', { id: this.flow.id, itemUnitIds: [id] })
    flowItem.finishedQuantity = (flowItem.finishedQuantity || 0) + 1
    const expected = this.flow.items.reduce((sum, item) => sum + (item.applyQuantity || 0), 0)
    const finished = this.flow.items.reduce((sum, item) => sum + (item.finishedQuantity || 0), 0)
    this.setData({ 'context.progress': `${finished} / ${expected}`, lastMessage: inbound ? `单体 #${id} 已入库` : `单体 #${id} 已领取` })
    wx.showToast({ title: inbound ? '入库成功' : '领取成功', icon: 'success' })
  },
  async loadInspectUnit(id) {
    const unit = await get('/mes/item-unit/get', { id })
    if (unit.engineeringOrderId !== this.data.context.id) throw new Error('单体不属于当前工程单')
    if (unit.qualityStatus !== QUALITY_PENDING) throw new Error('只能检测待检测单体')
    this.setData({ unit, showScanButton: false })
  },
  async submitQuality(event) {
    const qualityStatus = Number(event.currentTarget.dataset.quality)
    this.setData({ submitting: true })
    try {
      await post('/mes/item-unit/status/update', {
        id: this.data.unit.id,
        stockStatus: this.data.unit.stockStatus || STOCK_OUT,
        qualityStatus
      })
      wx.showToast({ title: qualityStatus === 2 ? '已标记合格' : '已标记不合格', icon: 'success' })
      this.setData({ unit: null, showScanButton: true, lastMessage: `单体 #${this.data.unit.id} 检测完成` })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  },
  async loadDetail(kind, id) {
    const config = {
      FLOW: ['/mes/inventory-flow/get', '流转单'],
      ITEM_UNIT: ['/mes/item-unit/get', '物料单体'],
      ENGINEERING_ORDER: ['/mes/engineering-order/get', '工程单']
    }[kind]
    if (!config) throw new Error('不支持的 MES 码')
    const value = await get(config[0], { id })
    const rows = []
    if (kind === 'FLOW') rows.push({ label: '进度', value: `${(value.items || []).reduce((s, x) => s + (x.finishedQuantity || 0), 0)} / ${(value.items || []).reduce((s, x) => s + (x.applyQuantity || 0), 0)}` })
    if (kind === 'ITEM_UNIT') rows.push({ label: '物料', value: value.itemName || `#${value.itemId}` }, { label: '库存状态', value: value.stockStatus === 1 ? '在库' : value.stockStatus === 2 ? '预留' : '不在库' }, { label: '质量状态', value: value.qualityStatus === 2 ? '合格' : value.qualityStatus === 3 ? '不合格' : '待检测' })
    if (kind === 'ENGINEERING_ORDER') rows.push({ label: '生产进度', value: `${value.producedQuantity || 0} / ${value.expectedQuantity || 0}` }, { label: '合格 / 不合格', value: `${value.qualifiedQuantity || 0} / ${value.unqualifiedQuantity || 0}` })
    this.setData({ detail: { id, kindLabel: config[1], name: value.name || value.itemName || config[1], rows }, showScanButton: false })
  },
  reset() {
    this.flow = null
    this.setData({ context: null, unit: null, detail: null, showScanButton: true, buttonText: '开始扫码', lastMessage: '' })
  },
  confirm(content) {
    return new Promise((resolve, reject) => wx.showModal({ title: '操作确认', content, confirmText: '确认', success: result => result.confirm ? resolve() : reject() }))
  }
})
