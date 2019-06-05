export default {
  name: 'Footer',
  data() {
    return {
      date: new Date().getFullYear()
    }
  },
  methods: {
    _images: function (name) {
      var _images = require.context('@/_images/', false, /\.(png|jpe?g|gif|svg)$/)
      return _images('./' + name)
    }
  }
}