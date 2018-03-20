let Images = require.context('@/_images/constant/', false, /\.(png|jpg)$/);
export default {
  name: 'Diamonds',
  data: function () {
    return {
      diamonds: [],
      featured: 'NO',
      shape: [],
      color: [],
      fancyColor: false,
      clarity: [],
      cut: [],
      polish: [],
      sym: [],
      certi: [],
      fluo: [],
      place: [],
      weight_from: '',
      weight_to: '',
      price_from: '',
      price_to: '',
      sorting: '',
      sorting_weight_direction: 'ASC',
      sorting_color_direction: 'ASC',
      sorting_clarity_direction: 'ASC',
      sorting_cut_direction: 'ASC',
      sorting_price_direction: 'ASC',
      sorting_direction: 'ASC',
      vat_include: 'NO',
      crr_page: 1
    }
  },
  computed: {
    getImage: function () {
      return Images('./ico-stones.png')
    },
    thelettertoagents: function () {
      // TODO is agent and accept order online
      return true
    }
  },
  methods: {
    update: function () {
      $('#loadingtxt').html('载入中...')
      $('div#loadingcover').fadeIn('fast')
      t = setTimeout('stopLoadingContent()', 7899)
      updatingstatus = 'updating'
      var formData = new FormData()
      formData.append('shape', hash_shape)
      formData.append('color', hash_color)
      formData.append('clarity', hash_clarity)
      formData.append('cut ', hash_cut)
      formData.append('sym', hash_sym)
      formData.append('fluo', hash_fluo)
      formData.append('certi', hash_certi)
      formData.append('weight_from', hash_weight_from)
      formData.append('weight_to', hash_weight_to)
      formData.append('price_from', hash_price_from)
      formData.append('price_to', hash_price_to)
      formData.append('sorting', hash_sorting)
      formData.append('sorting_direction', hash_sorting_direction)
      formData.append('crr_page', hash_crr_page)
      formData.append('vat_choice', hash_vat_choice)
      formData.append('place', hash_place)

      this.$http.post(
        this.$userURL + '/products/filter/diamonds',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        }
      ).then(response => {
        if (response.status === 200) {
          // token
          updatingstatus = 'updated'
          clearTimeout(t)
          if (response.body == 'NO-DIAMOND-FOUND') {
            $('div#diamondlist').html('<p id="nodiafoundfeedback">抱歉，没有找到符合您所选条件的钻石</p>')
            $('span#diapagenavi').html(' - ')
          } else {
            diamonds = response.body
            $('div#diamondlist').html(data)
            var howmanyrecords = $("div#howmanyrecords").html()
            $('#loadingtxt').html('找到<strong>' + howmanyrecords + '</strong>颗')
            diamondlistpagenavi(howmanyrecords)

          }
          $('div#loadingcover').delay(880).fadeOut('fast')
        }
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    stopLoadingContent: function () {
      if (updatingstatus == 'updating') {
        $('div#diamondlist').html('<p id="nodiafoundfeedback">抱歉，系统繁忙，请稍后重试</p>')
        $('span#diapagenavi').html(' - ')
        xhr.abort()
        $('div#loadingcover').fadeOut('fast')
      }
    },
    filter_shape: function (theshape) {
      this.shape = []
      let vm = this
      $('li#filter_shape' + theshape).toggleClass('btn-active')
      $('.fileber_shape_outer>li.btn-active').each(function () {
        var crr_obj = $(this)
        vm.shape.push(crr_obj.attr('title'))
      });
      console.log(this.shape)
    },
    filter_color: function (thecolor) {
      this.color = []
      let vm = this
      if (thecolor !== 'Fancy') {
        vm.fancyColor = false
        $('li#filter_color' + thecolor).toggleClass('btn-active')
        $('li#filter_colorFancy').removeClass('btn-active')
        $('.filter_color_outer>li.btn-active').each(function () {
          var crr_obj = $(this)
          var crrcolor = crr_obj.attr('title')
          if (crrcolor == 'N-Z') {
            vm.color.push("N")
            vm.color.push("O")
            vm.color.push("P")
            vm.color.push("Q")
            vm.color.push("R")
            vm.color.push("S")
            vm.color.push("T")
            vm.color.push("U")
            vm.color.push("V")
            vm.color.push("W")
            vm.color.push("X")
            vm.color.push("Y")
            vm.color.push("Z")
          } else {
            vm.color.push(crr_obj.attr('title'))
          }
        })
      } else {
        vm.fancyColor = true
        $('li.filter_color').not('li[title="Fancy"]').removeClass('btn-active')
        $('li#filter_colorFancy').toggleClass('btn-active')
        if ($('li#filter_colorFancy').hasClass('btn-active')) {
          vm.color.push(' color LIKE "F%" AND color <> "F" ')
        }
      }
      console.log(this.color)
    },
    filter_clarity: function (theclarity) {
      this.clarity = []
      let vm = this
      $('li#filter_clarity' + theclarity).toggleClass('btn-active')
      $('.filter_clarity_outer>li.btn-active').each(function () {
        vm.clarity.push(crr_obj.attr('title'))
      })
      console.log(this.clarity);
    },
    filter_cut: function (thegrade) {
      this.cut = []
      let vm = this
      $('li#filter_cut' + thegrade).toggleClass('btn-active')
      $('.filter_cut_outer>li.btn-active').each(function () {
        var crr_obj = $(this)
        vm.cut.push(crr_obj.attr('title'))
      })
      console.log(this.cut)
    },
    filter_polish: function (thegrade) {
      this.polish = []
      let vm = this
      $('li#filter_polish' + thegrade).toggleClass('btn-active')
      $('.filter_polish_outer>li.btn-active').each(function () {
        var crr_obj = $(this)
        vm.polish.push(crr_obj.attr('title'))
      })
      console.log(this.polish)
    },
    filter_sym: function (thegrade) {
      this.sym = []
      let vm = this
      $('li#filter_sym' + thegrade).toggleClass('btn-active')
      $('.filter_sym_outer>li.btn-active').each(function () {
        var crr_obj = $(this)
        vm.sym.push(crr_obj.attr('title'))
      })
      console.log(this.sym);
    },
    filter_certi: function (thelab) {
      this.certi = []
      let vm = this
      $('li#filter_certi' + thelab).toggleClass('btn-active')
      $('.filter_certi_outer>li.btn-active').each(function () {
        var crr_obj = $(this)
        vm.certi.push(crr_obj.attr('title'))
      })
      console.log(this.certi);
    },
    filter_fluo: function (thegrade) {
      this.fluo = []
      let vm = this
      $('li#filter_fluo' + thegrade).toggleClass('btn-active')
      $('.filter_fluo_outer>li.btn-active').each(function () {
        var crr_obj = $(this);
        var fluovalue = crr_obj.attr('title')
        if (fluovalue == 'VST') {
          vm.fluo.push("VST")
          vm.fluo.push("Very Strong")
        } else if (fluovalue == 'STG') {
          vm.fluo.push("STG")
          vm.fluo.push("Strong")
        } else if (fluovalue == 'MED') {
          vm.fluo.push("MED")
          vm.fluo.push("Medium")
        } else if (fluovalue == 'FNT') {
          vm.fluo.push("FNT")
          vm.fluo.push("SLT")
          vm.fluo.push("VSL")
          vm.fluo.push("Faint")
          vm.fluo.push("Very Slight")
          vm.fluo.push("Slight")
        } else if (fluovalue == 'NON') {
          vm.fluo.push("NON")
          vm.fluo.push("None")
        }
      })
      console.log(this.fluo);
    },
    filter_place: function () {
      this.place = []
      let vm = this
      var chosenPlace = $('select#placechooser').val();
      if (chosenPlace == 'CHINA') {
        this.place.push(' country LIKE "China" OR country LIKE "SZ"  OR country LIKE "HK" OR country LIKE "HSTHK" OR country LIKE "Shenzhen" OR country LIKE "Hongkong" OR country LIKE "Hong kong" ')
      } else if (chosenPlace == 'CZ') {
        this.place.push(' country LIKE "SZ" OR country LIKE "Shenzhen" ')
      } else if (chosenPlace == 'HK') {
        this.place.push(' country LIKE "HK" OR country LIKE "HSTHK" OR country LIKE "Hongkong" OR country LIKE "Hong kong" ')
      } else if (chosenPlace == 'BELGIUM') {
        this.place.push(' country LIKE "Belg%" OR country LIKE "Antwerp%" OR country = "" OR country IS NULL ')
      } else if (chosenPlace == 'INDIA') {
        this.place.push(' country LIKE "IND%" ')
      } else if (chosenPlace == 'OTHER') {
        this.place.push(' country NOT LIKE "China" AND country NOT LIKE "SZ" AND country NOT LIKE "HK" AND country NOT LIKE "HSTHK" AND country NOT LIKE "Belg%" AND country NOT LIKE "IND%"  AND country <> ""  AND country IS NOT NULL ')
      } else if (chosenPlace == 'ALL') {
        this.place = [];
      }
      console.log(this.place);
    },
    filter_weight: function () {
      this.weight_from = $('#weight_from').val()
      this.weight_to = $('#weight_to').val()
      $("#slider-range-weight").slider("values", 0, this.weight_from)
      $("#slider-range-weight").slider("values", 1, this.weight_to)

      $('#tooltip-weight').stop(true).delay(528).fadeOut('slow', function () {
        $('#tooltip-weight').removeAttr('style')
      })
    },
    filter_price: function () {
      this.price_from = $('#price_from').val()
      this.price_to = $('#price_to').val()
      if (this.price_from < 9800) {
        var crr_v_slider_l = Math.round((parseFloat(this.price_from) + 400) / 50)
      } else {
        var crr_v_slider_l = Math.round(Math.sqrt((parseFloat(this.price_from) - 9050) / 68) + 188)
      }
      $("#slider-range").slider("values", 0, crr_v_slider_l)
      if (this.price_to < 9800) {
        var crr_v_slider_r = Math.round((parseFloat(this.price_to) + 400) / 50)
      } else {
        var crr_v_slider_r = Math.round(Math.sqrt((parseFloat(this.price_to) - 9050) / 68) + 188)
      }
      $("#slider-range").slider("values", 1, crr_v_slider_r)

      $('#tooltip-price').stop(true).delay(528).fadeOut('slow', function () {
        $('#tooltip-price').removeAttr('style')
      })
    },
    sortingby: function (ref) {
      $('li.sortingbtn').removeClass('active')
      $('li#sortingby_' + ref).addClass('active').toggleClass('DESC')
      this.sorting = ref
      this.sorting_direction = $('li#sortingby_' + ref).hasClass('DESC') ? "DESC" : "ASC"
    },
    vatChoice: function (ref) {
      this.vatChoice = $('input#vat-choice').prop("checked") ? "YES" : "NO"
    },
    diamondlistpagenavi: function (howmanyrecords) {
      $('span#found-diamonds-num').html(howmanyrecords)
      //$('#dia-num-bg').stop(true).remove()
      //$('span#found-diamonds-num-container').append('<span id="dia-num-bg"></span>')
      $('#dia-num-bg').stop(true).css({ 'opacity': 1, 'display': 'inline-block' }).delay(58).fadeOut(2000)
      $('span#diapagenavi').empty()
      var totalrecords = parseFloat(howmanyrecords)
      var totalpages = Math.ceil(totalrecords / 28)
      if (totalpages < $crr_page) {
        $crr_page = 1
      }
      if (totalpages <= 13) {
        var loop_beginpage = 1
        var loop_endpage = totalpages
      } else {
        if (($crr_page + 10) < totalpages) {
          var loop_beginpage = $crr_page
          var loop_endpage = $crr_page + 10
        } else {
          var loop_beginpage = totalpages - 10
          var loop_endpage = totalpages
        }
      }
      if (loop_beginpage > 1) {
        $('span#diapagenavi').append('<span class="dia-page-btn" onclick="choosethispage(1)">1</span>')
      }
      if (loop_beginpage > 2) {
        $('span#diapagenavi').append('<span class="somethinginbetween-indi"> ... </span>')
      }

      for (var i = loop_beginpage; i <= loop_endpage; i++) {
        if (i == $crr_page) {
          $('span#diapagenavi').append('<span class="dia-page-btn" id="crr_page">' + i + '</span>');
        } else {
          $('span#diapagenavi').append('<span class="dia-page-btn" onclick="choosethispage(' + i + ')">' + i + '</span>');
        }
      }

      if ((totalpages - loop_endpage) > 1) {
        $('span#diapagenavi').append('<span class="somethinginbetween-indi"> ... </span>');
      }
      if ((totalpages - loop_endpage) >= 1) {
        $('span#diapagenavi').append('<span class="dia-page-btn" onclick="choosethispage(' + totalpages + ')">' + totalpages + '</span>');
      }
    },
    searchbyRef: function (theRef) {
      var theRef = $('input#stockreftosearch').val()
      if (theRef == '') {
        alert('请输入要查询的商品编号！')
        return
      }
      var formData = new FormData()
      formData.append('ref', theRef)
      this.$http.post(
        this.$userURL + '/products/search/diamonds',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        }
      ).then(response => {
        if (response.status === 200) {
          // token
          updatingstatus = 'updated'
          clearTimeout(t)
          if (response.body == 'NO-DIAMOND-FOUND') {
            $('div#diamondlist').html('<p id="nodiafoundfeedback">抱歉，没有找到符合您所选条件的钻石</p>')
            $('span#diapagenavi').html(' - ')
          } else {
            diamonds = response.body
            $('div#diamondlist').html(data)
            var howmanyrecords = $("div#howmanyrecords").html()
            $('#loadingtxt').html('找到<strong>' + howmanyrecords + '</strong>颗')
            diamondlistpagenavi(howmanyrecords)

          }
          $('div#loadingcover').delay(880).fadeOut('fast')
        }
      }, err => { console.log(err); alert('error:' + err.body) })
    },
    choosethispage: function (page) {
      this.crr_page = page
    },
    mobileFunctions: function () {
      //alert($('#slider-range-weight span.ui-slider-handle').length);
      var sliderHandelInitialed = $('#slider-range-weight span.ui-slider-handle').length
      var sliderHandelInitialed_price = $('#slider-range span.ui-slider-handle').length
      if (sliderHandelInitialed == 0 || sliderHandelInitialed_price == 0) {
        t_mobileFunctions = setTimeout(mobileFunctions, 200)
        return
      }
      clearTimeout(t_mobileFunctions)

      $('#slider-range-weight span.ui-slider-handle').eq(0).addClass('weight_from_handel')
      $('#slider-range-weight span.ui-slider-handle').eq(1).addClass('weight_to_handel')
      $('#slider-range span.ui-slider-handle').eq(0).addClass('price_from_handel')
      $('#slider-range span.ui-slider-handle').eq(1).addClass('price_to_handel')
      $('#slider-range-weight span.ui-slider-handle, #slider-range span.ui-slider-handle').off()
      $(document).on("vmousemove", '#slider-range-weight span.weight_from_handel', function (event) {
        var beginLeftNum = $('#slider-range-weight').offset().left
        var distanceToLeft = event.pageX - beginLeftNum
        $("#slider-range-weight").slider("values", 0, (distanceToLeft / 24))
        $('#weight_from').val(Math.round(distanceToLeft / 24 * 10) / 10)
        $('#tooltip-weight').html(Math.round(distanceToLeft / 24 * 10) / 10)
        $('#tooltip-weight').css({ 'left': (distanceToLeft + 15), 'display': 'block' })
        clearTimeout(t_update)
        t_update = setTimeout(this.filter_weight(), 288)
      })

      $(document).on("vmousemove", '#slider-range-weight span.weight_to_handel', function (event) {
        var beginLeftNum = $('#slider-range-weight').offset().left
        var distanceToLeft = event.pageX - beginLeftNum
        $("#slider-range-weight").slider("values", 1, (distanceToLeft / 24))
        $('#weight_to').val(Math.round(distanceToLeft / 24 * 10) / 10)
        $('#tooltip-weight').html(Math.round(distanceToLeft / 24 * 10) / 10)
        $('#tooltip-weight').css({ 'left': (distanceToLeft + 15), 'display': 'block' })
        clearTimeout(t_update)
        t_update = setTimeout(this.filter_weight(), 288)
      })

      $(document).on("vmousemove", '#slider-range span.price_from_handel', function (event) {
        var beginLeftNum = $('#slider-range').offset().left
        var distanceToLeft = event.pageX - beginLeftNum
        $("#slider-range").slider("values", 0, distanceToLeft)
        if (distanceToLeft < 188) {
          var pricevaluefrom = distanceToLeft * 50 - 400
        } else {
          var pricevaluefrom = (distanceToLeft - 188) * (distanceToLeft - 188) * 68 + 9450 - 400
        }
        $('#price_from').val(Math.round(pricevaluefrom))
        $('#tooltip-price').html(Math.round(pricevaluefrom) + '€')
        $('#tooltip-price').css({ 'left': (distanceToLeft - 3), 'display': 'block' })

        clearTimeout(t_update)
        t_update = setTimeout(this.filter_price(), 288)
      })

      $(document).on("vmousemove", '#slider-range span.price_to_handel', function (event) {
        var beginLeftNum = $('#slider-range').offset().left
        var distanceToLeft = event.pageX - beginLeftNum
        $("#slider-range").slider("values", 1, distanceToLeft)

        if (distanceToLeft < 188) {
          var pricevalueto = distanceToLeft * 50 - 400
        } else {
          var pricevalueto = (distanceToLeft - 188) * (distanceToLeft - 188) * 68 + 9450 - 400
        }

        $('#price_to').val(Math.round(pricevalueto))

        $('#tooltip-price').html(Math.round(pricevalueto) + '€')
        $('#tooltip-price').css({ 'left': (distanceToLeft - 3), 'display': 'block' })

        clearTimeout(t_update)
        t_update = setTimeout(this.filter_price(), 288)
      })
    },
    gotoMountingsfor: function () {
      var mountingURL = '/products/ringfordiamond?ref=' + crr_ordered_diamond + '&ordered=yes';
      this.$route.replace(mountingURL)
    },
    xtheSfdb: function () {
      $('div#feedbackcover').fadeOut(50);
    },
    agentLetterOk: function () {
      $.post(
        "/_content/ajax/agent-letter-read.php",
        { read: "YES" },
        function (data) {
          console.log(data);
        }
      )
      $('body').removeClass('no-overflow')
      $('div#thelettertoagents-container').fadeOut(588, function () {
        $('div#thelettertoagents-container').remove();
      })
    }
  },
  created() {
    var theHashSTR = location.hash.replace('#', '')
    var refinHash = "ref"
    if(theHashSTR.indexOf(refinHash) > -1) {
      var crr_freg_array = theHashSTR.split('=')
      var crr_freg_value = crr_freg_array[1]
      $('input#stockreftosearch').val(crr_freg_value)
      searchbyRef()
    } else {
      parseTheHash()
    }

    $('input#weight_from, input#weight_to, input#price_from, input#price_to').focus(function () {
      crr_active_input_field_id = $(this).attr('id')
      crr_active_input_field_value = $(this).val()
    })

    $('input#weight_from').blur(function () {
      var crr_value = $('input#weight_from').val()

      if (!$.isNumeric(crr_value)) {
        alert('您输入的数字格式不正确，请重试')
        return
      }

      if (crr_value != $weight_from) {
        //$( "#slider-range-weight" ).slider( "values", 0, crr_value )
        filter_weight()
      }  
    })

    $('input#weight_to').blur(function () {
      var crr_value = $('input#weight_to').val()
      if (!$.isNumeric(crr_value)) {
        alert('您输入的数字格式不正确，请重试')
        return
      }
      if (crr_value != $weight_to) {
        filter_weight()
      }
    })

    $('input#price_from').blur(function () {
      var crr_value = $('input#price_from').val()
      if (!$.isNumeric(crr_value)) {
        alert('您输入的数字格式不正确，请重试')
        return
      }
      if (crr_value != $price_from) {
        filter_price()
      }
    })

    $('input#price_to').blur(function () {
      var crr_value = $('input#price_to').val()
      if (!$.isNumeric(crr_value)) {
        alert('您输入的数字格式不正确，请重试')
        return
      }
      if (crr_value != $price_to) {
        filter_price()
      }
    })

    //enter key event when input the stock ref number
    $('#stockreftosearch').keypress(function (e) {
      var key = e.which
      if (key == 13)  // the enter key code
      {
        $('button#stockrefbtn').click()
        return false
      }
    })

    if (/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)) {
      mobileFunctions()
    }
  },
  beforeCreate() {
    this.$emit('getCurrentPage', 'diamonds')
  },
  beforeDestroy() {
    this.$emit('getCurrentPage', '')
  }
}

var $featured = 'NO';
var $shape = '';
var $color = '';
var $clarity = '';
var $cut = '';
var $sym = '';
var $polish = '';
var $certi = '';
var $fluo = '';
var $place = '';
var $weight_from = '';
var $weight_to = '';
var $price_from = '';
var $price_to = '';
var $sorting = '';
var $sorting_weight_direction = 'ASC';
var $sorting_color_direction = 'ASC';
var $sorting_clarity_direction = 'ASC';
var $sorting_cut_direction = 'ASC';
var $sorting_price_direction = 'ASC';
var $sorting_direction = 'ASC';
var $vat_include = 'NO';
var $crr_page = 1;
//for the url to record current conditions
var hash_shape = '';
var hash_color = '';
var hash_clarity = '';
var hash_cut = '';
var hash_sym = '';
var hash_polish = '';
var hash_certi = '';
var hash_fluo = '';
var hash_place = '';
var hash_weight_from = '';
var hash_weight_to = '';
var hash_price_from = '';
var hash_price_to = '';
var hash_sorting = '';
var hash_sorting_weight_direction = 'ASC';
var hash_sorting_color_direction = 'ASC';
var hash_sorting_clarity_direction = 'ASC';
var hash_sorting_cut_direction = 'ASC';
var hash_sorting_price_direction = 'ASC';
var hash_sorting_direction = 'ASC';
var hash_vat_include = 'NO';
var hash_crr_page = 1;
var t_update = 0;
var updatingstatus = 'updated';
var detailOpened = 0;
var crr_ordered_diamond = '';
function parseTheHash() {
  var theHashSTR = location.hash.replace('#', '');
  var theHashArray = theHashSTR.split("&");
  $.each(theHashArray, function (index, value) {
    //console.log( index + ": " + value );
    var crr_freg_condition = value;
    var crr_freg_array = crr_freg_condition.split('=');
    var crr_freg_name = crr_freg_array[0];
    var crr_freg_value = crr_freg_array[1];
    console.log(crr_freg_name + ' -> ' + crr_freg_value);
    if (crr_freg_value != '') {
      if (crr_freg_name == 'shape') {
        var crr_shape_array = crr_freg_value.split('+');
        var shapeChoiceCounter = 0;
        $.each(crr_shape_array, function (i, shapecrr) {
          $('li.filter_shape[title="' + shapecrr + '"]').addClass('btn-active');

          if (shapeChoiceCounter > 0) {
            $shape += ' OR ';
            hash_shape += '+';
          }
          $shape += ' shape = "' + shapecrr + '" ';
          hash_shape += shapecrr;
          shapeChoiceCounter++;
        });
      } else if (crr_freg_name == 'color') {
        if (crr_freg_value != 'Fancy') {
          var crr_color_array = crr_freg_value.split('+');
          var colorChoiceCounter = 0;
          $.each(crr_color_array, function (i, colorcrr) {
            $('li.filter_color[title="' + colorcrr + '"]').addClass('btn-active');
            if (colorcrr != 'Fancy') {
              if (colorChoiceCounter > 0) {
                $color += ' OR ';
                hash_color += '+';
              }
              if (colorcrr == 'N-Z') {
                $color += ' color = "N" OR color="O" OR color="P" OR color="Q" OR color="R" OR color="S" OR color="T" OR color="U" OR color="V" OR color="W" OR color="X" OR color="Y" OR color="Z" ';
                hash_color += 'N-Z';
              } else {
                $color += ' color = "' + colorcrr + '" ';
                hash_color += colorcrr;
              }
              colorChoiceCounter++;
            }
          });
        } else {
          $('li#filter_colorFancy').addClass('btn-active');
          $color = ' color LIKE "F%" AND color <> "F" ';
          hash_color = 'Fancy';
        }
      } else if (crr_freg_name == 'clarity') {
        var crr_clarity_array = crr_freg_value.split('+');
        var clarityChoiceCounter = 0;
        $.each(crr_clarity_array, function (i, claritycrr) {
          $('li.filter_clarity[title="' + claritycrr + '"]').addClass('btn-active');
          if (clarityChoiceCounter > 0) {
            $clarity += ' OR ';
            hash_clarity += '+';
          }
          $clarity += ' clarity = "' + claritycrr + '" ';
          hash_clarity += claritycrr;
          clarityChoiceCounter++;
        });
      } else if (crr_freg_name == 'cut') {
        var crr_cut_array = crr_freg_value.split('+');
        var cutChoiceCounter = 0;
        $.each(crr_cut_array, function (i, cutcrr) {
          $('li.filter_cut[title="' + cutcrr + '"]').addClass('btn-active');
          if (cutChoiceCounter > 0) {
            $cut += ' OR ';
            hash_cut += '+';
          }
          $cut += ' cut_grade = "' + cutcrr + '" ';
          hash_cut += cutcrr;
          cutChoiceCounter++;
        });
      } else if (crr_freg_name == 'polish') {
        var crr_polish_array = crr_freg_value.split('+');
        var ChoiceCounter = 0;
        $.each(crr_polish_array, function (i, polishcrr) {
          $('li.filter_polish[title="' + polishcrr + '"]').addClass('btn-active');
          if (ChoiceCounter > 0) {
            $polish += ' OR ';
            hash_polish += '+';
          }
          $polish += ' polish = "' + polishcrr + '" ';
          hash_polish += polishcrr;
          ChoiceCounter++;
        });
      } else if (crr_freg_name == 'sym') {
        var crr_sym_array = crr_freg_value.split('+');
        var ChoiceCounter = 0;
        $.each(crr_sym_array, function (i, symcrr) {
          $('li.filter_sym[title="' + symcrr + '"]').addClass('btn-active');
          if (ChoiceCounter > 0) {
            $sym += ' OR ';
            hash_sym += '+';
          }
          $sym += ' symmetry = "' + symcrr + '" ';
          hash_sym += symcrr;
          ChoiceCounter++;
        });

      } else if (crr_freg_name == 'fluo') {
        var crr_fluo_array = crr_freg_value.split('+');
        var ChoiceCounter = 0;
        $.each(crr_fluo_array, function (i, fluocrr) {
          $('li.filter_fluo[title="' + fluocrr + '"]').addClass('btn-active');
          if (ChoiceCounter > 0) {
            $fluo += ' OR ';
            hash_fluo += '+';
          }
          if (fluocrr == 'VST') {
            $fluo += ' fluorescence_intensity = "VST" OR fluorescence_intensity = "Very Strong" ';
          } else if (fluocrr == 'STG') {
            $fluo += ' fluorescence_intensity = "STG" OR fluorescence_intensity = "Strong" ';
          } else if (fluocrr == 'MED') {
            $fluo += ' fluorescence_intensity = "MED" OR fluorescence_intensity = "Medium" ';
          } else if (fluocrr == 'FNT') {
            $fluo += ' fluorescence_intensity = "FNT" OR fluorescence_intensity = "SLT"  OR fluorescence_intensity = "VSL" OR fluorescence_intensity = "Faint" OR fluorescence_intensity = "Very Slight" OR fluorescence_intensity = "Slight" ';
          } else if (fluocrr == 'NON') {
            $fluo += ' fluorescence_intensity = "NON" OR fluorescence_intensity = "None"';
          }
          hash_fluo += fluocrr;
          ChoiceCounter++;
        });
      } else if (crr_freg_name == 'place') {

        if (crr_freg_value == 'CHINA') {
          $place += ' country LIKE "China" OR country LIKE "SZ"  OR country LIKE "HK" OR country LIKE "HSTHK" OR country LIKE "Shenzhen" OR country LIKE "Hongkong" OR country LIKE "Hong kong" ';
        } else if (crr_freg_value == 'CZ') {
          $place += ' country LIKE "SZ" OR country LIKE "Shenzhen" ';
        } else if (crr_freg_value == 'HK') {
          $place += ' country LIKE "HK" OR country LIKE "HSTHK" OR country LIKE "Hongkong" OR country LIKE "Hong kong" ';
        } else if (crr_freg_value == 'BELGIUM') {
          $place += ' country LIKE "Belg%" OR country LIKE "Antwerp%" ';
        } else if (crr_freg_value == 'INDIA') {
          $place += ' country LIKE "IND%" ';
        } else if (crr_freg_value == 'OTHER') {
          $place += ' country NOT LIKE "China" AND country NOT LIKE "SZ" AND country NOT LIKE "HK" AND country NOT LIKE "HSTHK" AND country NOT LIKE "Belg%" AND country NOT LIKE "IND%" ';
        } else if (crr_freg_value == 'ALL') {
          $place = '';
        }

      } else if (crr_freg_name == 'certi') {
        var crr_certi_array = crr_freg_value.split('+');
        var certiChoiceCounter = 0;
        $.each(crr_certi_array, function (i, certicrr) {
          $('li.filter_certi[title="' + certicrr + '"]').addClass('btn-active');
          if (certiChoiceCounter > 0) {
            $certi += ' OR ';
            hash_certi += '+';
          }
          $certi += ' grading_lab = "' + certicrr + '" ';
          hash_certi += certicrr;
          certiChoiceCounter++;
        });
      } else if (crr_freg_name == 'weight_from') {
        hash_weight_from = crr_freg_value;
        $weight_from = crr_freg_value;
      } else if (crr_freg_name == 'weight_to') {
        hash_weight_to = crr_freg_value;
        $weight_to = crr_freg_value;
      } else if (crr_freg_name == 'price_from') {
        $price_from = crr_freg_value;
        hash_price_from = crr_freg_value;
        $('input#price_from').val($price_from);

      } else if (crr_freg_name == 'price_to') {
        $price_to = crr_freg_value;
        hash_price_to = crr_freg_value;
        $('input#price_to').val($price_to);

      } else if (crr_freg_name == 'sorting') {
        $sorting = crr_freg_value;
        hash_sorting = crr_freg_value;
        $('li#sortingby_' + crr_freg_value).addClass('active');
      } else if (crr_freg_name == 'sorting_direction') {
        $sorting_direction = 'crr_freg_value';
        hash_sorting_direction = crr_freg_value;
        $('li.sortingbtn').removeClass('DESC');
        $('li.sortingbtn.active').addClass(crr_freg_value);
      } else if (crr_freg_name == 'crr_page') {
        $crr_page = parseInt(crr_freg_value);
        hash_crr_page = crr_freg_value;
      } else if (crr_freg_name == 'vat_choice') {
        if (crr_freg_value == "YES") {
          $('input#vat-choice').prop('checked', true);
          $vat_include = 'YES';
          hash_vat_include = 'YES';
        } else {
          $('input#vat-choice').prop('checked', false);
          $vat_include = 'NO';
          hash_vat_include = 'NO';
        }
      }
    }
  });
}//////////////////////// function parseTheHash() FINI

function diamondlistpagenavi(howmanyrecords) {
  $('span#found-diamonds-num').html(howmanyrecords);
  //$('#dia-num-bg').stop(true).remove();
  //$('span#found-diamonds-num-container').append('<span id="dia-num-bg"></span>');
  $('#dia-num-bg').stop(true).css({ 'opacity': 1, 'display': 'inline-block' }).delay(58).fadeOut(2000);

  $('span#diapagenavi').empty();
  var totalrecords = parseFloat(howmanyrecords);
  var totalpages = Math.ceil(totalrecords / 28);

  if (totalpages < $crr_page) {
    $crr_page = 1;
  }

  if (totalpages <= 13) {
    var loop_beginpage = 1;
    var loop_endpage = totalpages;
  } else {
    if (($crr_page + 10) < totalpages) {
      var loop_beginpage = $crr_page;
      var loop_endpage = $crr_page + 10;
    } else {
      var loop_beginpage = totalpages - 10;
      var loop_endpage = totalpages;
    }
  }

  if (loop_beginpage > 1) {
    $('span#diapagenavi').append('<span class="dia-page-btn" onclick="choosethispage(1)">1</span>');
  }
  if (loop_beginpage > 2) {
    $('span#diapagenavi').append('<span class="somethinginbetween-indi"> ... </span>');
  }

  for (var i = loop_beginpage; i <= loop_endpage; i++) {
    if (i == $crr_page) {
      $('span#diapagenavi').append('<span class="dia-page-btn" id="crr_page">' + i + '</span>');
    } else {
      $('span#diapagenavi').append('<span class="dia-page-btn" onclick="choosethispage(' + i + ')">' + i + '</span>');
    }
  }

  if ((totalpages - loop_endpage) > 1) {
    $('span#diapagenavi').append('<span class="somethinginbetween-indi"> ... </span>');
  }
  if ((totalpages - loop_endpage) >= 1) {
    $('span#diapagenavi').append('<span class="dia-page-btn" onclick="choosethispage(' + totalpages + ')">' + totalpages + '</span>');
  }
}

if (/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)) {
  // TODO jquery ui vue
//   $(function () {
//     $("#slider-range-weight").slider({
//       range: true,
//       step: 0.1,
//       min: 0,
//       max: 12,
//       values: [0.5, 3]
//     });
//   });

//   $(function () {
//     $("#slider-range").slider({
//       range: true,
//       step: 1,
//       min: 10,
//       max: 288,
//       values: [11, 260],
//     });
//   });
// } else {// if it's mobile device, event listener should be different than desktop
//   $(function () {
//     $("#slider-range-weight").slider({
//       range: true,
//       step: 0.1,
//       min: 0,
//       max: 12,
//       values: [0.5, 3],
//       slide: function (event, ui) {
//         $('#weight_from').val(ui.values[0]);
//         $('#weight_to').val(ui.values[1]);
//       },
//       stop: function (event, ui) {
//         //console.log('stopped');
//         filter_weight();
//       }
//     });
//   });

//   $(function () {
//     $("#slider-range").slider({
//       range: true,
//       step: 1,
//       min: 10,
//       max: 288,
//       values: [11, 260],
//       slide: function (event, ui) {
//         var slidervalue_from = ui.values[0];
//         var slidervalue_to = ui.values[1];
//         if (slidervalue_from < 188) {
//           var pricevaluefrom = slidervalue_from * 50 - 400;
//         } else {
//           var pricevaluefrom = (slidervalue_from - 188) * (slidervalue_from - 188) * 68 + 9450 - 400;
//         }
//         if (slidervalue_to < 188) {
//           var pricevalueto = slidervalue_to * 50 - 400;
//         } else {
//           var pricevalueto = (slidervalue_to - 188) * (slidervalue_to - 188) * 68 + 9450 - 400;
//         }

//         $('#price_from').val(pricevaluefrom);
//         $('#price_to').val(pricevalueto);
//       },
//       stop: function (event, ui) {
//         //console.log('stopped');
//         filter_price();
//       }
//     });
//   });
}
