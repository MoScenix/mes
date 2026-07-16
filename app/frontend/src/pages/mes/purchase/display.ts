import dayjs from 'dayjs'
import {
  FLOW_TYPE_IN,
  FLOW_TYPE_OUT,
  QUALITY_STATUS_PENDING,
  QUALITY_STATUS_QUALIFIED,
  QUALITY_STATUS_UNQUALIFIED,
  STOCK_STATUS_IN_STOCK,
  STOCK_STATUS_OUT_STOCK,
} from '@/api/mesController'

export const formatTime = (t?: string) => (t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-')

export const flowTypeLabel = (s?: number) =>
  s === FLOW_TYPE_IN ? '入库' : s === FLOW_TYPE_OUT ? '出库' : '未知'

export const flowBusinessLabel = (s?: number) =>
  s === 1 ? '采购入库' : s === 2 ? '申请货物' : s === 3 ? '生产入库' : '未知'

export const flowStatusColor = (s?: number) => {
  if (s === 1) return 'default'
  if (s === 2) return 'blue'
  if (s === 3) return 'green'
  if (s === 4) return 'red'
  return 'default'
}

export const flowStatusLabel = (s?: number) => {
  if (s === 1) return '草稿'
  if (s === 2) return '待处理'
  if (s === 3) return '已通过'
  if (s === 4) return '已拒绝'
  return '未知'
}

export const stockLabel = (s?: number) =>
  s === STOCK_STATUS_IN_STOCK ? '在库' : s === STOCK_STATUS_OUT_STOCK ? '不在库' : '未知'

export const qualityLabel = (s?: number) =>
  s === QUALITY_STATUS_PENDING
    ? '待检测'
    : s === QUALITY_STATUS_QUALIFIED
      ? '合格'
      : s === QUALITY_STATUS_UNQUALIFIED
        ? '不合格'
        : '未知'
