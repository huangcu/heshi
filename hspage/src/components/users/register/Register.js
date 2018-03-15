export default {
  name: 'LoginPage',
  data: function () {
    return {
      registerErrors: [],
      email: null,
      username: null,
      password: null,
      password2: null,
      realname: null,
      wechatid: null,
      cellphone: null,
      address: null,
      additionalinfo: null,
      recommendedby: null,
      icon: null
    }
  },
  methods: {
    checkForm: function (e) {
      var formData = new FormData()
      this.registerErrors = []
      if (!this.email) {
        this.registerErrors.push('请输入邮箱.')
      } else if (!this.validEmail(this.email)) {
        this.registerErrors.push('邮箱格式不正确.')
      } else {
        formData.append('email', this.email)
      }
      if (!this.password) {
        this.registerErrors.push('请输入密码.')
      } else {
        if (!this.password2) {
          this.registerErrors.push('请输入确认密码.')
        } else if (this.password !== this.password2) {
          this.registerErrors.push('两次输入的密码不匹配')
        } else {
          formData.append('password', this.password)
        }
      }
      if (this.username) {
        if (!this.validUsername(this.username)) {
          this.registerErrors.push('请输入正确的用户名格式')
        } else {
          formData.append('username', this.username)
        }
      }
      if (this.realname) formData.append('real_name', this.realname)
      if (this.wechatid) formData.append('wechat_id', this.wechatid)
      if (this.cellphone) formData.append('cellphone', this.cellphone)
      if (this.address) formData.append('address', this.address)
      if (this.additionalinfo) formData.append('additional_info', this.additionalinfo)
      if (this.recommendedby) formData.append('recommended_by', this.recommendedby)
      if (this.icon) formData.append('icon', this.icon)
      if (!this.registerErrors.length) {
        this.$http.post(
          this.$userURL + '/users',
          formData,
          {
            headers: {
              'X-Auth-Token': 'Jbm6XfXQj/KqmMTqz6c4GQWl9U6JMLQ/T4LzPWIEi2W2Q23GDkuIfxvbUC/rar8ZJIWWSVo68fZ/hv6n0oAeXaQKEfhKmGUZ8m8JHm5TteBZwqZuqXAbOeowTJVBn8aaUhfSfZbmgNnXwDEnhjZ1DZ8jG2Khy9uzoHu5ogwbVHQ=',
              'Content-Type': 'multipart/form-data'
            }
          }).then(response => {
            this.$cookies.set('_account', response.body)
            this.$router.replace('/login');
          }, err => {
            // error callback
            console.log(err.body)
          })
        return true
      } else {
        e.preventDefault()
      }
    },
    validEmail: function (email) {
      // eslint-disable-next-line
      var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
      return re.test(email)
    },
    validUsername: function (username) {
      return true
    }
  }
}
