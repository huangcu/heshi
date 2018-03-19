<template>
  <div class="navbar-wrapper">
    <div class="container">
      <a class="logobox" href="/">
        <img id="site-logo" alt="比利时钻石" src="../../_images/constant/beyoudiamond-logo.png"/>
      </a>
      <div class="small-buttons-container">
        <a id="shoppingcartbtn" href="/shoppingList">
          <span class="glyphicon glyphicon-heart" aria-hidden="true"></span> 我的收藏
          <span id="num-shoppinglist-container" style="display:none;">
            <span id="num-shoppinglist">{{ interestedItems.length }}</span>
          </span>
        </a>
        <a v-if="accountID != ''" id="shoppingcartbtn-confirmed" href="/shoppinglistconfirmed">
          <span class="glyphicon glyphicon-list" aria-hidden="true"></span> 预定单
        </a>
        <a v-if="accountID != ''" id="shoppingcartbtn-confirmed" href="/orders">
          <span class="glyphicon glyphicon-th-list" aria-hidden="true"></span> 已下单
        </a>
        <a v-if="accountID === ''" id="myaccountbtn" href="/login">
          <span class="glyphicon glyphicon-user" aria-hidden="true"></span>
            登录 | 注册
        </a>
        <a v-else id="myaccountbtn" :href="'/myaccount?id=' + accountID">
          <span class="glyphicon glyphicon-user" aria-hidden="true"></span>
            我的账户
        </a>
        <a id="mypresentbtn" href="/mypresent">
          <img id="gifticon" src="../../_images/constant/gift-icon.png" /> 我的礼包
        </a>
      </div>
      <div class="small-buttons-container-mobile">
        <a id="shoppingcartbtn-mobile" href="/shoppingList">
          <span class="glyphicon glyphicon-heart" aria-hidden="true"></span> 收藏
          <span id="num-shoppinglist-container" style="display:none;">
            <span id="num-shoppinglist">{{ interestedItems.length }}</span>
          </span>
        </a>
        <a v-if="accountID != ''" id="shoppingcartbtn-confirmed-mobile" href="/shoppinglistconfirmed">
          <span class="glyphicon glyphicon-list" aria-hidden="true"></span> 预定单
        </a>
        <a v-if="accountID != ''" id="shoppingcartbtn-confirmed-mobile" href="/orders">
          <span class="glyphicon glyphicon-th-list" aria-hidden="true"></span> 已下单
        </a>
        <a v-if="accountID === ''" id="myaccountbtn-mobile" href="/login">
          <span class="glyphicon glyphicon-user" aria-hidden="true"></span>
            登录 | 注册
        </a>
        <a v-else  id="myaccountbtn-mobile" :href="'/myaccount?id=' + accountID">
          <span class="glyphicon glyphicon-user" aria-hidden="true"></span>
            账户
        </a>
        <a id="mypresentbtn-mobile" href="/mypresent">
          <img id="gifticon" src="../../_images/constant/gift-icon.png" /> 礼包
        </a>
      </div>
      <nav class="navbar navbar-expand-lg navbar-light bg-light">
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#nav-content" aria-controls="nav-content" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div id="nav-content" class="navbar-collapse collapse">
          <ul class="navbar-nav">
            <li class="nav-item"><a class="text-sm-center nav-link "  v-bind:class="{ active: isIndexPage }" href="/"><span class="glyphicon glyphicon-home"></span>首页</a></li>
            <li class="nav-item"><a class="text-sm-center nav-link "  v-bind:class="{ active: isDiamond }" href="/product/diamonds">钻石</a></li>
            <li class="nav-item"><a class="text-sm-center nav-link "  v-bind:class="{ active: isJewelryNeedMounting }" href="/product/jewelry?class=mounting">空托</a></li>
            <li class="nav-item"><a class="text-sm-center nav-link "  v-bind:class="{ active: isJewelry }" href="/product/jewelry">首饰</a></li>
            <li class="nav-item"><a class="text-sm-center nav-link "  v-bind:class="{ active: isGem }" href="/product/colored-gems">彩宝</a></li>
            <li class="nav-item"><a class="text-sm-center nav-link "  v-bind:class="{ active: isCustomizedJewelry }" href="/jewelry-custom">定制</a></li>
            <li class="nav-item"><a class="text-sm-center nav-link "  v-bind:class="{ active: isKnowledge }" href="/press/category/knowledge/">珠宝知识</a></li>
            <li class="nav-item" id="appobtnbox"><a class="text-sm-center nav-link "  v-bind:class="{ active: isAppointment }" href="/appointment">在线预约</a></li>
            <li class="nav-item" id="contactbtnbox"><a class="text-sm-center nav-link "  v-bind:class="{ active: isContact }" href="/contact">联系我们</a></li>
            <li class="nav-item" id="brandstorybtnbox"><a class="text-sm-center nav-link "  v-bind:class="{ active: isContact }" href="/brandstory">关于我们</a></li>
          </ul>
        </div>
      </nav>
    </div>
  </div>
</template>
<style src='./Header.css' scoped></style>
<script src='./Header.js'></script>
