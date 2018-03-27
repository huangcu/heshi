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
