export default {
  name: 'newAdmin',
  data: function () {
    return {
      registerErrors:null,
      adminprofile:null,
      username:null,
      realName:null,
      password:null,
      password2:null,
      cellphone:null,
      adminLevel:null,
      iconFile: null
    }
  },
  computed: {

  },
  methods: {
    validateIconFile: function (fileList) {
      var vm = this
      if (fileList.length > 1) {
        alert("please select only one file")
        vm.clearIconFile()
        return
      }
      var f = fileList[0]
      if (f.size < 50000) {
        var _URL = window.URL || window.webkitURL;
        var img = new Image();
        var imgwidth = 0;
        var imgheight = 0;
        var maxwidth = 320;
        var maxheight = 320;

        img.src = _URL.createObjectURL(f);
        img.onload = function () {
          imgwidth = this.width;
          imgheight = this.height;
          if (imgwidth > maxwidth || imgheight > maxheight) {
            alert('Sorry, image width and heigh can not exceed 320 pixel')
            vm.clearIconFile()
          } else {
            vm.iconFile = f
          }
        }
      } else {
        alert('Sorry, file size can not exceed 50KB!')
        vm.clearIconFile()
        return
      }
      vm.iconFile = f
    },
    clearIconFile: function () {
      document.getElementById("image").value = ""
    },
    newAdminUser: function () {
      var formData = new FormData()
      // if (!this.email) {
      //   this.registerErrors.push('请输入邮箱.')
      // } else if (!this.validEmail(this.email)) {
      //   this.registerErrors.push('邮箱格式不正确.')
      // } else {
      //   formData.append('email', this.email)
      // }
      if (!this.cellphone) {
        this.registerErrors.push('请输入电话号码.')
      } else if (!this.validateCellphone(this.cellphone)) {
        this.registerErrors.push('电话号码格式不正确.')
      } else {
        formData.append('cellphone', this.cellphone)
      }
      if (!this.password) {
        this.registerErrors.push('请输入密码.')
      } else {
        if (!this.password2) {
          this.registerErrors.push('请输入确认密码.')
        } else if (this.password !== this.password2) {
          this.registerErrors.push('两次输入的密码不匹配')
        } else {
          formData.append('password', this.password)
        }
      }
      if (this.username) {
        if (!this.validateUsername(this.username)) {
          this.registerErrors.push('请输入正确的用户名格式')
        } else {
          formData.append('username', this.username)
        }
      }
      if (this.realName) {
        formData.append('real_name', this.realName)
      }
      if (!this.adminLevel) {
        this.registerErrors.push('请输入管理员级别.')
      } else if (!this.validateadminLevel(this.adminLevel)) {
        this.registerErrors.push('电话号码格式不正确.')
      } else {
        formData.append('level', this.adminLevel)
      }
      if (this.iconFile !== null) {
        formData.append('images', this.iconFile)
      }
      formData.append('user_type','ADMIN')
      this.$http.post(
        this.$adminURL + '/users',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        }).then(response => {
          if (response.status === 200) {
            var registerResult = JSON.parse(response.bodyText)
            console.log(registerResult)
            if (registerResult.code!==undefined) {
              this.registerErrors = registerResult.message
              alert(registerResult.message)
              return
            }
            alert('添加成功')
            // TODO refresh page
            this.$forceUpdate()
          }
        }, err => { console.log(err); alert('error:' + err.body) })
    },
    validateUsername: function (username) {
      return true
    },
    validateCellphone: function (cellphone) {
      return true
    },
    validateadminLevel: function (adminLevel) {
      return (adminLevel>0 && adminLevel<10)
    }
  },
  created() {

  },
  mounted() {
    if (this.$cookies.isKey('adminprofile')) {
      this.adminprofile = JSON.parse(this.$cookies.get('adminprofile'))
      if (this.adminprofile.admin.admin_level!==1){
        alert('NOT ALLOWED!')
        this.$router.replace('/manage/admins')
      }
    } else {
      this.$router.replace('/manage/login')
    }
  }
}
$(document).ready(function () {
  $('input.image_selecting').change(function (e) {
    e.preventDefault();
    //check first if the file extension is correct
    var crr_form = $(this).parent('form');
    var filefullpath = $(this).val();
    var dotposition = filefullpath.lastIndexOf('.');
    var extension = filefullpath.substring(dotposition);
    if (extension == '.jpg' || extension == '.jpe' || extension == '.jpeg' || extension == '.JPG' || extension == '.JPEG' || extension == '.JPE' || extension == '.gif' || extension == '.GIF' || extension == '.png' || extension == '.PNG') {
      $('div.progress').fadeIn('fast');
      crr_form.submit();
    } else {
      alert("图片格式不支持，请上传. JPEG,. GIF 或. PNG 图片.");
    }
  });

  (function () {
    var bar = $('.progress .bar');
    var percent = $('.progress .percent');
    var status = $('#status');
    $('form.uploadImage').ajaxForm({
      beforeSend: function () {
        status.empty();
        var percentVal = '0%';
        bar.width(percentVal)
        percent.html(percentVal);
      },
      uploadProgress: function (event, position, total, percentComplete) {

        var percentVal = percentComplete + '%';
        bar.width(percentVal)
        percent.html(percentVal);
      },
      complete: function (xhr) {
        $('div.progress').fadeOut('fast');

        $('button').removeAttr('disabled');
        status.html(xhr.responseText);
        console.log(xhr.responseText);

        var feedback = status.find("p.message").html();

        if (feedback == 'OK') {
          var imagename = status.find("p#imagename").html();
          //var imagewhere=status.find("p#imagewhere").html();
          //$("#img"+imagewhere).attr('src',('/_images/jewelry/thumbs/'+imagename));
          $('div.picture-preview-box').append('<img class="uploadedpreviewpic" src="/_images/admins/thumbs/' + imagename + '" />');
          $("#icon").val(imagename);
        } else {
          console.log(feedback);
          alert("Er is een onbekende fout opgetreden. Probeer het opnieuw.");
        }
      }
    });
  })();
});

function addUploadedPics() {
  var imgcounting = 0;
  $('#status p.uploadedimage').each(function () {
    var crre = $(this);
    var crr_pic_num = crre.attr('title');
    var crr_pic_name = crre.html();
    imgcounting++;

    $('input#value_img_' + crr_pic_num).val(crr_pic_name);
    $('div.picture-preview-box').append('<img class="uploadedpreviewpic" src="/_images/jewelry/thumbs/' + crr_pic_name + '" />');
  });
  $('input#startwith').val(imgcounting);
}

function submittheform() {
  $('form#addnewadmin').submit();
}

// require_once('admin_authorizing_part.php');
// 		$message_db="添加成功";
// 		//1添加客服
// 		$kfURL="https://api.weixin.qq.com/customservice/kfaccount/add?access_token=".$theaccesstoken;
// 		//2 添加头像				
// 			$iconURL="https://api.weixin.qq.com/customservice/kfaccount/uploadheadimg?access_token=".$theaccesstoken.'&kf_account='.$username.'@admin';


// <script type="text/javascript" src="tinymce/tinymce.min.js"></script>
// <script type="text/javascript" src="formplugin.js"></script>
// <script type="text/javascript">
// tinymce.init({
//     selector: "textarea#content",
// 	theme: "modern",
//     width: 500,
//     height: 260,
//     plugins: [
//          "advlist autolink link image lists charmap print preview hr anchor pagebreak spellchecker",
//          "searchreplace wordcount visualblocks visualchars code fullscreen insertdatetime media nonbreaking",
//          "save table contextmenu directionality emoticons template paste textcolor"
//     ],
//     toolbar: "insertfile undo redo | fontsizeselect | styleselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | link image | print preview media fullpage | forecolor backcolor emoticons"
//  });
// </script>
