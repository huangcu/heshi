export default {
  name: 'Login',
  data: function () {
    return {
      qrCodeSrc: '',
      login_feedback: '',
      upgrade_feedback: '',
      userprofile: '',
      previewmode: '',
      sceneID: '',
      QRCodeHandle: null,
      QRCodeStatusHandle: null
    }
  },
  props: {
    loginResult: null
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
        }
      ).then(response => {
        if (response.status === 200) {
          // token
          this.loginResult = JSON.parse(response.body)
          this.$cookies.set('beyou', this.loginResult.id)
          this.$cookies.set('token', this.loginResult.token)
        }
        return response.body
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    logout: function () {
      this.$http.post(
        this.$userURL + '/logout'
      ).then(response => {
        console.log('logout')
        this.$router.replace('/')
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    getQRCode: function () {
      //MAX 1 - 4294967295
      var sceneID = Math.floor(Math.random() * 1000 * 1000 *1000)
      this.$http.get(
        this.$wechatURL + '/temp_qrcode?sceneID='+sceneID,
      ).then(response => {
        this.sceneID = sceneID
        this.qrCodeSrc = response.body
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    isQRCodeScanned: function () {
      if (this.sceneID !=='') {
        this.$http.get(
          this.$wechatURL + '/status?sceneID=' + this.sceneID
        ).then(response => {
          if (response.body !== '') {
            this.$cookies.set("wechatopenID", response.body, 60 * 30)
            this.$router.replace('/qrsign/'+response.body+'?referer='+this.referer)
          }
        }, err => { console.log(err); alert('error:' + err.body) })
      }
    }
  },
  mounted() {
    if (this.$cookies.isKey('userprofile')) {
      this.userprofile = JSON.parse(this.$cookies.get('userprofile'))
    } else {
      this.userprofile = ''
      this.getQRCode()
      this.QRCodeHandle = setInterval(function () {
        this.getQRCode();
      }.bind(this), 2 * 60 * 1000);
      this.QRCodeStatusHandle = setInterval(function () {
        this.isQRCodeScanned();
      }.bind(this), 5000);
    }
    // this.getInterestedItems()
    // this.getOrders()
  },
  beforeDestroy() {
    clearInterval(this.QRCodeHandle)
    clearInterval(this.QRCodeStatusHandle)
  }
}
