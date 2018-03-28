<template>
  <div>
    <p style='text-align:center;'><img :src="ImgURL('logo.png')" style="width:128px;" /></p>
    <h1 style="margin-top:20px; text-align:center; color:#FFF;">合适钻石 管理界面</h1>
    <p v-if="wrongmessage" class='warning'>{{ wrongmessage }}</p>
    <form id="loginform" @submit.prevent="login()">
      <p style="text-align:center;">
        <input type="text" name="user" id="user" v-model="username" v-validate="'required'" placeholder="用户名"><br />
        <input type="password" name="pwd" id="pwd" v-model="password" v-validate="'required'" placeholder="密码"><br />
        <input name="denglu" type="submit" id="login" value="登陆" style="background-color:#000; margin-top:55px;">
      </p>
    </form>
  </div>
</template>
<script src="./Login.js"></script>
<style src="./Login.css" scoped></style>
