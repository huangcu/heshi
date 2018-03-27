let Images = require.context('@/_images/constant/', false, /\.(png|jpg)$/);
let diamondImage = require.context('@/.image/diamond/', false, /\.(png|jpg)$/);
let diamondImageThumbs = require.context('@/.image/diamond/thumbs/', false, /\.(png|jpg)$/);

import diamondMixin from '../diamondcommon.js'
import agentPriceMixin from '@/util/agentprice.js'
import accountPriceMixin from '@/util/accountprice.js'
import currentCaculatorMixin from '@/util/currency_caculator.js'

export default {
  name: 'DiamondsData',
  props: {
    rate: {
      type: Object,
      required: true
    },
    vat_choice: String,
    diamonds: Array
  },
  mixins: [diamondMixin, accountPriceMixin, agentPriceMixin, currentCaculatorMixin],
  data: function () {
    return {
      userprofile: '',
      the_order_btn: '<span class="glyphicon glyphicon-heart"></span> 收藏',
      the_order_btn_expl: '<span class="btn-expl-words interested"><span class="glyphicon glyphicon-hand-left"></span> 先比较，再决定</span>',
      the_order_confirm_btn: '<span class="glyphicon glyphicon-check"></span> 预定',
      the_order_confirm_btn_expl: '<span class="btn-expl-words ordered"><span class="glyphicon glyphicon-hand-left"></span> 决定购买了请点这里</span>',
      checklisttxt: '查看已收藏钻石 &raquo;',
      checklist_ordered_txt: '查看已预定钻石 &raquo;',
      inthelistword_saved: '您已收藏该钻',
      inthelistword_ordered: '您已预定该钻',
      vat_status_txt: '(不含税)',
      if_agent: false,
      dia_items_shoppinglist: [],
      dia_items_shoppinglist_confirmed: [],
      pic_where: 0,
      pic_name: 'ico-stones.png',
      shape_TXT: '圆形',
      color_TXT:''
    }
  },
  computed: {
    hrdGradingLabImage: function () {
      return Images('./hrd-label.jpg')
    },
    giaGradingLabImage: function () {
      return Images('./gia-label.jpg')
    },
    igiGradingLabImage: function () {
      return Images('./igi-label.jpg')
      
    },
    getImage: function () {
      // return Images('./ico-stones.png')
      return Images('./' + this.pic_name)
    },
    
    imagePosition: function () {
      return  this.pic_where
    },
    in_dia_items_shoppinglist_confirmed: function () {
      // TODO sudo function
      return false
    },
    in_dia_items_shoppinglist: function () {
      // TODO sudo function
      return false
    }
  }, 
  methods: {
    getDiamondImageThumb: function (theName) {
      return diamondImageThumbs("./"+theName)
    },
    getDiamondImage: function (theName) {
      return diamondImage("./"+theName)
    },
    getData: function () {
      if (userprofile!=='') {
        //get from cookie, if it is agent 
        if (userprofile.user_type === 'AGENT') {
          this.the_order_confirm_btn = '<span class="glyphicon glyphicon-check"></span> 预定'
          this.the_order_confirm_btn_expl = '<span class="btn-expl-words ordered"><span class="glyphicon glyphicon-hand-left"></span> 决定购买了请点这里</span>'
        }
      }
      //if vat_choice is selected 
      if (vat_choice!==undefined) {
        if (vat_choice !== 'YES') {
          this.vat_status_txt = '(含税)'
        }
      }
    },
    getInterestedItems: function () {
      if (this.userprofile === undefined) {
        return
      }
      if (this.userprofile === '') {
        if (this.$cookies.isKey('interestedItems')) {
          this.interestedItems = this.$cookies.get('interestedItems')
        }
      } else {
        // TODO post
        this.$http.post(this.$userURL + '/interestedItems/' + this.userprofile.id).then(response => {
          return response.bodyText
        }, err => { console.log(err); alert('error:' + err.bodyText) })
      }
    },
    showOrderSuccessFeedback: function () {
      $('div#feedbackcover').fadeIn(500)
    },
    makeorder: function (theRef) {
      // TODO
      var itemtype = 'DIAMOND'
      $('div#loadingcover').fadeIn('fast')
      $.post(
        '/_includes/functions/addtoshoppinglist.php',
        { id: theRef, item_type: itemtype },
        function (data) {
          console.log('order feedback: ' + data)
          if ($.trim(data) == 'OK') {
            //alert('ordered')
            $('button.btnfororder.interested[title="' + theRef + '"]').html('<span class="glyphicon glyphicon-ok"></span> 已收藏').addClass('itemordered').attr('disabled', 'disabled')
            $('div.details_' + theRef + ' span.btn-expl-words.interested').fadeOut('fast')
            $('a#addringfordiamondbtn_' + theRef).fadeIn('fast')
            $('a#checkmyorderlistbtn_' + theRef).fadeIn('fast')
            var crr_list_num = parseInt($('#num-shoppinglist').html())
            var new_list_num = crr_list_num + 1
            $('#num-shoppinglist, #num-shoppinglist-mobile').html(new_list_num)
            $('#num-shoppinglist-container, #num-shoppinglist-container-mobile').removeAttr('style')
          } else {
            alert('Server is busy, please try later!')
          }
          $('div#loadingcover').fadeOut('fast')
        }
      )
    },
    makeorder_confirmed: function (theRef) {
      // TODO
      var itemtype = 'DIAMOND'
      $('div#loadingcover').fadeIn('fast')
      $.post(
        '/_includes/functions/addtoshoppinglist.php',
        { id: theRef, item_type: itemtype, confirm_for_check: 'YES' },
        function (data) {
          console.log('order feedback: ' + data)
          if ($.trim(data) == 'OK') {
            //alert('ordered')
            $('button.btnfororder.ordered[title="' + theRef + '"]').html('<span class="glyphicon glyphicon-ok"></span> 已预定').addClass('itemordered').attr('disabled', 'disabled')
            $('div.details_' + theRef + ' span.btn-expl-words.ordered').fadeOut('fast')
            $('p#interested-btn-box_' + theRef).fadeOut('fast')
            $('a#addringfordiamondbtn_confirmed_' + theRef).fadeIn('fast')
            $('a#checkmyorderlistbtn_confirmed_' + theRef).fadeIn('fast')
            var crr_list_num = parseInt($('#num-shoppinglist-c').html())
            var new_list_num = crr_list_num + 1
            $('#num-shoppinglist-c, #num-shoppinglist-c-mobile').html(new_list_num)
            $('#num-shoppinglist-c-container, #num-shoppinglist-c-container-mobile').removeAttr('style')

            crr_ordered_diamond = theRef
            showOrderSuccessFeedback()
          } else if ($.trim(data) == 'NEED-VERIFICATION') {
            window.location.href = "/myaccount.php?r=order"
          } else {
            alert('Server is busy, please try later!')
          }
          $('div#loadingcover').fadeOut('fast')
        }
      )
    },
    openDiaDetail: function (theRef) {
      // detailOpened++
      $('div.details').not(".details_" + theRef).slideUp()
      $('div.details_' + theRef).slideDown()
      //dia-piece-box
      $('div.dia-piece-box').removeClass('detailisopen')
      $('div.details_' + theRef).parent('div.dia-piece-box').addClass('detailisopen')
      var crr_browser_width = $(window).width()
      if (crr_browser_width > 680) {
        var thedevice = 'PC'
      } else {
        var thedevice = 'MOBILE'
      }
      // $.post(
        // TODO post back -record use viewed this, and interested in this diamonds
        // '/_includes/functions/userusingrecord.php',
        // { ref: theRef, device: thedevice }
      // )
      // TODO
      // if ($('div#thelettertoagents-container').length > 0) {
      //   $('body').addClass('no-overflow')
      //   $('div#thelettertoagents-container').fadeIn(1888)
      // }
    },
    closeDiaDetail: function (theRef) {
      $('div.details_' + theRef).slideUp()
      $('div.dia-piece-box').removeClass('detailisopen')
    }
  },
  mounted() {
    if (this.$cookies.isKey('userprofile')) {
      this.userprofile = this.$cookies.get('userprofile')
    } else {
      this.userprofile = ''
    }
    // this.getInterestedItems()
    // this.getOrders()
  }
}
