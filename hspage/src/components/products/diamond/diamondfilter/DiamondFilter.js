let Images = require.context('@/_images/constant/', false, /\.(png|jpg)$/);
export default {
  name: 'DiamondFilter',
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
        }).then(response => {
          if (response.status === 200) {
            // token
            updatingstatus = 'updated'
            clearTimeout(t)
            if (data == 'NO-DIAMOND-FOUND') {
              $('div#diamondlist').html('<p id="nodiafoundfeedback">抱歉，没有找到符合您所选条件的钻石</p>')
              $('span#diapagenavi').html(' - ')
            } else {
              $('div#diamondlist').html(data)
              var howmanyrecords = $("div#howmanyrecords").html()
              $('#loadingtxt').html('找到<strong>' + howmanyrecords + '</strong>颗')
              diamondlistpagenavi(howmanyrecords)

            }
            $('div#loadingcover').delay(880).fadeOut('fast')            
          }
        }, err => { console.log(err); alert('error:' + err.bodyText) })
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
    }
  }
}