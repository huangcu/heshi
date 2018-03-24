<template>
<div>
  <div v-for="(diamond,index) in diamonds" id="dia-data-box" :key="diamond.id">
    <div class="dia-piece-box">
      <div v-if="diamond.supplier==='BEYOU'">
        <span v-if="diamond.profitable==='NO'" class="r-indi non-profit"><span class="glyphicon glyphicon-certificate" style="color:#ff5;"></span><span class="reco-txt">特惠</span></span>
        <span v-else class="r-indi"><span class="glyphicon glyphicon-thumbs-up"></span><span class="reco-txt">推荐</span></span>
      </div>
      <div :class="index" class="generalinfobox">
        <span class="shapedesc-box">
          <!-- :style="{ 'background-position:'' + imagePosition + 'px 0px; background-image': 'url(' + getImage + ')' }" -->
          <span class="shape-txt">{{ detail_forDiamond_byShape(diamond.shape, "NAMECN") }}</span>
          <span class="shapeiconcontainer" :style="{ 'background-position': ' '+ imagePosition + 'px 0px'}"></span>
          <!-- <span class="shapeiconcontainer" :style="{ 'background-image': 'url(' + getImage + ');' + 'background-position': + imagePosition + 'px 0px' }"></span> -->
        </span>
        <span class="valuetxt value_carat">
          <span class="thevalue"> {{ diamond.carat }}</span>克拉
        </span>
        <span class="valuetxt value_color">颜色 <span class="thevalue"><span v-html="diamondColor(diamond.color)"></span></span>
        </span>
        <span class="valuetxt value_clarity">净度 <span class="thevalue"> {{ diamond.clarity }}</span>
        </span>
        <span class="valuetxt value_certificate">证书 <span class="thevalue">{{ diamond.grading_lab }}</span>
        </span>
        <span class="valuetxt value_cut">切工 <span class="thevalue">{{ diamondCutGrade(diamond.cut_grade) }}</span>
        </span>
        <span class="valuetxt value_polish">抛光 <span class="thevalue">{{ diamond.polish}}</span>
        </span>
        <span class="valuetxt value_symmetry">对称 <span class="thevalue">{{ diamond.symmetry }}</span>
        </span>
        <span v-if= "userprofile!==''" class="valuetxt value_price">
          <div v-if="userprofile.user_level ==-1">
            参考价{{ vat_status_txt }}
            <span  class="thevalue">{{ DollarToEuro(diamond.price_retail*userprofile.user_discount) }}</span> €
          </div>
          <div v-else>
            参考价{{ vat_status_txt }}
            <span class="thevalue">{{ DollarToEuro(diamond.price_retail*userprofile.user_discount) }}</span> €
          </div>
        </span>
        <span v-else>
          <a href="/login">登录查看价格</a>
        </span>
        <span class="detail-btn" v-on:click="openDiaDetail(diamond.id)">
          <span class="glyphicon glyphicon-eye-open"></span> 详情
        </span>
      </div><!-- end generalinfobox -->
      <div :id="'d_'+index" class="details" :class ="'details_'+diamond.id"><!-- detail box -->
        <div class="detailcol detailcol1"><!-- TODO handle fancycolor -->
          <img v-if="diamond.color==='FANCY'" class="dia-demo-pic" src="/_images/constant/fancycolor.jpg"/>
          <img v-else class="dia-demo-pic" :src="getImage"/>
          <img v-if="diamond.grading_lab==='HRD'" class="certi-demo-pic" :src="hrdGradingLabImage" />
          <img v-else-if="diamond.grading_lab==='GIA'" class="certi-demo-pic" :src="gidGradingLabImage" />
          <img v-else-if="diamond.grading_lab==='IGI'" class="certi-demo-pic" :src="igiGradingLabImage" />
        </div>
        <div class="detailcol detailcol2">
          <span>荧光强度: {{ diamond.fluorescence_intensity }}</span>
          <span>所在地: {{ diamondPlace(diamond.country) }}</span>
          <span>证书编号: {{ diamond.certificate_number }}</span>
          <span>
            <a v-if="diamond.grading_lab==='HRD'" class="certi_linker" target="_blank" :href="'https://my.hrdantwerp.com/?L=&record_number=' + diamond.certificate_number+'&certificatetype=MC'">查看证书
              <span class="glyphicon glyphicon-new-window"></span>
            </a>
            <a v-else-if="diamond.grading_lab==='GIA'" class="certi_linker" target="_blank" :href="'http://www.gia.edu/cs/Satellite?pagename=GST%2FDispatcher&childpagename=GIA%2FPage%2FReportCheck&c=Page&cid=1355954554547&reportno='+ diamond.certificate_number">查看证书 <span class="glyphicon glyphicon-new-window"></span></a>
            <a v-else-if="diamond.grading_lab==='IGI'" class="certi_linker" target="_blank" :href="'http://www.igiworldwide.com/verify.php?r=' + diamond.certificate_number">查看证书 <span class="glyphicon glyphicon-new-window"></span></a>
          </span>
          <p v-if="diamond.recommend_words!==''" class="commentbox">
            <label class="commentwords-label">
              <span class="glyphicon glyphicon-asterisk"></span> 公司检验结论
            </label>
            <span class="thecomment">{{ diamond.recommend_words }}</span>
          </p><!-- TODO handle images -->
          <label v-if="diamond.images!=0" class="realpics-label"><span class="glyphicon glyphicon-camera"></span> 实拍照片(点击放大)</label>
          <p v-for="(image, imageNumber) in diamond.images" class="dia-pics-container" :key="imageNumber">
            <a :href="'/_images/diamond/'+image" class="biggerimglinker fancybox" :rel="diamond.id"><!-- TODO handle thumbs images(small image) -->
            <span class="dia-detail-pic-holder" :style="'background-image:url(/_images/diamond/thumbs/' + image"></span>
            </a>
          </p>
          <p v-if="diamond.extra_words!==''" class="commentbox">
            <label class="commentwords-label extrwords">备注</label>
            <span class="thecomment extrwords"> {{ diamond.extra_words }}</span>
          </p>
        </div>
        <div class="detailcol detailcol3"><!-- TODO handle diamonds page with query in url -->
          <span>钻石ID: <a target="_blank" style="color:#000; text-decoration:underline;" :href="'/product/diamond?&ref='+diamond.id">{{ diamond.stock_ref }}</a></span>
          <div v-if="diamond.profitable==='YES'">
            <div v-if="userprofile!==''">
              <div v-if="userprofile.user_level===-1">
                参考价{{ vat_status_txt}} {{ DollarToEuro(diamond.price_retail*reference_price_ratio) }} €
              </div>
              <div v-else>
                <div v-if="userprofile.user_type==='AGENT'">
                  <span class="btnforprice">合适价{{ vat_status_txt }}: {{ DollarToEuro(diamond.price_retail) }}<span class="currencyunit">欧元</span> ({{ priceretail(diamond.price_retail) }}<span class="currencyunit">美元</span>, {{ DollarToYuan(diamond.retail_price) }}<span class="currencyunit">人民币</span>)</span>
                  <span class="btnforprice">您的代理价{{ vat_status_txt }}: <span class="specialprice">{{ DollarToEuro(priceForAgent(agentLevel, diamond.price_retail)) }} </span><span class="currencyunit">欧元</span> ({{ priceForAgent(agentlevel, diamond.price_retail) }}<span class="currencyunit">美元</span>, {{ DollarToYuan(priceForAgent(agentLevel, diamond.price_retail)) }}<span class="currencyunit">人民币</span>)</span>
                </div>
                <div v-else>
                  <span class="btnforprice">合适价{{ vat_status_txt }}: {{ DollarToEuro(diamond.price_retail) }}<span class="currencyunit">欧元</span> ({{ priceretail(diamond.price_retail) }}<span class="currencyunit">美元</span>, {{ DollarToYuan(diamond.retail_price) }}<span class="currencyunit">人民币</span>)</span>
                  <span class="btnforprice">参考价{{ vat_status_txt }}: <span class="specialprice">{{ DollarToEuro(priceForAccount(accountLevel, diamond.price_retail)) }}</span><span class="currencyunit">欧元</span> ({{ priceForAccount(agentlevel, diamond.price_retail) }}<span class="currencyunit">美元</span>, {{ DollarToYuan(priceForAccount(accountLevel, diamond.price_retail)) }}<span class="currencyunit">人民币</span>)</span>
                </div>
              </div>
            </div>
            <div v-else>
              <span class="btnforprice">价格:  <a href="/login">请先登录／注册</a></span>
            </div>

            <div v-if="diamond.status==='SOLD'">
              <button type="button" class="btnfororder disabledbtn itemordered" title="chosenitem" disabled="disabled">
                <span class="glyphicon glyphicon-globe"></span> 该钻已售
              </button>
            </div>
            <div v-else>
              <!-- todo in dia_items_shoppinglist_confirmed list -->
              <div v-if="in_dia_items_shoppinglist_confirmed">
                <p class="actionbuttonsbox">
                  <button type="button" class="btnfororder disabledbtn itemordered" title="chosenitem" disabled="disabled">
                    <span class="glyphicon glyphicon-ok"></span> {{ inthelistword_ordered }}
                  </button>
                </p>
              </div>
              <!-- todo in dia_items_shoppinglist_confirmed list -->
              <div v-else-if="in_dia_items_shoppinglist">
                <p  class="actionbuttonsbox interested-btn-box" :id="'interested-btn-box_'+diamond.stock_ref">
                  <button type="button" class="btnfororder disabledbtn itemordered" title="chosenitem" disabled="disabled">
                    <span class="glyphicon glyphicon-ok"></span>  {{ inthelistword_saved }}
                  </button>
                </p>
                <p class="actionbuttonsbox ordered-bt-box" :id="'ordered-btn-box_'+ diamond.stock_ref">
                  <button type="button" class="btnfororder" :title="diamond.stock_ref" v-on:click="makeorder_confirmed(diamond.id)"><span v-html="the_order_confirm_btn"></span></button>
                    <!-- todo addringfordiamond.php list -->
                    <a class="addringfordiamondbtn confirmed" :id="'addringfordiamondbtn_confirmed_'+diamond.stock_ref" :href="'/product/ringfordiamond?ref='+diamond.stock_ref+'&ordered=yes'">选空托</a>
                    <!-- todo shoppinglist-confirmed.php list -->
                    <a class="checkmyorderlistbtn" :id="'checkmyorderlistbtn_confirmed_'+ diamond.stock_ref" href="/shoppinglist-confirmed">{{ checklist_ordered_txt }}</a>
                  <span v-html="the_order_btn_expl"></span>
                </p>
              </div>
              <div v-else>
                <p class="actionbuttonsbox interested-btn-box" :id="'interested-btn-box_'+ diamond.stock_ref">
                  <button type="button" class="btnfororder interested" :title="diamond.stock_ref" v-on:click="makeorder(diamond.stock_ref)"><span v-html="the_order_btn"></span></button>
                  <a class="addringfordiamondbtn interested" :id="'addringfordiamondbtn_'+diamond.stock_ref" :href="'/product/ringfordiamond?ref='+ diamond.stock_ref">选空托</a>
                  <a class="checkmyorderlistbtn" :id="'checkmyorderlistbtn_'+diamond.stock_ref" href="/shoppinglist.php">{{ checklisttxt }}</a>
                  <span v-if="the_order_btn_expl!=''" v-html="the_order_btn_expl"></span>
                </p>
                <p class="actionbuttonsbox ordered-bt-box" :id="'ordered-btn-box_'+ diamond.stock_ref">
                  <button type="button" class="btnfororder ordered" :title="diamond.stock_ref" v-on:click="makeorder_confirmed( diamond.stock_ref)"><span v-html="the_order_confirm_btn"></span></button>
                  <a class="addringfordiamondbtn confirmed" :id="'addringfordiamondbtn_confirmed_'+  diamond.stock_ref" :href="'/product/ringfordiamond?ref='+  diamond.stock_ref + '&ordered=yes'">选空托</a>
                  <a class="checkmyorderlistbtn" :id="'checkmyorderlistbtn_confirmed_'+ diamond.stock_ref" href="/shoppinglist-confirmed.php"> {{ checklist_ordered_txt }}</a>
                  <span v-if="the_order_btn_expl!=''" v-html="the_order_confirm_btn_expl"></span>
                </p>
              </div>
            </div>
          </div><!-- end of profitable -->
          <!-- not profitable -->
          <div v-else>
            <a class="specialofferlinker" href="/product/diamondoftheweek">
            <span class="glyphicon glyphicon-gift"></span> 特惠详情 &raquo;</a>
          </div>
          <p class="closebtnplaceholder"></p>
        </div>
        <br class="clear" />
        <button class="closeDetailBtn" v-on:click="closeDiaDetail(diamond.id)">
          <span class="glyphicon glyphicon-arrow-up"></span> 收起
        </button>
      </div><!-- end details -->
    </div><!-- end dia-pice-box -->
  </div><!-- end of for loop -->
    <div id="howmanyrecords" style="display:none;">{{ diamonds.length }}
    </div>
  </div>
</template>
<script src='./DiamondsData.js'></script>
