export default {
  name: 'Jewelrys',
  data: function () {
    return {
      jewelrys: []
    }
  },
  beforeCreate() {
    this.$emit('getCurrentPage', 'jewelry')
  },
  beforeDestroy() {
    this.$emit('getCurrentPage', '')
  }
}
