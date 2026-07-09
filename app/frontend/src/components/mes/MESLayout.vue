<template>
  <a-layout class="mes-shell">
    <button class="mobile-menu-button" type="button" aria-label="打开功能栏" @click="drawerOpen = true">
      <MenuOutlined />
    </button>

    <aside class="mes-sidebar">
      <div class="sidebar-title">
        <span class="title-mark"></span>
        <div>
          <strong>MES 工作台</strong>
        </div>
      </div>
      <nav class="side-nav" aria-label="MES 功能">
        <button
          v-for="item in visibleNavItems"
          :key="item.key"
          class="side-nav-item"
          :class="{ active: isActive(item) }"
          type="button"
          @click="go(item)"
        >
          <component :is="item.icon" />
          <span>{{ item.label }}</span>
        </button>
      </nav>
    </aside>

    <a-drawer
      v-model:open="drawerOpen"
      placement="left"
      width="288"
      :closable="false"
      class="mes-mobile-drawer"
    >
      <div class="drawer-head">
        <strong>MES 工作台</strong>
        <a-button shape="circle" size="small" @click="drawerOpen = false">
          <CloseOutlined />
        </a-button>
      </div>
      <nav class="side-nav mobile" aria-label="MES 移动功能">
        <button
          v-for="item in visibleNavItems"
          :key="item.key"
          class="side-nav-item"
          :class="{ active: isActive(item) }"
          type="button"
          @click="go(item)"
        >
          <component :is="item.icon" />
          <span>{{ item.label }}</span>
        </button>
      </nav>
    </a-drawer>

    <a-layout class="mes-main-layout">
      <main class="mes-content">
        <router-view />
      </main>
    </a-layout>

    <FloatingAssistant v-if="route.path !== '/mes/assistant'" />
  </a-layout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  CloseOutlined,
  MenuOutlined,
  InboxOutlined,
  ShoppingCartOutlined,
  ToolOutlined,
  AppstoreOutlined,
  DatabaseOutlined,
  FileAddOutlined,
  FileSearchOutlined,
  FormOutlined,
  SafetyCertificateOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { useLoginUserStore } from '@/stores/loginUser'
import FloatingAssistant from '@/components/mes/FloatingAssistant.vue'
import { normalizeMesRole, type MesRoleKey } from '@/utils/mesRole'

type NavItem = {
  key: string
  path: string
  label: string
  panel?: string
  view: string
  scanMode?: string
  roles: MesRoleKey[]
  icon: typeof ShoppingCartOutlined
}

const navItems: NavItem[] = [
  {
    key: 'purchase-add',
    path: '/mes/purchase',
    panel: 'items',
    view: 'catalog',
    label: '物料',
    roles: ['purchase'],
    icon: ShoppingCartOutlined,
  },
  {
    key: 'purchase-units',
    path: '/mes/purchase',
    panel: 'itemUnits',
    view: 'units',
    label: '库存单体',
    roles: ['purchase'],
    icon: ToolOutlined,
  },
  {
    key: 'purchase-scan',
    path: '/mes/scan',
    view: 'inbound',
    scanMode: 'inbound',
    label: '扫描入库',
    roles: ['purchase'],
    icon: InboxOutlined,
  },
  {
    key: 'purchase-flow',
    path: '/mes/purchase',
    panel: 'flows',
    view: 'purchase',
    label: '流转单',
    roles: ['purchase'],
    icon: FileSearchOutlined,
  },
  {
    key: 'worker-add',
    path: '/mes/worker',
    panel: 'itemUnits',
    view: 'units',
    label: '新增单品',
    roles: ['worker'],
    icon: ToolOutlined,
  },
  {
    key: 'worker-receive',
    path: '/mes/scan',
    view: 'receive',
    scanMode: 'receive',
    label: '领取货物',
    roles: ['worker', 'sales'],
    icon: InboxOutlined,
  },
  {
    key: 'worker-inspect',
    path: '/mes/scan',
    view: 'inspect',
    scanMode: 'inspect',
    label: '检验单品',
    roles: ['worker'],
    icon: SafetyCertificateOutlined,
  },
  {
    key: 'process-engineer-processes',
    path: '/mes/processes',
    panel: 'processes',
    view: 'processes',
    label: '工艺管理',
    roles: ['process_engineer'],
    icon: AppstoreOutlined,
  },
  {
    key: 'leader-engineering',
    path: '/mes/leader',
    panel: 'engineering',
    view: 'engineering',
    label: '工程单',
    roles: ['leader'],
    icon: AppstoreOutlined,
  },
  {
    key: 'leader-workorder',
    path: '/mes/leader',
    panel: 'workOrders',
    view: 'workOrders',
    label: '发工单',
    roles: ['leader'],
    icon: FormOutlined,
  },
  {
    key: 'warehouse-audit',
    path: '/mes/warehouse',
    panel: 'audit',
    view: 'audit',
    label: '审批流转单',
    roles: ['warehouse_admin'],
    icon: InboxOutlined,
  },
  {
    key: 'warehouse-inventory',
    path: '/mes/warehouse',
    panel: 'inventory',
    view: 'inventory',
    label: '物资情况',
    roles: ['warehouse_admin'],
    icon: DatabaseOutlined,
  },
  {
    key: 'warehouse-workorder',
    path: '/mes/warehouse',
    panel: 'workOrders',
    view: 'workOrders',
    label: '发工单',
    roles: ['warehouse_admin'],
    icon: FormOutlined,
  },
  {
    key: 'sales-apply',
    path: '/mes/sales',
    panel: 'flows',
    view: 'flows',
    label: '流转单',
    roles: ['sales'],
    icon: FileAddOutlined,
  },
  {
    key: 'admin-users',
    path: '/mes/admin/users',
    view: 'admin-users',
    label: '员工管理',
    roles: ['admin'],
    icon: UserOutlined,
  },
]

const router = useRouter()
const route = useRoute()
const loginUserStore = useLoginUserStore()
const drawerOpen = ref(false)

const normalizedRole = computed(() => normalizeMesRole(loginUserStore.loginUser.userRole))

const visibleNavItems = computed(() => {
	const role = normalizedRole.value
	if (role === 'admin') return navItems
	return navItems.filter((item) => item.roles.includes(role))
})

const isActive = (item: NavItem) => {
  if (route.path !== item.path) return false
  if (item.scanMode) return String(route.query.mode || '') === item.scanMode
  return String(route.query.panel || '') === (item.panel || '') && String(route.query.view || '') === item.view
}

const go = async (item: NavItem) => {
  drawerOpen.value = false
  const query = item.scanMode ? { mode: item.scanMode } : item.panel ? { panel: item.panel, view: item.view } : { view: item.view }
  const active = item.scanMode
    ? route.path === item.path && String(route.query.mode || '') === item.scanMode
    : route.path === item.path && String(route.query.panel || '') === (item.panel || '') && String(route.query.view || '') === item.view
  if (!active) {
    await router.push({ path: item.path, query })
  }
}

</script>

<style scoped>
.mes-shell {
  min-height: calc(100vh - 52px);
  display: flex;
  flex-direction: row;
  align-items: stretch;
  background: #fafafa;
  color: #1d1d1f;
  position: relative;
}

.mes-sidebar {
  width: 200px;
  flex: 0 0 200px;
  min-height: calc(100vh - 52px);
  padding: 20px 10px;
  background: #fff;
  border-right: 1px solid var(--border);
}

.sidebar-title {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 10px 18px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 12px;
}

.title-mark {
  width: 4px;
  height: 20px;
  border-radius: 2px;
  background: var(--primary);
}

.sidebar-title strong,
.drawer-head strong {
  display: block;
  font-size: 14px;
  line-height: 1.3;
  font-weight: 600;
  color: var(--foreground);
}

.side-nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.side-nav-item {
  width: 100%;
  min-height: 36px;
  display: flex;
  align-items: center;
  gap: 10px;
  border: 0;
  border-radius: 6px;
  padding: 0 10px;
  background: transparent;
  color: var(--muted-foreground);
  font-size: 13px;
  line-height: 1.2;
  text-align: left;
  cursor: pointer;
  transition: all 0.15s ease;
}

.side-nav-item:hover {
  background: var(--muted);
  color: var(--foreground);
}

.side-nav-item.active {
  background: var(--muted);
  color: var(--primary);
  font-weight: 500;
}

.side-nav-item span {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.mes-main-layout {
  flex: 1 1 auto;
  min-width: 0;
  background: transparent;
}

.mes-content {
  min-width: 0;
  padding: 20px 24px 96px;
}

.mobile-menu-button {
  display: none;
}

.drawer-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
}

.mobile.side-nav {
  gap: 4px;
}

@media (max-width: 768px) {
  .mes-shell {
    min-height: calc(100vh - 48px);
  }

  .mes-sidebar {
    display: none;
  }

  .mobile-menu-button {
    position: fixed;
    left: 12px;
    top: 60px;
    z-index: 30;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border: 1px solid var(--border);
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.92);
    color: #1d1d1f;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
  }

  .mes-content {
    padding: 14px 12px 80px;
  }
}
</style>
