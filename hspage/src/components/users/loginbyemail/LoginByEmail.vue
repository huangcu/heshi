<template>
  <div>
    <vue-title title="合适钻石  BEYOU DIAMOND  用户登录"></vue-title>
    <app-header></app-header>
    <div id="maincontentbox">
      <div id="userloginform">
        <h4>我的合适账户</h4>
        <p v-if="login_feedback"> {{ login_feedback }} </p>
        <p v-else-if="upgrade_feedback"> {{ upgrade_feedback }}</p>
        <div v-if="appointment !==''">
          <p class="loginforappointmentwords">欢迎您和我们预约看钻石和首饰<br />请先登录您的帐户</p>
        </div>
        <div v-if="account!==''">
          <form action="" @submit.prevent="login()">
            <input name="email" type="text" class="inputtextforloginform" v-model="username" v-validate="'required'" placeholder="电子邮箱" />
            <input name="accountpassword" type="password" class="inputtextforloginform" v-model="password" v-validate="'required'" placeholder="密码" />
            <p v-if="referer!==''" id="recommend-ticket-box">
              <label><span class="glyphicon glyphicon-certificate"></span> 高级用户推荐码 <span class="glyphicon glyphicon-certificate"></span></label>
              <input type="text" disabled="disabled" id="therecommendticketcode" name="therecommendticketcode" :value="referer" />
            </p>
            <input name="loginformsubmit" type="submit" id="loginbtn" value="登录" />
          </form>
          <p class="getbackpassword">忘记密码了？<a href="retrievepass.php">点击找回 &raquo;</a></p>
          <p class="registerindibox">还没有合适账户? <a id="registeropenbtn" href="/register">立即注册</a></p>
        </div>
        <div v-else>
          <p>您已经登录。<a href="'/myaccount">回到我的帐户 &raquo;</a></p>
        </div>
      </div>
    </div>
  </div>
</template>
<script src='./LoginByEmail.js'></script>
<style src='./LoginByEmail.css' scoped></style>
