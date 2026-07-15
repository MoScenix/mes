const { get } = require('./utils/request')

App({
  globalData: {
    user: null
  },

  async onLaunch() {
    try {
      const user = await get('/user/get/login')
      this.globalData.user = user
    } catch (error) {
      this.globalData.user = null
    }
  }
})
