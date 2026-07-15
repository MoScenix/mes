const { baseUrl } = require('./config')

let cookie = wx.getStorageSync('mesCookie') || ''

function request(path, method = 'GET', data = {}) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: `${baseUrl}${path}`,
      method,
      data,
      header: {
        'content-type': 'application/json',
        ...(cookie ? { Cookie: cookie } : {})
      },
      success(response) {
        const setCookie = response.header['Set-Cookie'] || response.header['set-cookie']
        if (setCookie) {
          cookie = setCookie.split(';')[0]
          wx.setStorageSync('mesCookie', cookie)
        }
        const body = response.data || {}
        if (response.statusCode === 401 || body.code === 40100) {
          cookie = ''
          wx.removeStorageSync('mesCookie')
          reject(new Error('请先登录'))
          return
        }
        if (response.statusCode < 200 || response.statusCode >= 300 || body.code !== 0) {
          reject(new Error(body.message || `请求失败 (${response.statusCode})`))
          return
        }
        resolve(body.data)
      },
      fail(error) {
        reject(new Error(error.errMsg || '网络连接失败'))
      }
    })
  })
}

module.exports = {
  get(path, data) { return request(path, 'GET', data) },
  post(path, data) { return request(path, 'POST', data) },
  clearSession() {
    cookie = ''
    wx.removeStorageSync('mesCookie')
  }
}
