// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import VueResource from 'vue-resource'
import VueLocalStorage from 'vue-localstorage'
import VueFilter from 'vue-filter'
import VueCookies from 'vue-cookies'
import vueHeadful from 'vue-headful'
import VeeValidate from 'vee-validate'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
// import 'jquery-ui-css/jquery-ui.min.css'
import 'jquery-ui/themes/base/all.css'
// import 'jquery-ui/themes/base/slider.css'

Vue.use(VueCookies)
Vue.use(VueResource)
Vue.use(VueLocalStorage)
Vue.use(VueFilter)
Vue.use(VeeValidate)
Vue.use(BootstrapVue)
Vue.http.interceptors.push((request, next) => {
  if (VueCookies.isKey('token')) {
    request.headers.set('Authorization', 'Bearer ' + VueCookies.get('token'))
  }
  request.headers.set('X-Auth-Token', 'Jbm6XfXQj/KqmMTqz6c4GQWl9U6JMLQ/T4LzPWIEi2W2Q23GDkuIfxvbUC/rar8ZJIWWSVo68fZ/hv6n0oAeXaQKEfhKmGUZ8m8JHm5TteBZwqZuqXAbOeowTJVBn8aaUhfSfZbmgNnXwDEnhjZ1DZ8jG2Khy9uzoHu5ogwbVHQ=')
  request.credentials = true
  next()
})
Vue.config.productionTip = true
Vue.prototype.$userURL = (Vue.config.productionTip) ? 'https://localhost:8443/api' : 'http://localhost:8080/api'
Vue.prototype.$customerURL = (Vue.config.productionTip) ? 'https://localhost:8443/api/customer' : 'http://localhost:8080/api/customer'
Vue.prototype.$wechatURL = (Vue.config.productionTip) ? 'https://localhost:8443/api/wechat' : 'http://localhost:8080/api/wechat'
Vue.prototype.$adminURL = (Vue.config.productionTip) ? 'https://localhost:8443/api/admin' : 'http://localhost:8080/api/admin'

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})

Vue.component('vue-headful', vueHeadful)
// TODO bootstrap glyphicon css
// TODO favicon.ico not working
