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
      qrCodeSrc: ''
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
        this.$cookies.set('sceneID', sceneID, 60 * 2)
        window.sessionStorage.setItem("HSSESSIONID", sceneID)
        this.qrCodeSrc = response.body
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    isQRCodeScanned: function () {
      if (this.$cookies.isKey('sceneID')) {
        this.$http.get(
          this.$wechatURL + '/status?sceneID=' + this.$cookies.get('sceneID')
        ).then(response => {
          console.log(response.body)
          if (response.body !== '') {
            this.cookies.set("openID", response.body, 60 * 30)
            this.$route.replace('/')
          }
        }, err => { console.log(err); alert('error:' + err.body) })
      }
    }
  },
  created() {
    if (!this.$cookies.isKey('_account')) {
      // this.$router.replace('/login')
    }
    this.getQRCode()
    setInterval(function () {
      this.getQRCode();
    }.bind(this), 2 * 60 * 1000);
    setInterval(function () {
      this.isQRCodeScanned();
    }.bind(this), 5000);
  },
  mounted() {
    if ($('#qrcode-box').length > 0) {
      // qrCounter = setTimeout(checkQRresult, 2000);
    }
    // parseTheHash();

  }
}

function parseTheHash() {
  var theHashSTR = location.hash.replace('#section-', '');
  // goToSection(theHashSTR);
}

function checkQRresult() {
  $.get("/_content/ajax/weixin-link-status.php", function (data) {
    //$('#notificationpad').html(data);
    console.log(data);
    if (data == 'LINKED') {
      $('img#the_account_qrcode').attr('src', '/_images/constant/success.jpg')
      $('#qrcodebg').delay(1288).fadeOut('normal', function () {
        $('#qrcodebg').remove();
      });
      clearTimeout(qrCounter);
    } else {
      clearTimeout(qrCounter);
      qrCounter = setTimeout(checkQRresult, 1280);
    }
  });
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
