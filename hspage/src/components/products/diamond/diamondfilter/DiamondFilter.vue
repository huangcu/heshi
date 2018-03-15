<template>
  <div>
    <div id="filter_box">
      <a id="fancycolorbtn" href="recommend-diamonds.php"><span class="glyphicon glyphicon-thumbs-up"></span> 合适推荐 &raquo;</a>
      <div id="weightandpricebox">
        <div class="filter_line_dependent" id="filter_line_weight">
          <span class="filter_title">钻重</span>
          <div id="slider-range-weight" style="width:288px;"></div>
          <input type="number" id="weight_from" value="0.2" step=".1" /> - <input type="number" id="weight_to" value="20" step=".1" /> ct <button type="button" id="btn_weight" onclick="filter_weight()"><span class="glyphicon glyphicon-refresh"></span> 刷新</button>
          <span id="tooltip-weight"></span>
        </div>

        <div class="filter_line_dependent" id="filter_line_price">
          <span class="filter_title">价格</span>
          <div id="slider-range" style="width:288px;"></div>
          <input type="number" id="price_from" value="100" /> - <input type="number" id="price_to" value="899999" /> € <button type="button" id="btn_price" onclick="filter_price()"><span class="glyphicon glyphicon-refresh"></span> 刷新</button>
          <span class="vatchoicebox"><input type="checkbox" id="vat-choice" class="vat-choice" name="vat-choice" onchange="vatChoice()" /> 价格含税</span>
          <span id="tooltip-price"></span>
        </div>
        <br class="clear"/>
      </div>

      <div id="filterboxinner-left">
        <!-- shape -->
        <div class="filter_line" id="filterline-shape">
          <span class="filter_title" id="filter_title_shape">形状</span>
          <ul class="fileber_shape_outer">
            <li class="filter_shape" id="filter_shapeBR" title="BR" onclick="filter_shape('BR')">
              <span class="shapeiconcontainer" title="圆形" style="background-image:url(/_images/constant/ico-stones.png); background-position:0px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapePS" title="PS" onclick="filter_shape('PS')">
              <span class="shapeiconcontainer" title="梨形" style="background-image:url(/_images/constant/ico-stones.png); background-position:-87px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapePR" title="PR" onclick="filter_shape('PR')">
              <span class="shapeiconcontainer" title="公主方" style="background-image:url(/_images/constant/ico-stones.png); background-position:-66px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapeHS" title="HS" onclick="filter_shape('HS')">
              <span class="shapeiconcontainer" title="心形" style="background-image:url(/_images/constant/ico-stones.png); background-position:-154px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapeMQ" title="MQ" onclick="filter_shape('MQ')">
              <span class="shapeiconcontainer" title="橄榄形" style="background-image:url(/_images/constant/ico-stones.png); background-position:-132px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapeOV" title="OV" onclick="filter_shape('OV')">
              <span class="shapeiconcontainer" title="椭圆形" style="background-image:url(/_images/constant/ico-stones.png); background-position:-110px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapeEM" title="EM" onclick="filter_shape('EM')">
              <span class="shapeiconcontainer" title="祖母绿" style="background-image:url(/_images/constant/ico-stones.png); background-position:-176px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapeRAD" title="RAD" onclick="filter_shape('RAD')">
              <span class="shapeiconcontainer" title="雷蒂恩" style="background-image:url(/_images/constant/ico-stones.png); background-position:-44px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapeCU" title="CU" onclick="filter_shape('CU')">
            <span class="shapeiconcontainer" title="枕形" style="background-image:url(/_images/constant/ico-stones.png); background-position:-198px 0px;"></span>
            </li><li class="filter_shape" id="filter_shapeAS" title="AS" onclick="filter_shape('AS')">
              <span class="shapeiconcontainer" title="Asscher" style="background-image:url(/_images/constant/ico-stones.png); background-position:-22px 0px;"></span>
            </li>
          </ul>
        </div>
        <!-- color -->
        <div class="filter_line" id="filterline-color">
          <span class="filter_title">颜色</span>
          <ul class="filter_color_outer">
            <li class="filter_color" id="filter_colorD" title="D" onclick="filter_color('D')">D
            </li><li class="filter_color" id="filter_colorE" title="E" onclick="filter_color('E')">E
            </li><li class="filter_color" id="filter_colorF" title="F" onclick="filter_color('F')">F
            </li><li class="filter_color" id="filter_colorG" title="G" onclick="filter_color('G')">G
            </li><li class="filter_color" id="filter_colorH" title="H" onclick="filter_color('H')">H
            </li><li class="filter_color" id="filter_colorI" title="I" onclick="filter_color('I')">I
            </li><li class="filter_color" id="filter_colorJ" title="J" onclick="filter_color('J')">J
            </li><li class="filter_color" id="filter_colorK" title="K" onclick="filter_color('K')">K
            </li><li class="filter_color" id="filter_colorL" title="L" onclick="filter_color('L')">L
            </li><li class="filter_color" id="filter_colorM" title="M" onclick="filter_color('M')">M
            </li><li class="filter_color" id="filter_colorN" title='N-Z' onclick="filter_color('N')">N-Z
            </li><li class="filter_color" id="filter_colorFancy" title='Fancy' onclick="filter_color('Fancy')">彩钻
            </li>
          </ul>
        </div>
        <!-- clarity -->
        <div class="filter_line" id="filterline-clarity">
          <span class="filter_title">净度</span>
          <ul class="filter_clarity_outer">
            <li class="filter_clarity" id="filter_clarityFL" title="FL" onclick="filter_clarity('FL')">FL
            </li><li class="filter_clarity" id="filter_clarityIF" title="IF" onclick="filter_clarity('IF')">IF
            </li><li class="filter_clarity" id="filter_clarityVVS1" title="VVS1" onclick="filter_clarity('VVS1')">VVS1
            </li><li class="filter_clarity" id="filter_clarityVVS2" title="VVS2" onclick="filter_clarity('VVS2')">VVS2
            </li><li class="filter_clarity" id="filter_clarityVS1" title="VS1" onclick="filter_clarity('VS1')">VS1
            </li><li class="filter_clarity" id="filter_clarityVS2" title="VS2" onclick="filter_clarity('VS2')">VS2
            </li><li class="filter_clarity" id="filter_claritySI1" title="SI1" onclick="filter_clarity('SI1')">SI1
            </li><li class="filter_clarity" id="filter_claritySI2" title="SI2" onclick="filter_clarity('SI2')">SI2
            </li><li class="filter_clarity" id="filter_clarityI1" title="I1" onclick="filter_clarity('I1')">I1</li>
          </ul>
        </div>
        <!-- lab -->
        <div class="filter_line" id="filterline-lab">
          <span class="filter_title">证书</span>
          <ul class="filter_certi_outer">
            <li class="filter_certi" id="filter_certiIGI" title="IGI" onclick="filter_certi('IGI')">IGI
            </li><li class="filter_certi" id="filter_certiHRD" title="HRD" onclick="filter_certi('HRD')">HRD
            </li><li class="filter_certi" id="filter_certiGIA" title="GIA" onclick="filter_certi('GIA')">GIA</li>
          </ul>
        </div>
      </div><!-- filterboxinner-left -->
      <div id="filterboxinner-right">
        <!-- cut -->
        <div class="filter_line" id="filterline-cut">
          <span class="filter_title">切工</span>
          <ul class="filter_cut_outer">
            <li class="filter_cut" id="filter_cutEX" title="EX" onclick="filter_cut('EX')">EX
            </li><li class="filter_cut" id="filter_cutVG" title="VG" onclick="filter_cut('VG')">VG
            </li><li class="filter_cut" id="filter_cutG" title="G" onclick="filter_cut('G')">G
            </li><li class="filter_cut" id="filter_cutF" title="F" onclick="filter_cut('F')">F</li>
          </ul>
        </div>
        <!-- polish -->
        <div class="filter_line" id="filterline-polish">
        <span class="filter_title">抛光</span>
          <ul class="filter_polish_outer">
            <li class="filter_polish" id="filter_polishEX" title="EX" onclick="filter_polish('EX')">EX
            </li><li class="filter_polish" id="filter_polishVG" title="VG" onclick="filter_polish('VG')">VG
            </li><li class="filter_polish" id="filter_polishG" title="G" onclick="filter_polish('G')">G
            </li><li class="filter_polish" id="filter_polishF" title="F" onclick="filter_polish('F')">F</li>
          </ul>
        </div>
        <!-- symmetry -->
        <div class="filter_line" id="filterline-symmetry">
          <span class="filter_title">对称</span>
          <ul class="filter_sym_outer">
            <li class="filter_sym" id="filter_symEX" title="EX" onclick="filter_sym('EX')">EX
            </li><li class="filter_sym" id="filter_symVG" title="VG" onclick="filter_sym('VG')">VG
            </li><li class="filter_sym" id="filter_symG" title="G" onclick="filter_sym('G')">G
            </li><li class="filter_sym" id="filter_symF" title="F" onclick="filter_sym('F')">F</li>
          </ul>
        </div>
        <!-- fluo -->
        <div class="filter_line" id="filterline-fluo">
          <span class="filter_title">荧光</span>
          <ul class="filter_fluo_outer">
            <li class="filter_fluo" id="filter_fluoVST" title="VST" onclick="filter_fluo('VST')">极强
            </li><li class="filter_fluo" id="filter_fluoSTG" title="STG" onclick="filter_fluo('STG')">强
            </li><li class="filter_fluo" id="filter_fluoMED" title="MED" onclick="filter_fluo('MED')">中
            </li><li class="filter_fluo" id="filter_fluoFNT" title="FNT" onclick="filter_fluo('FNT')">弱
            </li><li class="filter_fluo" id="filter_fluoNON" title="NON" onclick="filter_fluo('NON')">无</li>
          </ul>
        </div>
      </div><!-- filterboxinner-right -->
      <br class="clear" />
      <div id="filterinner-diaid">
        <p class="placechoosercontainer">
          <span class="filter_title">地点</span>
          <select id="placechooser" onchange="filter_place()">
            <option value="ALL">全部</option>
            <option value="BELGIUM">比利时</option>
            <option value="CHINA">中国</option>
            <option value="CZ">-深圳</option>
            <option value="HK">-香港</option>
            <option value="INDIA">印度</option>
            <option value="OTHER">其他地区</option>
          </select>
        </p>
        <div class="filter_line_dependent" id="filterline-stockref">
          <span id="therefsearcherbox">
            <input name="stockreftosearch" id="stockreftosearch" placeholder="钻石ID／证书编号" />
            <button id="stockrefbtn" onclick="searchbyRef()"><span class="glyphicon glyphicon-search"></span></button>
          </span>
          <div id="found-diamonds-num-box">为您找到<span id="found-diamonds-num-container"><span id="dia-num-bg"></span><span id="found-diamonds-num"></span></span>颗钻石</div>
        </div>
        <br class="clear" />
    </div>
    <br class="clear" />
    </div>
    <div id="diamondsorting-box">
      <label class="filter_title">排序</label>
        <ul class="sortingbox">
        <li class="sortingbtn" id="sortingby_weight" onclick="sortingby('weight')">重量 <span class="glyphicon glyphicon-menu-up" aria-hidden="true"></span> <span class="glyphicon glyphicon-menu-down" aria-hidden="true"></span></li>
        <li class="sortingbtn" id="sortingby_price" onclick="sortingby('price')">价格 <span class="glyphicon glyphicon-menu-up" aria-hidden="true"></span> <span class="glyphicon glyphicon-menu-down" aria-hidden="true"></span></li>
        <li class="sortingbtn" id="sortingby_color" onclick="sortingby('color')">颜色 <span class="glyphicon glyphicon-menu-up" aria-hidden="true"></span> <span class="glyphicon glyphicon-menu-down" aria-hidden="true"></span></li>
        <li class="sortingbtn" id="sortingby_clarity" onclick="sortingby('clarity')">净度 <span class="glyphicon glyphicon-menu-up" aria-hidden="true"></span> <span class="glyphicon glyphicon-menu-down" aria-hidden="true"></span></li>
      </ul>
      <br class="clear" />
    </div>
  </div>
</template>
