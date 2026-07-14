<template>
  <a-layout-header class="header">
    <a-row class="header-row" :wrap="false" align="middle">
      <!-- 左侧：Logo和标题 -->
      <a-col class="header-brand-col" flex="160px">
        <RouterLink to="/mes" class="header-left">
          <img class="logo" src="@/assets/logo.png" alt="Logo" />
          <h1 class="site-title">MES</h1>
        </RouterLink>
      </a-col>
      <!-- 中间：导航菜单 -->
      <a-col class="header-center-col" flex="auto">
        <a-menu
          class="desktop-menu"
          v-model:selectedKeys="selectedKeys"
          mode="horizontal"
          :items="menuItems"
          @click="handleMenuClick"
        />
      </a-col>
      <!-- 右侧：扫码 + 邮件 + 用户操作 -->
      <a-col class="header-action-col" flex="none">
        <div class="header-actions">
          <ScanButton />
          <WorkOrderInbox />
          <div v-if="loginUserStore.loginUser.id">
            <a-dropdown>
              <a-space class="user-trigger">
                <a-avatar :src="loginUserStore.loginUser.userAvatar" :size="28" />
                <span class="user-name">{{ loginUserStore.loginUser.userName ?? '无名' }}</span>
              </a-space>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="router.push('/user/center')">
                    <UserOutlined />
                    个人中心
                  </a-menu-item>
                  <a-menu-divider />
                  <a-menu-item @click="doLogout" danger>
                    <LogoutOutlined />
                    退出登录
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
          <div v-else>
            <a-button type="primary" href="/user/login" class="login-btn">登录</a-button>
          </div>
        </div>
      </a-col>
    </a-row>
    <div class="mobile-module-dropdown">
      <a-dropdown :trigger="['click']">
        <button class="mobile-module-trigger" type="button">
          <component :is="mobileModuleIcon" />
          <span>{{ mobileModuleTitle }}</span>
          <DownOutlined />
        </button>
        <template #overlay>
          <a-menu :items="mobileMenuItems" @click="handleMobileMenuClick" />
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { type MenuProps, message } from 'ant-design-vue'
import { useLoginUserStore } from '@/stores/loginUser.ts'
import { userLogout } from '@/api/userController.ts'
import { DownOutlined, HomeOutlined, LogoutOutlined, UserOutlined } from '@ant-design/icons-vue'
import ScanButton from '@/components/ScanButton.vue'
import WorkOrderInbox from '@/components/mes/WorkOrderInbox.vue'
import { normalizeMesRole } from '@/utils/mesRole'
import { mesNavTargetFor, visibleMesNavItems, type MesNavItem } from '@/components/mes/mesNav'

const loginUserStore = useLoginUserStore()
const router = useRouter()
const route = useRoute()
const selectedKeys = ref<string[]>(['/mes'])

router.afterEach((to) => {
  selectedKeys.value = [to.path.startsWith('/mes') ? '/mes' : to.path]
})

const menuItems: MenuProps['items'] = [
  {
    key: '/mes',
    icon: () => h(HomeOutlined),
    label: '工作台',
    title: '工作台',
  },
]

const normalizedRole = computed(() => normalizeMesRole(loginUserStore.loginUser.userRole))

const visibleMobileNavItems = computed(() => visibleMesNavItems(normalizedRole.value))

const targetFor = (item: MesNavItem) => mesNavTargetFor(item, normalizedRole.value)

const isActiveMobileItem = (item: MesNavItem) => {
  const target = targetFor(item)
  if (route.path !== target.path) return false
  if (target.scanMode) return String(route.query.mode || '') === target.scanMode
  return (
    String(route.query.panel || '') === (target.panel || '') &&
    String(route.query.view || '') === target.view
  )
}

const mobileMenuItems = computed<MenuProps['items']>(() =>
  visibleMobileNavItems.value.map((item) => ({
    key: item.key,
    icon: () => h(item.icon),
    label: item.label,
    title: item.label,
  })),
)

const mobileModuleTitle = computed(() => {
  const activeItem = visibleMobileNavItems.value.find(isActiveMobileItem)
  return activeItem?.label || '功能'
})

const mobileModuleIcon = computed(() => {
  const activeItem = visibleMobileNavItems.value.find(isActiveMobileItem)
  return activeItem?.icon || HomeOutlined
})

const handleMenuClick: MenuProps['onClick'] = (e) => {
  const key = e.key as string
  selectedKeys.value = [key]
  if (key.startsWith('/')) {
    router.push(key)
  }
}

const handleMobileMenuClick: MenuProps['onClick'] = (e) => {
  const key = e.key as string
  const item = visibleMobileNavItems.value.find((navItem) => navItem.key === key)
  if (!item) return
  const target = targetFor(item)
  const query = target.scanMode
    ? { mode: target.scanMode }
    : target.panel
      ? { panel: target.panel, view: target.view }
      : { view: target.view }
  if (!isActiveMobileItem(item)) {
    router.push({ path: target.path, query })
  }
}

const doLogout = async () => {
  const res = await userLogout()
  if (res.data.code === 0) {
    loginUserStore.setLoginUser({ userName: '未登录' })
    message.success('退出登录成功')
    await router.push('/user/login')
  } else {
    message.error('退出登录失败，' + res.data.message)
  }
}
</script>

<style scoped>
.header {
  height: 52px;
  background: #fff;
  padding: 0 18px;
  border-bottom: 1px solid var(--border);
  line-height: 52px;
  position: relative;
}

.header-row {
  width: 100%;
  min-width: 0;
}

.header-brand-col {
  flex: 0 0 160px;
  min-width: 0;
}

.header-center-col {
  min-width: 0;
  overflow: hidden;
}

.header-action-col {
  min-width: max-content;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo {
  height: 28px;
  width: 28px;
}

.site-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1d1d1f;
  letter-spacing: -0.01em;
}

.desktop-menu:deep(.ant-menu-horizontal),
.desktop-menu {
  border-bottom: none !important;
  background: transparent;
}

:deep(.ant-menu-item) {
  padding: 0 12px;
  font-size: 13px;
}

:deep(.ant-menu-light.ant-menu-horizontal > .ant-menu-item-selected) {
  color: #0066cc;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.mobile-module-dropdown {
  display: none;
}

.mobile-module-trigger {
  height: 48px;
  max-width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  border: 0;
  padding: 0 4px;
  background: transparent;
  color: #1d1d1f;
  font-size: 13px;
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
  cursor: pointer;
}

.mobile-module-trigger :deep(svg) {
  width: 13px;
  height: 13px;
  color: var(--muted-foreground);
}

.user-trigger {
  cursor: pointer;
  padding: 0 8px;
  height: 36px;
  border-radius: var(--radius);
  transition: background 0.15s;
}

.user-trigger:hover {
  background: var(--muted);
}

.user-name {
  font-size: 13px;
  color: var(--foreground);
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.login-btn {
  height: 32px;
  font-size: 13px;
  border-radius: var(--radius);
}

@media (max-width: 768px) {
  .header {
    height: 48px;
    padding: 0 10px;
    line-height: 48px;
  }

  .header-brand-col {
    flex: 0 0 66px !important;
  }

  .header-center-col {
    flex: 1 1 auto !important;
    overflow: visible;
  }

  .header-action-col {
    flex: 0 0 auto !important;
  }

  .header-left {
    gap: 6px;
  }

  .logo {
    height: 24px;
    width: 24px;
  }

  .site-title {
    font-size: 14px;
  }

  .desktop-menu {
    display: none;
  }

  .mobile-module-dropdown {
    position: absolute;
    top: 0;
    left: 50%;
    z-index: 2;
    display: flex;
    height: 48px;
    max-width: min(180px, calc(100vw - 196px));
    transform: translateX(-50%);
    justify-content: center;
  }

  .mobile-module-trigger {
    max-width: 100%;
  }

  .mobile-module-trigger span {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .header-actions {
    gap: 0;
  }

  .header-actions :deep(.ant-btn) {
    padding-inline: 6px;
  }

  .user-name {
    display: none;
  }
}
</style>
