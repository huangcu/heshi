export default {
  name: 'Jewelrys',
  data: function () {
    return {
      jewelrys: []
    }
  },
  beforeCreate() {
    this.$emit('getCurrentPage', 'jewelrys')
  },
  beforeDestroy() {
    this.$emit('getCurrentPage', '')
  }
}
