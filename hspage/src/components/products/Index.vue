<template>
  <div>
    <app-header :currentPage="currentPage" :userprofile='userprofile'></app-header>
    <router-view :rate='rate' @getCurrentPage='getCurrentPage' @updateAccount='updateAccount'/>
    <app-footer></app-footer>
  </div>
</template>
<script src='./Index.js'></script>
