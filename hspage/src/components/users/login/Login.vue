<template>
  <div>
    <vue-title title="BEYOU DIAMOND|合适钻石 登录"></vue-title>
    <div id="maincontentbox">
      <p v-if="login_feedback !== ''"  class="loginfeedbackwords">{{ login_feedback }}</p>
      <p v-if="upgrade_feedback !== ''" class="registerfeedbackmessage">{{ upgrade_feedback }}</p>
      <div v-if="previewmode !== ''" id="userloginform" class="forpreviewer">
        <ul>
          <li class="share-instruction-line"><strong>使用微信发送：</strong><br />点击右上角 <span class="glyphicon glyphicon-option-horizontal"></span> 或 <span class="glyphicon glyphicon-option-vertical"></span> 打开菜单，选择“传送给朋友”或“分享到朋友圈”</li>
          <li class="share-instruction-line"><strong>使用电脑网页发送：</strong><br />复制网页上面的网址，用邮件、QQ、微信等方式发给你的朋友</li>
        </ul>
        <p class="share-info">你的朋友将会看到一个二维码。长按或扫描该二维码，注册后即成为您的客人。您的客人购买了合适商品，合适会按照我们的奖励规则，给您奖励。</p>
      </div>
      <div id="userloginform">
        <h4>我的合适账户</h4>
        <div v-if="userprofile == ''">
          <p v-if="appointment!==''" class="loginforappointmentwords">欢迎您和我们预约看钻石和首饰<br />请先登录您的帐户</p>
          <p class="weixin-sign-title">
            <img id="weixin-indi" src="../../../_images/constant/weixin2.png" /> 使用微信注册、登录</p>
          <p class="weixin-sign-expl">请用微信扫(或长按)二维码</p>
          <div id="qrcode-img-container">
            <span id="loading-txt">生成二维码中...</span>
            <img id="the_account_qrcode" :src="qrCodeSrc"/>
          </div>
          <p class="expire_expl">两分钟内有效</p>
          <p v-if="referer!==''" id="recommend-ticket-box">
            <label style="font-size:10px;">推荐人:</label>
            <input type="text" readonly="readonly" id="therecommendticketcode" name="therecommendticketcode-readonly" :value="referer" />
          </p>
          <!-- <form action="qr_sign.php" method="post" id="sign-form">
            <input type="hidden" name="wechatopenID" id="wechatopenID" value="" />
            <input type="hidden" name="therecommendticketcode" value="<>" />
          </form> -->
          <p class="registerindibox">
            <a href="/users/loginbyemail" id="registeropenbtn"> 电子邮箱登录 &raquo;</a>
          </p>
        </div>
        <div v-else>
          <p>您已经登录。<a href="/users/myaccount" >回到我的帐户 &raquo;</a></p>
          <div class="logoutbtnbox" style="text-align:center;">
            <form method="post" @submit="logout" action="/users/myaccount" id="logoutform">
              <input type="hidden" name="logout" value="confirmed" />
              <input type="submit" value="登出账号" name="logoutbtn" id="logoutbtn" />
            </form>
          </div>
        </div>
      </div><!-- maincontentbox -->
      <div id="sign_status" style="display:none;">
        <p id="status"></p>
      </div>
    </div>
  </div>
</template>

<script src="./Login.js"></script>
<style src="./Login.css" scoped></style>
