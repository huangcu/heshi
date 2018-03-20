import Vue from 'vue'
import Router from 'vue-router'
import Header from '@/components/header/Header.vue'
import titleComponent from '@/components/title.component.vue'
import diamondFilter from '@/components/products/diamond/diamondfilter/DiamondFilter.vue'
import diamondsData from '@/components/products/diamond/diamondsdata/DiamondsData.vue'
import pageNotFound from '@/components/page.not.found.vue'
import currencyCaculator from '@/util/currency_caculator.js'
import agentPrice from '@/util/agentprice.js'
import accountPrice from '@/util/accountprice.js'
Vue.use(Router)

// TODO global mixin - post to server to log user activity
Vue.mixin({
  methods: {
    capitalizeFirstLetter: str => str.charAt(0).toUpperCase() + str.slice(1)
  }
})

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
        appointment: route.query.appointment
      })
    },
    {
      path: '/loginbyemail',
      name: 'LoginByEmail',
      components: {
        default: () => import('@/components/users/loginbyemail/LoginByEmail.vue')
      }
    },
    {
      path: '/qrsign',
      name: 'QRSign',
      components: {
        default: () => import('@/components/users/qrsign/QRSign.vue')
      }
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
      mixins: [currencyCaculator, agentPrice, accountPrice],
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
      path: '/product/ringfordiamond',
      name: 'RingForDiamond',
      components: {
        default: () => import('@/components/products/diamond/ringfordiamond/RingForDiamond.vue')
      }
    },
    {
      path: '/product/diamondoftheweek',
      name: 'DiamondOfTheWeek',
      components: {
        default: () => import('@/components/products/diamond/diamondsoftheweek/DiamondsOfTheWeek.vue')
      }
    },
    {
      path: '/product/recommenddiamonds',
      name: 'RecommendDiamonds',
      components: {
        default: () => import('@/components/products/diamond/recommenddiamonds/RecommendDiamonds.vue')
      }
    },
    {
      path: '/product/jewelrys',
      name: 'Jewelry',
      components: {
        default: () => import('@/components/products/jewelry/jewelrys/Jewelrys.vue')
      }
    },
    {
      path: '*',
      component: pageNotFound
    }
  ]
})

Vue.component('app-header', Header)
Vue.component('vue-title', titleComponent)
Vue.component('diamond-filter', diamondFilter)
Vue.component('diamonds-data', diamondsData)
