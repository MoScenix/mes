const { post } = require('../../utils/request')

Page({
  data: { userAccount: '', userPassword: '', checkPassword: '', submitting: false },
  onAccount(event) { this.setData({ userAccount: event.detail.value }) },
  onPassword(event) { this.setData({ userPassword: event.detail.value }) },
  onCheckPassword(event) { this.setData({ checkPassword: event.detail.value }) },
  async submit() {
    const { userAccount, userPassword, checkPassword } = this.data
    if (!userAccount || userPassword.length < 8) {
      wx.showToast({ title: '请填写有效账号和密码', icon: 'none' })
      return
    }
    if (userPassword !== checkPassword) {
      wx.showToast({ title: '两次输入密码不一致', icon: 'none' })
      return
    }
    this.setData({ submitting: true })
    try {
      await post('/user/register', { userAccount, userPassword, checkPassword })
      wx.showToast({ title: '注册成功', icon: 'success' })
      setTimeout(() => wx.redirectTo({ url: '/pages/login/login' }), 500)
    } catch (error) {
      wx.showToast({ title: error.message, icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  },
  goLogin() { wx.navigateBack() }
})
