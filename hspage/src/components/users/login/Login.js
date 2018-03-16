export default {
  name: 'Login',
  data: function () {
    return {
      qrCodeSrc: '',
      login_feedback: '',
      upgrade_feedback: '',
      account: '',
      previewmode: ''
    }
  },
  props: {
    // referer: String,
    for: String,
    loginResult: null
  },
  computed: {
    referer: function () {
      // `this` points to the vm instance
      console.log(this.$route.params.referer)
      if (this.$route.params.referer === undefined) {
        return ''
      }
      return ''
    },
    appointment: function () {
      console.log(this.$route.params.referer)
      if (this.$route.params.referer === undefined) {
        return ''
      }
      return ''
    }
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
        this.$userURL + '/login',
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
        this.$userURL + '/logout'
      ).then(response => {
        console.log('logout')
        this.$route.replace('/')
      }, err => { console.log(err); alert('error:' + err.bodyText) })
    },
    getQRCode: function () {
      //MAX 1 - 4294967295
      var sceneID = Math.floor(Math.random() * 1000 * 1000 *1000)
      this.$http.get(
        this.$wechatURL + '/temp_qrcode?sceneID='+sceneID,
      ).then(response => {
        this.$cookies.set('sceneID', sceneID, 60 * 2)
        window.sessionStorage.setItem("HSSESSIONID", sceneID)
        this.qrCodeSrc = response.body
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    isQRCodeScanned: function () {
      if (this.$cookies.isKey('sceneID')) {
        this.$http.get(
          this.$wechatURL + '/status?sceneID=' + this.$cookies.get('sceneID')
        ).then(response => {
          if (response.body !== '') {
            this.cookies.set("openID", response.body, 60 * 30)
            this.$route.replace('/')
          }
        }, err => { console.log(err); alert('error:' + err.body) })
      }
    }
  },
  mounted() {
    if (this.$cookies.isKey('_account')) {
      this.account = this.$cookies.get('_account')
    } else {
      this.account = ''
      this.getQRCode()
      setInterval(function () {
        this.getQRCode();
      }.bind(this), 2 * 60 * 1000);
      setInterval(function () {
        this.isQRCodeScanned();
      }.bind(this), 5000);
    }
    // this.getInterestedItems()
    // this.getOrders()
  },
  created() {
    this.$currentPage = 'LOGIN'
  }
}
