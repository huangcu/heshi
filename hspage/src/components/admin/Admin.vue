<template>
  <div>
    <div class="pageheader">
      <img id="logo" :src="imgURL" />
      <h3>合适钻石 系统管理</h3>
      <div class="navibox">
        <ul class="navi">
          <!-- admin level 4 and 5 -->
          <li v-if="userprofile.admin.admin_level >= 4" class="homepagenavi">
            <a href="index.php"><span class="glyphicon glyphicon-file"></span>销售流程</a>
            <ul class="subnavi-containner">
              <li class="homepagenavi subnavi">
              <a href="verify-request.php">新预定</a>
              </li>
              <li class="homepagenavi subnavi">
              <a href="orders-transaction.php">订单处理</a>
              </li>
              <li class="homepagenavi subnavi">
              <a href="pointhistory.php">返点</a>
              </li>
            </ul>
          </li>
          <!--
          <li><a href="orders.php"><span class="glyphicon glyphicon-list-alt"></span> 订单处理</a></li>
          -->
          <li v-if="userprofile.admin.admin_level >= 4" class="accountsnavi">
            <a href="userlist.php"><span class="glyphicon glyphicon-user"></span> 用户管理</a>
            <ul class="subnavi-containner">
              <li class="homepagenavi subnavi">
              <a href="userlist.php">客人管理</a>
              </li>
              <li class="homepagenavi subnavi">
              <a href="admins.php">内部人员管理</a>
              </li>
              <li class="homepagenavi subnavi">
              <a href="message-list.php">聊天纪录</a>
              </li>
            </ul>
          </li>

          <li class="contentnavi">
            <a href="import_excel.php">钻石货单</a>
            <ul class="subnavi-containner">
              <li class="homepagenavi subnavi">
                <a href="import_excel.php">管理货单</a>
              </li>
              <li v-if="userprofile.admin.admin_level >=4" class="homepagenavi subnavi">----</li>
              <li v-if="userprofile.admin.admin_level >=4" class="homepagenavi subnavi">
                <a href="pricemanagement.php">钻石价格规则</a>
              </li>
            </ul>
          </li>

          <li class="contentnavi"><a href="gemstonemanage.php">彩宝管理</a></li>
          <li class="contentnavi">
            <a href="jewelrymanage.php">首饰货单</a>
            <ul class="subnavi-containner">
              <li class="homepagenavi subnavi">
              <a href="new_jewelry.php">添加新首饰</a>
              </li>
              <li class="homepagenavi subnavi">
              <a href="jewelrymanage.php">编辑首饰内容</a>
              </li>
            </ul>
          </li>
          <li v-if="userprofile.admin.admin_level >=4" class="accountsnavi">
            <a href="user-analytics.php">分析报告</a>
          </li>
          <!--
          <li class="contentnavi"><a href="homepage.php">首页促销</a></li>
          -->
          <!--
          <li><a href="analyse.php">分析报告</a></li>
          -->
          <li><a target="_blank" href="../">回到网站 <span class="glyphicon glyphicon-new-window"></span></a></li>
        </ul>
      </div>
    </div>
    <router-view></router-view>
    <div id="notificationpad">
      <a v-if="userprofile.admin.admin_level >=4" href="/administrator/verify-request.php" id="checkrequestbtn">新预定 <span class="newmessagenum" id="newdiacheck_num">{{ num_verify }}</span></a>
      <a v-if="userprofile.admin.admin_level >=4" href="/administrator/orders-transaction.php" id="neworderbtn">处理中订单 <span class="newmessagenum" id="newdiacheck_num">{{ new_order_transaction_num }}</span></a>
      <a href="/administrator/message-list.php" id="newmessagebtn">客户消息 <span class="newmessagenum" id="newmessage_num"> {{ num_newmessage }}</span></a>
    </div>
    <form class="logout" action="" method="post" style="position:absolute; top:5px; right:0px; width:80px;">
        <input type="submit" name="logout" id="logout" value="logout">
    </form>
  </div>
</template>
<script src="./Admin.js"></script>
<style src="./Admin.css" scoped></style>
