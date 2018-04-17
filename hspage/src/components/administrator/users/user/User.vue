<template>
  <div>
    <vue-title title="合适钻石 用户管理 历史纪录"></vue-title>
    <h3 class="pagetitle">用户详情</h3>
    <div id="userinfobox">
      <img style="height:168px; width:auto; display:inline-block; margin:18px;" src="u_icon" />
      <p id="nameofuser">u_name
        <span id="actionbtns">
          <a id="manageuserbtn" href="users.php?userid=usertoviewid"><span class="glyphicon glyphicon-option-horizontal"></span> 管理用户  &raquo;</a>
          <a id="sendmessagtbtn" href="message.php?user=usertoviewid"><span class="glyphicon glyphicon-comment"></span> 发消息 &raquo;</a>
        </span>
      </p>
      <p id="account-level">用户类型: 
      <!-- <?php
      if($u_agent=='YES'){	
        echo $u_agent_level.'级代购';
      }else{
        if(!empty($u_belong_to_agent)){
          if($u_relation_level==1){
            echo '代购'.$agentname.'的客人';
          }else{
            echo '普通客人，'.$u_account_level.'级折扣';
          }
        }else{
          echo '普通客人，'.$u_account_level.'级折扣';
        }
      }
      ?> -->
      </p>
      <p v-if="recommendername" class="recommender">
        推荐人：<a href="/administrator/user.php?user='.$u_recommenedby.'">recommendername&raquo;</a>
      </p>
      <p v-else>无推荐人
      </p>
      <p id="addressofuser">电话：u_tel</p>
      <p id="addressofuser">微信：u_wechatid</p>
      <p id="addressofuser">地址：u_address</p>
      <p id="balance">余额: u_balance 欧元</p>
    </div>
    <div v-if="recommendercustomerofthisuser" class="userbybox">
      <h4>该用户推荐的客人</h4>
      <ul v-for="(user, index) in users" v-bind:key="index">
        <li class="byuser"><a href="/administrator/user.php?user=by_id">by_name> &raquo;</a></li>
      </ul>
    </div>
    <!-- history -->
    <div class="history">
      <h3>购买历史纪录</h3>
      <div class="historybox history-diamond">
        <h4>钻石购买历史纪录</h4>
        <ul>
          <li class="dia-order-history">
            <span class="picfororder" style="background-image:url('crr_order_picture')"></span>
            <span class="dia-description">detail_forDiamond_byShape($crr_dia_shape, 'NAMECN') 钻石: crr_dia_weight 克拉，crr_dia_color 色， 净度 crr_dia_clarity</span>
            <span class="jewelry-description">配饰:
              <!-- <?php
              if(isset($jewelry_name)){
                echo $jewelry_name;
              }else{
                echo '无';
              }
              ?> -->
            </span>
            <span class="totalpricecrr">
              sold_price_total 欧元 <!-- (EuroToYuan($sold_price_total) 元人们币, EuroToDollar($sold_price_total) ?>美元) -->
            </span>
            <span class="orderedtime">ordered_at</span>
            <br class="clear" />
          </li>
        </ul>
      </div>
      <div class="historybox history-diamond">
        <h4>首饰购买历史记录</h4>
        <ul>
          <li class="dia-order-history">
            <span class="picfororder" style="background-image:url('crr_order_picture');"></span>
            <span class="jewelry-name">crr_jew_category</span>
            <span class="jewelry-description">crr_jew_name</span>
            <span class="totalpricecrr">
            总价格: sold_price_total欧元 (EuroToYuan($sold_price_total)元人民币, EuroToDollar($sold_price_total) 美元)
            </span>
            <span class="orderedtime">ordered_at</span>
            <br class="clear" />
          </li>
        </ul>
      </div>
      <p class="totalboughtamount">该用户的全部购买金额为 total_price欧元, 合 EuroToYuan(total_price) 元人民币, EuroToDollar(total_price) 美元</p>
    </div>
    <div v-if="isagent" class="history history-subseller">
      <h3>返点商品历史记录（由该代理的客人所购买）</h3>
      <div class="historybox history-diamond">
        <h4>钻石购买</h4>
        <ul>
          <li class="dia-order-history">
          <p class="subseller">来自用户: $usersname </p>
          <span class="picfororder" style="background-image:url('$crr_order_picture ');"></span>
          <span class="dia-description">detail_forDiamond_byShape($crr_dia_shape, 'NAMECN') 钻石: $crr_dia_weight 克拉，$crr_dia_color 色， 净度 $crr_dia_clarity </span>
          <span class="jewelry-description">配饰:
          <!-- <?php
          if(isset($jewelry_name)){
            echo $jewelry_name;
          }else{
            echo '无';
          }
          ?> -->
          </span>
          <span class="totalpricecrr">
          总价格: $sold_price_total 欧元 <!-- (EuroToYuan($sold_price_total) 元人们币, EuroToDollar($sold_price_total) ?>美元) -->
          </span>
          <span class="orderedtime">$ordered_at </span>
          <br class="clear" />
          </li>
        </ul>
      </div>
      <div class="historybox history-diamond">
        <h4>首饰购买</h4>
        <ul>
          <li class="dia-order-history">
            <p class="subseller">来自用户: usersname</p>
            <span class="picfororder" style="background-image:url('crr_order_picture');"></span>
            <span class="jewelry-name">crr_jew_category</span>
            <span class="jewelry-description">crr_jew_name</span>
            <span class="totalpricecrr">
            总价格: sold_price_total欧元 (EuroToYuan($sold_price_total)元人们币, EuroToDollar($sold_price_total) 美元)
            </span>
            <span class="orderedtime">ordered_at</span>
          </li>
        </ul>
      </div>
    </div>
    <div id="indication" style="position:fixed; width:100%; height:100%; background-color:rgba(255,255,255, 0.88); top:0; left:0; z-index:28; display:none;">
      <div class="indiinner" style="position:relative; width:200px; background-color:#0CF; margin:150px auto; padding:20px; text-align:center;">
        正在处理。。。
      </div>
    </div>
  </div>
</template>
<script src="./User.js"></script>
<style src="./User.css" scoped></style>
