<template>
  <a-layout-header class="header">
    <a-row :wrap="false" align="middle">
      <!-- 左侧：Logo和标题 -->
      <a-col flex="160px">
        <RouterLink to="/mes" class="header-left">
          <img class="logo" src="@/assets/logo.png" alt="Logo" />
          <h1 class="site-title">MES</h1>
        </RouterLink>
      </a-col>
      <!-- 中间：导航菜单 -->
      <a-col flex="auto">
        <a-menu v-model:selectedKeys="selectedKeys" mode="horizontal" :items="menuItems" @click="handleMenuClick" />
      </a-col>
      <!-- 右侧：扫码 + 邮件 + 用户操作 -->
      <a-col>
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
  </a-layout-header>
</template>

<script setup lang="ts">
import { h, ref } from 'vue'
import { useRouter } from 'vue-router'
import { type MenuProps, message } from 'ant-design-vue'
import { useLoginUserStore } from '@/stores/loginUser.ts'
import { userLogout } from '@/api/userController.ts'
import { LogoutOutlined, HomeOutlined, UserOutlined } from '@ant-design/icons-vue'
import ScanButton from '@/components/ScanButton.vue'
import WorkOrderInbox from '@/components/mes/WorkOrderInbox.vue'

const loginUserStore = useLoginUserStore()
const router = useRouter()
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

const handleMenuClick: MenuProps['onClick'] = (e) => {
  const key = e.key as string
  selectedKeys.value = [key]
  if (key.startsWith('/')) {
    router.push(key)
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
  padding: 0 20px;
  border-bottom: 1px solid var(--border);
  line-height: 52px;
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

:deep(.ant-menu-horizontal) {
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
    padding: 0 12px;
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

  :deep(.ant-menu) {
    display: none;
  }

  .user-name {
    display: none;
  }
}
</style>
