<template>
  <div>
    <vue-title title="合适钻石  BEYOU DIAMOND  用户登录"></vue-title>
    <app-header></app-header>
    <div id="maincontentbox">
      <div id="userloginform">
      <h4><span class="glyphicon glyphicon-ok-circle"></span> 恭喜！您成功创建了合适帐户</h4>
      <p v-if="login_feedback" class="loginfeedbackwords">{{ login_feedback }}</p>
      <p v-else-if="upgrade_feedback"  class="registerfeedbackmessage">{{ upgrade_feedback }}</p>
      <p v-if="appointment" class="loginforappointmentwords">欢迎您和我们预约看钻石和首饰<br />请先登录您的帐户</p>
      <div v-if="account">
        <button style="display:block;" id="showqrsignup_btn" v-on:click="showqrSignUp()" type="button">
          <span class="glyphicon glyphicon-user"></span> 进入我的帐户 &raquo;
        </button>
        <p class="link-old-account-box">
          <span class="glyphicon glyphicon-info-sign"></span> 如果您之前已经用邮箱注册过，点这里
          <button v-if="account" id="showqrsignin_btn" v-on:click="showqrSignIn()" type="button">
            <span class="glyphicon glyphicon-lock"></span> 连接已有帐户
          </button>
        </p>
        <form id="signin-form" style="display:none;">
          <input type="hidden" name="wechatopenID" id="wechatopenID" :value="wechatopenID" />
          <input name="email" type="text" class="inputtextforloginform" placeholder="电子邮箱" />
          <input name="accountpassword" type="password" class="inputtextforloginform" placeholder="密码" />
          <p v-if="therecommendticketcode" id="recommend-ticket-box">
            <label>
              <span class="glyphicon glyphicon-certificate"></span> 高级用户推荐码
              <span class="glyphicon glyphicon-certificate"></span>
            </label>
            <input type="text" disabled="disabled" id="therecommendticketcode" name="therecommendticketcode" value="<?php echo $therecommendticketcode; ?>" />
          </p>
          <button id="qrsignin_btn" v-on:click="qrSignIn()" type="button">登录</button>
        </form>
        <form id="signup-form">
          <input type="hidden" name="wechatopenID" id="wechatopenID_new" :value="wechatopenID" />
          <input name="newaccountemail" type="hidden" class="inputtextforloginform" id="newaccountemail" value="<?php echo $wechatopenID; ?>@beyoudiamond.com" />
          <input name="newaccountpassword" type="hidden" class="inputtextforloginform" id="newaccountpassword" value="<?php echo str_replace('-','', str_replace('_', '', $wechatopenID)); ?>" />
          <input v-if="therecommendticketcode" type="hidden" id="therecommendticketcode" name="therecommendticketcode" value="<?php echo $therecommendticketcode; ?>" />
        </form>
      </div>
      <div v-else>
        <p>您已经登录。<a href="/myaccount">回到我的帐户 &raquo;</a></p>
      </div>
    </div>
  </div>
  <div id="sign_status" style="display:none;"></div>
  </div>
</template>
<script src='./QRSign.js'></script>
<style src='./QRSign.css' scoped></style>
