<template>
  <div>
    <vue-title title="合适钻石 用户管理 历史纪录"></vue-title>
    <h3 class="pagetitle">用户详情</h3>
    <div id="userinfobox">
      <img style="height:168px; width:auto; display:inline-block; margin:18px;" src="u_icon" />
      <div class="userinfotxtbox">
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
      </div>
      <br style="clear:both;" />
      <p><a id="addnewbtn" href="addhistory.php?user=userid">
        <span class="glyphicon glyphicon-plus"></span> 添加购买纪录</a>
      </p>
    </div>
    <!-- history -->
    <div class="history">
      <h3>购买历史纪录</h3>
      <div class="historybox history-diamond">
        <h4>钻石购买历史纪录</h4>
        <ul>
          <li v-if="found" v-for="(d,index) in drecords" v-bind:key="index" class="dia-order-history" v-bind:class="{oneyear_or_no}" id="history_record_id_dia">
            <span v-if="ordertimenotpass1year" class="orderedtime">
              <span class="glyphicon glyphicon-ok">积分有效</span>
              <span class="profitable-box">
                <input type="checkbox" class="profitablecheckbox" id="profitablechekbox-crr_dia_stockref" checked onchange="updateSoldDiaProfittingStatus(crr_dia_stockref)" /> 可返点 
                <span class="indi-update-profitting-status" id="indi-update-profitting-status-crr_dia_stockref">
                  <span class="glyphicon glyphicon-ok"></span> 已更新
                </span>
              </span>
            </span>
          <span v-else class="orderedtime">
            <span>超过一年，不积分</span>
            <span class="profitable-box no">
              <input type="checkbox" class="profitablecheckbox" id="profitablechekbox-crr_dia_stockref" onchange="updateSoldDiaProfittingStatus(crr_dia_stockref)" /> 可返点 
              <span class="indi-update-profitting-status" id="indi-update-profitting-status-crr_dia_stockref">
                <span class="glyphicon glyphicon-ok"></span> 已更新
              </span>
            </span>
          </span>

          <span v-if="idofuser!=$usertoviewid" class="buyertxt">
            <span class="glyphicon glyphicon-user"></span>由'.$u_name.'的客人 
            <a class="clientofthisuserlinker" href="user.php?user='.$idofuser.'">'.$nameofuser.'</a> 购买
          </span>
          <span class="picfororder" style="background-image:url('crr_order_picture');"></span>
          <span class="dia-description">detail_forDiamond_byShape($crr_dia_shape, 'NAMECN')钻石: crr_dia_weight ct，crr_dia_color色， crr_dia_clarity, 荧光:crr_dia_fluo, crr_dia_lab(crr_dia_certi_num)</span>
          <span v-if="jewelry_name" class="jewelry-description">配饰: jewelry_name</span>
          <span v-else  class="jewelry-description">配饰: 无</span>
          <span class="totalpricecrr shouldpay">预定时客人价格: ($jewellery_price+$diamond_price)欧元</span>
          <span class="totalpricecrr">最终付款: sold_price_total欧元 
            <span v-if="extra_info==added_manually" class="totaldiscountbox" id="totaldiscountbox-record_id_dia">
              <span v-if="final_price_ratio_empty" class="glyphicon glyphicon-warning-sign"></span>
              实际折扣：
              <select class="actualdiscountchooser" id="actualdiscountchooser-record_id_dia" onchange="updateOrderActualDiscount(record_id_dia)">
                <option value="">请选择...</option>
                <!-- <option value="1"<?php if($final_price_ratio==1){ echo ' selected="selected"'; } ?>>无折扣</option>
                <option value="0.98"<?php if($final_price_ratio==0.98){ echo ' selected="selected"'; } ?>>98折</option>
                <option value="0.95"<?php if($final_price_ratio==0.95){ echo ' selected="selected"'; } ?>>95折</option>
                <option value="0.93"<?php if($final_price_ratio==0.93){ echo ' selected="selected"'; } ?>>93折</option>
                <option value="0.9"<?php if($final_price_ratio==0.9){ echo ' selected="selected"'; } ?>>9折</option>
                <option value="0.87"<?php if($final_price_ratio==0.87){ echo ' selected="selected"'; } ?>>87折</option>
                <option value="0.85"<?php if($final_price_ratio==0.85){ echo ' selected="selected"'; } ?>>85折</option>
                <option value="0.83"<?php if($final_price_ratio==0.83){ echo ' selected="selected"'; } ?>>83折</option>
                <option value="0.80"<?php if($final_price_ratio==0.8){ echo ' selected="selected"'; } ?>>83折以下</option> -->
              </select>
            </span>
            <span v-else> 实际折扣:round($final_price_ratio,3)
              <span class="totaldiscountbox" id="totaldiscountbox-record_id_dia">
                <span class="glyphicon glyphicon-warning-sign"></span> 实际折扣：
                <select class="actualdiscountchooser" id="actualdiscountchooser-record_id_dia" onchange="updateOrderActualDiscount(record_id_dia)">
                  <option value="">请选择...</option>
                  <option value="1">无折扣</option>
                  <option value="0.98">98折</option>
                  <option value="0.95">95折</option>
                  <option value="0.93">93折</option>
                  <option value="0.9">9折</option>
                  <option value="0.87">87折</option>
                  <option value="0.85">85折</option>
                  <option value="0.83">83折</option>
                  <option value="0.80">83折以下</option>
                </select>
              </span>
            </span>
          </span>

          <span class="specialnotecontainer">
            <textarea class="specialnote" id="specialnote-record_id_dia" placeholder="无特别标注,请添加">special_notice</textarea>
            <span class="oldspecialnot" id="oldspecialnote-record_id_dia" style="display:none;">special_notice</span>
            <button class="specialnoteBTN" id="specialnoteBTNrecord_id_dia" onclick="specialnoteSave(record_id_dia)">保存</button>
          </span>
          <div v-if="extra_info=='ADDED MANUALLY'">
            <span class="manuallyaddedrecord">人工添加的纪录</span>
            <button type="button" class="deleteorderbtn" onclick="deleteOrder(record_id_dia)">
              <span class="glyphicon glyphicon-trash"></span> 删除</button>
          </div>
        </li>
        <li v-else class="dia-order-history no-record">还没有呢</li>
        <li id="noprofittitle">&darr; 不能积分购买纪录 &darr;</li>
        <li class="dia-order-history no-profit" id="history_na_record_id_dia">
          <span class="orderedtime"> ordered_at 
            <span class="glyphicon glyphicon-ban-circle"></span> 不能积分 
          </span>
          <span v-if="idofuser!=$usertoviewid" class="buyertxt"><span class="glyphicon glyphicon-user"></span>由'.$u_name.'的客人
            <a class="clientofthisuserlinker" href="user.php?user='.$idofuser.'">'.$nameofuser.'</a> 购买
          </span>
          <span class="picfororder" style="background-image:url('crr_order_picture');"></span>
          <span class="dia-description">detail_forDiamond_byShape($crr_dia_shape, 'NAMECN')钻石: crr_dia_weight ct，crr_dia_color色， crr_dia_clarity, 荧光:crr_dia_fluo, crr_dia_lab(crr_dia_certi_num)</span>
          <span v-if="jewelry_name" class="jewelry-description">配饰:jewelry_name</span>
          <span v-else  class="jewelry-description">配饰: 无</span>
          <span class="totalpricecrr shouldpay">应付款: ($jewellery_price+$diamond_price)欧元</span>
          <span class="totalpricecrr">最终付款: sold_price_total欧元</span>
          <button v-if="extra_info=='ADDED MANUALLY'" type="button" class="deleteorderbtn" onclick="deleteOrderNA(record_id_dia)">
          <span class="glyphicon glyphicon-trash"></span> 删除</button>
        </li>
      </ul>
    </div>
    <div class="historybox history-diamond">
      <h4>首饰购买历史记录</h4>
      <ul>
        <li v-if="found" v-for="(d,index) in drecords" v-bind:key="index" class="dia-order-history$listclass">
          <span v-if="ordertimenotpass1year" class="orderedtime">
            <span class="glyphicon glyphicon-ok">积分有效</span>
            <span class="profitable-box">
              <input type="checkbox" class="profitablecheckbox" id="profitablechekbox-crr_dia_stockref" checked onchange="updateSoldDiaProfittingStatus(crr_dia_stockref)" /> 可返点 
              <span class="indi-update-profitting-status" id="indi-update-profitting-status-crr_dia_stockref">
                <span class="glyphicon glyphicon-ok"></span> 已更新
              </span>
            </span>
          </span>
          <span v-else class="orderedtime" >
            <span>超过一年，不积分</span>
            <span class="profitable-box no">
              <input type="checkbox" class="profitablecheckbox" id="profitablechekbox-crr_dia_stockref" onchange="updateSoldDiaProfittingStatus(crr_dia_stockref)" /> 可返点 
              <span class="indi-update-profitting-status" id="indi-update-profitting-status-crr_dia_stockref">
                <span class="glyphicon glyphicon-ok"></span> 已更新
              </span>
            </span>
          </span>

          <span v-if="idofuser!=$usertoviewid" class="buyertxt">
            <span class="glyphicon glyphicon-user"></span>由'.$u_name.'的客人
            <a class="clientofthisuserlinker" href="user.php?user='.$idofuser.'">'.$nameofuser.'</a> 购买
          </span>
          <span class="picfororder" style="background-image:url('$crr_order_picture');"></span>
          <span class="jewelry-name">$crr_jew_category</span>
          <span class="jewelry-description">$crr_jew_name</span>
          <span class="totalpricecrr shouldpay" style="display:none;">应付款: $jewellery_price欧元</span>
          <span class="totalpricecrr">最终付款: $sold_price_total 欧元; (实际折扣: round($crr_jew_final_price_ratio,3))</span>
        </li>
        <li v-else class="dia-order-history no-record">还没有呢</li>
      </ul>
      <div id="indication" style="position:fixed; width:100%; height:100%; background-color:rgba(255,255,255, 0.88); top:0; left:0; z-index:28; display:none;">
        <div class="indiinner" style="position:relative; width:200px; background-color:#0CF; margin:150px auto; padding:20px; text-align:center;">
          正在处理。。。
        </div>
      </div>
  </div>
</template>
<script src="./UserBuyRecords.js"></script>
<style src="./UserBuyRecords.css" scoped></style>
