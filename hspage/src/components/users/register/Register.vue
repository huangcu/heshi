<template>
  <div>
    <vue-title title="合适钻石  BEYOU DIAMOND  用户注册"></vue-title>
    <app-header></app-header>
    <div id="signupform">
      <h4>注册合适账户</h4>
      <h6 class="memeberbenefitexpl"><span class="glyphicon glyphicon-piggy-bank"></span> 立即享受会员折扣和积分</h6>
      <form id="the_signup_form" @submit.prevent="checkForm" novalidate="true">
        <p v-if="registerErrors.length">
          <b>Please correct the following error(s):</b>
          <ul>
            <li v-bind:key="error" v-for="error in registerErrors">{{ error }}</li>
          </ul>
        </p>
        <p>
          <label for="email">电子邮箱</label>
          <input type="email" name="email" id="newaccountemail" class="inputtextforloginform" v-model="email" placeholder="电子邮箱">
        </p>
        <p>
          <label for="username">用户名</label>
          <input type="text" name="username" id="newaccountusername" class="inputtextforloginform" v-model="username" placeholder="用户名">
        </p>
        <p>
          <label for="password">密码</label>
          <input type="password" name="username" id="newaccountpassword" class="inputtextforloginform" v-model="password" placeholder="密码">
        </p>
        <p>
          <label for="password2">密码确认</label>
          <input type="password" name="username" id="newaccountpasswordagain" class="inputtextforloginform" v-model="password2" placeholder="确认密码">
        </p>
        <p>
          <label for="realname">姓名</label>
          <input type="text" name="realname" id="realname" class="inputtextforloginform" v-model="realname" placeholder="姓名">
        </p>
        <p>
          <label for="wechatid">微信ID(可选)</label>
          <input type="text" name="wechatid" id="wechatid" class="inputtextforloginform" v-model="wechatid" placeholder="微信ID">
        </p>
        <p>
          <label for="tel">手机号码(可选)</label>
          <input type="text" name="tel" id="tel" class="inputtextforloginform" v-model="cellphone" placeholder="手机号码">
        </p>
        <p>
          <input type="submit" id="registerformsubmitbtn" value="Submit">
        </p>
      </form>
    </div>
  </div>
</template>
<script src='./Register.js'></script>
<style src='./Register.css' scoped></style>
