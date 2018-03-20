<template>
<div>
  <app-header :currentPage="currentPage"></app-header>
  <div id="app">
    <router-view @getCurrentPage='getCurrentPage'/>
  </div>
  <app-footer></app-footer>

</div>
</template>

<script>
import Header from '@/components/header/Header.vue'
import Footer from '@/components/footer/Footer.vue'

export default {
  name: 'App',
  data: function () {
    return {
      currentPage: ''
    }
  },
  components: {
    'app-header': Header,
    'app-footer': Footer
  },
  methods: {
    getCurrentPage: function (arg) {
      this.currentPage = arg
    }
  }
}
</script>

<style src='./App.css'></style>
