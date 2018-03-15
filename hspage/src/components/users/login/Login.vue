<template>
  <div>
    <vue-title title="BEYOU DIAMOND|合适钻石 登录"></vue-title>
    <div id="userloginform">
      <h4>我的合适账户</h4>
      <p class="loginforappointmentwords">欢迎您和我们预约看钻石和首饰<br />请先登录您的帐户</p>
      <p class="weixin-sign-title"><img id="weixin-indi" src="../../../_images/constant/weixin2.png" /> 使用微信注册、登录</p>
      <p class="weixin-sign-expl">请用微信扫(或长按)二维码</p>
      <div id="qrcode-img-container">
        <span id="loading-txt">生成二维码中...</span>
        <img id="the_account_qrcode" :src="qrCodeSrc"/>
      </div>
      <p class="expire_expl">两分钟内有效</p>
      <form action="qr_sign.php" method="post" id="sign-form">
        <input type="hidden" id="signID" value="<?php echo $the_sign_id; ?>" />
        <input type="hidden" name="wechatopenID" id="wechatopenID" value="" />
        <p id="recommend-ticket-box">
        <label style="font-size:10px;">推荐人:</label>
        <input type="text" readonly="readonly" id="therecommendticketcode" name="therecommendticketcode-readonly" value="<>" />
        </p>
        <input type="hidden" name="therecommendticketcode" value="<>" />
      </form>
      <p class="registerindibox">
        <a href="/loginbyemail" id="registeropenbtn"> 电子邮箱登录 &raquo;</a>
      </p>
      <p>您已经登录。<a href="/myaccount?_account" >回到我的帐户 &raquo;</a></p>
      <div class="logoutbtnbox" style="text-align:center;">
        <form method="post" @submit="logout" action="/MyAccount" id="logoutform">
          <input type="hidden" name="logout" value="confirmed" />
          <input type="submit" value="登出账号" name="logoutbtn" id="logoutbtn" />
        </form>
      </div>
    </div><!-- maincontentbox -->
    <div id="sign_status" style="display:none;">
    </div>
  </div>
</template>

<script src="./Login.js"></script>
<style src="./Login.css" scoped></style>
