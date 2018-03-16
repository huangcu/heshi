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
Vue.http.interceptors.push((request, next) => {
  request.headers.set('X-Auth-Token', 'Jbm6XfXQj/KqmMTqz6c4GQWl9U6JMLQ/T4LzPWIEi2W2Q23GDkuIfxvbUC/rar8ZJIWWSVo68fZ/hv6n0oAeXaQKEfhKmGUZ8m8JHm5TteBZwqZuqXAbOeowTJVBn8aaUhfSfZbmgNnXwDEnhjZ1DZ8jG2Khy9uzoHu5ogwbVHQ=')
  next()
})
// Vue.use(BootstrapVue)
Vue.config.productionTip = false
Vue.prototype.$currentPage = 'Index'
Vue.prototype.$userURL = 'https://localhost:8443/api'
Vue.prototype.$wechatURL = 'https://localhost:8443/api/wechat'
Vue.prototype.$adminURL = 'https://localhost:8443/api/admin'

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router: router,
  components: { App },
  template: '<App/>'
})

Vue.component('vue-headful', vueHeadful)
