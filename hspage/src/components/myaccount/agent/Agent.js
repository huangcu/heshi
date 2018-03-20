import agentPrice from '@util/agentprice.js'
export default {
  name: 'Agent',
  data: function () {
    return {
      agentLevel: 0,
      name: ''
    }
  },
  computed: {
    fromAgentLevel: function () {
      var str =''
      for (i=0; i< this.agentLevel; i++) {
        str = str + '<span class="glyphicon glyphicon-certificate"></span>'
      }
      return str
    }
  },
  methods: {
    goToSection: function (sectionID) {
      $('div#generalinfo, div#clientorders, div#mypoints, div#myclients, div.history-recommenedusers.extendedclients, div#coupon, div.generalinfo.heshibi-box, div.accountinfobox, div#personalinfo, div#mywebsite').not('#' + sectionID).css('display', 'none')
      $('#' + sectionID).fadeIn('fast')
      $('button.inpage-navi').removeClass('active')
      $('button#inpage-navi_' + sectionID).addClass('active')
      if (sectionID == 'myclients') {
        $('div.history-recommenedusers.extendedclients').fadeIn('fast')
      }
      if (sectionID == 'coupon') {
        $('div.generalinfo.heshibi-box').fadeIn('fast')
      }
      location.hash = 'section-' + sectionID
    }
  }
}
