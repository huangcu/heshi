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
    }
  }
}
export default diamondMixin
