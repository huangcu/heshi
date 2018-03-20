let _images = require.context('@/_images/', false, /\.png$/)
export default {
  name: 'Header',
  props: {
    currentPage: String,
    jobType: String,
    platform: String,
    deviceType: String
  },
  data: function () {
    return {
      isLocalStorageSaved: false,
      accountID: '',
      interestedItems: [],
      orders: []
    }
  },
  computed: {
    isIndexPage: function () {
      return this.currentPage === 'index'
    },
    isDiamond: function () {
      return this.currentPage === 'diamonds'
    },
    isJewelry: function () {
      return this.currentPage === 'jewelry'
    },
    isJewelryNeedMounting: function () {
      return this.currentPage === 'jewelry' && this.$route.query.class === "mounting"
    },
    isGem: function () {
      return this.currentPage === 'gem'
    },
    isCustomizedJewelry: function () {
      return this.currentPage === 'customizedJewelry'
    },
    isKnowledge: function () {
      return this.currentPage === 'knowledge'
    },
    isAppointment: function () {
      return this.currentPage === 'appointment'
    },
    isContact: function () {
      return this.currentPage === 'contact'
    },
    isBrandstory: function () {
      return this.currentPage === 'brandstory'
    }
  },
  methods: {
    getInterestedItems: function () {
      if (this.accountID === undefined) {
        return
      }
      if (this.accountID === '') {
        if (this.$cookies.isKey('interestedItems')) {
          this.interestedItems = this.$cookies.get('interestedItems')
        }
      } else {
        // TODO post
        this.$http.post(this.$userURL + '/interestedItems/' + this.accountID).then(response => {
          return response.bodyText
        }, err => { console.log(err); alert('error:' + err.bodyText) })
      }
    },
    // TODO post
    getOrders: function () {
      if (this.accountID === undefined) {
        return
      }
      if (this.accountID !== '') {
        this.$http.post(this.$userURL + '/interestedItems/' + this.accountID).then(response => {
          return response.bodyText
        }, err => { console.log(err); alert('error:' + err.bodyText) })
      }
    },
    gotoMyAccount() {
      this.$router.replace('/myaccount');
    },
    gotoMyPresent() {
      this.$router.replace('/');
    },
    gotoShoppingList() {
      this.$router.replace('/');
    },
    gotoShoppingListConfirmed() {
      this.$router.replace('/');
    },
    gotoOrders() {
      this.$router.replace('/');
    },
    _images: function (name) {
      return _images('./' + name + ".png")
    }
  },
  mounted() {
    if (this.$cookies.isKey('_account')) {
      this.accountID = this.$cookies.get('_account')
    } else {
      this.accountID = ''
    }
    // this.getInterestedItems()
    // this.getOrders()
  }
}
