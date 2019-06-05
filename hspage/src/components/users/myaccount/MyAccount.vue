<template>
  <div>
    <vue-title title="合适钻石 | BEYOU DIAMOND | 我的账户"></vue-title>
    <div class="contentmainwrapper">
      <h2>
      <span id="thetxtofpagetitle">我的合适账户</span>
      </h2>
      <div v-if="isAgent">
        <agent></agent>
      </div>
      <div v-else>
        <div v-if ="currentUserProfile.user_level==6">
          <!-- _includes/documentparts/myaccount-recommender.php -->
        </div>
        <div v-else>
          <!-- _includes/documentparts/myaccount-highleveluser.php -->
        </div>
      </div>
      <!-- account info -->
      <div class="accountinfobox" id="personalinfo">
        <h3>个人信息</h3>
        <p v-if="currentUserProfile.email==''" id="updatewarning">
          <span class="glyphicon glyphicon-exclamation-sign"></span> 为了交易顺利，请更新这里的信息
        </p>
        <div>
          <form action="" method="post">
            <p><label>姓名</label><input type="text" name="name_new" id="account_name" :value="currentUserProfile.name"></p>
            <p><label>电话</label><input type="text" name="tel_new" id="account_tel" :value="currentUserProfile.cellphone"></p>
            <p><label>微信ID</label><input type="text" name="wechatid_new" id="account_wechatid" :value="currentUserProfile.wechat_id"></p>
            <p><label>地址</label><input type="text" name="address_new" id="account_address" :value="currentUserProfile.address"></p>
            <p><label>邮箱</label><input type="text" name="email" :value="currentUserProfile.email"></p>
            <!-- <p><label>密码</label><input type="text" name="website_password" value="website_password"></p> -->
            <p><input type="submit" id="changeinfosavebtn" value="保存修改"  /></p>
          </form>
          <p><a href="changepassword.php">修改登录密码</a></p>
        </div>
      </div><!-- end of account info -->
      <div class="logoutbtnbox">
        <form action="" @submit.prevent="logout()" id="logoutform">
          <input type="hidden" name="logout" value="confirmed" />
          <input type="submit" value="登出账号" name="logoutbtn" id="logoutbtn" />
        </form>
      </div>
    </div> <!-- end of contentmainwrapper -->

    <!-- show the qr code for who haven't linked their wechat account -->
    <div v-if ="currentUserProfile.wechat_id==''" id="qrcodebg">
      <div id="qrcode-box">
        <h3>关注我们的微信客服</h3><p class="qrexpl"><img id="weixin-indi" :src="ourWXCode" />请<strong style="font-size:16px;">微信扫描(或长按)</strong>下面的二维码，<br />以验证您的帐户和享有完整的客户服务</p>
        <!-- get_account_qrcode.php -->
        <div v-if ="justsignedup">
          <h3>注册成功!</h3><p class="qrexpl"><img id="weixin-indi" :src="ourWXCode" />微信扫一扫，加关注，您的合适帐户即开通</p>
        </div>
        <div id="qrcode-img-container">
          <span id="loading-txt">生成二维码中，稍等...</span>
          <img id="the_account_qrcode" :src="qrCodeSrc"/>
        </div>
        <p class="expire_expl">有效期为两分钟，过期请刷新页面</p>
        <p class="no_weixin">没有微信? 请<a href="contact.php">联系我们</a></p>
        <div class="logoutbtnbox">
          <form action="" @submit.prevent="logout()" id="logoutform-inqrstatus">
            <input type="hidden" name="logout" value="confirmed" />
            <input type="submit" value="登出账号" name="logoutbtn" id="logoutbtn-qrstates" />
          </form>
        </div>
      </div>
    </div>
    <div v-else-if="wechat_open_idwechatnameicon">
      <div id="qrcodebg-watch">
        <div id="qrcode-box">
          <h3>关注我们的微信公众号</h3><p class="qrexpl"><img id="weixin-indi" :src="ourWXCode" />亲，您上次扫码注册时似乎没有关注我们的公众号，为了您数据的安全，请再扫一下，然后别忘了关注我们！</p>
          我们的公众号：(请扫码或长按)
          <div id="qrcode-img-container">
            <img id="the_account_qrcode_w" :src="ourQRCode" />
          </div>
          <p class="no_weixin">没有微信? 请<a href="contact.php">联系我们</a></p>
        </div>
        <button id="x-it" v-on:click="x()" style="position:absolute; bottom:25px; right:25px;">关闭</button>
      </div>
    </div>
    <div v-if="updatefeedback">
      <div id="updatedinindication" style="position:fixed; width:100%; height:58px; top:0; left:0; text-align:center;">
        <span id="updatedindi-inner" style="position:relative; display:inline-block; padding:3px 15px; background-color:#090; color:#FFF; font-size:16px;">{{ userinfoupdatedfeedback }}</span>
      </div>
    </div>
    <div id="upgratedinindication" style="position:fixed; width:100%; height:58px; top:0; left:0; text-align:center; display:none;">
      <span id="upgratedindi-inner" style="position:relative; display:inline-block; padding:3px 15px; background-color:#090; color:#FFF; font-size:16px;">恭喜您成功升级了账户!</span>
    </div>
  </div>
</template>

<script src="./MyAccount.js"></script>
<style src="./MyAccount.css" scoped></style>
