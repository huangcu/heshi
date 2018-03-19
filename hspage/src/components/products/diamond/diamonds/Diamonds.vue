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
      <diamond-filter></diamond-filter>
      <div id="diamondlist"></div>
      <div id="list-page-navi">第 <span id="diapagenavi"></span> 页</div>
      </div>
      <div id="feedbackcover">
        <div id="feedbackbox">
          <div class="bgimgcover" style="position:absolute; width:100%; height:100%; left:0; top:0; background-color:rgba(151,240,84,0.85); z-index:0;">
          </div>
          <p class="ordersuccwords"><span class="glyphicon glyphicon-thumbs-up"></span> 您成功预定了钻石!</p>
          <p style="z-index:2; position:relative;"><a id="ordersuccshoppinglistlinker" href="/shoppinglist-confirmed.php"><span class="glyphicon glyphicon-th-list"></span> 查看预定单</a> <button id="ordersuccmountinglinker" onclick="gotoMountingsfor()"><span></span> 配戒托</button> <button id="ordersuccsfdbcXbtn" onclick="xtheSfdb()">继续选购 &raquo;</button></p>
        </div>
      </div>
    </div>
</template>
<script src='./Diamonds.js'></script>
<style src='./Diamonds.css'></style>
