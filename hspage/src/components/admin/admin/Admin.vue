<template>
  <div>
    <vue-title title="管理界面"></vue-title>
    <h3 style="text-align:center;">管理首页</h3>
    <h4 style="text-align:center;">管理员 {{userprofile.real_name}}
      <img v-if="userprofile.icon!==''" :src="userproifle.icon" class="adminiconpic" />
      <a href="admin-edit.php" style="display:inline-block; position:relative; width:1px; overflow:visible; left:20px; font-size:14px;">
        <span style="display:inline-block; white-space:nowrap;">修改管理员信息 &raquo;</span>
      </a>
    </h4>
    <div id="maincontent">
      <!-- SELECT COUNT(*) AS num_verify_t FROM interested_items, diamonds WHERE item_id = stock_ref AND confirmed_for_check = "YES" AND item_type = "DIAMOND" AND (interested_items.available = "TOBECHECKED" OR (interested_items.available = "YES" AND diamonds.recommend_words IS NULL)) -->
      <div v-if="num_verify!==0" class="procedurebtnbox">
        <div class="step" id="step1">
          <a class="verifyrequestbtn" href="/admin/verify-request.php">验货请求 <span class="num">{{ num_verify }}</span></a>
        </div>
      </div>
      <br style="clear:both;" />
    </div> <!-- end of main content -->
    <div id="indication" style="position:fixed; width:100%; height:100%; background-color:rgba(255,255,255, 0.88); top:0; left:0; z-index:28; display:none;">
      <div id="indiinner" style="position:relative; width:200px; background-color:#0CF; margin:150px auto; padding:20px; text-align:center;">
        正在存储。。。
      </div>
    </div>
  </div> 
</template>
<script src="./Admin.js"></script>
<style src="./Admin.css" scoped></style>
