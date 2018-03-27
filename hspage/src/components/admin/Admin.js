let ImagesConstant = require.context('@/_images/constant/', false, /\.(png|jpg)$/);

export default {
  name: 'Admin',
  data: function () {
    return {
      userprofile: {
        additional_info: '',
        address: '',
        admin: {
          admin_level:        1,
          created_by:'system_dev_admin',
          wechat_kefu:'',
        },
        agent: {
          agent_discount: 0,
          agent_level: 0,
        },
        created_by: '',
        cellphone: '15864540221',
        email: '',
        icon: 'beyourdiamond.jpg',
        id: 'dae0b1b5-1f35-4f1a-a10a-01af466bd123',
        invitation_code: '-L8aeUxPZxsaLgLnYsMW',
        point: 0,
        real_name: '',
        recommended_by: '',
        status: 'ACTIVE',
        total_purchase_amount: 0,
        user_discount: 0.98,
        user_level: 1,
        user_type: 'ADMIN',
        username: 'hsadmin',
        wechat_id: '',
        wechat_name: '',
        wechat_qr: ''
      },
      num_newmessage:''
    }
  },
  computed: {
    imgURL: function () {
      return ImagesConstant('./logo.png')
    }
  }, 
  mounted() {
    if (this.$cookies.isKey('userprofile')) {
      this.userprofile = this.$cookies.get('userprofile')
    }
  }
}

var crr_dia_to_check_num = 0;
var crr_message_num = 0;
var t = -20;
// $(document).ready(function () {
//   audioElement = document.createElement('audio');
//   audioElement.setAttribute('src', '/_assets/sounds-766-graceful.mp3');
//   audioElement.setAttribute('autoplay', 'autoplay');
//   $.get();
//   audioElement.addEventListener("load", function () {
//     //audioElement.play();
//   }, true);

//   crr_dia_to_check_num = $('#newdiacheck_num').html();
//   crr_message_num = $('#newmessage_num').html();
//   t = setTimeout(retrievMessageNum, 20000);

//   $('.accountsnavi, .homepagenavi, .contentnavi').mouseenter(function () {
//     $(this).children('ul.subnavi-containner').stop(true).css({ 'display': 'block', 'opacity': '1' });
//   });
//   $('.accountsnavi, .homepagenavi, .contentnavi').mouseleave(function () {
//     $(this).children('ul.subnavi-containner').stop(true).fadeOut();
//   });

// });

function retrievMessageNum() {
  //console.log('retrieve now');
  $.get("newmessagesnum.php", function (data) {
    $('#notificationpad').html(data);
  });

  var new_dia_to_check_num = $('#newdiacheck_num').html();
  var new_message_num = $('#newmessage_num').html();

  if (new_dia_to_check_num > crr_dia_to_check_num || new_message_num > crr_message_num) {
    audioElement.play();
  }

  crr_dia_to_check_num = new_dia_to_check_num;
  crr_message_num = new_message_num;

  if (t > 0) {
    clearTimeout(t);
    t = setTimeout(retrievMessageNum, 20000);
  }
}
