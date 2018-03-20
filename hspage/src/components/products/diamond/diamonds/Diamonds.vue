<template>
  <div>
    <vue-title title="合适钻石 | BEYOU DIAMOND | 挑选钻石"></vue-title>
    <app-header></app-header>
    <div class="contentmainwrapper">
      <h2 class="pagetitle">合适钻石</h2>
      <h3 class="pagesubtitle">挑选适合您的钻石</h3>
      <p class="subtitle-expl">
        <span class="glyphicon glyphicon-hand-right"></span> 将您选好的钻石放到预订单中，提交后可以预约到公司来看钻石。
      </p>
      <!-- start of <diamond-filter -->
      <div id="filter_box">
        <a id="fancycolorbtn" href="/product/recommenddiamonds"><span class="glyphicon glyphicon-thumbs-up"></span> 合适推荐 &raquo;</a>
        <div id="weightandpricebox">
          <div class="filter_line_dependent" id="filter_line_weight">
            <span class="filter_title">钻重</span>
            <div id="slider-range-weight" style="width:288px;"></div>
            <input type="number" id="weight_from" value="0.2" step=".1" /> - <input type="number" id="weight_to" value="20" step=".1" /> ct <button type="button" id="btn_weight" v-on:click="filter_weight()"><span class="glyphicon glyphicon-refresh"></span> 刷新</button>
            <span id="tooltip-weight"></span>
          </div>

          <div class="filter_line_dependent" id="filter_line_price">
            <span class="filter_title">价格</span>
            <div id="slider-range" style="width:288px;"></div>
            <input type="number" id="price_from" value="100" /> - <input type="number" id="price_to" value="899999" /> € <button type="button" id="btn_price" v-on:click="filter_price()"><span class="glyphicon glyphicon-refresh"></span> 刷新</button>
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
              <li class="filter_shape" id="filter_shapeBR" title="BR" v-on:click="filter_shape('BR')">
                <span class="shapeiconcontainer" title="圆形" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:0px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapePS" title="PS" v-on:click="filter_shape('PS')">
                <span class="shapeiconcontainer" title="梨形" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-87px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapePR" title="PR" v-on:click="filter_shape('PR')">
                <span class="shapeiconcontainer" title="公主方" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-66px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapeHS" title="HS" v-on:click="filter_shape('HS')">
                <span class="shapeiconcontainer" title="心形" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-154px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapeMQ" title="MQ" v-on:click="filter_shape('MQ')">
                <span class="shapeiconcontainer" title="橄榄形" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-132px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapeOV" title="OV" v-on:click="filter_shape('OV')">
                <span class="shapeiconcontainer" title="椭圆形" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-110px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapeEM" title="EM" v-on:click="filter_shape('EM')">
                <span class="shapeiconcontainer" title="祖母绿" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-176px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapeRAD" title="RAD" v-on:click="filter_shape('RAD')">
                <span class="shapeiconcontainer" title="雷蒂恩" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-44px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapeCU" title="CU" v-on:click="filter_shape('CU')">
              <span class="shapeiconcontainer" title="枕形" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-198px 0px;'></span>
              </li><li class="filter_shape" id="filter_shapeAS" title="AS" v-on:click="filter_shape('AS')">
                <span class="shapeiconcontainer" title="Asscher" :style="{ 'background-image': 'url(' + getImage + ')' }" style='background-position:-22px 0px;'></span>
              </li>
            </ul>
          </div>
          <!-- color -->
          <div class="filter_line" id="filterline-color">
            <span class="filter_title">颜色</span>
            <ul class="filter_color_outer">
              <li class="filter_color" id="filter_colorD" title="D" v-on:click="filter_color('D')">D
              </li><li class="filter_color" id="filter_colorE" title="E" v-on:click="filter_color('E')">E
              </li><li class="filter_color" id="filter_colorF" title="F" v-on:click="filter_color('F')">F
              </li><li class="filter_color" id="filter_colorG" title="G" v-on:click="filter_color('G')">G
              </li><li class="filter_color" id="filter_colorH" title="H" v-on:click="filter_color('H')">H
              </li><li class="filter_color" id="filter_colorI" title="I" v-on:click="filter_color('I')">I
              </li><li class="filter_color" id="filter_colorJ" title="J" v-on:click="filter_color('J')">J
              </li><li class="filter_color" id="filter_colorK" title="K" v-on:click="filter_color('K')">K
              </li><li class="filter_color" id="filter_colorL" title="L" v-on:click="filter_color('L')">L
              </li><li class="filter_color" id="filter_colorM" title="M" v-on:click="filter_color('M')">M
              </li><li class="filter_color" id="filter_colorN" title='N-Z' v-on:click="filter_color('N')">N-Z
              </li><li class="filter_color" id="filter_colorFancy" title='Fancy' v-on:click="filter_color('Fancy')">彩钻
              </li>
            </ul>
          </div>
          <!-- clarity -->
          <div class="filter_line" id="filterline-clarity">
            <span class="filter_title">净度</span>
            <ul class="filter_clarity_outer">
              <li class="filter_clarity" id="filter_clarityFL" title="FL" v-on:click="filter_clarity('FL')">FL
              </li><li class="filter_clarity" id="filter_clarityIF" title="IF" v-on:click="filter_clarity('IF')">IF
              </li><li class="filter_clarity" id="filter_clarityVVS1" title="VVS1" v-on:click="filter_clarity('VVS1')">VVS1
              </li><li class="filter_clarity" id="filter_clarityVVS2" title="VVS2" v-on:click="filter_clarity('VVS2')">VVS2
              </li><li class="filter_clarity" id="filter_clarityVS1" title="VS1" v-on:click="filter_clarity('VS1')">VS1
              </li><li class="filter_clarity" id="filter_clarityVS2" title="VS2" v-on:click="filter_clarity('VS2')">VS2
              </li><li class="filter_clarity" id="filter_claritySI1" title="SI1" v-on:click="filter_clarity('SI1')">SI1
              </li><li class="filter_clarity" id="filter_claritySI2" title="SI2" v-on:click="filter_clarity('SI2')">SI2
              </li><li class="filter_clarity" id="filter_clarityI1" title="I1" v-on:click="filter_clarity('I1')">I1</li>
            </ul>
          </div>
          <!-- lab -->
          <div class="filter_line" id="filterline-lab">
            <span class="filter_title">证书</span>
            <ul class="filter_certi_outer">
              <li class="filter_certi" id="filter_certiIGI" title="IGI" v-on:click="filter_certi('IGI')">IGI
              </li><li class="filter_certi" id="filter_certiHRD" title="HRD" v-on:click="filter_certi('HRD')">HRD
              </li><li class="filter_certi" id="filter_certiGIA" title="GIA" v-on:click="filter_certi('GIA')">GIA</li>
            </ul>
          </div>
        </div><!-- filterboxinner-left -->
        <div id="filterboxinner-right">
          <!-- cut -->
          <div class="filter_line" id="filterline-cut">
            <span class="filter_title">切工</span>
            <ul class="filter_cut_outer">
              <li class="filter_cut" id="filter_cutEX" title="EX" v-on:click="filter_cut('EX')">EX
              </li><li class="filter_cut" id="filter_cutVG" title="VG" v-on:click="filter_cut('VG')">VG
              </li><li class="filter_cut" id="filter_cutG" title="G" v-on:click="filter_cut('G')">G
              </li><li class="filter_cut" id="filter_cutF" title="F" v-on:click="filter_cut('F')">F</li>
            </ul>
          </div>
          <!-- polish -->
          <div class="filter_line" id="filterline-polish">
          <span class="filter_title">抛光</span>
            <ul class="filter_polish_outer">
              <li class="filter_polish" id="filter_polishEX" title="EX" v-on:click="filter_polish('EX')">EX
              </li><li class="filter_polish" id="filter_polishVG" title="VG" v-on:click="filter_polish('VG')">VG
              </li><li class="filter_polish" id="filter_polishG" title="G" v-on:click="filter_polish('G')">G
              </li><li class="filter_polish" id="filter_polishF" title="F" v-on:click="filter_polish('F')">F</li>
            </ul>
          </div>
          <!-- symmetry -->
          <div class="filter_line" id="filterline-symmetry">
            <span class="filter_title">对称</span>
            <ul class="filter_sym_outer">
              <li class="filter_sym" id="filter_symEX" title="EX" v-on:click="filter_sym('EX')">EX
              </li><li class="filter_sym" id="filter_symVG" title="VG" v-on:click="filter_sym('VG')">VG
              </li><li class="filter_sym" id="filter_symG" title="G" v-on:click="filter_sym('G')">G
              </li><li class="filter_sym" id="filter_symF" title="F" v-on:click="filter_sym('F')">F</li>
            </ul>
          </div>
          <!-- fluo -->
          <div class="filter_line" id="filterline-fluo">
            <span class="filter_title">荧光</span>
            <ul class="filter_fluo_outer">
              <li class="filter_fluo" id="filter_fluoVST" title="VST" v-on:click="filter_fluo('VST')">极强
              </li><li class="filter_fluo" id="filter_fluoSTG" title="STG" v-on:click="filter_fluo('STG')">强
              </li><li class="filter_fluo" id="filter_fluoMED" title="MED" v-on:click="filter_fluo('MED')">中
              </li><li class="filter_fluo" id="filter_fluoFNT" title="FNT" v-on:click="filter_fluo('FNT')">弱
              </li><li class="filter_fluo" id="filter_fluoNON" title="NON" v-on:click="filter_fluo('NON')">无</li>
            </ul>
          </div>
        </div><!-- filterboxinner-right -->
        <br class="clear" />
        <div id="filterinner-diaid">
          <p class="placechoosercontainer">
            <span class="filter_title">地点</span>
            <select id="placechooser" v-on:change="filter_place()">
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
              <button id="stockrefbtn" v-on:click="searchbyRef()"><span class="glyphicon glyphicon-search"></span></button>
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
          <li class="sortingbtn" id="sortingby_weight" v-on:click="sortingby('weight')">重量 <span class="glyphicon glyphicon-menu-up" aria-hidden="true"></span> <span class="glyphicon glyphicon-menu-down" aria-hidden="true"></span></li>
          <li class="sortingbtn" id="sortingby_price" v-on:click="sortingby('price')">价格 <span class="glyphicon glyphicon-menu-up" aria-hidden="true"></span> <span class="glyphicon glyphicon-menu-down" aria-hidden="true"></span></li>
          <li class="sortingbtn" id="sortingby_color" v-on:click="sortingby('color')">颜色 <span class="glyphicon glyphicon-menu-up" aria-hidden="true"></span> <span class="glyphicon glyphicon-menu-down" aria-hidden="true"></span></li>
          <li class="sortingbtn" id="sortingby_clarity" v-on:click="sortingby('clarity')">净度 <span class="glyphicon glyphicon-menu-up" aria-hidden="true"></span> <span class="glyphicon glyphicon-menu-down" aria-hidden="true"></span></li>
        </ul>
        <br class="clear" />
      </div><!-- end of </diamond-filter> -->
      <div id="diamondlist">
        <diamonds-data :diamonds="diamonds">
        </diamonds-data>
      </div>
      <div id="list-page-navi">第 <span id="diapagenavi"></span> 页</div>
      </div>
      <div id="feedbackcover">
        <div id="feedbackbox">
          <div class="bgimgcover" style="position:absolute; width:100%; height:100%; left:0; top:0; background-color:rgba(151,240,84,0.85); z-index:0;">
          </div>
          <p class="ordersuccwords"><span class="glyphicon glyphicon-thumbs-up"></span> 您成功预定了钻石!</p>
          <p style="z-index:2; position:relative;"><a id="ordersuccshoppinglistlinker" href="/shoppinglist-confirmed.php"><span class="glyphicon glyphicon-th-list"></span> 查看预定单</a> <button id="ordersuccmountinglinker" v-on:click="gotoMountingsfor()"><span></span> 配戒托</button> <button id="ordersuccsfdbcXbtn" v-on:click="xtheSfdb()">继续选购 &raquo;</button></p>
        </div>
      </div>
      <!-- if Agent yes and $accept_order_online=='TOBECHECKED' -->
      <div v-if="thelettertoagents" id="thelettertoagents-container">
        <img id="thelettertoagents" src="_images/constant/lettertoagents.jpg" />
        <p class="agentletter-actionbox">
          <button id="agentletterokbtn" v-on:click="agentLetterOk()">好的，知道了</button>
        </p>
      </div>
    </div>
</template>
<script src='./Diamonds.js'></script>
<style src='./Diamonds.css'></style>
