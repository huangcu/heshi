export default {
  methods: {
    DollarToEuro: function (thePrice, rateEUR) {
      return (rateEUR * thePrice).toFixed(2)
    },
    DollarToYuan: function (thePrice, rateCNY) {
      return (rateCNY * thePrice).toFixed(2)
    },
    EuroToYuan: function (thePrice, rateEURToCNY) {
      return (thePrice * rateEURToCNY).toFixed(2)
    },
    YuanToEuro: function (thePrice, rateCNYToEUR) {
      return (thePrice * rateCNYToEUR).toFixed(2)
    },
    EuroToDollar: function (thePrice, rateEUR) {
      return (thePrice / rateEUR).toFixed(2)
    },
    YuanToDollar: function (thePrice, rateCNY) {
      return (thePrice / rateCNY).toFixed(2)
    }
  }
}
