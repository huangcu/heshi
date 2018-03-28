import Vue from 'vue'
import Router from 'vue-router'
import Header from '@/components/header/Header.vue'
import titleComponent from '@/components/title.component.vue'
import pageNotFound from '@/components/page.not.found.vue'
// import RoutingGuard from './routerguard.js'
Vue.use(Router)

// TODO global mixin - post to server to log user activity
Vue.mixin({
  methods: {
    capitalizeFirstLetter: str => str.charAt(0).toUpperCase() + str.slice(1)
  }
})

Vue.component('app-header', Header)
Vue.component('vue-title', titleComponent)

const router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      components: {
        default: () => import('@/components/products/Index.vue')
      },
      redirect: '/home',
      children: [
        {
          path: 'home',
          name: 'Home',
          components: {
            default: () => import('@/components/products/home/Home.vue')
          },
          alias: ['index']
        },
        {
          path: 'product/diamonds',
          name: 'Diamonds',
          components: {
            default: () => import('@/components/products/diamond/diamonds/Diamonds.vue')
          }
        },
        {
          path: 'product/diamond/:id',
          name: 'Diamond',
          components: {
            default: () => import('@/components/products/diamond/diamond/Diamond.vue')
          }
        },
        {
          path: 'product/ringfordiamond/:id',
          name: 'RingForDiamond',
          components: {
            default: () => import('@/components/products/diamond/ringfordiamond/RingForDiamond.vue')
          }
        },
        {
          path: 'product/diamondoftheweek',
          name: 'DiamondOfTheWeek',
          components: {
            default: () => import('@/components/products/diamond/diamondoftheweek/DiamondOfTheWeek.vue')
          }
        },
        {
          path: 'product/recommenddiamonds',
          name: 'RecommendDiamonds',
          components: {
            default: () => import('@/components/products/diamond/recommenddiamonds/RecommendDiamonds.vue')
          }
        },
        {
          path: 'product/jewelrys',
          name: 'Jewelrys',
          components: {
            default: () => import('@/components/products/jewelry/jewelrys/Jewelrys.vue')
          }
        }
      ]
    },
    {
      path: '/users',
      name: 'users',
      components: {
        default: () => import('@/components/users/Users.vue')
      },
      redirect: '/users/login',
      children: [
        {
          path: 'login',
          name: 'Login',
          components: {
            default: () => import('@/components/users/login/Login.vue')
          }
        },
        {
          path: 'loginbyemail',
          name: 'LoginByEmail',
          components: {
            default: () => import('@/components/users/loginbyemail/LoginByEmail.vue')
          }
        },
        {
          path: 'qrsign/:wechatopenID',
          name: 'QRSign',
          components: {
            default: () => import('@/components/users/qrsign/QRSign.vue')
          }
        },
        {
          path: 'register',
          name: 'Register',
          components: {
            default: () => import('@/components/users/register/Register.vue')
          }
        },
        {
          path: 'myaccount',
          name: 'MyAccount',
          components: {
            default: () => import('@/components/users/myaccount/MyAccount.vue')
          },
          meta: {requiresAuth: true, role: ['CUSTOMER', 'AGENT']},
          props: (route) => ({
            _account: route.query._account
          })
        }
      ]
    },
    {
      path: '/manage',
      name: 'Manage',
      components: {
        default: () => import('@/components/admin/Admin.vue')
      },
      redirect: '/manage/login',
      meta: {requiresAuth: true, role: ['ADMIN']},
      children: [
        {
          path: 'login',
          name: 'adminLogin',
          components: {
            default: () => import('@/components/admin/login/Login.vue')
          }
        },
        {
          path: 'admins',
          name: 'Admins',
          components: {
            default: () => import('@/components/admin/admins/Admins.vue')
          }
        },
        {
          path: 'newadmin',
          name: 'newAdmin',
          components: {
            default: () => import('@/components/admin/newadmin/NewAdmin.vue')
          }
        }
      ]
    },
    {
      path: '*',
      component: pageNotFound
    }
  ]
})

// router.beforeResolve(RoutingGuard)

export default router
