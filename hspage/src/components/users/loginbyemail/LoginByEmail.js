export default {
  name: 'Login',
  data: function (params) {
    return {
      login_feedback: '',
      upgrade_feedback: ''
    }
  },
  props: {
    _ref: String,
    _for: String,
    _account: String,
    loginResult: null
  },
  methods: {
    reference: function () {
      if (this.referer) {
        if (this.for) {
          if (this.for === 'appointment') {
            this.$http.headers.common['Location'] = 'appointment'
          } else {
            this.$http.headers.common['Location'] = 'myaccount'
          }
        } else {
          this.$http.headers.common['Location'] = 'myaccount'
        }
      }
    },
    login: function (user, password) {
      var formData = new FormData()
      formData.append('username', user)
      formData.append('password', password)

      this.$http.post(
        this.$userURL + 'login',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        }).then(response => {
          if (response.status === 200) {
            // token
            this.loginResult = JSON.parse(response.bodyText)
            this.$cookies.set('beyou', this.loginResult.id)
            this.$cookies.set('token', this.loginResult.token)
          }
          return response.bodyText
        }, err => { console.log(err); alert('error:' + err.bodyText) })
    },
    logout: function () {
      this.$http.post(
        this.$userURL + 'logout'
      ).then(response => {
        console.log('logout')
        this.$route.router.go('/')
      }, err => { console.log(err); alert('error:' + err.bodyText) })
    },
    getQRCode: function () {
      //MAX 1 - 4294967295
      var sceneID = Math.floor(Math.random() * 1000 * 1000 * 1000)
      this.$http.get(
        this.$wechatURL + '/temp_qrcode?sceneID=' + sceneID
      ).then(response => {
        return response.body
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    returnPic: function () {
      return 'https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=gQGo8DwAAAAAAAAAAS5odHRwOi8vd2VpeGluLnFxLmNvbS9xLzAyR3UtUEVDclFleGoxSWxDRnhxMTAAAgSd5alaAwR4AAAA'
    }
  }
}