let constantImage = require.context('@/_images/constant/', false, /\.(png|jpg)$/)

export default {
  name: 'DiamondOfTheWeek',
  data: function() {
    return {
      userprofile: {
        id: ''
      },
      jewelrys: [],
      diamonds: [],
      promo_type: 'FREE_ACC',
      promo_price: 0
    }
  },
  computed: {
    bigSaleImage: function () {
      return constantImage('./bigsale.png')
    },
    soldOutImage: function () {
      return constantImage('./soldout.png')
    }
  },
  methods: {
    getConstantImage: function (theName) {
      return constantImage('./' + theName)
    },
    getDiamondsOfTheWeek: function () {
      this.$http.get(
        this.$userURL + '/products/diamonds'
      ).then(response => {
        if (response.status === 200) {
          this.diamonds = response.body
          if (this.diamonds.length > 0) {
            // if (this.diamonds[0].promo_type == "FREE_ACC") {
            this.getJewelryForDiamond()
            // }
          }
        }
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    getJewelryForDiamond: function () {
      if (this.diamonds.length !== 0) {
        // SELECT * FROM jewelry WHERE free_acc = "YES" AND 
        // dia_shape LIKE "%'.$shape.'%" AND 
        // dia_size_min <= '.$carat.' AND 
        // dia_size_max >= '.$carat.' 
        // AND need_diamond = "YES" AND online = "YES" AND stock_quantity > 1
        var formData = new FormData()
        formData.append('free_acc', 'YES')
        formData.append('dia_shape', this.diamonds[0].shape)
        formData.append('size', this.diamonds[0].carat)
        formData.append('online', 'YES')
        formData.append('stock_quantity', 1)
        this.$http.post(
          this.$userURL + '/products/filter/jewelrys?class=mounting',
          formData
        ).then(response => {
          if (response.status === 200) {
            this.jewelrys = response.body
          }
        }, err => { console.log(err); alert('error:' + err.body) })
      }
    },
    chooseThisJew: function (jewid) {
      $('div.jewelrybox').removeClass('chosen')
      $('div#jewelrybox_' + jewid).addClass('chosen')
      if ($('input[name="jewID"]').length > 0) {
        $('input[name="jewID"]').val(jewid)
      }
    },
    submitOrder: function () {
      var r = true
      var jewsnum = $('div.jewelrybox').length
      if (jewsnum > 0 && $('input[name="jewID"]').val() == '') {
        r = confirm('您没有选空托，确定只够买钻石吗？')
      }
      if (r) {
        $('form#specialofferorderform').submit()
      }
    },
    diamondFluo: function (theFluo) {
      switch (theFluo) {
        case 'NON', 'NONE':
          return '无荧光'
        case 'FNT':
          return '弱荧光'
        case 'MED':
          return '中荧光'
        case 'STG','VST':
          return '强荧光'
        case 'SLT':
          return '弱荧光'
        case 'VSL':
          return '弱荧光'
        default:
          break;
      }
    },
    diamondShapeTxtPic: function(theShape, thetype) {
      // "BR","PS","PR","HS","MQ","OV","EM","CU","AS","RAD","RBC","RCRB","RC","PE","HT","CMB"
      if (thetype==="TXT") {
        switch (theShape) {
          case 'BR':
            return '圆形'
          case 'PS':
            return '公主方'
          case 'PR':
            return '梨形'
          case 'HS':
            return '心形'
          case 'MQ':
            return '橄榄形'
          case 'OV':
            return '椭圆形'
          case 'EM':
            return '祖母绿'
          case 'RA':
            return '雷蒂恩'
          case 'CU':
            return '枕形'
          case 'AS':
            return 'Asscher'
          default:
            return '圆形'
        }
      } else {
        switch (theShape) {
          case 'BR':
            return constantImage('./round.png')
          case 'PS':
          // TODO no pic
            return constantImage('./pear.png')
          case 'PR':
            return constantImage('./pear.png')
          case 'HS':
            return constantImage('./heart.png')
          case 'MQ':
            return constantImage('./marquise.png')
          case 'OV':
            return constantImage('./oval.png')
          case 'EM':
            return constantImage('./emerald.png')
          case 'RA':
            return constantImage('./radiant.png')
          case 'CU':
            return constantImage('./cushion.png')
          case 'AS':
            return constantImage('./assher.png')
          default:
            return '圆形'
        }
      }
    }
  },
  created() {
    this.getDiamondsOfTheWeek()
  }
}
$(document).ready(function () {
  $("a.fancybox").fancybox({
    maxWidth: 728,
    maxHeight: 588,
    margin: 20,
    padding: 0,
    helpers: {
      thumbs: {
        width: 50,
        height: 50
      }
    }
  })
})
