<template>
  <div>
    <vue-title title="管理界面:添加历史购买纪录"></vue-title>
    <div id="maincontent">
      <h3>为用户<strong> username</strong>添加购买历史记录</h3>
      <p><a href="userhistory.php?user=userid">&laquo; 返回  username的历史纪录</a></p>
      <div>
        <p v-if="errormessage" class="errormessage">errormessage</p>
        <p v-if="successmessage" class="successmessage">successmessage</p>
        <!-- shape -->
        <div class="filter_line" id="filterline-shape">
          <h3 style="margin-bottom:25px;">购买的钻石：</h3>
          <p><label>钻重 carat</label>
            <input type="number" id="weight" />
          </p>
          <span class="filter_title" id="filter_title_shape">形状: <span id="theshapetxt"></span></span> 
          <ul class="fileber_shape_outer">
            <li class="filter_shape" id="filter_shapeBR" title="BR" onclick="filter_shape('BR')">
              <span class="shapeiconcontainer" title="round" style="background-image:url(/_images/constant/ico-stones.png); background-position:0px 0px;"></span>
            </li>
            <li class="filter_shape" id="filter_shapePS" title="PS" onclick="filter_shape('PS')">
              <span class="shapeiconcontainer" title="pear" style="background-image:url(/_images/constant/ico-stones.png); background-position:-87px 0px;"></span>
            </li>
            <li class="filter_shape" id="filter_shapePR" title="PR" onclick="filter_shape('PR')">
              <span class="shapeiconcontainer" title="princess" style="background-image:url(/_images/constant/ico-stones.png); background-position:-66px 0px;"></span>
            </li>
            <li class="filter_shape" id="filter_shapeHS" title="HS" onclick="filter_shape('HS')">
              <span class="shapeiconcontainer" title="heart" style="background-image:url(/_images/constant/ico-stones.png); background-position:-154px 0px;"></span>
            </li>
            <li class="filter_shape" id="filter_shapeMQ" title="MQ" onclick="filter_shape('MQ')">
            <span class="shapeiconcontainer" title="marquise" style="background-image:url(/_images/constant/ico-stones.png); background-position:-132px 0px;"></span>
            </li>
            <li class="filter_shape" id="filter_shapeOV" title="OV" onclick="filter_shape('OV')">
              <span class="shapeiconcontainer" title="oval" style="background-image:url(/_images/constant/ico-stones.png); background-position:-110px 0px;"></span>
            </li>
            <li class="filter_shape" id="filter_shapeEM" title="EM" onclick="filter_shape('EM')">
              <span class="shapeiconcontainer" title="emerald" style="background-image:url(/_images/constant/ico-stones.png); background-position:-176px 0px;"></span>
            </li>
            <li class="filter_shape" id="filter_shapeRAD" title="RAD" onclick="filter_shape('RAD')">
              <span class="shapeiconcontainer" title="radiant" style="background-image:url(/_images/constant/ico-stones.png); background-position:-44px 0px;"></span>
            </li>
            <li class="filter_shape" id="filter_shapeCU" title="CU" onclick="filter_shape('CU')">
              <span class="shapeiconcontainer" title="cushion" style="background-image:url(/_images/constant/ico-stones.png); background-position:-198px 0px;"></span>
            </li>
            <li class="filter_shape defined">
              <input id="shape_defined" type="text" placeholder="其它，请添写" />
            </li>
          </ul>
        </div>
        <!-- color -->
        <div class="filter_line" id="filterline-color">
          <span class="filter_title">颜色: </span>
          <ul class="filter_color_outer">
            <li class="filter_color" id="filter_colorD" title="D" onclick="filter_color('D')">D</li>
            <li class="filter_color" id="filter_colorE" title="E" onclick="filter_color('E')">E</li>
            <li class="filter_color" id="filter_colorF" title="F" onclick="filter_color('F')">F</li>
            <li class="filter_color" id="filter_colorG" title="G" onclick="filter_color('G')">G</li>
            <li class="filter_color" id="filter_colorH" title="H" onclick="filter_color('H')">H</li>
            <li class="filter_color" id="filter_colorI" title="I" onclick="filter_color('I')">I</li>
            <li class="filter_color" id="filter_colorJ" title="J" onclick="filter_color('J')">J</li>
            <li class="filter_color" id="filter_colorK" title="K" onclick="filter_color('K')">K</li>
            <li class="filter_color" id="filter_colorL" title="L" onclick="filter_color('L')">L</li>
            <li class="filter_color" id="filter_colorM" title="M" onclick="filter_color('M')">M</li>
            <li class="filter_color" id="filter_colorN" title='N' onclick="filter_color('N')">N</li>
            <li class="filter_color" id="filter_colorO" title='O' onclick="filter_color('O')">O</li>
            <li class="filter_color" id="filter_colorP" title='P' onclick="filter_color('P')">P</li>
            <li class="filter_color" id="filter_colorQ" title='Q' onclick="filter_color('Q')">Q</li>
            <li class="filter_color" id="filter_colorR" title='R' onclick="filter_color('R')">R</li>
            <li class="filter_color" id="filter_colorS" title='S' onclick="filter_color('S')">S</li>
            <li class="filter_color" id="filter_colorT" title='T' onclick="filter_color('T')">T</li>
            <li class="filter_color" id="filter_colorU" title='U' onclick="filter_color('U')">U</li>
            <li class="filter_color" id="filter_colorV" title='V' onclick="filter_color('V')">V</li>
            <li class="filter_color" id="filter_colorW" title='W' onclick="filter_color('W')">W</li>
            <li class="filter_color" id="filter_colorX" title='X' onclick="filter_color('X')">X</li>
            <li class="filter_color" id="filter_colorY" title='Y' onclick="filter_color('Y')">Y</li>
            <li class="filter_color" id="filter_colorZ" title='Z' onclick="filter_color('Z')">Z</li>
            <li class="filter_color fancycolor defined">
              <input id="fancycolorinput" type="text" placeholder="彩色或其它，请添写" style="padding-left:12px;" />
            </li>
          </ul>
        </div>
        <!-- clarity -->
        <div class="filter_line" id="filterline-clarity">
          <span class="filter_title">净度</span>
          <ul class="filter_clarity_outer">
            <li class="filter_clarity" id="filter_clarityFL" title="FL" onclick="filter_clarity('FL')">FL</li>
            <li class="filter_clarity" id="filter_clarityIF" title="IF" onclick="filter_clarity('IF')">IF</li>
            <li class="filter_clarity" id="filter_clarityVVS1" title="VVS1" onclick="filter_clarity('VVS1')">VVS1</li>
            <li class="filter_clarity" id="filter_clarityVVS2" title="VVS2" onclick="filter_clarity('VVS2')">VVS2</li>
            <li class="filter_clarity" id="filter_clarityVS1" title="VS1" onclick="filter_clarity('VS1')">VS1</li>
            <li class="filter_clarity" id="filter_clarityVS2" title="VS2" onclick="filter_clarity('VS2')">VS2</li>
            <li class="filter_clarity" id="filter_claritySI1" title="SI1" onclick="filter_clarity('SI1')">SI1</li>
            <li class="filter_clarity" id="filter_claritySI2" title="SI2" onclick="filter_clarity('SI2')">SI2</li>
            <li class="filter_clarity" id="filter_clarityI1" title="I1" onclick="filter_clarity('I1')">I1</li>
            <li class="filter_clarity" id="filter_clarityI2" title="I2" onclick="filter_clarity('I2')">I2</li>
            <li class="filter_clarity" id="filter_clarityI3" title="I3" onclick="filter_clarity('I3')">I3</li>
            <li class="filter_clarity defined">
              <input id="clarity_defined" type="text" placeholder="其它，请添写" />
            </li>
          </ul>
        </div>
        <!-- cut -->
        <div class="filter_line" id="filterline-cut">
          <span class="filter_title">切工</span>
          <ul class="filter_cut_outer">
            <li class="filter_cut" id="filter_cutEX" title="EX" onclick="filter_cut('EX')">EX</li>
            <li class="filter_cut" id="filter_cutVG" title="VG" onclick="filter_cut('VG')">VG</li>
            <li class="filter_cut" id="filter_cutG" title="G" onclick="filter_cut('G')">G</li>
            <li class="filter_cut" id="filter_cutF" title="F" onclick="filter_cut('F')">F</li>
            <li class="filter_cut" id="filter_cut-" title="-" onclick="filter_cut('-')">Not Available</li>
            <li class="filter_cut defined">
              <input id="cut_defined" type="text" placeholder="其它，请添写" />
            </li>
          </ul>
        </div>
        <!-- polish -->
        <div class="filter_line" id="filterline-polish">
          <span class="filter_title">抛光</span>
          <ul class="filter_polish_outer">
            <li class="filter_polish" id="filter_polishEX" title="EX" onclick="filter_polish('EX')">EX</li>
            <li class="filter_polish" id="filter_polishVG" title="VG" onclick="filter_polish('VG')">VG</li>
            <li class="filter_polish" id="filter_polishG" title="G" onclick="filter_polish('G')">G</li>
            <li class="filter_polish" id="filter_polishF" title="F" onclick="filter_polish('F')">F</li>
            <li class="filter_polish defined">
              <input id="polish_defined" type="text" placeholder="其它，请添写" />
            </li>
          </ul>
        </div> 
        <!-- symmetry -->
        <div class="filter_line" id="filterline-symmetry">
          <span class="filter_title">对称性</span>
          <ul class="filter_sym_outer">
            <li class="filter_sym" id="filter_symEX" title="EX" onclick="filter_sym('EX')">EX</li>
            <li class="filter_sym" id="filter_symVG" title="VG" onclick="filter_sym('VG')">VG</li>
            <li class="filter_sym" id="filter_symG" title="G" onclick="filter_sym('G')">G</li>
            <li class="filter_sym" id="filter_symF" title="F" onclick="filter_sym('F')">F</li>
            <li class="filter_sym defined">
              <input id="symmetry_defined" type="text" placeholder="其它，请添写" />
            </li>
          </ul>
        </div>
        <!-- fluo -->
        <div class="filter_line" id="filterline-fluo">
          <span class="filter_title">荧光</span>
          <ul class="filter_fluo_outer">
            <li class="filter_fluo" id="filter_fluoVST" title="VST" onclick="filter_fluo('VST')">very strong</li>
            <li class="filter_fluo" id="filter_fluoSTG" title="STG" onclick="filter_fluo('STG')">strong</li>
            <li class="filter_fluo" id="filter_fluoMED" title="MED" onclick="filter_fluo('MED')">medium</li>
            <li class="filter_fluo" id="filter_fluoFNT" title="FNT" onclick="filter_fluo('FNT')">faint</li>
            <li class="filter_fluo" id="filter_fluoNON" title="NON" onclick="filter_fluo('NON')">none</li>
            <li class="filter_fluo defined">
              <input id="fluo_defined" type="text" placeholder="其它，请添写" />
            </li>  
          </ul>
        </div>
        <!-- lab -->  
        <div class="filter_line" id="filterline-lab">
          <span class="filter_title">证书</span>
          <ul class="filter_certi_outer">
            <li class="filter_certi" id="filter_certiIGI" title="IGI" onclick="filter_certi('IGI')">IGI</li>
            <li class="filter_certi" id="filter_certiHRD" title="HRD" onclick="filter_certi('HRD')">HRD</li>
            <li class="filter_certi" id="filter_certiGIA" title="GIA" onclick="filter_certi('GIA')">GIA</li>  
            <li class="filter_certi defined">
              <input id="certi_defined" type="text" placeholder="其它，请添写" />
            </li>
          </ul>
        </div>
        <form action="" method="post" id="theform">
          <input type="hidden" name="carat" id="fordb_carat" />
          <input type="hidden" name="shape" id="fordb_shape" />
          <input type="hidden" name="color" id="fordb_color" />
          <input type="hidden" name="clarity" id="fordb_clarity" />
          <input type="hidden" name="cut" id="fordb_cut" />
          <input type="hidden" name="polish" id="fordb_polish" />
          <input type="hidden" name="symmetry" id="fordb_symmetry" />
          <input type="hidden" name="fluo" id="fordb_fluo" />
          <input type="hidden" name="lab" id="fordb_lab" />
          <div class="filter_line">
            <p style="margin-bottom:15px;"><label>证书编号:</label>
              <input type="text" class="inputfield" name="certi_num" id="certi_num" style="background-color:#FFF;" /></p>
            <p><label>钻石最终付款<span style="color:#C00">(必填)</span>:</label>
              <input type="number" class="inputfield" name="amount" id="amount" /> 欧元&nbsp; &nbsp; 
              <input type="checkbox" name="d_no_point" id="d_no_point" />
              <strong style="font-size:16px; color:#C60;">该钻石不能积分, 只做纪录</strong>
            </p>
            <p><label>空托最终付款:</label>
              <input type="number" class="inputfield" name="acc_amount" id="acc_amount" /> 欧元 &nbsp; &nbsp; 
              <input type="checkbox" name="acc_point" id="acc_point" />
              <strong style="font-size:16px; color:#C60;">空托算做积分</strong>
            </p>
            <p><label>购买时实际折扣:</label> 
              <select name="actual_price_ratio">
              <!-- <option value="1"<?php if($user_price_ratio==1){ echo ' selected="selected"'; } ?>>无折扣</option>
              <option value="0.98"<?php if($user_price_ratio==0.98){ echo ' selected="selected"'; } ?>>98折</option>
              <option value="0.95"<?php if($user_price_ratio==0.95){ echo ' selected="selected"'; } ?>>95折</option>
              <option value="0.93"<?php if($user_price_ratio==0.93){ echo ' selected="selected"'; } ?>>93折</option>
              <option value="0.9"<?php if($user_price_ratio==0.9){ echo ' selected="selected"'; } ?>>9折</option>
              <option value="0.87"<?php if($user_price_ratio==0.87){ echo ' selected="selected"'; } ?>>87折</option>
              <option value="0.85"<?php if($user_price_ratio==0.85){ echo ' selected="selected"'; } ?>>85折</option>
              <option value="0.83"<?php if($user_price_ratio==0.83){ echo ' selected="selected"'; } ?>>83折</option> -->
              </select>
            </p>
            <p><label>购买日期<span style="color:#C00">(必填)</span>:</label>
              <input type="datetime" class="inputfield" name="thetime" id="thetime" placeholder="点击选择日期" /></p>
            <p><button type="button" id="submitbtn" onclick="submitTheForm()">保存</button></p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
<script src="./AddBuyRecords.js"></script>
<style src="./AddBuyRecords.css" scoped></style>
