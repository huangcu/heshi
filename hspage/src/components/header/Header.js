let _images = require.context('@/_images/', false, /\.png$/)
export default {
  name: 'Header',
  props: {
    jobType: String,
    platform: String,
    deviceType: String
  },
  data: function () {
    return {
      isIndexPage: this.$currentPage === 'Index',
      isDiamond: this.$currentPage === 'diamonds',
      isJewelry: this.$currentPage === 'jewelry',
      isJewelryNeedMounting: this.$currentPage === 'needMounting',
      isGem: this.$currentPage === 'gem',
      isCustomizedJewelry: this.$currentPage === 'customizedJewelry',
      isKnowledge: this.$currentPage === 'knowledge',
      isAppointment: this.$currentPage === 'appointment',
      isContact: this.$currentPage === 'contact',
      isBrandstory: this.$currentPage === 'brandstory',
      isLocalStorageSaved: false,
      accountID: '',
      interestedItems: [],
      orders: []
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
