
$(function () {
  $("#thetime").datepicker({
    numberOfMonths: 2,
    dayNames: ["周日", "周一", "周二", "周三", "周四", "周五", "周六"],
    dayNamesMin: ["日", "一", "二", "三", "四", "五", "六"],
    monthNames: ["一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"],
    buttonText: "Select date"
  });
  $("#thetime").datepicker("option", "dateFormat", 'yy-mm-dd');
});

function filter_shape(theshape) {
  $('#shape_defined').val('');
  $('li.filter_shape').removeClass('btn-active');
  $('li#filter_shape' + theshape).addClass('btn-active');
  $('span#theshapetxt').html($('li#filter_shape' + theshape).attr('title'));
}
function filter_color(thecolor) {
  $('#fancycolorinput').val('');
  $('li.filter_color').removeClass('btn-active');
  $('li#filter_color' + thecolor).addClass('btn-active');
}
function filter_clarity(theclarity) {
  $('#clarity_defined').val('');
  $('li.filter_clarity').removeClass('btn-active');
  $('li#filter_clarity' + theclarity).addClass('btn-active');
}
function filter_cut(thegrade) {
  $('#cut_defined').val('');
  $('li.filter_cut').removeClass('btn-active');
  $('li#filter_cut' + thegrade).addClass('btn-active');
}
function filter_polish(thegrade) {
  $('#polish_defined').val('');
  $('li.filter_polish').removeClass('btn-active');
  $('li#filter_polish' + thegrade).addClass('btn-active');
}
function filter_sym(thegrade) {
  $('#symmetry_defined').val('');
  $('li.filter_sym').removeClass('btn-active');
  $('li#filter_sym' + thegrade).addClass('btn-active');
}
function filter_certi(thelab) {
  $('#certi_defined').val('');
  $('li.filter_certi').removeClass('btn-active');
  $('li#filter_certi' + thelab).addClass('btn-active');
}
function filter_fluo(thegrade) {
  $('#fluo_defined').val('');
  $('li.filter_fluo').removeClass('btn-active');
  $('li#filter_fluo' + thegrade).addClass('btn-active');
}

function submitTheForm() {
  $('#fordb_carat').val($('#weight').val());


  if ($.trim($('#shape_defined').val()) != '') {
    var shape = $.trim($('#shape_defined').val());
  } else {
    var shape = $('.fileber_shape_outer>li.btn-active').attr('title');
  }
  $('#fordb_shape').val(shape);

  if ($.trim($('#fancycolorinput').val()) != '') {
    var color = $.trim($('#fancycolorinput').val());
  } else {
    var color = $('.filter_color_outer>li.btn-active').attr('title');
  }
  $('#fordb_color').val(color);

  if ($.trim($('#clarity_defined').val()) != '') {
    var clarity = $.trim($('#clarity_defined').val());
  } else {
    var clarity = $('.filter_clarity_outer>li.btn-active').attr('title');
  }
  $('#fordb_clarity').val(clarity);

  if ($.trim($('#cut_defined').val()) != '') {
    var cut = $.trim($('#cut_defined').val());
  } else {
    var cut = $('.filter_cut_outer>li.btn-active').attr('title');
  }
  $('#fordb_cut').val(cut);

  if ($.trim($('#symmetry_defined').val()) != '') {
    var symm = $.trim($('#symmetry_defined').val());
  } else {
    var symm = $('.filter_sym_outer>li.btn-active').attr('title');
  }
  $('#fordb_symmetry').val(symm);

  if ($.trim($('#polish_defined').val()) != '') {
    var polish = $.trim($('#polish_defined').val());
  } else {
    var polish = $('.filter_polish_outer>li.btn-active').attr('title');
  }
  $('#fordb_polish').val(polish);

  if ($.trim($('#fluo_defined').val()) != '') {
    var fluo = $.trim($('#fluo_defined').val());
  } else {
    var fluo = $('.filter_fluo_outer>li.btn-active').attr('title');
  }
  $('#fordb_fluo').val(fluo);
  if ($.trim($('#certi_defined').val()) != '') {
    var certi = $.trim($('#certi_defined').val());
  } else {
    var certi = $('.filter_certi_outer>li.btn-active').attr('title');
  }
  $('#fordb_lab').val(certi);
  $('form#theform').submit();
}
