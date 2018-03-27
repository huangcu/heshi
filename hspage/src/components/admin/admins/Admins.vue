<template>
  <div>
    <vue-title title="合适钻石 用户管理"></vue-title>
    <h3 class="pagetitle">内部人员管理</h3>
    <div v-if="activeAdminUsers === null">
      <p style="text-align:center;">没有找到用户</p>
    </div>
    <div v-else>
      <ul id="userlist" class="userlist">
        <li v-for="(user, index) in activeAdminUsers" :key="index" class="user-record">
          <span class="name">{{ user.real_name }} </span>
          <span class="icon">👦 头像： <img class="iconpic" :src="user.icon" /></span>
          <span class="username">用户名: {{ user.username}}</span>
          <span class="password">密码: xxxxxxxx</span>
          <span v-if="user.admin.admin_level===1" class="admin_level superadmin">
            <span class="glyphicon glyphicon-star"></span>
            超级管理员
          </span>
          <span v-else class="admin_level normaladmin">
            网站内容管理员
          </span>
          <a class="editbtn" :href="'editadmin?id='+userprofile.id">修改</a>
          <button class="deletebtn" v-on:click="deleteAdmin(user.id)">删除</button>
        </li>
      </ul>
      <p style="text-align:center;">
        <a id="new-btn" :href="'newadmin?id='+userprofile.id">
          <span class="glyphicon glyphicon-plus"></span> 添加新管理人员
        </a>
      </p>
    </div>
  </div>
</template>
<script src="./Admins.js"></script>
<style src="./Admins.css" scoped></style>
