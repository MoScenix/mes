import { useLoginUserStore } from '@/stores/loginUser'
import { message } from 'ant-design-vue'
import router from '@/router'

const publicPathPrefixes = ['/user/login', '/user/register']
const protectedPathPrefixes = ['/mes', '/admin', '/app', '/user/center']

const matchesPathPrefix = (path: string, prefixes: string[]) =>
  prefixes.some((prefix) => path === prefix || path.startsWith(`${prefix}/`))

/**
 * 全局权限校验
 */
router.beforeEach(async (to, from, next) => {
  if (matchesPathPrefix(to.path, publicPathPrefixes)) {
    next()
    return
  }

  const loginUserStore = useLoginUserStore()
  let loginUser = loginUserStore.loginUser

  if (matchesPathPrefix(to.path, protectedPathPrefixes) && !loginUser.id) {
    try {
      await loginUserStore.fetchLoginUser()
    } catch (error) {
      console.warn('获取登录用户失败', error)
    }
    loginUser = loginUserStore.loginUser
  }

  const toUrl = to.fullPath
  if (toUrl.startsWith('/admin') || toUrl.startsWith('/mes/admin')) {
    if (!loginUser || loginUser.userRole !== 'admin') {
      message.error('没有权限')
      next(`/user/login?redirect=${to.fullPath}`)
      return
    }
  }
  next()
})
