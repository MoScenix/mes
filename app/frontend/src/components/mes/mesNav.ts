import {
  AppstoreOutlined,
  DatabaseOutlined,
  DashboardOutlined,
  FileSearchOutlined,
  InboxOutlined,
  SafetyCertificateOutlined,
  ShoppingCartOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import type { MesRoleKey } from '@/utils/mesRole'

export type MesNavTarget = {
  path: string
  panel?: string
  view: string
  scanMode?: string
  businessType?: string
}

export type MesNavItem = {
  key: string
  label: string
  target: MesNavTarget | ((role: MesRoleKey) => MesNavTarget)
  roles: MesRoleKey[]
  icon: typeof ShoppingCartOutlined
  group: MesNavGroup
}

export type MesNavGroup =
  | '原料采购'
  | '生产管理'
  | '仓库管理'
  | '工艺管理'
  | '扫码作业'
  | '系统管理'

export const mesNavItems: MesNavItem[] = [
  {
    key: 'production-dashboard',
    target: { path: '/mes/dashboard', view: 'dashboard' },
    label: '智能生产看板',
    roles: ['leader'],
    icon: DashboardOutlined,
    group: '生产管理',
  },
  {
    key: 'purchase-add',
    target: (role) => {
      if (role === 'process_engineer')
        return { path: '/mes/processes', panel: 'items', view: 'catalog' }
      return { path: '/mes/purchase', panel: 'items', view: 'catalog' }
    },
    label: '物料',
    roles: ['purchase', 'process_engineer'],
    icon: ShoppingCartOutlined,
    group: '原料采购',
  },
  {
    key: 'purchase-units',
    target: { path: '/mes/purchase', panel: 'itemUnits', view: 'units' },
    label: '库存单体',
    roles: ['purchase'],
    icon: DatabaseOutlined,
    group: '原料采购',
  },
  {
    key: 'purchase-scan',
    target: { path: '/mes/scan', view: 'inbound', scanMode: 'inbound' },
    label: '扫描入库',
    roles: ['purchase', 'worker'],
    icon: InboxOutlined,
    group: '扫码作业',
  },
  {
    key: 'inventory-flow',
    target: { path: '/mes/purchase', panel: 'flows', view: 'purchase', businessType: '1' },
    label: '采购入库',
    roles: ['purchase'],
    icon: FileSearchOutlined,
    group: '原料采购',
  },
  {
    key: 'material-request',
    target: { path: '/mes/leader', panel: 'flows', view: 'flows', businessType: '2' },
    label: '申请货物',
    roles: ['leader'],
    icon: FileSearchOutlined,
    group: '生产管理',
  },
  {
    key: 'production-inbound',
    target: { path: '/mes/leader', panel: 'flows', view: 'flows', businessType: '3' },
    label: '成品入库',
    roles: ['leader'],
    icon: InboxOutlined,
    group: '生产管理',
  },
  {
    key: 'worker-receive',
    target: { path: '/mes/scan', view: 'receive', scanMode: 'receive' },
    label: '领取货物',
    roles: ['worker', 'sales'],
    icon: InboxOutlined,
    group: '扫码作业',
  },
  {
    key: 'worker-inspect',
    target: { path: '/mes/scan', view: 'inspect', scanMode: 'inspect' },
    label: '检验单品',
    roles: ['worker'],
    icon: SafetyCertificateOutlined,
    group: '扫码作业',
  },
  {
    key: 'process-engineer-processes',
    target: { path: '/mes/processes', panel: 'processes', view: 'processes' },
    label: '工艺管理',
    roles: ['process_engineer'],
    icon: AppstoreOutlined,
    group: '工艺管理',
  },
  {
    key: 'leader-engineering',
    target: { path: '/mes/leader', panel: 'engineering', view: 'engineering' },
    label: '生产计划',
    roles: ['leader'],
    icon: AppstoreOutlined,
    group: '生产管理',
  },
  {
    key: 'warehouse-audit',
    target: { path: '/mes/warehouse', panel: 'audit', view: 'audit' },
    label: '审批流水',
    roles: ['warehouse_admin'],
    icon: InboxOutlined,
    group: '仓库管理',
  },
  {
    key: 'warehouse-flow',
    target: { path: '/mes/warehouse', panel: 'flows', view: 'flows' },
    label: '流水审计',
    roles: ['warehouse_admin'],
    icon: FileSearchOutlined,
    group: '仓库管理',
  },
  {
    key: 'warehouse-inventory',
    target: { path: '/mes/warehouse', panel: 'inventory', view: 'inventory' },
    label: '库存总览',
    roles: ['warehouse_admin'],
    icon: DatabaseOutlined,
    group: '仓库管理',
  },
  {
    key: 'admin-users',
    target: { path: '/mes/admin/users', view: 'admin-users' },
    label: '员工管理',
    roles: ['admin'],
    icon: UserOutlined,
    group: '系统管理',
  },
]

export const mesNavGroups: MesNavGroup[] = [
  '原料采购',
  '生产管理',
  '仓库管理',
  '工艺管理',
  '扫码作业',
  '系统管理',
]

export const visibleMesNavItems = (role: MesRoleKey) => {
  if (role === 'admin') return mesNavItems
  return mesNavItems.filter((item) => item.roles.includes(role))
}

export const mesNavTargetFor = (item: MesNavItem, role: MesRoleKey) => {
  return typeof item.target === 'function' ? item.target(role) : item.target
}
