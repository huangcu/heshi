let diamondImage = require.context('@/.image/diamond/', false, /\.(png|jpg)$/);
import diamondMixin from '../diamondcommon.js'
import agentPriceMixin from '@/util/agentprice.js'
import accountPriceMixin from '@/util/accountprice.js'
import currentCaculatorMixin from '@/util/currency_caculator.js'

export default {
  name: 'Diamond',
  mixins: [diamondMixin, agentPriceMixin, accountPriceMixin, currentCaculatorMixin],
  props: {
    rate: Object,
    diamondID: String,
  },
  data: function () {
    return {
      diamond: '',
      userprofile:''
    }
  },
  methods: {
    getDiamondImage: function (theName) {
      return diamondImage("./" + theName)
    },
    getDiamond: function (theID) {
      this.$http.get(
        this.$userURL + '/products/diamonds/' + theID
      ).then(response => {
        if (response.status === 200) {
          var diamonds = response.body
          if (diamonds !== null && diamonds.length ==1) {
            this.diamond = diamonds[0]
          } else {
            alert('error: ' + diamonds)
          }
        }
      }, err => { console.log(err); alert('error:' + err.body) })
    },
  },
  created() {
    var did = this.$route.params.id
    if (did !== null && did !== undefined && did !=='') {
      this.getDiamond(did)
    }
  }
}
