export default {
  name: 'Login',
  data: function () {
    return {
      qrCodeSrc: ''
    }
  },
  props: {
    referer: String,
    for: String,
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
        {
          headers: {
            'X-Auth-Token': 'Jbm6XfXQj/KqmMTqz6c4GQWl9U6JMLQ/T4LzPWIEi2W2Q23GDkuIfxvbUC/rar8ZJIWWSVo68fZ/hv6n0oAeXaQKEfhKmGUZ8m8JHm5TteBZwqZuqXAbOeowTJVBn8aaUhfSfZbmgNnXwDEnhjZ1DZ8jG2Khy9uzoHu5ogwbVHQ=',
          }
        }
      ).then(response => {
        this.$cookies.set('sceneID', sceneID, 60 * 2)
        window.sessionStorage.setItem("HSSESSIONID", sceneID)
        this.qrCodeSrc = response.body
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    isQRCodeScanned: function () {
      if (this.$cookies.isKey('sceneID')) {
        this.$http.get(
          this.$wechatURL + '/status?sceneID=' + this.$cookies.get('sceneID'),
          {
            headers: {
              'X-Auth-Token': 'Jbm6XfXQj/KqmMTqz6c4GQWl9U6JMLQ/T4LzPWIEi2W2Q23GDkuIfxvbUC/rar8ZJIWWSVo68fZ/hv6n0oAeXaQKEfhKmGUZ8m8JHm5TteBZwqZuqXAbOeowTJVBn8aaUhfSfZbmgNnXwDEnhjZ1DZ8jG2Khy9uzoHu5ogwbVHQ=',
            }
          }
        ).then(response => {
          console.log(response.body)
          if (response.body !== '') {
            this.cookies.set("openID", response.body, 60 * 30)
            this.$route.replace('/')
          }
        }, err => { console.log(err); alert('error:' + err.body) })
      }
    }
  },
  created: function () {
    this.getQRCode()
    setInterval(function () {
      this.getQRCode();
    }.bind(this), 2 * 60 * 1000);
    setInterval(function () {
      this.isQRCodeScanned();
    }.bind(this), 5000);
  }
}
