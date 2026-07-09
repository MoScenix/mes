export type MesRoleKey =
  | 'purchase'
  | 'worker'
  | 'process_engineer'
  | 'leader'
  | 'warehouse_admin'
  | 'sales'
  | 'admin'

const roleAliases: Record<string, MesRoleKey> = {
  purchase: 'purchase',
  采购专员: 'purchase',
  worker: 'worker',
  普通工人: 'worker',
  process_engineer: 'process_engineer',
  工艺工程师: 'process_engineer',
  leader: 'leader',
  组长: 'leader',
  warehouse: 'warehouse_admin',
  warehouse_admin: 'warehouse_admin',
  仓库管理员: 'warehouse_admin',
  sales: 'sales',
  销售: 'sales',
  admin: 'admin',
  管理员: 'admin',
}

export const normalizeMesRole = (role?: unknown): MesRoleKey => {
  const rawRole = String(role || '').trim()
  if (!rawRole) return 'worker'
  return roleAliases[rawRole] || roleAliases[rawRole.toLowerCase()] || 'worker'
}

export const mesRoleHomePath = (role?: unknown) => {
  const normalized = normalizeMesRole(role)
  if (normalized === 'purchase') return '/mes/purchase'
  if (normalized === 'process_engineer') return '/mes/processes'
  if (normalized === 'leader') return '/mes/leader'
  if (normalized === 'warehouse_admin') return '/mes/warehouse'
  if (normalized === 'sales') return '/mes/sales'
  if (normalized === 'admin') return '/mes/purchase'
  return '/mes/worker'
}
