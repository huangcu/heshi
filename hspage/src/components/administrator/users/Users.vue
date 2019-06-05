<template>
  <div>
    <vue-title title="合适钻石 管理界面"></vue-title>
    <h3 class="pagetitle">用户管理页面</h3>
    <p class="totalnumbox">总注册用户数：totaluser</p>
    <form method="post" action="userlist.php" id="searchuserform">
      <input name="searchkey" type="text" placeholder="姓名／微信名／email" />
      <input type="submit" value="搜索用户" />
    </form>
    <div class="usersfilter"> 
      <a class="filterbtn" href="userlist.php?order=registered">按注册日期排序</a>
      <a class="filterbtn" href="userlist.php?order=transactiontime">按交易日期排序</a>
      <a class="filterbtn" id="agentbtn" href="userlist.php?agent=yes">只显示代购用户 agentlist</a>
    </div>
    <div v-if="usernotfound">
    </div>
    <div v-else>
      <!-- $agentclass=' agent'; -->
      <ul v-for="(user,index)in users" v-bind:key="index"  id="userlist" class="userlist">
        <li class="user-record" v-bind:class="{isagentclass} " id="user_crr_user_id">
          <div class="userdetail-col1">
            <p class="userwechatinfo">
              <img class="usericon" src="r['icon']" /> 
              <span class="wechatname">wechat_name</span>
            </p>
            <h4 class="realname">
            <a class="detaillinker" href="/administrator/user.php?user=crr_user_id">
              name 
            </a>
            <span class="user-account-status">
              " *代购*  or not"
              推荐人：<a href="/administrator/user.php?user='.$recommenedby.'">rec_name&raquo; </a>
              无推荐人
            </span>
            <span class="regitime">
              <label>注册: </label> subscribed_time
            </span>
          </div>
          <div class="userdetail-col2">
            <span class="email">email</span>
            <span class="tel">电话: tel</span>
            <span class="we-chat_id">wechatID: wechat_id</span>
            <span class="address">address</span>
          </div>
          <div class="userdetail-col3">
            <span><label>总购买金额(一年内)</label>total_amount</span>
            <span><label>总钻石购买颗数(一年内):</label>bought_dia_num</span>
            <span><label>系统生成用户级别</label>account_level</span>
            <a class="addhistorybtn" href="/administrator/userhistory.php?user=crr_user_id">查看、添加购买记录</a>
          </div>
          <div v-if="notagent" class="userdetail-action">
            <p v-bind:style="account_level_assigned==-1, style='display:none'"> id="userlevelbox_crr_user_id">
              <label>指定用户级别</label>
              <select class="assigned_account_level" name="assigned_account_level" id="assigned_account_level_crr_user_id" onchange="assignUserLevel(crr_user_id)">
                <!-- <option value="0" <?php if($account_level_assigned==0){ echo 'selected="selected"';} ?> <?php if($r['account_level']>0){ echo 'disabled="disabled"'; } ?>>系统自动</option>
                <option value="1" <?php if($account_level_assigned==1){ echo 'selected="selected"';} ?> class="normaluseroption">1级(98折)</option>
                <option value="2" <?php if($account_level_assigned==2){ echo 'selected="selected"';} ?>>2级(95折)</option>
                <option value="3" <?php if($account_level_assigned==3){ echo 'selected="selected"';} ?>>3级(93折)</option>
                <option value="4" <?php if($account_level_assigned==4){ echo 'selected="selected"';} ?>>4级(90折)</option>
                <option value="5" <?php if($account_level_assigned==5){ echo 'selected="selected"';} ?>>5级(85折)</option>
                <option value="6" <?php if($account_level_assigned==6){ echo 'selected="selected"';} ?> class="recommenderoption">6级(推介人)</option> -->
              </select>
              <span class="userlevelchanged_inid" id="userlevelchanged_indi_crr_user_id">
                <span class="glyphicon glyphicon-ok"></span> 成功</span>
            </p>
            <p v-bind:style="account_level_assigned==-1, style='display:none'"> id="rpbox_crr_user_id">
              <label>推荐客人得的返点</label> 
              <select class="pointasrecommender" id="point_crr_user_id" onchange="changeRecommenderPoint(crr_user_id)">
                <!-- <option value="1" <?php if($point_as_recommender==1){ echo 'selected="selected"';} ?> class="normalpoint">1 个点</option>
                <option value="2" <?php if($point_as_recommender==2){ echo 'selected="selected"';} ?>>2 个点</option>
                <option value="3" <?php if($point_as_recommender==3){ echo 'selected="selected"';} ?>>3 个点</option>
                <option value="4" <?php if($point_as_recommender==4){ echo 'selected="selected"';} ?>>4 个点</option>
                <option value="5" <?php if($point_as_recommender==5){ echo 'selected="selected"';} ?>>5 个点</option>
                <option value="6" <?php if($point_as_recommender==6){ echo 'selected="selected"';} ?>>6 个点</option>
                <option value="7" <?php if($point_as_recommender==7){ echo 'selected="selected"';} ?>>7 个点</option>
                <option value="8" <?php if($point_as_recommender==8){ echo 'selected="selected"';} ?>>8 个点</option>
                <option value="9" <?php if($point_as_recommender==9){ echo 'selected="selected"';} ?>>9 个点</option>
                <option value="10" <?php if($point_as_recommender==10){ echo 'selected="selected"';} ?>>10 个点</option>
                <option value="11" <?php if($point_as_recommender==11){ echo 'selected="selected"';} ?>>11 个点</option>
                <option value="12" <?php if($point_as_recommender==12){ echo 'selected="selected"';} ?>>12 个点</option> -->
              </select>
                <span class="pointchanged_inid" id="pointchanged_indi_crr_user_id">
                  <span class="glyphicon glyphicon-ok"></span> 成功</span>
            </p>
            <p style="color:#900;">
              <!-- <?php
              $referencepriceattr='';
              $ratio_attr=' disabled';
              if($account_level_assigned==-1){
                $referencepriceattr=' checked="checked"';
                $ratio_attr='';
              }
              ?> -->
              <span class="glyphicon glyphicon-eye-close"></span> 只能看参考价
                <input type="checkbox"  class="referenceprice" id="referenceprice_crr_user_id" onclick="thisUserCompare(crr_user_id)" referencepriceattr/>
              <span class="referenceprice_indi" id="referenceprice_indi_crr_user_id"><span class="glyphicon glyphicon-ok"></span> 成功</span>
              </p>
              <p class="rrboxratio_attr" id="rrbox_crr_user_id">
              <label>参考价系数：</label>
                <input type="number" class="referencepriceratio" id="referencepriceratio_crr_user_id" value="reference_price_ratio" step=".01" ratio_attr/> 
                <button type="button" class="referencepriceratiosavebtn" id="referencepriceratiosavebtn_crr_user_id" onclick="saveReferencePriceRatio(crr_user_id)" ratio_attr>保存</button>
              <span class="referencepriceratio_indi" id="referencepriceratio_indi_crr_user_id">
                <span class="glyphicon glyphicon-ok"></span> 成功
              </span>
            </p>
          </div>
          <div v-else>
            <p>代购级别请在下面设定 &darr;</p>
          </div>
          <div class="userdetail-action action2">
            <button v-if="notagent" class="assignresellerbtn" type="button" id="assignresellerbtn_crr_user_id" onclick="assignagent(crr_user_id)"><span class="glyphicon glyphicon-certificate"></span> 指定为代购</button>
            <button v-else class="resignresellerbtn" type="button" id="resignresellerbtn_crr_user_id" onclick="resignagent(crr_user_id)"><span class="glyphicon glyphicon-minus-sign"></span> 解除代购</button>
          </div>
          <div class="userdetail-action action3">
            <a class="historybtn" href="/administrator/user.php?user=crr_user_id">
              <span class="glyphicon glyphicon-th-list"></span> 用户详情
              <span v-if="totalpointnum!=0" class="newpointnote">(<span class="glyphicon glyphicon-bell"></span>totalpointnum)</span>
            </a>
            <a class="heshibibtn" href="/administrator/heshibi.php?user=crr_user_id">发／扣 合适币</a>
            <a class="discountticketbtn" href="/administrator/discountticket.php?user=crr_user_id">发打折券</a>
          </div>
          <br class="clear" />
          <div v-if="isagent" class="agent-contact-info-box">
            <p class="levelinputbox">代购级别：
              <!--
              <input type="text" name="agent_level_crr_user_id" id="agent_level_crr_user_id" value="crr_user_agent_level" />
              <button id="agent_level_btn_crr_user_id" type="button" onclick="updateAgentInfo('level','crr_user_id')" class="agent_level_btn">设定</button>
              -->
              <select class="agentleveloption" name="agent_level_crr_user_id" id="agent_level_crr_user_id" onchange="updateAgentInfo('level','crr_user_id')">
              <!-- <option value="1" <?php if($crr_user_agent_level==1){ echo 'selected="selected"';} ?>>1级</option>
              <option value="2" <?php if($crr_user_agent_level==2){ echo 'selected="selected"';} ?>>2级</option>
              <option value="3" <?php if($crr_user_agent_level==3){ echo 'selected="selected"';} ?>>3级</option> -->
              </select>
              <span class="agentlevelsetindi" id="agentlevelsetindi_crr_user_id">
                <span class="glyphicon glyphicon-ok"></span> 成功
              </span>
              
              <!-- if($agent_level_locked=='YES'){
                $crr_checked_txt=' checked="checked"';
                $crr_checked_words='<span class="glyphicon glyphicon-lock"></span> 锁定状态';
              }else{
                $crr_checked_txt='';
                $crr_checked_words='未锁定';
              } -->
              <label>锁定代购级别</label> 
              <input type="checkbox" class="agentlevellocked" name="agentlevellocked_crr_user_id" id="agentlevellocked_crr_user_id" onclick="agentLevelLockStatus(crr_user_id)" crr_checked_txt />
              <span class="agentlevellockedindiwords" id="agentlevellockedindiwords_crr_user_id">crr_checked_words</span>
            </p>
            <h4>代理联系信息</h4>
            <form action="uploadingqrimage.php" method="post" enctype="multipart/form-data" id="uploadImage2" class="uploadImage"> 
              <label>微信二维码图片:</label>
              <input type="hidden" name="picturewhere" value="<?php echo $crr_user_id; ?>" />
              <img class="qrimage" id="qrimage_<?php echo $crr_user_id; ?>" width="58" src="<?php echo $qrimage_path; ?>" />
              <input type="file" name="image" id="image_selecting_<?php echo $crr_user_id; ?>" class="image_selecting">                
            </form>
            <div class="agentinfoboxinner">
              <p>
              <label>电话:</label>
              <input type="text" name="agent_phone_<?php echo $crr_user_id; ?>" id="agent_phone_<?php echo $crr_user_id; ?>" value="<?php echo $phone; ?>" />
              <button id="agent_phone_btn_<?php echo $crr_user_id; ?>" type="button" onclick="updateAgentInfo('phone','<?php echo $crr_user_id; ?>')">更新</button>
              </p>
              <p>
              <label>邮箱:</label>
              <input type="text" name="agent_email_<?php echo $crr_user_id; ?>" id="agent_email_<?php echo $crr_user_id; ?>" value="<?php echo $email; ?>" />
              <button id="agent_email_btn_<?php echo $crr_user_id; ?>" type="button" onclick="updateAgentInfo('email','<?php echo $crr_user_id; ?>')">更新</button>
              </p>
            </div>
            <br class="clear" />
          </div>
        </li>
      </ul>
      <div class="progress" style="width:300px; position:relative; left:488px;">
        <p>上传中</p>
        <div class="bar"></div >
        <div class="percent">0%</div >
      </div>
      <p id="indi" style="position:fixed; top:50px; left:50px; width:150px; text-align:center; padding: 12px 30px; border-style:solid; border-width:3px; border-color:#000; background-color:#FF9; font-size:24px; font-family:'Lucida Sans Unicode', 'Lucida Grande', sans-serif; display:none;"> 处理中 ... </p>
      <p style="height:50px;">&nbsp;</p>
      <div id="status" style="display:none"></div>
    </div>
  </div> 
</template>
<script src="./Users.js"></script>
<style src="./Users.css" scoped></style>
