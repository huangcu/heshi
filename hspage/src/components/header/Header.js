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
      appToken: '',
      accountId: '',
      interestedItems: [],
      orders: []
    }
  },
  methods: {
    getAccount() {
      if (this.$cookies.isKey('beyou')) {
        return this.$cookies('beyou')
      } else {
        return ''
      }
    },
    getInterestedItems: function () {
      if (this.accountId === '') {
        if (this.$cookies.isKey('interestedItems')) {
          return this.$cookies('interestedItems')
        }
      } else {
        // TODO post
        this.$http.post(this.$userURL + 'interestedItems/' + this.accountId).then(response => {
          return response.bodyText
        }, err => { console.log(err); alert('error:' + err.bodyText) })
      }
    },
    // TODO post
    getOrders: function () {
      if (this.accountId !== '') {
        this.$http.post(this.$userURL + 'interestedItems/' + this.accountId).then(response => {
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
      var _images = require.context('@/_images/', false, /\.png$/)
      return _images('./' + name + ".png")
    }
  },
  created() {
    this.getInterestedItems()
    this.getAccount()
    this.getOrders()
  }
}
