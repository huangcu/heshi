<template>
<div>
  <app-header></app-header>
  <div id="app">
    <router-view/>
  </div>
  <app-footer></app-footer>

</div>
</template>

<script>
import Header from '@/components/header/Header.vue'
import Footer from '@/components/footer/Footer.vue'

export default {
  name: 'App',
  components: {
    'app-header': Header,
    'app-footer': Footer
  }
}
</script>

<style src='./App.css'></style>
