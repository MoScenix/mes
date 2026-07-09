import { createRouter, createWebHistory } from 'vue-router'
import UserLoginPage from '@/pages/user/UserLoginPage.vue'
import UserRegisterPage from '@/pages/user/UserRegisterPage.vue'
import UserManagePage from '@/pages/admin/UserManagePage.vue'
import AppChatPage from '@/pages/app/AppChatPage.vue'
import AppEditPage from '@/pages/app/AppEditPage.vue'
import UserCenterPage from '@/pages/user/UserCenterPage.vue'
import MESLayout from '@/components/mes/MESLayout.vue'
import PurchaseWorkspace from '@/pages/mes/PurchaseWorkspace.vue'
import WorkerWorkspace from '@/pages/mes/WorkerWorkspace.vue'
import LeaderWorkspace from '@/pages/mes/LeaderWorkspace.vue'
import WarehouseWorkspace from '@/pages/mes/WarehouseWorkspace.vue'
import SalesWorkspace from '@/pages/mes/SalesWorkspace.vue'
import ProcessWorkspace from '@/pages/mes/ProcessWorkspace.vue'
import MESAssistantPage from '@/pages/mes/MESAssistantPage.vue'
import MESHomePage from '@/pages/mes/MESHomePage.vue'
import MESDetailPage from '@/pages/mes/MESDetailPage.vue'
import MESWorkOrdersPage from '@/pages/mes/MESWorkOrdersPage.vue'
import MESCreatePage from '@/pages/mes/MESCreatePage.vue'
import MESScanPage from '@/pages/mes/MESScanPage.vue'
import MESProcessEngOrdersPage from '@/pages/mes/MESProcessEngOrdersPage.vue'
import { useLoginUserStore } from '@/stores/loginUser'

async function requireLogin(to: any) {
  const loginUserStore = useLoginUserStore()
  if (!loginUserStore.loginUser.id) {
    try {
      await loginUserStore.fetchLoginUser()
    } catch (error) {
      console.warn('获取登录用户失败', error)
    }
  }
  if (!loginUserStore.loginUser.id) {
    return `/user/login?redirect=${encodeURIComponent(to.fullPath)}`
  }
  return true
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: '主页',
      redirect: '/mes',
    },
    {
      path: '/user/login',
      name: '用户登录',
      component: UserLoginPage,
    },
    {
      path: '/user/register',
      name: '用户注册',
      component: UserRegisterPage,
    },
    {
      path: '/admin/userManage',
      redirect: '/mes/admin/users',
    },
    {
      path: '/app/chat/:id',
      name: '应用对话',
      component: AppChatPage,
    },
    {
      path: '/app/edit/:id',
      name: '编辑应用',
      component: AppEditPage,
    },
    {
      path: '/user/center',
      name: '个人中心',
      component: UserCenterPage,
    },
    {
      path: '/mes',
      name: 'MES 工作台',
      component: MESLayout,
      beforeEnter: requireLogin,
      children: [
        {
          path: '',
          name: 'MES 工作台入口',
          component: MESHomePage,
        },
        {
          path: 'purchase',
          name: 'MES 采购',
          component: PurchaseWorkspace,
        },
        {
          path: 'worker',
          name: 'MES 普通员工',
          component: WorkerWorkspace,
        },
        {
          path: 'processes',
          name: 'MES 工艺管理',
          component: ProcessWorkspace,
        },
        {
          path: 'leader',
          name: 'MES 组长',
          component: LeaderWorkspace,
        },
        {
          path: 'warehouse',
          name: 'MES 仓库',
          component: WarehouseWorkspace,
        },
        {
          path: 'sales',
          name: 'MES 销售',
          component: SalesWorkspace,
        },
        {
          path: 'assistant',
          name: 'MES 助手',
          component: MESAssistantPage,
        },
        {
          path: 'detail',
          name: 'MES 详情',
          component: MESDetailPage,
        },
        {
          path: 'scan',
          name: 'MES 扫码',
          component: MESScanPage,
        },
        {
          path: 'workorders',
          name: 'MES 工单',
          component: MESWorkOrdersPage,
        },
        {
          path: 'process-eng-orders',
          name: '工艺关联工程单',
          component: MESProcessEngOrdersPage,
        },
        {
          path: 'create',
          name: 'MES 新建',
          component: MESCreatePage,
        },
        {
          path: 'admin/users',
          name: 'MES 员工管理',
          component: UserManagePage,
        },
      ],
    },
  ],
})

export default router
