export default {
  name: 'Diamond',
  data: function () {
    return {
      diamond: '',
      userprofile:''
    }
  },
  props: {
    diamondID: String,
  },
  methods: {
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
          console.log(this.diamond)
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
