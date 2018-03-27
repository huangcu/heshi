import agentPrice from '@util/agentprice.js'
export default {
  name: 'Agent',
  data: function () {
    return {
      agentLevel: 0,
      name: '',
      account: '',
      offer_ticket: ''
    }
  },
  computed: {
    fromAgentLevel: function () {
      var str =''
      for (i=0; i< this.agentLevel; i++) {
        str = str + '<span class="glyphicon glyphicon-certificate"></span>'
      }
      return str
    }
  },
  methods: {
    goToSection: function (sectionID) {
      $('div#generalinfo, div#clientorders, div#mypoints, div#myclients, div.history-recommenedusers.extendedclients, div#coupon, div.generalinfo.heshibi-box, div.accountinfobox, div#personalinfo, div#mywebsite').not('#' + sectionID).css('display', 'none')
      $('#' + sectionID).fadeIn('fast')
      $('button.inpage-navi').removeClass('active')
      $('button#inpage-navi_' + sectionID).addClass('active')
      if (sectionID == 'myclients') {
        $('div.history-recommenedusers.extendedclients').fadeIn('fast')
      }
      if (sectionID == 'coupon') {
        $('div.generalinfo.heshibi-box').fadeIn('fast')
      }
      location.hash = 'section-' + sectionID
    },
    showRegulation: function() {
      $('div.ticket-answer').slideDown()
      $('div#reco-contentbox').slideUp()
    },
    openRecoContent: function () {
      $('div#reco-contentbox').slideDown()
      $('div.ticket-answer').slideUp()
    },
    checkdomainavailability: function (){
      var crr_domainchoice = $('input#sitedomain').val() + '.' + $('select#domainsuffix').val()
      var domaincheckurl = 'https://sg.godaddy.com/zh/domains/searchresults.aspx?checkAvail=1&tmskey=&domainToCheck=' + crr_domainchoice
      window.open(domaincheckurl)
    },
    completeTransaction: function (transactionid) {
      $('button#transactioncompletebtn_' + transactionid).html('<span class="glyphicon glyphicon-refresh"></span> 执行中...稍候').attr('disabled', 'disabled');
      $.post(
        "/_includes/functions/completetransaction.php",
        { id: transactionid },
        function (data) {
          //alert(data);
          console.log('order feedback: ' + data)

          if ($.trim(data) == 'OK') {
            //alert('ordered');
            $('button#transactioncompletebtn_' + transactionid).html('<span class="glyphicon glyphicon-ok"></span> 交易完成').addClass('complete').attr('disabled', 'disabled');
            /*
            $.post(
              "/_includes/functions/balancecheck.php", 
              {id: "<?php echo $accountID; ?>"}, 
              function(data){
                //alert(data);
                console.log('order feedback: '+data)
              	
                if($.isNumeric(data)){
                  $('span#heshibi-balance-value').html(data);
                }else{
                  //alert('网络异常，请重试!');
                }
              	
              }
            );
            */
          } else {
            alert('网络异常，请重试!');
          }
        }
      );
    },
    archiveTransaction: function (transactionid) {
      $('button#transactionarchive_' + transactionid).html('<span class="glyphicon glyphicon-refresh"></span> 执行中...稍候').attr('disabled', 'disabled');
      $.post(
        "/_includes/functions/archiveransaction.php",
        { id: transactionid },
        function (data) {
          //alert(data);
          console.log('order feedback: ' + data)

          if ($.trim(data) == 'OK') {
            //alert('ordered');
            $('button#transactionarchive_' + transactionid).html('<span class="glyphicon glyphicon-ok"></span> 已存档').addClass('complete').attr('disabled', 'disabled');
            $('li#transactionpiece_' + transactionid).delay(500).fadeOut('normal');
          } else {
            alert('网络异常，请重试!');
          }

        }
      );
    },
    agent_level_words: function () {
    //   if ($agent_level == 1) {
    //     if ($total_price_one_year < 40000 && $total_piece_one_year < 12) {
    //       $agent_level_words = '您再继续购买'.(12 - $total_piece_one_year).'件合适商品，并且这些商品的价值超过'.(40000 - $total_price_one_year).'欧元，您即可升级为二级代理，且得到一年消费总额的5%作为返点。';
    //     } else if ($total_price_one_year < 40000 && $total_piece_one_year >= 12) {
    //       $agent_level_words = '您再继续消费'.(40000 - $total_price_one_year).'欧元，您即可升级为二级代理，享受85折价格，且得到一年消费总额的5%作为返点。';
    //     } else if ($total_price_one_year >= 40000 && $total_piece_one_year < 12) {
    //       $agent_level_words = '您再继续购买'.(12 - $total_piece_one_year).'件合适商品，您即可升级为二级代理，享受85折价格，且得到一年消费总额的5%作为返点。';
    //     } else if ($total_piece_one_year >= 12 && $total_price_one_year >= 40000 && ($total_piece_one_year < 20 || $total_price_one_year < 100000)) {
    //       $agent_level_words = '<span class="glyphicon glyphicon-thumbs-up"></span> 您已经可以升级为二级代理，享受85折价格，且得到一年消费总额的5%(须根据购买时的实际折扣核算)作为返点。您在几个小时内会收到确认消息。';
    //     } else if ($total_piece_one_year >= 20 && $total_price_one_year >= 100000) {
    //       $agent_level_words = '<span class="glyphicon glyphicon-thumbs-up"></span> 您已经可以升级为三级代理，享受83折价格，且得到一年消费总额的7%(须根据购买时的实际折扣核算)作为返点。您在几个小时内会收到确认消息。';
    //     }
    //   } else if ($agent_level == 2) {
    //     if ($total_price_one_year < 40000 || $total_piece_one_year < 12) {
    //       $agent_level_words = '<span class="glyphicon glyphicon-bell"></span> 您一年内的购买量不足二级代理的最低限度，您即将成为一级代理。';
    //     } else if ($total_piece_one_year >= 12 && $total_price_one_year >= 40000 && ($total_piece_one_year < 20 || $total_price_one_year < 100000)) {
    //       if ($total_piece_one_year >= 20) {
    //         $agent_level_words = '您再继续消费'.(100000 - $total_price_one_year).'欧元，您即可升级为三级代理，享受83折价格，且得到一年消费总额的2%作为返点。';
    //       } else if ($total_price_one_year >= 100000) {
    //         $agent_level_words = '您再继续购买'.(20 - $total_piece_one_year).'件合适商品，您即可升级为三级代理，享受83折价格，且得到一年消费总额的2%作为返点。';
    //       } else {
    //         $agent_level_words = '您再继续购买'.(20 - $total_piece_one_year).'件合适商品，并且金额超过'.(100000 - $total_price_one_year).'，您即可升级为三级代理，享受83折价格，且得到一年消费总额的2%作为返点。';
    //       }
    //     } else if ($total_piece_one_year >= 20 && $total_price_one_year >= 100000) {
    //       $agent_level_words = '<span class="glyphicon glyphicon-thumbs-up"></span> 您已经可以升级为三级代理，享受83折价格，且得到一年消费总额的2%(须根据购买时的实际折扣核算)作为返点。我们会尽快联系您和您确认。';
    //     }
    //   } else if ($agent_level == 3) {
    //     if ($total_price_one_year < 40000 || $total_piece_one_year < 12) {
    //       $agent_level_words = '<span class="glyphicon glyphicon-bell"></span> 您一年内的购买量不足三级代理和二级代理的最低限度，您即将成为一级代理。';
    //     } else if ($total_piece_one_year >= 12 && $total_price_one_year >= 40000 && ($total_piece_one_year < 20 || $total_price_one_year < 100000)) {
    //       $agent_level_words = '<span class="glyphicon glyphicon-bell"></span> 您一年内的购买量不足三级代理的最低限度，您即将成为二级代理。';
    //     } else if ($total_piece_one_year >= 20 && $total_price_one_year >= 100000) {
    //       $agent_level_words = '<span class="glyphicon glyphicon-thumbs-up"></span> 您目前是三级代理，享受83折价格（最高折扣），感谢您带来的傲人成绩。';
    //     }
    //   }

    // }else{
    //   $agent_level_words='您的代理级别被设定为恒定'.$agent_level.'级。永久享有折扣'.priceforagent($agent_level, 1).'。';
    // }//if($agent_level_locked=='NO'){
    //   ?>
    }
  }
}
