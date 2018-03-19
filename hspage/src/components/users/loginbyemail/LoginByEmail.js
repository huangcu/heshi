export default {
  name: 'Login',
  data: function (params) {
    return {
      login_feedback: '',
      upgrade_feedback: '',
      username: '',
      password: '',
      account: '',
    }
  },
  computed: {
    referer: function () {
      // `this` points to the vm instance
      if (this.$route.query.referer === undefined) {
        return ''
      }
      return this.$route.query.referer
    },
    appointment: function () {
      if (this.$route.query.for === undefined) {
        return ''
      }
      return this.$route.query.for
    }
  },
  methods: {
    reference: function () {
      if (this.referer) {
        if (this.appointment) {
          if (this.appointment === 'appointment') {
            this.$http.headers.common['Location'] = 'appointment'
          } else {
            this.$http.headers.common['Location'] = 'myaccount'
          }
        } else {
          this.$http.headers.common['Location'] = 'myaccount'
        }
      }
    },
    login: function () {
      var formData = new FormData()
      formData.append('username', this.username)
      formData.append('password', this.password)

      this.$http.post(
        this.$userURL + '/login',
        formData,
        {
          headers: {
              'Content-Type': 'multipart/form-data'
            }
        }).then(response => {
          if (response.status === 200) {
            var loginResult = JSON.parse(response.bodyText)
            if (loginResult.code !== 200 ) {
              this.login_feedback = loginResult.message
              return
            }
            this.$cookies.set('token', loginResult.token)
            var userprofile = JSON.parse(loginResult.userprofile)
            this.$cookies.set('_account', userprofile.id)
            this.$cookies.set('userprofile', loginResult.userprofile)
            this.$router.replace('/home')
          }
        }, err => { console.log(err); alert('error:' + err.body) })
    },
    logout: function () {
      this.$http.post(
        this.$userURL + '/logout'
      ).then(response => {
        console.log('logout')
        this.$route.router.go('/')
      }, err => { console.log(err); alert('error:' + err.body) })
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
  },
  created()  {
    this.$currentPage = 'LOGIN'
  },
  mounted() {
    if (this.$cookies.isKey('_account')) {
      this.account = this.$cookies.get('_account')
    } else {
      this.account = ''
    }
  }
}
