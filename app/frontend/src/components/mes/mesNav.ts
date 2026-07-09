import {
  AppstoreOutlined,
  DatabaseOutlined,
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
}

export type MesNavItem = {
  key: string
  label: string
  target: MesNavTarget | ((role: MesRoleKey) => MesNavTarget)
  roles: MesRoleKey[]
  icon: typeof ShoppingCartOutlined
}

export const mesNavItems: MesNavItem[] = [
  {
    key: 'purchase-add',
    target: (role) => {
      if (role === 'process_engineer') return { path: '/mes/processes', panel: 'items', view: 'catalog' }
      return { path: '/mes/purchase', panel: 'items', view: 'catalog' }
    },
    label: '物料',
    roles: ['purchase', 'process_engineer'],
    icon: ShoppingCartOutlined,
  },
  {
    key: 'purchase-units',
    target: { path: '/mes/purchase', panel: 'itemUnits', view: 'units' },
    label: '库存单体',
    roles: ['purchase'],
    icon: DatabaseOutlined,
  },
  {
    key: 'purchase-scan',
    target: { path: '/mes/scan', view: 'inbound', scanMode: 'inbound' },
    label: '扫描入库',
    roles: ['purchase'],
    icon: InboxOutlined,
  },
  {
    key: 'inventory-flow',
    target: (role) => {
      if (role === 'sales') return { path: '/mes/sales', panel: 'flows', view: 'flows' }
      if (role === 'leader') return { path: '/mes/leader', panel: 'flows', view: 'flows' }
      return { path: '/mes/purchase', panel: 'flows', view: 'purchase' }
    },
    label: '流转单',
    roles: ['purchase', 'sales', 'leader'],
    icon: FileSearchOutlined,
  },
  {
    key: 'worker-receive',
    target: { path: '/mes/scan', view: 'receive', scanMode: 'receive' },
    label: '领取货物',
    roles: ['worker', 'sales'],
    icon: InboxOutlined,
  },
  {
    key: 'worker-inspect',
    target: { path: '/mes/scan', view: 'inspect', scanMode: 'inspect' },
    label: '检验单品',
    roles: ['worker'],
    icon: SafetyCertificateOutlined,
  },
  {
    key: 'process-engineer-processes',
    target: { path: '/mes/processes', panel: 'processes', view: 'processes' },
    label: '工艺管理',
    roles: ['process_engineer'],
    icon: AppstoreOutlined,
  },
  {
    key: 'leader-engineering',
    target: { path: '/mes/leader', panel: 'engineering', view: 'engineering' },
    label: '工程单',
    roles: ['leader'],
    icon: AppstoreOutlined,
  },
  {
    key: 'warehouse-audit',
    target: { path: '/mes/warehouse', panel: 'audit', view: 'audit' },
    label: '审批流转单',
    roles: ['warehouse_admin'],
    icon: InboxOutlined,
  },
  {
    key: 'warehouse-flow',
    target: { path: '/mes/warehouse', panel: 'flows', view: 'flows' },
    label: '流转单',
    roles: ['warehouse_admin'],
    icon: FileSearchOutlined,
  },
  {
    key: 'warehouse-inventory',
    target: { path: '/mes/warehouse', panel: 'inventory', view: 'inventory' },
    label: '物资情况',
    roles: ['warehouse_admin'],
    icon: DatabaseOutlined,
  },
  {
    key: 'admin-users',
    target: { path: '/mes/admin/users', view: 'admin-users' },
    label: '员工管理',
    roles: ['admin'],
    icon: UserOutlined,
  },
]

export const visibleMesNavItems = (role: MesRoleKey) => {
  if (role === 'admin') return mesNavItems
  return mesNavItems.filter((item) => item.roles.includes(role))
}

export const mesNavTargetFor = (item: MesNavItem, role: MesRoleKey) => {
  return typeof item.target === 'function' ? item.target(role) : item.target
}
