let Images = require.context('@/_images/articles/', false, /\.(png|jpg)$/);
export default {
  name: 'Home',
  computed: {
    backgroudImgUrl: function () {
      return Images('./beyoudiamond-20160609_144555.jpg');
    },
    backgroudImgUrl2: function () {
      return Images('./beyoudiamond-20160122_130349.jpg');
    }
  },
  created() {
    this.$currentPage = 'Index'
  }
}
