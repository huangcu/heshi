export default {
  name: 'qrsign',
  data: function () {
    return {
      login_feedback: '',
      upgrade_feedback: '',
      appointment: '',
      account: 'test',
      therecommendticketcode: '',
      wechatopenID: ''
    }
  },
  methods: {
    showqrSignUp: function () {
      $('button#showqrsignup_btn').attr('disabled', 'disabled');
      // TODO
      // $('form#signup-form').submit();
    },
    showqrSignIn: function () {
      $('#showqrsignup_btn, #showqrsignin_btn').css('display', 'none');
      $('#signin-form').fadeIn();
    },
    qrSignIn: function () {
      $('#qrsignin_btn').attr('disabled', 'disabled');
      // TODO
      // $('form#signin-form').submit();
    }
  },
  created() {
    if (!this.$cookies.isKey('wechatopenID')) {
      // this.$router.replace('/login')
    }
  }
}
