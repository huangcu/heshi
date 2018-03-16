import Vue from 'vue'
import Router from 'vue-router'
import Header from '@/components/header/Header.vue'
import titleComponent from '@/components/title.component.vue'
import diamondFilter from '@/components/products/diamond/diamondfilter/DiamondFilter.vue'
Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Home',
      components: {
        default: () => import('@/components/home/Home.vue')
      },
      alias: ['/index', '/home']
    },
    {
      path: '/login',
      name: 'Login',
      components: {
        default: () => import('@/components/users/login/Login.vue')
      },
      props: (route) => ({
        referer: route.query.referer,
        for: route.query.for
      })
    },
    {
      path: '/loginbyemail',
      name: 'LoginByEmail',
      components: {
        default: () => import('@/components/users/loginbyemail/LoginByEmail.vue')
      },
      props: (route) => ({
        _ref: route.query._ref,
        _for: route.query._for
      })
    },
    {
      path: '/register',
      name: 'Register',
      components: {
        default: () => import('@/components/users/register/Register.vue')
      }
    },
    {
      path: '/myaccount',
      name: 'MyAccount',
      components: {
        default: () => import('@/components/myaccount/MyAccount.vue')
      },
      props: (route) => ({
        _account: route.query._account
      })
    },
    {
      path: '/product/diamonds',
      name: 'Diamonds',
      components: {
        default: () => import('@/components/products/diamond/diamonds/Diamonds.vue')
      }
    },
    {
      path: '/product/diamond',
      name: 'Diamond',
      components: {
        default: () => import('@/components/products/diamond/diamond/Diamond.vue')
      }
    },
    {
      path: '/product/jewelry',
      name: 'Jewelry',
      components: {
        default: () => import('@/components/products/jewelry/jewelrys/Jewelrys.vue')
      }
    }
  ]
})

Vue.component('app-header', Header)
Vue.component('vue-title', titleComponent)
Vue.component('diamond-filter', diamondFilter)
