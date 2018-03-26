var agentPriceMixin = {
  methods: {

    priceForAgent: function (theAgentLevel, thePrice) {
      switch (theAgentLevel) {
        case '0':
          return thePrice * 1
        case '1':
          return thePrice * 0.9
        case '2':
          return thePrice * 0.85
        case '3':
          return thePrice * 0.83
        default:
          return thePrice * 1
      }
    },
    toFixedTwo: function (price) {
      return Math.round(price * 10) / 10
    },
    priceForAgentJewelry: function (theAgentLevel, thePrice) {
      switch (theAgentLevel) {
        case '0':
          return thePrice * 1
        case '1':
          return thePrice * 0.9
        case '2':
          return thePrice * 0.85
        case '3':
          return thePrice * 0.83
        default:
          return thePrice * 1
      }
    }
  }
}
export default agentPriceMixin
