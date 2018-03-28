let ImagesConstant = require.context('@/_images/constant/', false, /\.(png|jpg)$/);

export default {
  name: 'adminLogin',
  data: function () {
    return {
      wrongmessage:'',
      adminprofile: '',
      username:'',
      password:''
    }
  },
  methods: {
    ImgURL: function(theName) {
      return ImagesConstant('./'+theName)
    },
    login: function () {
      var formData = new FormData()
      formData.append('username', this.username)
      formData.append('password', this.password)
      formData.append('user_type', 'ADMIN')

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
              alert(loginResult.message)
              return
            }
            this.$cookies.set('token', loginResult.token, 60 * 30)
            var adminprofile = JSON.parse(loginResult.userprofile)
            this.$cookies.set('adminprofile', loginResult.userprofile, 60 * 30)
            this.$emit('updateProfile', adminprofile)
            this.$router.replace('/manage/admins')
          }
        }, err => { console.log(err); alert('error:' + err.body) })
    }
  },
  created() {
    if (this.$cookies.isKey('adminprofile')) {
      this.adminprofile = JSON.parse(this.$cookies.get('adminprofile'))
      this.$route.replace('/manage/admins')
    } else {
      this.adminprofile = ''
    }
  }
}
