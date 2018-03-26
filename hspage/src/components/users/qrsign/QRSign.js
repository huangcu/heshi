export default {
  name: 'qrsign',
  data: function () {
    return {
      login_feedback: '',
      upgrade_feedback: '',
      appointment: '',
      userprofile: '',
      username: '',
      password: '',
      therecommendticketcode: '',
    }
  },
  computed: {
    wechatopenID: function () {
      var wechatopenID = this.$route.params.wechatopenID
      if (wechatopenID === null || wechatopenID === undefined) {
        return ''
      } else {
        return wechatopenID
      }
    },
    referer: function () {
      var referer = this.$route.query.get('referer')
      if (referer === null  || referer === undefined) {
        return ''
      } else {
        return referer
      }
    }
  },
  methods: {
    showqrSignUp: function () {
      $('button#showqrsignup_btn').attr('disabled', 'disabled');
      // TODO
      // $('form#signup-form').submit();
    },
    showqrSignIn: function () {
      $('#showqrsignup_btn, #showqrsignin_btn').css('display', 'none');
      $('#signin-form').fadeIn();
    },
    qrSignIn: function () {
      $('#qrsignin_btn').attr('disabled', 'disabled');
      var formData = new FormData()
      formData.append('username', this.username)
      formData.append('password', this.password)
      formData.append('wechat_id', this.wechatopenID)
      // TODO in backend server
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
            if (loginResult.code !== 200) {
              this.login_feedback = loginResult.message
              return
            }
            this.$cookies.set('token', loginResult.token, 60 * 30)
            var userprofile = JSON.parse(loginResult.userprofile)
            this.$cookies.set('userprofile', loginResult.userprofile, 60 * 30)
            this.$emit('updateAccount', userprofile.id)
            this.$router.replace('/home')
          }
        }, err => { console.log(err); alert('error:' + err.body) })
    }
  },
  created() {
    if (!this.$cookies.isKey('wechatopenID')) {
      // this.$router.replace('/login')
    }
  },
  mounted() {
    if (this.$cookies.isKey('userprofile')) {
      this.userprofile = this.$cookies.get('userprofile')
    } else {
      this.userprofile = ''
    }
  }
}
