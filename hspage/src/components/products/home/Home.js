let Images = require.context('@/_images/articles/', false, /\.(png|jpg)$/);
let ImagesConstant = require.context('@/_images/constant/', false, /\.(png|jpg)$/);
export default {
  name: 'Home',
  computed: {
    backgroudImgUrl: function () {
      return Images('./beyoudiamond-20160609_144555.jpg');
    },
    backgroudImgUrl2: function () {
      return Images('./beyoudiamond-20160122_130349.jpg');
    },
    mapImgUrl: function () {
      return ImagesConstant('./antwerp-diamond-district.jpg')
    }
  },
  beforeCreate() {
    this.$emit('getCurrentPage', 'index')
  },
  beforeDestroy() {
    this.$emit('getCurrentPage', '')
  }
}
