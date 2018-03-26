<template>
  <div>
    <vue-title title="合适钻石 | BEYOU DIAMOND | 推荐详情"></vue-title>
    <div class="contentmainwrapper">
    <h2 class="jewelry-title">合适推荐钻石</h2>
    <div v-if="diamond.status !== 'AVAILABLE'">
      <p>抱歉，该钻石已经售出。</p>
    </div>
    <div v-else class="jewely-description">
      <p>商品编号: <strong>{{ diamond.stock_ref }}></strong></p>
      <span class="dia-description">{{ diamond.carat }}克拉</span>
      <span class="dia-description">{{ diamondShape(diamond.shape)}}</span>
      <span class="dia-description">{{ diamond.color}}色</span>
      <span class="dia-description">净度: {{ diamond.clarity}}</span>
      <span class="dia-description">切工: {{ diamond.cut_grade}}／ 抛光: {{ diamond.polish}}／ 对称: {{ diamond.symmetry}}</span>
      <span class="dia-description">荧光强度: {{ diamond.fluorescence_intensity }}</span>
      <a v-if="diamond.grading_lab==='HRD'" class="certi_linker" target="_blank" :href="'https://my.hrdantwerp.com/?L=&record_number=' + diamond.certificate_number+'&certificatetype=MC'">查看证书
        <span class="glyphicon glyphicon-new-window"></span>
      </a>
      <a v-else-if="diamond.grading_lab==='GIA'" class="certi_linker" target="_blank" :href="'http://www.gia.edu/cs/Satellite?pagename=GST%2FDispatcher&childpagename=GIA%2FPage%2FReportCheck&c=Page&cid=1355954554547&reportno='+ diamond.certificate_number">查看证书 <span class="glyphicon glyphicon-new-window"></span></a>
      <a v-else-if="diamond.grading_lab==='IGI'" class="certi_linker" target="_blank" :href="'http://www.igiworldwide.com/verify.php?r=' + diamond.certificate_number">查看证书 <span class="glyphicon glyphicon-new-window"></span></a>
      <div v-if="diamond.recommend_words!==''" id="reco-words">
        {{ diamond.recommend_words }}
      </div>
    </div>
    <div class="pricebox">
      <div v-if="userprofile!==''">
        <div v-if="userprofile.user_type === 'AGENT'">
          <span class="price">合适价: {{priceforagent(userprofile.agent.level, diamond.retail_price)}} 欧元</span>
          <span class="price-special"> 您的代理价: {{priceforagent(userprofile.agent.level, diamond.retail_price)}} 美元</span>
        </div>
        <div v-else-if="userprofile.user_level === -1">
          <span class="price">合适价:{{ priceforaccount(userprofile.user_level, diamond.retail_price)}} 欧元</span>
          <span class="price-special">您的会员价: {{ priceforaccount(userprofile.user_level, diamond.retail_price)}} 美元</span>
        </div>
        <div v-else>
          <span class="price">合适价: {{ diamond.retail_price }}美元</span>
        </div>
      </div>
      <div v-else>
        <span class="btnforprice">价格:  <a href="/login">请先登录／注册</a></span>
      </div>
    </div>
    <div v-if="diamond.images!== null" class="jewelry-pics">
      <img class="jewelry-img" v-if="diamond.images.length!==0" v-for="(image, index) in diamond.images" :key="index" v-bind:src="image"/>
    </div>
      <p class="choosebtnbox">
        <a class="choosebtn" :href="'shoppinglist-confirmed.php?type=diamond&id='+diamond.id+'&action=add'">
          <span class="glyphicon glyphicon-th-list"></span> 加入预订单</a>
      </p>
    </div>
    <p class="gobackbtnbox">
    <span class="gobackindisign"><span class="glyphicon glyphicon-arrow-left"></span> 返回:</span>
      <a class="gobackbtn" href="/product/recommenddiamonds"><span class="glyphicon glyphicon-th"></span> 全部推荐钻石 </a>
    </p>
  </div>
</template>
<style src='./Diamond.css'></style>
<script src='./Diamond.js'></script>
