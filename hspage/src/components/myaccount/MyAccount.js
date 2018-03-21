let Images = require.context('@/_images/constant/', false, /\.(png|jpg)$/);

export default {
  name: 'MyAccount',
  props: {
    account: String,
  },
  data: function () {
    return {
      agent: '',
      accountLevel: 0,
      emailnotbeyoudiamond: '',
      wechat_open_id: '',
      wechat_open_idwechatnameicon: '',
      updatefeedback: '',
      qrCodeSrc: '',
      sceneID: '',
      QRCodeHandle: null,
      QRCodeStatusHandle: null
    }
  },
  computed: {
    ourWXCode: function () {
      return Images('./weixin2.png');
    },
    ourQRCode: function () {
      return Images('./qr_code.jpg');
    }
  },
  methods: {
    getQRCode: function () {
      //MAX 1 - 4294967295
      var sceneID = Math.floor(Math.random() * 1000 * 1000 * 1000)
      this.$http.get(
        this.$wechatURL + '/temp_qrcode?sceneID=' + sceneID,
      ).then(response => {
        this.sceneID = sceneID
        window.sessionStorage.setItem("HSSESSIONID", sceneID)
        this.qrCodeSrc = response.body
      }, err => { console.log(err); alert('error:' + err.bodyText) })
    },
    isQRCodeScanned: function () {
      if (this.sceneID !=='') {
        this.$http.get(
          this.$wechatURL + '/status?sceneID=' + this.sceneID
        ).then(response => {
          if (response.body !== '') {
            this.cookies.set("openID", response.body, 60 * 30)
            this.$route.replace('/')
          }
        }, err => { console.log(err); alert('error:' + err.body) })
      }
    },
    logout: function () {
      this.$http.post(
        this.$customerURL + '/logout',
        '',
        {
          headers: {
            'Authorization': 'Bearer ' + this.$cookies.get('token')
          }
        }
      ).then(response => {
        console.log('logout')
        this.$cookies.remove('token')
        this.$cookies.remove('_account')
        this.$cookies.remove('userprofile')
        this.$cookies.remove('SESSIONID')
        this.$emit('updateAccount','')
        this.$router.replace('/')
      }, err => { console.log(err); alert('error:' + err.bodyText) })
    }
  },
  created() {
    if (!this.$cookies.isKey('_account')) {
      // this.$router.replace('/login')
    }
  },
  mounted() {
    if (!this.$cookies.isKey('_account')) {
      // this.$router.replace('/login')
    }
    this.getQRCode()
    this.QRCodeHandle = setInterval(function () {
      this.getQRCode()
    }.bind(this), 2 * 60 * 1000)
    this.QRCodeStatusHandle = setInterval(function () {
      this.isQRCodeScanned()
    }.bind(this), 5000)
    if ($('#qrcode-box').length > 0) {
      // qrCounter = setTimeout(checkQRresult, 2000);
    }
    // parseTheHash();
  },
  beforeDestroy() {
    clearInterval(this.QRCodeHandle)
    clearInterval(this.QRCodeStatusHandle)
  }
}

function parseTheHash() {
  var theHashSTR = location.hash.replace('#section-', '');
  // goToSection(theHashSTR);
}

function x() {
  $('div#qrcodebg-watch').fadeOut('fast', function () {
    $('div#qrcodebg').remove();
  });
}

function openRecoContent() {
  $('div#reco-contentbox').slideDown();
  $('div.ticket-answer').slideUp();
}
function showRegulation() {
  $('div.ticket-answer').slideDown();
  $('div#reco-contentbox').slideUp();
}

function goToSection(sectionID) {
  $('div#generalinfo, div#clientorders, div#mypoints, div#myclients, div.history-recommenedusers.extendedclients, div#coupon, div.generalinfo.heshibi-box, div.accountinfobox, div#personalinfo, div#mywebsite').not('#' + sectionID).css('display', 'none');
  $('#' + sectionID).fadeIn('fast');
  $('button.inpage-navi').removeClass('active');
  $('button#inpage-navi_' + sectionID).addClass('active');
  if (sectionID == 'myclients') {
    $('div.history-recommenedusers.extendedclients').fadeIn('fast');
  }
  if (sectionID == 'coupon') {
    $('div.generalinfo.heshibi-box').fadeIn('fast');
  }
  location.hash = 'section-' + sectionID;
}
