function saveCurrency() {
  var eur_usd = $('#eur_usd').val();
  var eur_cny = $('#eur_cny').val();
  var usd_cny = $('#usd_cny').val();
  if (isNaN(eur_usd)) {
    alert('欧元－美元 汇率填写有误');
    return;
  }
  if (isNaN(eur_cny)) {
    alert('欧元－人民币 汇率填写有误');
    return;
  }
  if (isNaN(usd_cny)) {
    alert('美元－人民币 汇率填写有误');
    return;
  }
  $('#indication').fadeIn();
  $.post(
    "updatecurrency.php",
    { eur_usd: eur_usd, eur_cny: eur_cny, usd_cny: usd_cny },
    function (data) {
      //if the returned data indicates ok, then do the actions
      if (data == "OK") {
        $('#currencysavebtn').html('<span class="glyphicon glyphicon-ok"></span> 已更新').attr('disabled', 'disabled');
        $('#indication').fadeOut('fast');
      } else {
        alert('未知错误，请重试');
      }
    }
  );
}

function approvePoint(pointID) {
  $('#indication').fadeIn();

  $.post(
    "approvepoint.php",
    { id: pointID },
    function (data) {
      //if the returned data indicates ok, then do the actions
      console.log(data);
      if (data == "OK") {
        $('#approve_' + pointID).html('<span class="glyphicon glyphicon-ok"></span> 已核准').attr('disabled', 'disabled');
        $('#indication').fadeOut('fast');
      } else {
        alert('未知错误，请重试');
      }
    }
  );
}

function cashedPoint(pointID) {
  $('#indication').fadeIn();

  $.post(
    "avoidpoint.php",
    { id: pointID },
    function (data) {
      //if the returned data indicates ok, then do the actions
      console.log(data);
      if (data == "OK") {
        $('#cashisback_' + pointID).html('<span class="glyphicon glyphicon-ok"></span> 已标注').attr('disabled', 'disabled');
        $('#indication').fadeOut('fast');
      } else {
        alert('未知错误，请重试');
      }
    }
  );
}
