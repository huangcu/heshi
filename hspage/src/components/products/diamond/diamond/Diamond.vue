<template>
  <div>
    <vue-title title="合适钻石 | BEYOU DIAMOND | 推荐详情"></vue-title>
    <div class="contentmainwrapper">
    <h2 class="jewelry-title">合适推荐钻石</h2>
    <div v-if="status !== 'AVAIABLE'">
      <p>抱歉，该钻石已经售出。</p>
    </div>
    <div v-else class="jewely-description">
      <p>商品编号: <strong>{{ ref }}></strong></p>
      <span class="dia-description">'.$carat.'克拉</span>
      <span class="dia-description">'.$shape_txt.'</span>
      <span class="dia-description">'.$color.'色</span>
      <span class="dia-description">净度:'.$clarity.'</span>
      <span class="dia-description">切工: '.$cut_grade.' ／ 抛光: '.$polish.' ／ 对称: '.$symmetry.'</span>
      <span class="dia-description">'.$fluo_txt.'</span>
      <div v-if = "getLab == 'HRD'">
        <p class="certi-linker">
        <span class="dia-description">HRD</span>
        <a class="certi_linker" target="_blank" href="https://my.hrdantwerp.com/?L=&record_number=<?php echo $certificate_number; ?>&certificatetype=MC">查看证书 <span class="glyphicon glyphicon-new-window"></span></a></p>
      </div>
      <div v-else-if="getLab == 'GIA'">
        <p class="certi-linker">
        <span class="dia-description">GIA</span>
        <a class="certi_linker" target="_blank" href="http://www.gia.edu/cs/Satellite?pagename=GST%2FDispatcher&childpagename=GIA%2FPage%2FReportCheck&c=Page&cid=1355954554547&reportno=<?php echo $certificate_number; ?>">查看证书 <span class="glyphicon glyphicon-new-window"></span></a></p>
      </div>
      <div v-else-if="getLab == 'IGI'">
        <p class="certi-linker">
        <span class="dia-description">IGI</span>
        <a class="certi_linker" target="_blank" href="http://www.igiworldwide.com/verify.php?r=<?php echo $certificate_number; ?>">查看证书 <span class="glyphicon glyphicon-new-window"></span></a></p>
      </div>
      <div v-if="reco-words" id="reco-words">
      {{ recowords }}
      </div>
    </div>
    <div class="pricebox">
      <div v-if="isAgent()">
        <span class="price">合适价: {{agentPrice}} 欧元</span>
        <span class="price-special"> 您的代理价: {{agentPrice}}美元)</span>
      </div>
      <div v-else-if="accountLevel()">
        <span class="price">合适价:{{ customerPrice}} 欧元</span>
        <span class="price-special">您的会员价: {{ customerPrice}} 美元)</span>
      </div>
      <div v-else>
        <span class="price">合适价: {{ price }}美元)</span>
      </div>
    </div>
    <div v-if="imgs.length" class="jewelry-pics">
      <!-- <img class="jewelry-img" v-for="image in imgs" v-bind:src="image"/> -->
    </div>
      <p class="choosebtnbox">
        <a class="choosebtn" href="shoppinglist-confirmed.php?type=diamond&id=<?php echo $_GET['ref']; ?>&action=add"><span class="glyphicon glyphicon-th-list"></span> 加入预订单</a>
      </p>
    </div>
    <p class="gobackbtnbox">
    <span class="gobackindisign"><span class="glyphicon glyphicon-arrow-left"></span> 返回:</span> <a class="gobackbtn" href="/recommend-diamonds.php"><span class="glyphicon glyphicon-th"></span> 全部推荐钻石 </a>
    </p>
  </div>
</template>
<script src='./Diamond.js'></script>
<style src='./Diamond.css'></style>
