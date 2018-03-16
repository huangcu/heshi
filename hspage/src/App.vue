<template>
<div>
  <div id="app">
    <router-view/>
  </div>
  <app-footer></app-footer>

</div>
</template>

<script>
import Footer from '@/components/footer/Footer.vue'

export default {
  name: 'App',
  components: {
    'app-footer': Footer
  }
}
</script>

<style src='./App.css'></style>
