/**
 * 环境变量配置
 */

// 应用部署域名
export const DEPLOY_DOMAIN = import.meta.env.VITE_DEPLOY_DOMAIN || window.location.origin
export const DEPLOY_BASE_PATH = import.meta.env.VITE_DEPLOY_BASE_PATH || '/static/deploy'

// API 基础地址
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

// 静态资源地址
export const STATIC_BASE_URL = import.meta.env.VITE_STATIC_BASE_URL || '/static'

// 获取部署应用的完整URL
export const getDeployUrl = (deployKey: string) => {
  const origin = DEPLOY_DOMAIN.replace(/\/$/, '')
  const basePath = DEPLOY_BASE_PATH.replace(/^\/?/, '/').replace(/\/$/, '')
  return `${origin}${basePath}/${deployKey}/`
}

// 获取 history 文件静态资源 URL
export const getHistoryDocumentUrl = (historyId: string, fileId: string, filename: string) => {
  return `${STATIC_BASE_URL}/document/${historyId}/${fileId}/${encodeURIComponent(filename)}`
}
