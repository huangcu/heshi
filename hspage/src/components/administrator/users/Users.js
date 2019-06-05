let ImagesConstant = require.context('@/_images/constant/', false, /\.(png|jpg)$/)

export default {
  name: 'Admin',
  data: function () {
    return {
      adminprofile: '',
      num_newmessage: ''
    }
  },
  computed: {
  },
  methods: {
    assignUserLevel: function (userid) {
      //var crruid=userid
      var theassignedlevel = $('#assigned_account_level_' + userid).val()
      if (theassignedlevel == '') {
        alert('请指定级别！')
        return
      }
      if (theassignedlevel == 6) {
        alert('不能直接定为6级，请先设定“推荐得的返点”')
        return
      }
      $('#indi').fadeIn()
      $.post(
        'assignaccountlevel.php',
        { user: userid, level: theassignedlevel },
        function (data) {
          console.log(data)
          if (data == 'OK') {
            $('#indi').fadeOut()
            $('#userlevelchanged_indi_' + userid).fadeIn().delay(1588).fadeOut()
            //alert($('#point_'+crruid+' option.normalpoint').attr('name'))
            $('#point_' + userid + ' option.normalpoint').attr('selected', 'selected')
            //$('#assigned_account_level_'+userid+' option.normaluseroption').attr('selected','selected')
          } else if (data == '6') {
            alert('不能直接定为6级，请先设定“推荐得的返点”')
          } else {
            alert('未知错误，请重试')
          }
        }
      )
    },
    assignagent: function(userid) {
      $('#indi').fadeIn()
      $.post(
        'assignagent.php',
        { user: userid },
        function (data) {
          console.log(data)
          if (data == 'OK') {
            $('#indi').fadeOut()
            $('button#assignresellerbtn_' + userid).attr('disabled', 'disabled').html('指定成功')
            location.reload()
          } else {
            alert('未知错误，请重试')
          }
        }
      )
    },
    resignagent: function(userid) {
      $('#indi').fadeIn()
      $.post(
        'resignagent.php',
        { user: userid },
        function (data) {
          if (data == 'OK') {
            $('#indi').fadeOut()
            $('button#resignresellerbtn_' + userid).attr('disabled', 'disabled').html('解除成功')
            location.reload()
          } else {
            alert('未知错误，请重试')
          }
        }
      )
    },
    updateAgentInfo: function(action, userid) {
      var crr_value = $('#agent_' + action + '_' + userid).val()
      if (action == 'level') {
        if (crr_value != 1 && crr_value != 2 && crr_value != 3) {
          alert('代理级别只能为1 或 2 或 3！')
          return
        }
      }

      $('#indi').fadeIn()
      $.post(
        'updateagentinfo.php',
        { action: action, userid: userid, value: crr_value },
        function (data) {
          //alert(data)
          console.log(data)
          if (data == 'OK') {
            $('#indi').fadeOut()
            $('#agentlevelsetindi_' + userid).fadeIn().delay(1280).fadeOut()
          } else {
            alert('未知错误，请重试')
          }
        }
      )
    },
    agentLevelLockStatus: function(userid) {
      //alert($('#agentlevellocked_'+userid).prop('checked'))
      if ($('#agentlevellocked_' + userid).prop('checked')) {
        var crr_value = 'YES'
      } else {
        var crr_value = 'NO'
      }

      $('#indi').fadeIn()
      $.post(
        'updateagentstatus.php',
        { userid: userid, locked: crr_value },
        function (data) {
          //alert(data)
          if (data == 'OK') {
            $('#indi').fadeOut()
            if (crr_value == "YES") {
              $('span#agentlevellockedindiwords_' + userid).html('<span class="glyphicon glyphicon-lock"></span> 已锁定')
            } else {
              $('span#agentlevellockedindiwords_' + userid).html('已解除锁定')
            }

          } else {
            alert('未知错误，请重试')
          }
        }
      )
    },
   changeRecommenderPoint: function(userid) {
  //alert(userid)
    var crr_value = $('#point_' + userid).val()
    $('#indi').fadeIn()
    $.post(
      'updatepointforrecommender.php',
      { userid: userid, point: crr_value },
      function (data) {
        console.log(data)
        if (data == 'OK') {
          $('#indi').fadeOut()
          $('#pointchanged_indi_' + userid).fadeIn().delay(1500).fadeOut()
          if (crr_value > 1) {
            $('#assigned_account_level_' + userid + ' option.recommenderoption').attr('selected', 'selected')
          } else {
            $('#assigned_account_level_' + userid + ' option.normaluseroption').attr('selected', 'selected')
          }

        } else {
          alert('未知错误，请重试')
        }
      })
    },
    thisUserCompare: function(userid) {
      if ($('#referenceprice_' + userid).prop('checked')) {
        var thelevel = -1
      } else {
        var thelevel = 1
      }

      $('#indi').fadeIn()
      $.post(
        'assignaccountlevel.php',
        { user: userid, level: thelevel },
        function (data) {
          console.log(data)
          if (data == 'OK') {
            $('#indi').fadeOut()
            if (thelevel == -1) {
              $('#userlevelbox_' + userid).fadeOut()
              $('#rpbox_' + userid).fadeOut()
              $('input#referencepriceratio_' + userid).removeAttr('disabled')
              $('button#referencepriceratiosavebtn_' + userid).removeAttr('disabled')
              $('p#rrbox_' + userid).removeClass('disabled')
            } else {
              $('#userlevelbox_' + userid).fadeIn()
              $('#rpbox_' + userid).fadeIn()
            }
            $('#referenceprice_indi_' + userid).fadeIn().delay(580).fadeOut()
          } else {
            alert('未知错误，请重试')
          }
        }
      )
    },
    saveReferencePriceRatio: function(userid) {
      var theratio = $('input#referencepriceratio_' + userid).val()
      $('#indi').fadeIn()
      $('button#referencepriceratiosavebtn_' + userid).attr('disabled', 'disabled')
      $.post(
        'updateuserreferenceprice.php',
        { user: userid, priceratio: theratio },
        function (data) {
          console.log(data)
          if (data == 'OK') {
            $('#indi').fadeOut()
            $('button#referencepriceratiosavebtn_' + userid).removeAttr('disabled')
            $('#referencepriceratio_indi_' + userid).fadeIn().delay(580).fadeOut()
          } else {
            alert('未知错误，请重试')
          }
        }
      )
    },
    seemore: function(id) {
      $('li.user-record').removeClass('active')
      $('li#user_' + id).addClass('active')
    },
    updateProfile: function (arg) {
      this.adminprofile = arg
    },
    imgURL: function (theName) {
      return ImagesConstant('./' + theName)
    }
  },
  mounted() {
    if (this.$cookies.isKey('adminprofile')) {
      this.adminprofile = JSON.parse(this.$cookies.get('adminprofile'))
    }
  },
  ready() {
    $('input.image_selecting').change(function (e) {
      e.preventDefault()
      //check first if the file extension is correct
      var crr_form = $(this).parent('form')
      var filefullpath = $(this).val()
      var dotposition = filefullpath.lastIndexOf('.')
      var extension = filefullpath.substring(dotposition)
      if (extension == '.jpg' || extension == '.jpe' || extension == '.jpeg' || extension == '.JPG' || extension == '.JPEG' || extension == '.JPE' || extension == '.gif' || extension == '.GIF' || extension == '.png' || extension == '.PNG') {
        $('div.progress').fadeIn('fast')
        crr_form.submit()
      } else {
        alert("图片格式不支持，请上传.JPEG,.GIF 或.PNG 图片.")
      }
    })
    (function () {
      var bar = $('.progress .bar')
      var percent = $('.progress .percent')
      var status = $('#status')
      $('form.uploadImage').ajaxForm({
        beforeSend: function () {
          status.empty()
          var percentVal = '0%'
          bar.width(percentVal)
          percent.html(percentVal)
        },
        uploadProgress: function (event, position, total, percentComplete) {
          var percentVal = percentComplete + '%'
          bar.width(percentVal)
          percent.html(percentVal)
        },
        complete: function (xhr) {
          $('div.progress').fadeOut('fast')
          $('button').removeAttr('disabled')
          status.html(xhr.responseText)

          var feedback = status.find("p.message").html()
          if (feedback == 'OK') {
            var imagename = status.find("p#imagename").html()
            var imagewhere = status.find("p#imagewhere").html()
            $("#qrimage_" + imagewhere).attr('src', ('/_images/qrcodes/' + imagename))
            $('#user_' + imagewhere + ' .agent-contact-info-box form.uploadImage label').html('二维码已存')
          } else {
            alert("Er is een onbekende fout opgetreden. Probeer het opnieuw.")
          }
        }
      })
    })()
  }
}
