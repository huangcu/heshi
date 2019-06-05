export default {
  name: "admins",
  data: function () {
    return {
      users: [],
      adminprofile: ''
    }
  },
  computed: {
    activeAdminUsers: function () {
      if (this.users === null || this.users === '' || this.users.length === 0) {
        return ''
      }
      return this.users
    }
  },
  methods: {
    deleteAdmin: function (theID) {
      this.$http.delete(
        this.$adminURL + '/users/' + theID,
      ).then(response => {
        alert(response.body)
        for (var i=0; i<this.users.length; i++) {
          if (this.users[i].id === theID) {
            this.users[i].status = 'disabled'
            break
          }
        }
        var u = this.users
        if (u !== null) {
          this.users = u.filter(filterUsers)
        }
      }, err => { console.log(err); alert('error:' + err.bodyText) })
    },
    getAdmins: function () {
      this.$http.get(
        this.$adminURL + '/users?user_type=ADMIN',
      ).then(response => {
        if (response.status === 200) {
          this.users = response.body
        }
      }, err => { console.log(err); alert('error:' + err.bodyText) })
    }
  },
  created () {
  },
  mounted() {
    if (this.$cookies.isKey('adminprofile')) {
      this.adminprofile = JSON.parse(this.$cookies.get('adminprofile'))
      this.getAdmins()
    } else {
      this.$router.replace('/manage/login')
    }
  }
}

function filterUsers(theUser) {
  return theUser.status !== 'disabled'
}
