<template>
<div>
  <app-header :currentPage="currentPage" :account='account'></app-header>
  <div id="app">
    <router-view :rate='rate' @getCurrentPage='getCurrentPage' @updateAccount='updateAccount'/>
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
      currentPage: '',
      account: '',
      rate: {
        USD: 1,
        CNY: 6.332200,
        EUR: 0.811434,
        CAD: 1.309410,
        AUD: 1.298319,
        CHF: 0.949890,
        RUB: 57.917500,
        NZD: 1.383853
      }
    }
  },
  components: {
    'app-header': Header,
    'app-footer': Footer
  },
  methods: {
    getCurrentPage: function (arg) {
      this.currentPage = arg
    },
    updateAccount: function (arg) {
      this.account = arg
    },
    getCurrencyRate: function () {
      this.$http.get(
        this.$userURL + '/exchangerate'
      ).then(response => {
        var exchangeRate = JSON.parse(response.bodyText)
        if (exchangeRate.code === undefined) {
          this.rate = exchangeRate.rates
        } else {
          console.log(exchangeRate.message)
        }
      }, err => { console.log(err); alert('error:' + err.bodyText) })
    }
  },
  created () {
    this.getCurrencyRate()

    setInterval(function () {
      // poll server every 2 hrs
      this.getCurrencyRate()
    }.bind(this), 2 * 60 * 60 * 1000)
  }
}
</script>

<style src='./App.css'></style>
