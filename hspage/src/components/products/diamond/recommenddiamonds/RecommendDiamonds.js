let constantImage = require.context('@/_images/constant/', false, /\.(png|jpg)$/)
let diamondImage = require.context('@/_images/diamond/', false, /\.(png|jpg)$/)
let diamondImageThumbs = require.context('@/_images/diamond/thumbs/', false, /\.(png|jpg)$/)
export default {
  name: 'RecommendDiamonds',
  data: function () {
    return {
      diamonds: [],
      sortBy: 'carat',
      orderBy: 'up'
    }    
  },
  computed: {
    
  },
  methods: {
    getConstantImage: function (theName) {
      return constantImage('./' + theName)
    },
    getRecommendDiamonds: function () {
      this.$http.get(
        this.$userURL + '/products/diamonds'
      ).then(response => {
        if (response.status === 200) {
          this.diamonds = response.body
        }
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    sortOrderDiamonds: function () {
      var order_b = this.orderBy
      var sort_b = this.sortBy
      function compare(a, b) {
        if (order_b === 'up') {
          if (sort_b === 'carat') {
            if (a.carat < b.carat)
              return -1;
            if (a.carat > b.carat)
              return 1;
            return 0;
          }
          if (sort_b === 'price') {
            if (a.price_retail < b.price_retail)
              return -1;
            if (a.price_retail > b.price_retail)
              return 1;
            return 0;
          }
        }
        if (order_b === 'down') {
          if (sort_b === 'carat') {
            if (a.carat < b.carat)
              return 1;
            if (a.carat > b.carat)
              return -1;
            return 0;
          }
          if (sort_b === 'price') {
            if (a.price_retail < b.price_retail)
              return 1;
            if (a.price_retail > b.price_retail)
              return -1;
            return 0;
          }
        }
      }
      var d = this.diamonds 
      if (this.diamonds.length !== 0) {
        this.diamonds = d.sort(compare)
        this.$forceUpdate()
      }
    }
  },
  created() {
    this.getRecommendDiamonds()
  }
}
