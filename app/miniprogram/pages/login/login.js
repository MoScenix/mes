const { get, post } = require('../../utils/request')

Page({
  data: { userAccount: '', userPassword: '', submitting: false },
  async onShow() {
    try {
      const user = await get('/user/get/login')
      if (user && user.id) {
        getApp().globalData.user = user
        wx.reLaunch({ url: '/pages/home/home' })
      }
    } catch (error) {}
  },
  onAccount(event) { this.setData({ userAccount: event.detail.value }) },
  onPassword(event) { this.setData({ userPassword: event.detail.value }) },
  async submit() {
    const { userAccount, userPassword } = this.data
    if (!userAccount || userPassword.length < 8) {
      wx.showToast({ title: '请填写有效账号和密码', icon: 'none' })
      return
    }
    this.setData({ submitting: true })
    try {
      const user = await post('/user/login', { userAccount, userPassword })
      getApp().globalData.user = user
      wx.reLaunch({ url: '/pages/home/home' })
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  },
  goRegister() { wx.navigateTo({ url: '/pages/register/register' }) }
})
