const { get, post, clearSession } = require('../../utils/request')
const { scan } = require('../../utils/mes-code')

Page({
  data: { userName: 'MES 用户' },
  async onShow() {
    try {
      const user = await get('/user/get/login')
      getApp().globalData.user = user
      this.setData({ userName: user.userName || user.name || user.userAccount || 'MES 用户' })
    } catch (error) {
      wx.reLaunch({ url: '/pages/login/login' })
    }
  },
  openMode(event) {
    wx.navigateTo({ url: `/pages/scan/scan?mode=${event.currentTarget.dataset.mode}` })
  },
  async openGeneralScan() {
    try {
      const parsed = await scan()
      wx.navigateTo({ url: `/pages/scan/scan?mode=detail&kind=${parsed.kind}&id=${parsed.id}` })
    } catch (error) {
      if (error) wx.showToast({ title: error.message, icon: 'none' })
    }
  },
  async logout() {
    try { await post('/user/logout') } catch (error) {}
    clearSession()
    getApp().globalData.user = null
    wx.reLaunch({ url: '/pages/login/login' })
  }
})
