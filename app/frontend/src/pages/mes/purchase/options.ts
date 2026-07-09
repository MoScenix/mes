import {
  FlowStatus,
  QUALITY_STATUS_PENDING,
  QUALITY_STATUS_QUALIFIED,
  QUALITY_STATUS_UNQUALIFIED,
  STOCK_STATUS_IN_STOCK,
  STOCK_STATUS_OUT_STOCK,
} from '@/api/mesController'

export const flowStatusOptions = [
  { label: '草稿', value: FlowStatus.Draft },
  { label: '待处理', value: FlowStatus.Submitted },
  { label: '已通过', value: FlowStatus.Approved },
  { label: '已拒绝', value: FlowStatus.Rejected },
]

export const stockOptions = [
  { label: '在库', value: STOCK_STATUS_IN_STOCK },
  { label: '不在库', value: STOCK_STATUS_OUT_STOCK },
]

export const qualityOptions = [
  { label: '待检测', value: QUALITY_STATUS_PENDING },
  { label: '合格', value: QUALITY_STATUS_QUALIFIED },
  { label: '不合格', value: QUALITY_STATUS_UNQUALIFIED },
]

export const stockFilterOptions = [...stockOptions]
export const qualityFilterOptions = [...qualityOptions]
