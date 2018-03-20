export default {
  priceforaccount: function (theAgentLevel, thePrice) {
    switch (theAgentLevel) {
      case '0':
        return thePrice * 0.99
      case '1':
        return thePrice * 0.98
      case '2':
        return thePrice * 0.95
      case '3':
        return thePrice * 0.93
      case '4':
        return thePrice * 0.90
      case '5':
        return thePrice * 0.85
      case '6':
        return thePrice * 0.90
      default:
        return thePrice * 1
    }
    // _includes/functions/order_price_user.php'
  },
  priceforaccount_jewelry: function (theAgentLevel, thePrice) {
    switch (theAgentLevel) {
      case '0':
        return thePrice * 0.99
      case '1':
        return thePrice * 0.98
      case '2':
        return thePrice * 0.95
      case '3':
        return thePrice * 0.93
      case '4':
        return thePrice * 0.90
      case '5':
        return thePrice * 0.85
      case '6':
        return thePrice * 0.90
      default:
        return thePrice * 1
    }
  }
}
