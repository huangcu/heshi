<template>
  <div>
    <vue-title title="合适钻石 | BEYOU DIAMOND | 合适推荐钻石"></vue-title>
    <div class="contentmainwrapper">
      <h2 class="pagetitle">合适推荐钻石</h2>
      <p class="title-expl">以下皆是合适精挑细选的高性价比钻石， 均有现货</p>
      <div class="sortingbtncontainer">
        <form v-on:click.prevent="sortOrderDiamonds" action="" method="get" class="sortingform">
          <select v-model="sortBy" name="sort">
            <option value="carat">钻石重量</option>
            <option value="price">钻石价格</option>
          </select>
          <select v-model="orderBy" name="order">
            <option value="up">由低到高 &darr;</option>
            <option value="down">由高到低 &uarr;</option>
          </select>
          <!-- <input type="submit" value="排序" id="orderbtn" /> -->
        </form>
      </div>
      <div class="jewelry-container-box for-reco-dias">
        <div v-for="(diamond,index) in diamonds" :key="index" class="jewelrybox complete reco-dia">
          <a class="seedetailbtn-big demo-box jew-pic-reco-holder" :href="'/diamond?ref='+diamond.id">
            <span v-if="diamond.images" class="imageholder" :style="'background-image:url(/_images/diamond/thumbs/'+diamond.images[0]"></span>
            <span v-else class="imageholder" :style="{ 'background-image': 'url(' + getConstantImage('noimage.jpg') + ')' }"></span>
          </a>
          <span class="jewelryname reco-jew-title">{{ diamond.carat}} 克拉, {{ diamond.color }} 色，{{ diamond.clarity }}</span>
          <span class="jewelrypricereco-jew-price"> {{ diamond.price_retail }} € </span>
        </div>
        <br class="clear" />
      </div><!-- class jewelry-container-box -->
    </div><!-- end of contentmainwrapper -->
  </div>
</template>
<style src='./RecommendDiamonds.css' scoped></style>
<script src='./RecommendDiamonds.js'></script>
