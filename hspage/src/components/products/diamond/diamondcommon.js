var diamondMixin = {
  methods: {
    diamondShape (shape) {
      switch (shape) {
        case 'BR':
          return '圆形'
        case 'CU':
          return '枕形'
        case 'EM':
          return '祖母绿'
        case 'AS':
          return 'Asscher'
        case 'HS':
          return '心形'
        case 'MQ':
          return '橄榄形'
        case 'OV':
          return '椭圆形'
        case 'PR':
          return '公主方'
        case 'PS':
          return '梨形'
        case 'RAD':
          return '雷蒂恩'
        default:
          return shape
      }
    },
    diamondPosition (shape) {
      switch (shape) {
        case 'BR':
          return 0
        case 'CU':
          return -198
        case 'EM':
          return -176
        case 'AS':
          return -22
        case 'HS':
          return -154
        case 'MQ':
          return -132
        case 'OV':
          return -110
        case 'PR':
          return -66
        case 'PS':
          return -87
        case 'RAD':
          return -44
        default:
          return 0
      }
    },
    diamondImage (shape) {
      switch (shape) {
        case 'BR':
          return 'round.png'
        case 'CU':
          return 'cushion.png'
        case 'EM':
          return 'emerald.png'
        case 'AS':
          return 'asscher.png'
        case 'HS':
          return 'heart.png'
        case 'MQ':
          return 'marquise.png'
        case 'OV':
          return 'oval.png'
        case 'PR':
          return 'square.png'
        case 'PS':
          return 'pear.png'
        case 'RAD':
          return 'radiant.png'
        default:
          return 'round.png'
      }
    },
    diamondColor: function (color) {
      // TODO color map to backend db
      switch (color.toUpperCase()) {
        case 'FY':
          return '<span class="fancycolortxt">黄色</span>'
        case 'FANCY YELLOW':
          return '<span class="fancycolortxt">黄色</span>'
        case 'FLY':
          return '<span class="fancycolortxt">浅黄色</span>'
        case 'FANCY BROWNISH YELLOW':
          return '<span class="fancycolortxt">棕黄色</span>'
        case 'FBY':
          return '<span class="fancycolortxt">棕黄色</span>'
        case 'FANCY LIGHT BROWNISH YELLOW':
          return '<span class="fancycolortxt">浅棕黄</span>'
        case 'FLBY':
          return '<span class="fancycolortxt">浅棕黄</span>'
        case 'FANCY INTENSE YELLOW':
          return '<span class="fancycolortxt">浓彩黄</span>'
        case 'FIY':
          return '<span class="fancycolortxt">浓彩黄</span>'
        case 'FANCY VIVID YELLOW':
          return '<span class="fancycolortxt">艳黄色</span>'
        case 'FVY':
          return '<span class="fancycolortxt">艳黄色</span>'
        case 'FLBGY':
          return '<span class="fancycolortxt">浅棕灰</span>'
        default:
          return color.toUpperCase()
      }
    },
    diamondCutGrade: function (cutGrade) {
      if (cutGrade == null || cutGrade === '' || cutGrade === undefined) {
        return ''
      } else {
        return cutGrade
      }
    },
    diamondPlace: function (country) {
      var place = ''
      if (country === 'SZ') {
        place = '中国深圳'
      } else if (country === 'HK' || country === 'HSTHK') {
        place = '香港'
      } else if (country === 'IND') {
        place = '印度'
      } else if (country === 'Belgi' || country === 'Belgium' || country === '-' || country === undefined) {
        place = '安特卫普'
      } else if (country === 'China' || country === 'CN') {
        place = '中国'
      }
      return place
    }
  }
}
export default diamondMixin
