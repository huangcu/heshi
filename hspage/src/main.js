// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
// import BootstrapVue from 'bootstrap-vue'
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import VueResource from 'vue-resource'
import VueLocalStorage from 'vue-localstorage'
import VueFilter from 'vue-filter'
import VueCookies from 'vue-cookies'
import vueHeadful from 'vue-headful'
import VeeValidate from 'vee-validate'
// import 'bootstrap/dist/css/bootstrap.css'
// import 'bootstrap-vue/dist/bootstrap-vue.css'
require('bootstrap/dist/css/bootstrap.min.css')
Vue.use(VueCookies)
Vue.use(VueResource)
Vue.use(VueLocalStorage)
Vue.use(VueFilter)
Vue.use(VeeValidate)
// Vue.use(BootstrapVue)
Vue.config.productionTip = false
Vue.prototype.$currentPage = 'Index'
Vue.prototype.$userURL = 'http://localhost:8080/api'
Vue.prototype.$wechatURL = 'http://localhost:8080/api/wechat'
Vue.prototype.$adminURL = 'http://localhost:8080/api/admin'

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router: router,
  components: { App },
  template: '<App/>'
})

Vue.component('vue-headful', vueHeadful)
