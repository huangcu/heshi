var crr_user_id = "<?php echo $usertoviewid; ?>";


function deleteOrder(orderid) {
  var r = confirm('确定删除该纪录?');
  if (r) {
    $('#indication').fadeIn();
    $.post(
      "deletehistory.php",
      { recordid: orderid },
      function (data) {
        //if the returned data indicates ok, then do the actions

        if (data == "OK") {

          $('li#history_' + orderid).fadeOut('slow');

        } else {
          alert(data);

        }
        $('#indication').fadeOut('fast');
        /**/
      }
    );
  }
}

function deleteOrderNA(orderid) {
  var r = confirm('确定删除该纪录?');
  if (r) {
    $('#indication').fadeIn();
    $.post(
      "deletehistory-na.php",
      { recordid: orderid },
      function (data) {
        //if the returned data indicates ok, then do the actions

        if (data == "OK") {

          $('li#history_na_' + orderid).fadeOut('slow');

        } else {
          alert(data);

        }
        $('#indication').fadeOut('fast');
        /**/
      }
    );
  }
}
function specialnoteSave(orderID) {

  var newNote = $.trim($('textarea#specialnote-' + orderID).val());
  var oldNote = $.trim($('span#oldspecialnote-' + orderID).html());

  console.log(newNote);
  if (newNote != oldNote) {

    $('#indication').fadeIn();
    $('button#specialnoteBTN' + orderID).attr('disabled', 'disabled');
    $.post(
      "newspecialnotice.php",
      { recordid: orderID, note: newNote },
      function (data) {
        if (data == "OK") {
          $('span#oldspecialnote-' + orderID).html(newNote);
          $('button#specialnoteBTN' + orderID).html('<span class="glyphicon glyphicon-ok"></span> 已存');
          $('button#specialnoteBTN' + orderID).addClass('specialnotesaved');
        } else {
          alert(data);
        }
        $('#indication').fadeOut('fast');
        //$('button#specialnoteBTN'+orderID).removeAttr('disabled');
        /**/
      }
    );
  } else {
    alert('no change' + newNote);
  }
}




function updateSoldDiaProfittingStatus(diaID) {
  if ($('input#profitablecheckbox-' + diaID).prop('checked')) {
    var crr_option = 'YES';
  } else {
    var crr_option = 'NO';
  }

  $('#indication').fadeIn();
  //$('button#specialnoteBTN'+orderID).attr('disabled','disabled');
  $.post(
    "updatesolddiaprofittingstatus.php",
    { diaref: diaID, value: crr_option },
    function (data) {
      if (data == "OK") {
        $('span#indi-update-profitting-status-' + diaID).fadeIn('fast').delay(850).fadeOut(function () {
          $('span.indi-update-profitting-status').removeAttr('style');
        });
        if (crr_option == 'YES') {
          $('input#profitablecheckbox-' + diaID).parent('span.profitable-box').removeClass('no');
        } else {
          $('input#profitablecheckbox-' + diaID).parent('span.profitable-box').addClass('no');
        }
      } else {
        alert(data);
      }
      $('#indication').fadeOut('fast');
      //$('button#specialnoteBTN'+orderID).removeAttr('disabled');
      /**/
    }
  );

}


function updateOrderActualDiscount(orderID) {
  $('#indication').fadeIn();
  //$('button#specialnoteBTN'+orderID).attr('disabled','disabled');
  var crr_option = $('select#actualdiscountchooser-' + orderID).val();
  $.post(
    "updateorderactualdiscount.php",
    { recordid: orderID, value: crr_option },
    function (data) {
      if (data == "OK") {
        $('span#totaldiscountbox-' + orderID + ' span.glyphicon.glyphicon-warning-sign').remove();
        $('span#totaldiscountbox-' + orderID).addClass('statusupdated');
        $('span#totaldiscountbox-' + orderID).prepend('<span class="glyphicon glyphicon-ok"></span> ');
        $('select#actualdiscountchooser-' + orderID).attr('disabled', 'disabled');
      } else {
        alert(data);
      }
      $('#indication').fadeOut('fast');
      //$('button#specialnoteBTN'+orderID).removeAttr('disabled');
      /**/
    }
  );
}
