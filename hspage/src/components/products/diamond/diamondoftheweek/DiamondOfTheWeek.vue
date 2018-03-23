<template>
  <div>
    <vue-title title="本周特价钻 | 合适钻石"></vue-title>
    <div class="contentmainwrapper">
      <img id="bigsaleicon" :src="bigSaleImage" />
      <h2 class="jewelry-title">本周特惠钻石</h2>
      <div v-if="this.diamonds=== null || this.diamonds.length ===0">
        <p>抱歉，本周没有特惠钻石，请下周再来。</p>
      </div>
      <div v-else class="jewely-description">
        <p>商品编号: <strong>{{ diamonds[0].stock_ref }}</strong></p>
        <p v-if="diamonds[0].ordered_by!=='' && diamonds[0].ordered_by===userprofile.id" class="spreadwords">恭喜！您已经预定了这颗钻石！</p>
        <div class="jewelry-pics">
          <a v-if="diamonds[0].images !== null" v-for="(image, index) in diamonds[0].images" :key="index" class="bigpiclinker fancybox" rel="diapic" :href="'/_images/diamond/'+image" style="display:inline-block;">
            <img class="jewelry-img" :src="'/_images/diamond/'+image" />
          </a>
          <img v-if="diamonds[0].status!=='AVAILABLE'" :src="soldOutImage" id="soldouticon" />
        </div>
        <span class="dia-description"> {{ diamonds[0].carat }}克拉</span>
        <span class="dia-description"> {{ diamondShapeTxtPic(diamonds[0].shape, 'TXT') }}</span>
        <span class="dia-description"> {{ diamonds[0].color }}色</span>
        <span class="dia-description">净度: {{ diamonds[0].clarity }}</span>
        <span class="dia-description">切工:  {{ diamonds[0].cut_grade }} ／ 抛光:  {{ diamonds[0].polish }} ／ 对称:  {{ diamonds[0].symmetry }}</span>
        <span class="dia-description"> {{ diamondFluo(diamonds[0].fluorescence_intensity) }}</span>
        <p class="certi-linker">
          <span class="dia-description"> {{ diamonds[0].grading_lab }}</span>
          <a class="certi_linker" target="_blank" :href="diamonds[0].certificate_link">查看证书: <span class="glyphicon glyphicon-new-window"></span>
          </a>
        </p>
        <div v-if="diamonds[0].recommand_word" id="reco-words">{{ diamonds[0].recommand_word }}</div>
        <div class="pricebox">
          <div v-if="promo_type=='FREE_ACC'">
            <span class="title-free_acc">买特价钻石，而且免费送空托</span>
            <span class="price-special">钻石价格: {{ promo_price }} 欧元</span>
          </div>
          <div v-else-if="promo_type=='DISCOUNT'">
            <span class="price">合适价: <span style="text-decoration:line-through; color:#666;">{{ diamonds[0].price_retail }} 欧元</span>
            </span>
            <span class="price-special">特惠价:{{ promo_price }} 欧元</span>
          </div>
          <div v-else>
            <span class="specialofferexpl">
              <span class="glyphicon glyphicon-info-sign"></span>注: 该特惠商品已是底价，不能再使用折扣
            </span>
          </div>
        </div>

        <div v-if="promo_type=='FREE_ACC'" class="free_acc_box">
          <h3>免费附送精美空托，请任选一款</h3>
          <div class="free_acc_box_inner">
            <div v-if="this.jewelrys === null || this.jewelrys.length === 0">
              <p>没有找到适合的空托</p>
            </div>
            <div v-else>
              <div v-for="(jewelry, index) in jewelrys" :key="index" class="jewelrybox" :id="'jewelrybox_'+crr_id">
                <a v-if="jewelry.images !== null" v-for="(image, imageIndex) in jewelry.images" :key="imageIndex" class="seedetailbtn-big demo-box fancybox" :rel="crr_id" :href="'/_images/jewelry/'+image">
                  <span class="imageholder" :style="'background-image:url(/_images/jewelry/thumbs/'+image"></span>
                  <span class="jewelryname">{{ jewelry.name }}</span>
                  <!-- <span class="jewelryprice">原价:<?php echo $crr_price; ?> 欧元</span> -->
                </a>
                <!-- <a class="fancybox hiddenforpic" href="/_images/jewelry/<?php echo $crr_img2; ?>" rel="<?php echo $crr_id; ?>"></a> -->
                <p class="actionbox">
                  <button class="choosebtn" :id="'choosebtn_'+ jewelry.id" v-on:Click="chooseThisJew(jewelry.id)">
                    <span class="glyphicon glyphicon-gift"></span> 选择</button>
                </p>
              </div>
            </div>
            <br class="clear" />
          </div><!-- free_acc_box_inner -->
          <p v-if="diamonds[0].ordered_by!=='' && diamonds[0].ordered_by===userprofile.id" class="spreadwords">恭喜！您已经预定了这颗钻石！</p>
          <p v-if="diamonds[0].ordered_by!=='' && diamonds[0].ordered_by!==userprofile.id" class="spreadwords">这颗钻石已经被手快的朋友抢走了。别灰心，我们每周一都会推出一颗新的特价钻！</p>
          <p v-if="diamonds[0].ordered_by===''" class="spreadwords">无意购买该钻？将这个好消息分享给你的朋友 <span class="glyphicon glyphicon-bullhorn"></span> ，或许他们正在寻找！
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
<script src='./DiamondOfTheWeek.js'></script>
<style src='./DiamondOfTheWeek.css' scoped></style>
