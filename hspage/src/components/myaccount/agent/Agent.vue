<template>
  <div>
  <h3>合适代理
    <span class="accountlevelstars">
    {{ fromAgentLevel }}
    </span>
  </h3>
  <p class="inpage-navi-box">
    <button class="inpage-navi active" id="inpage-navi_generalinfo" v-on:click="goToSection('generalinfo')">基本信息</button>
    <span id="new-iconbox">
      <img id="newicon-mysite" src="/_images/constant/new-icon.png" />
    </span>
    <button class="inpage-navi mywebsitebtn" id="inpage-navi_mywebsite" v-on:click="goToSection('mywebsite')">我的网站</button>
    <button class="inpage-navi" id="inpage-navi_mypoints" v-on:click="goToSection('mypoints')">我的积分</button>
    <button class="inpage-navi" id="inpage-navi_myclients" v-on:click="goToSection('myclients')">我的客人</button>
    <button class="inpage-navi" id="inpage-navi_clientorders" v-on:click="goToSection('clientorders')">客人的预定</button>
    <button class="inpage-navi" id="inpage-navi_coupon" v-on:click="goToSection('coupon')">打折／代金券</button>
    <button class="inpage-navi" id="inpage-navi_personalinfo" v-on:click="goToSection('personalinfo')">个人信息</button>
    <a class="inpage-navi" id="listdownloadbtn" href="stocklist-csv.php">下载货单 &raquo;</a>
  </p>
  <!-- offerticket box -->
  <div class="generalinfo discountbox" id="generalinfo">
    <h3 class="g-info-title">您好，{{ name }}
      <form method="post" action="" id="logoutform">
        <input type="hidden" name="logout" value="confirmed" />
        <input type="submit" value="登出账号" name="logoutbtn" id="logoutbtn" />
      </form>
    </h3>
<?php
$theagent_price_ratio=priceforagent($agent_level, 1);
if($agent_level_locked=='NO'){
?>
<h4 class="g-info-content">您目前的代理级别: <?php echo $agent_level; ?>级</h4> 
<h4 class="g-info-content">购买折扣: <?php echo $theagent_price_ratio; ?></h4>

<?php
####################################### 代购的积分 ############################
$sql_total_points_dia_a='SELECT SUM(sold_price_total) AS totalsoldpricedia_a FROM diamonds_orders WHERE buyer_id = '.$accountID.' AND transaction_status = "SOLD" AND ordered_at > SUBDATE(NOW(), INTERVAL 1 YEAR)';
foreach($conn->query($sql_total_points_dia_a) as $rtpda){
	$totalpointsdia_a=$rtpda['totalsoldpricedia_a'];
}
$sql_total_points_dia_c='SELECT SUM(sold_price_total) AS totalsoldpricedia_c FROM diamonds_orders, users WHERE buyer_id = users.id AND users.recommenedby = '.$accountID.' AND transaction_status = "SOLD" AND ordered_at > SUBDATE(NOW(), INTERVAL 1 YEAR)';
foreach($conn->query($sql_total_points_dia_c) as $rtpdc){
	$totalpointsdia_c=$rtpdc['totalsoldpricedia_c'];
}


$sql_total_points_jew_a='SELECT SUM(sold_price_total) AS totalsoldpricejew_a FROM jewelry_orders WHERE buyer_id = '.$accountID.' AND transaction_status = "SOLD" AND ordered_at > SUBDATE(NOW(), INTERVAL 1 YEAR)';
foreach($conn->query($sql_total_points_jew_a) as $rtpja){
	$totalpointsjew_a=$rtpja['totalsoldpricejew_a'];
}
$sql_total_points_jew_c='SELECT SUM(sold_price_total) AS totalsoldpricejew_c FROM jewelry_orders, users WHERE buyer_id = users.id AND users.recommenedby = '.$accountID.' AND transaction_status = "SOLD" AND ordered_at > SUBDATE(NOW(), INTERVAL 1 YEAR)';
foreach($conn->query($sql_total_points_jew_c) as $rtpjc){
	$totalpointsjew_c=$rtpjc['totalsoldpricejew_c'];
}


$sql_total_dia_num_a='SELECT COUNT(*) AS totaldianum_a FROM diamonds_orders WHERE buyer_id = '.$accountID.' AND transaction_status = "SOLD" AND ordered_at > SUBDATE(NOW(), INTERVAL 1 YEAR)';
foreach($conn->query($sql_total_dia_num_a) as $rtdna){
	$totaldiamum_a=$rtdna['totaldianum_a'];
}
$sql_total_dia_num_c='SELECT COUNT(*) AS totaldianum_c FROM diamonds_orders, users WHERE buyer_id = users.id AND users.recommenedby = '.$accountID.' AND transaction_status = "SOLD" AND ordered_at > SUBDATE(NOW(), INTERVAL 1 YEAR)';
foreach($conn->query($sql_total_dia_num_c) as $rtdnc){
	$totaldiamum_c=$rtdnc['totaldianum_c'];
}


$sql_total_jew_num_a='SELECT COUNT(*) AS totaljewnum_a FROM jewelry_orders WHERE buyer_id = '.$accountID.' AND transaction_status = "SOLD" AND ordered_at > SUBDATE(NOW(), INTERVAL 1 YEAR)';
foreach($conn->query($sql_total_jew_num_a) as $rtjna){
	$totaljewmum_a=$rtjna['totaljewnum_a'];
}
$sql_total_jew_num_c='SELECT COUNT(*) AS totaljewnum_c FROM jewelry_orders, users WHERE buyer_id = users.id AND users.recommenedby = '.$accountID.' AND transaction_status = "SOLD" AND ordered_at > SUBDATE(NOW(), INTERVAL 1 YEAR)';
foreach($conn->query($sql_total_jew_num_c) as $rtjnc){
	$totaljewmum_c=$rtjnc['totaljewnum_c'];
}



$total_point=$totalpointsdia_a+$totalpointsdia_c+$totalpointsjew_a+$totalpointsjew_c;
$total_dia_num=$totaldiamum_a+$totaldiamum_c+$totaljewmum_a+$totaljewmum_c;
################################### END 代购的积分 ############################
?>

<h4>我的积分：(一年内) 总购买金额: <strong><?php echo $total_point; ?></strong>欧元；钻石购买总颗数: <strong style="display:inline-block; margin-right:15px;"><?php echo $total_dia_num; ?></strong> <button id="mypointdetailbtn" href="/myaccount.php#section-mypoints" onclick="goToSection('mypoints')">详情 &raquo;</button></h4>






<p class="ticket-question" onclick="showRegulation()">折扣、积分、返点规则 &raquo;</p>
<div class="ticket-answer">
合适代理折扣：
<ul>
<li>一级代理（享有 9 折价格）: 一年购买12颗钻石以下 或 总购买金额在4万欧元以下</li>
<li>二级代理（享有8.5折价格）: 一年内购买超过12颗钻石而且总购买金额超过4万欧元</li>
<li>三级代理（享有8.3折价格）: 一年内购买超过20颗钻石而且总购买金额超过10万欧元</li>
</ul>


代理积分、返点规则
<ul>
<li>
一级代理如果一年购买12颗钻石以上而且总购买金额超过4万欧元，将获得一年内总购买金额的5%的返点，并且升为二级代理
</li>
<li>
二级代理如果一年购买超过20颗钻石，将获得一年内（如果上次返点在一年内，则从上次返点开始到现在）总购买金额的2%的返点，并且升为三级代理。如果一年内购买低于12颗或总金额低于4万欧元，则降为1级代理。
</li>
<li>
三级级代理已享有最高折扣，不再返点。如果一年内购买低于20颗钻石或总金额低于10万欧元，相应的降为二级或一级代理。
</li>
</ul>
</div>
<?php
}else{
?>
<h4>您目前的代理级别: <?php echo $agent_level; ?>级 &nbsp; &nbsp; &nbsp; 购买折扣为: <?php echo $theagent_price_ratio; ?></h4>

<?php	
}
?>




<div id="myreferenceticketbox">
<p id="titleformyreference" onclick="openRecoContent();">
轻松推广 增加客人 &raquo;
</p>

<div id="reco-contentbox">
由您推荐的注册用户都是您的客人。您的客人仍然可以向下继续推荐客人。<br />
您的客人推荐的客人如果购买合适商品，您同样享有<strong>1％</strong>的返点。您客人推荐的客人再推荐客人，您依然享有1%的返点。以此类推，惊喜无限！<br />
<div class="share-ticket-tools">
<span class="glyphicon glyphicon-hand-right"></span> 推荐方法:
<p id="toolexpl">点击打开下面链接中的页面，然后把该页面分享到微信朋友圈(或发送给朋友)。也可以拷贝这个链接，邮件等方式发给朋友。就这么简单！</p>
<p><a target="_blank" href="http://beyoudiamond.com/login.php?ref=<?php echo $accountID.'X'.$offer_ticket; ?>" target="_blank" id="recommenduserurl">http://beyoudiamond.com/login.php?ref=<?php echo $accountID.'X'.$offer_ticket; ?></a></p>
</div>
</div><!-- end roco-contentbox -->

</div>
</div>
<!-- my website -->
<div class="mywebsite" id="mywebsite">
<?php
$sql_reseller='SELECT * FROM agent_contact_info WHERE user_id = '.$accountID;
foreach($conn->query($sql_reseller) as $rr){
	$r_phone=$rr['phone'];
	$r_email=$rr['email'];
	$r_wechat_qr=$rr['wechat_qr'];
	$r_address=$rr['address'];
	$r_qq=$rr['qq'];
	$r_other_contact_info=$rr['other_contact_info'];
}


$filename='reseller'.$accountID.'/wp-content/themes/flatsome-child/account-configuration.php';

if(!file_exists($filename)){
?>

<h3>您的销售利器，我们已经为您准备好！</h3>
<p><br />合适总部可以专门为代理配置属于代理自己的网站。</p><p>网站中不出现「合适钻石」的品牌信息，而使用代理自己的品牌。</p><p>钻石和首饰的库存和合适钻石保持一致，显示的价格可以由代理自己调整。</p><p>开通网站和详情请咨询合适总部<strong>何丽丽</strong><br /><br /></p>

<?php
}else{
?>
<div class="mysitemanagingbloc">
<?php
if(isset($feedbackmessageofrequest)){
?>
<h3><?php echo $feedbackmessageofrequest; ?></h3>
<?php
}
?>

<?php	
	require_once('reseller'.$accountID.'/wp-content/themes/flatsome-child/account-configuration.php');
	
	function dbConnect_rs($usertype='write', $connectionType = 'pdo', $crruseraccountID) {
		$host = 'localhost';
		$db = 'beyoudia_reseller'.$crruseraccountID;
		if ($usertype  == 'read') {
		$user = 'beyoudia_user';
		$pwd = 'Sihui8mp5e878';
		} elseif ($usertype == 'write') {
		$user = 'beyoudia_user';
		$pwd = 'Sihui8mp5e878';
		} else {
		exit('Unrecognized connection type');
		}
		if ($connectionType == 'mysqli') {
			
		return new mysqli($host, $user, $pwd, $db) or die ('Cannot open database');
		} else {
			try {
				return new PDO("mysql:host=$host;dbname=$db", $user, $pwd);
			} catch (PDOException $e) {
				echo 'Cannot connect to database!!!!!!!!!!!!!';
				exit;
			}
		}
	}
	
	$conn_rs=dbConnect_rs('write','pdo', $accountID);
  $conn_rs->query("SET NAMES 'utf8'");

	
	$sql_site_data='SELECT * FROM wp_options';
	foreach($conn_rs->query($sql_site_data) as $rrs){
		
		$crr_option=$rrs['option_name'];
		$crr_value=$rrs['option_value'];
		
		if($crr_option=='siteurl'){
			$active_siteurl=$crr_value;
		}else if($crr_option=='blogname'){
			$site_title=$crr_value;
		}
		
		$pos = strpos($active_siteurl, "beyoudiamond.com");
		
		
	}
	

if($pos === FALSE || $active_siteurl=='http://beyoudiamond.com/reseller1'){
	
	$domainstrtoreplace=array('http://','/');
	
	$domaintoshow=str_replace($domainstrtoreplace, '', $active_siteurl);
	
?>
<p class="sitestatuswords"><span class="glyphicon glyphicon-thumbs-up"></span>您的网站已经开通 (<?php echo $domaintoshow; ?>) <a target="_blank" href="<?php echo $active_siteurl; ?>">查看 <span class="glyphicon glyphicon-new-window"></span></a></p>


<form action="" method="post">
<p>
<label>商品价格：</label> <input type="number" step="0.01" name="resellerpriceratio" id="resellerpriceratio" value="<?php echo $resellerpriceratio; ?>" /> x 合适原始价格
</p>

<input type="submit" value="保存" />
</form>

<div class="sitemanageformbox">
<form name="loginform" id="loginform" action="<?php echo $active_siteurl; ?>/wp-login.php" method="post" target="_blank">
	<p style="display:none;">
		<input type="text" name="log" id="user_login" class="input" value="manager" size="20"></label>
	</p>
	<p style="display:none;">
		<label for="user_pass">Password<br>
		<input type="password" name="pwd" id="user_pass" class="input" value="BeyouR!2017" size="20"></label>
	</p>
		<p class="forgetmenot" style="display:none;"><label for="rememberme"><input name="rememberme" type="checkbox" id="rememberme" value="forever" checked="checked"> Remember Me</label></p>
	<p class="submit">
		<input type="submit" name="wp-submit" id="wp-submit" class="button button-primary button-large" value="管理网站内容">
		<input type="hidden" name="redirect_to" value="<?php echo $active_siteurl; ?>">
		<input type="hidden" name="testcookie" value="1">
	</p>
</form>
<p style="font-size:12px"><span class="glyphicon glyphicon-info-sign"></span> 如果需要登录，用户名: manager 密码: BeyouR!2017</p>
</div>

<?php
}else{
	if(!empty($userdefineddomainname)){//这个变量要放在 account-configuration.php 文件里
?>
<p class="sitestatuswords"><span class="glyphicon glyphicon-ok"></span> 我们已经收到您提交的信息，请稍等，我们会在24小时内为您开通。</p>

<p><label>网址(域名)选择：</label> <input type="text" name="sitedomain" disabled="disabled" id="sitedomain" value="<?php echo $userdefineddomainname; ?>" /></p>
<p>
<label>品牌名称：</label> <input type="text" name="sitetitle" id="sitetitle" disabled="disabled" value="<?php echo $site_title; ?>" />
</p>
<?php
	}else{
?>
<p class="sitestatuswords">请填写您的网站信息并提交，填好后我们会在24小时内为您开通。</p>

<form action="" method="post">

<p><label>网址(域名)选择：</label> <input type="text" name="sitedomain" id="sitedomain" value="" />
<select name="domainsuffix" id="domainsuffix">
<option value="com">.com</option>
<option value="be">.be</option>
<option value="info">.info</option>
<option value="net">.net</option>
<option value="org">.org</option>
<option value="online">.online</option>
<option value="biz">.biz</option>
<option value="me">.me</option>
</select>

<button type="button" id="domaincheckbtn" onclick="checkdomainavailability()">查询是否可用 <span class="glyphicon glyphicon-search"></span> </button>
<span class="note"><span class="glyphicon glyphicon-exclamation-sign"></span> 注：提交后不可更改</span>
</p>

<span class="domainstatuswords">
<?php

if(!empty($r_website_domain)){
	if($r_domain_active=="NO"){
		$domainstatuswords = '域名 “'.$r_website_domain.'” 申请中，请稍等';
	}else if($r_domain_active=="OCCU"){
		$domainstatuswords = '域名 “'.$r_website_domain.'” 已经被占用，请另选一个';
	}
}else{
	$domainstatuswords = '';
}

echo $domainstatuswords;
?>
</span>

<p>
<label>品牌名称：</label> <input type="text" name="sitetitle" id="sitetitle" value="<?php echo $site_title; ?>" />
<span class="note"><span class="glyphicon glyphicon-exclamation-sign"></span> 注：提交后不可更改</span>
</p>
<p>
<label>商品价格：</label> <input type="number" step="0.01" name="resellerpriceratio" id="resellerpriceratio" value="<?php echo $resellerpriceratio; ?>" /> x 合适原始价格
<span class="note"><span class="glyphicon glyphicon-info-sign"></span> 注：提交后可以随时更改</span>
</p>



<input type="submit" value="提交" />
</form>



<?php
	}
}
?>

</div>


<div class="tech_support_bloc">
<h4>技术支持</h4>

<h5>如何更改网站上的内容? </h5>
<p>操作非常简单，请看一下我们为您录制的两段几分钟的视频，您即能轻松掌握。</p>
<p>
<a target="_blank" href="http://beyoudiamond.com/press/%E4%BB%A3%E8%B4%AD%E7%BD%91%E7%AB%99%E4%BD%BF%E7%94%A8%E4%BB%8B%E7%BB%8D-%E9%A6%96%E9%A1%B5-2/">请点击这里，</a>查看视频解说 (上集)</p>
<p>
<a target="_blank" href="http://beyoudiamond.com/press/%E4%BB%A3%E8%B4%AD%E7%BD%91%E7%AB%99%E4%BD%BF%E7%94%A8%E4%BB%8B%E7%BB%8D-%E5%85%B6%E4%BB%96%E8%AF%B4%E6%98%8E-small/">请点击这里，</a>查看视频解说 (下集)</p>

<p>更多的问题，请点击右侧的“在线咨询”联系我们，我们会给您详尽的解答</p>

</div>

<br class="clear" />
<?php
}
?>



<script type="text/javascript">
function checkdomainavailability(){
	var crr_domainchoice=$('input#sitedomain').val()+'.'+$('select#domainsuffix').val();
	var domaincheckurl='https://sg.godaddy.com/zh/domains/searchresults.aspx?checkAvail=1&tmskey=&domainToCheck='+crr_domainchoice;
	window.open(domaincheckurl);
}
</script>
</div>









<!-- orders -->
<div class="orders" id="clientorders">
<h3>客人的订单</h3>

<?php
$sql_transactions='SELECT * FROM transactions WHERE belong_to_agent = '.$accountID.' AND archived = "NO"';
$stmt_transaction=$conn->query($sql_transactions);
$transactionfound=$stmt_transaction->rowCount();
if(!$transactionfound){
?>
<p>目前尚无订单可以显示</p>
<?php
}else{
?>
<ul>
<?php
foreach($stmt_transaction as $rt){
	$crr_transaction_id=$rt['id'];
	$crr_tr_buyer_id=$rt['buyer_id'];
	$crr_tr_created=$rt['created'];
	$crr_tr_agent_status=$rt['agent_status'];
	$crr_tr_heshi_status=$rt['heshi_status'];
	
	$crr_buyer_price_ratio=1;
	
	$sql_crr_user='SELECT name, email, account_level, account_level_assigned, agent, agent_level FROM users WHERE id = '.$crr_tr_buyer_id;
	foreach($conn->query($sql_crr_user) as $rcui){
		$crr_buyer_name=$rcui['name'];
		$crr_buyer_email=$rcui['email'];
		
		$crr_buyer_agent=$rcui['agent'];
		$crr_buyer_agent_level=$rcui['agent_level'];
		$crr_buyer_account_level=$rcui['account_level'];
		$crr_buyer_account_level_assigned=$rcui['account_level_assigned'];
		
		if($crr_buyer_account_level_assigned>$crr_buyer_account_level){
			$crr_buyer_account_level=$crr_buyer_account_level_assigned;
		}
		
		
		if(empty($crr_buyer_name)){
			$crr_buyer_name=$crr_buyer_email;
		}
	}
	
	if($crr_buyer_agent=='YES'){
		$crr_buyer_price_ratio=priceforagent($crr_buyer_agent_level, 1);
	}else{
		$crr_buyer_price_ratio=priceforaccount($crr_buyer_account_level, 1);
	}
	
	
	
	$thisTransactionisComplete=true;
	$thisTransactionDiaList='';
	$thisTransactionJewList='';
	
	$crr_tr_total_price=0;
	$crr_tr_heshibi_amount=0;
	//$crr_tr_ori_price=0;
	
	$crr_tr_agent_profit=0;
	
	$sql_dia_orders='SELECT * FROM diamonds_orders WHERE transaction_id = '.$crr_transaction_id;
	$stmt_dia_orders=$conn->query($sql_dia_orders);
	$dia_orders_found=$stmt_dia_orders->rowCount();
	
	if($dia_orders_found){
		//$crr_tr_diamonds_array=explode(',',$crr_tr_diamonds);
		foreach($stmt_dia_orders as $r_c_d){
			$crr_order_id=$r_c_d['id'];
			$crr_dia_id=$r_c_d['diamondid'];
			$crr_dia_status=$r_c_d['transaction_status'];
			$crr_dia_price=$r_c_d['diamond_price'];
			$crr_dia_acc_id=$r_c_d['jewellery_id'];
			$crr_dia_acc_price=$r_c_d['jewellery_price'];
			
			$sql_ori_price='SELECT price_retail FROM diamonds WHERE stock_ref = "'.$crr_dia_id.'"';
			foreach($conn->query($sql_ori_price) as $rop){
				$crr_dia_price_ori=$rop['price_retail'];
			}
			//$crr_tr_ori_price+=$crr_dia_price_ori;
			
			if(!empty($thisTransactionDiaList)){
				$thisTransactionDiaList.=', ';
			}
			//$crr_tr_total_price+=$crr_dia_price;
			$thisTransactionDiaList.=$crr_dia_id;
			$crr_acc_price_ori=0;
			if(!empty($crr_dia_acc_id)){
				$thisTransactionDiaList.='(空托'.$crr_dia_acc_id.')';
				//$crr_tr_total_price+=$crr_dia_acc_price;
				
				$sql_ori_acc_price='SELECT price FROM jewelry WHERE id = '.$crr_dia_acc_id;
				foreach($conn->query($sql_ori_acc_price) as $ropa){
					$crr_acc_price_ori=$ropa['price'];
				}
				//$crr_tr_ori_price+=$crr_acc_price_ori;
			}
			
			$crr_order_buyer_prices=userOrderPrice($crr_buyer_price_ratio, $crr_order_id, 'DIAMOND');
			$crr_order_buyer_price=$crr_order_buyer_prices['currentorderprice_euro'];
			$crr_order_agent_price=($crr_dia_price_ori*$theagent_price_ratio*$_currencies['DollarToEuro'])+($crr_acc_price_ori*$theagent_price_ratio);
			
			//echo '||||||||||||||||||||||||'.$crr_order_buyer_price.'||||||||||||||||||||||||'.$crr_order_agent_price.'||||||||||||||||||||||||';
			
			$crr_order_agent_profit=$crr_order_buyer_price-$crr_order_agent_price;
			
			$crr_tr_total_price+=$crr_order_buyer_price;
			$crr_tr_agent_profit+=$crr_order_agent_profit;
			
			
			if($crr_dia_status=="SOLD"){
				$thisTransactionDiaList.='<span class="glyphicon glyphicon-ok"></span>';
			}else{
				$thisTransactionDiaList.='<span class="glyphicon glyphicon-hourglass"></span>';
				$thisTransactionisComplete=false;
			}
				
			
		}
	}
	
	
	$sql_jew_orders='SELECT * FROM jewelry_orders WHERE transaction_id = '.$crr_transaction_id;
	$stmt_jew_orders=$conn->query($sql_jew_orders);
	$jew_orders_found=$stmt_jew_orders->rowCount();
	
	if($jew_orders_found){
	
		foreach($stmt_jew_orders as $r_c_j){
			$crr_order_id=$r_c_j['id'];
			$crr_jew_id=$r_c_j['jewellery_id'];
			$crr_jew_status=$r_c_j['transaction_status'];
			$crr_jew_price=$r_c_j['jewellery_price'];
			
			$sql_ori_jew_price='SELECT price FROM jewelry WHERE id = '.$crr_jew_id;
			foreach($conn->query($sql_ori_jew_price) as $ropj){
				$crr_jew_price_ori=$ropj['price'];
			}
			//$crr_tr_ori_price+=$crr_jew_price_ori;
			
			$crr_order_buyer_prices=userOrderPrice($crr_buyer_price_ratio, $crr_order_id, 'JEWELRY');
			$crr_order_buyer_price=$crr_order_buyer_prices['currentorderprice_euro'];
			$crr_order_agent_price=$crr_jew_price_ori*$theagent_price_ratio;
			$crr_order_agent_profit=$crr_order_buyer_price-$crr_order_agent_price;
			
			$crr_tr_total_price+=$crr_order_buyer_price;
			$crr_tr_agent_profit+=$crr_order_agent_profit;
			
			if(!empty($thisTransactionJewList)){
				$thisTransactionJewList.=', ';
			}
			
			$thisTransactionJewList.=$crr_jew_id;
			//$crr_tr_total_price+=$crr_jew_price;
			
			if($crr_jew_status=="SOLD"){
				$thisTransactionJewList.='<span class="glyphicon glyphicon-ok"></span>';
			}else{
				$thisTransactionJewList.='<span class="glyphicon glyphicon-hourglass"></span>';
				$thisTransactionisComplete=false;
			}
			
		}
	}
	
	
	$sql_heshibi='SELECT * FROM heshibi_in_use WHERE transaction_id = '.$crr_transaction_id;
	$stmt_heshibi=$conn->query($sql_heshibi);
	$heshibi_found=$stmt_heshibi->rowCount();
	if($heshibi_found){
		foreach($stmt_heshibi as $rh){
			$crr_heshibi_price=$rh['amount'];
			$crr_tr_heshibi_amount+=$crr_heshibi_price;
		}
	}
	
	$crr_tr_final_price=$crr_tr_total_price-$crr_tr_heshibi_amount;
	
	
	
	
	
	if(empty($thisTransactionDiaList)){
		$thisTransactionDiaList='无';
	}
	if(empty($thisTransactionJewList)){
		$thisTransactionJewList='无';
	}
	
	
	
?>
<li class="transactionpiece <?php echo $crr_tr_heshi_status; ?>" id="transactionpiece_<?php echo $crr_transaction_id; ?>">
<span class="buyername"><label><span class="glyphicon glyphicon-user"></span> 预定人</label> <a target="_blank" href="userhistory.php?id=<?php echo $crr_tr_buyer_id; ?>"><?php echo $crr_buyer_name; ?> &raquo; </a></span>
<span class="transactiondiamonds"><label>预定的钻石(ID)</label><?php echo $thisTransactionDiaList; ?></span>
<span class="transactionjewelry"><label>预定的首饰(ID)</label><?php echo $thisTransactionJewList; ?></span>
<span class="transactiontotalprice"><label>总价格</label><?php echo round($crr_tr_total_price,2); ?>€</span>
<span class="transactionheshibi"><label>使用合适币</label><?php echo $crr_tr_heshibi_amount; ?>€</span>
<span class="transactionfinalprice"><label>最终应付款</label><?php echo round($crr_tr_final_price,2); ?>€</span>

<span class="agentprice"><label>代理利润</label><?php echo round($crr_tr_agent_profit,2); ?>€</span>

<span class="transactiontime">下单时间: <?php echo $crr_tr_created; ?></span>
<p class="transactionstatusactionsbox">
<?php
if($crr_tr_heshi_status=='PROCESSING'){
?>
<span class="precessinginfo"><span class="glyphicon glyphicon-th-list"></span> 合适总部正在处理订单</span>
<?php
}else if($crr_tr_heshi_status=='COMPLETE'){
	if($crr_tr_agent_status=='PROCESSING'){
	?>
	<button class="transactioncompletebtn" id="transactioncompletebtn_<?php echo $crr_transaction_id; ?>" onclick="completeTransaction(<?php echo $crr_transaction_id; ?>)"> <span class="glyphicon glyphicon-check"></span> 标记为交易完成</button>
  
	<?php
	}else{
	?>
	<span class="tr_complete_words"><span class="glyphicon glyphicon-ok"></span> 交易已完成</span>
	<button class="archivebtn" id="transactionarchive_<?php echo $crr_transaction_id; ?>" onclick="archiveTransaction(<?php echo $crr_transaction_id; ?>)"><span class="glyphicon glyphicon-inbox"></span> 存档</button>
	<?php
	}
	
}
?>
</li>
<?php
}
?>

</ul>


<script type="text/javascript">
function completeTransaction(transactionid){
	$('button#transactioncompletebtn_'+transactionid).html('<span class="glyphicon glyphicon-refresh"></span> 执行中...稍候').attr('disabled','disabled');
	$.post(
		"/_includes/functions/completetransaction.php", 
		{id: transactionid}, 
		function(data){
			//alert(data);
			console.log('order feedback: '+data)
			
			if($.trim(data)=='OK'){
				//alert('ordered');
				$('button#transactioncompletebtn_'+transactionid).html('<span class="glyphicon glyphicon-ok"></span> 交易完成').addClass('complete').attr('disabled','disabled');
				
				/*
				$.post(
					"/_includes/functions/balancecheck.php", 
					{id: "<?php echo $accountID; ?>"}, 
					function(data){
						//alert(data);
						console.log('order feedback: '+data)
						
						if($.isNumeric(data)){
							$('span#heshibi-balance-value').html(data);
						}else{
							//alert('网络异常，请重试!');
						}
						
					}
				);
				*/
				
			}else{
				alert('网络异常，请重试!');
			}
			
		}
	);
}

function archiveTransaction(transactionid){
	$('button#transactionarchive_'+transactionid).html('<span class="glyphicon glyphicon-refresh"></span> 执行中...稍候').attr('disabled','disabled');
	$.post(
		"/_includes/functions/archiveransaction.php", 
		{id: transactionid}, 
		function(data){
			//alert(data);
			console.log('order feedback: '+data)
			
			if($.trim(data)=='OK'){
				//alert('ordered');
				$('button#transactionarchive_'+transactionid).html('<span class="glyphicon glyphicon-ok"></span> 已存档').addClass('complete').attr('disabled','disabled');
				$('li#transactionpiece_'+transactionid).delay(500).fadeOut('normal');
			}else{
				alert('网络异常，请重试!');
			}
			
		}
	);
}
</script>


<?php
}
?>

</div>







<!-- history -->
<div class="history" id="mypoints">
<!--
<h3>购买历史纪录</h3>
-->
<?php
$total_price=0;
$total_price_one_year=0;
$total_piece_one_year=0;



$current_date=date("Y-m-d");
$strtotime=strtotime($current_date);
$last_year=strtotime("-1 year",$strtotime);
?>
<div class="historybox history-diamond">
<h4>钻石购买历史纪录</h4>
<ul>
<?php
$total_price_diamond=0;

//$sql_history_diamond='SELECT * FROM diamonds_orders WHERE buyer_id = '.$accountID.' AND transaction_status = "SOLD" ORDER BY ordered_at DESC';
$sql_history_diamond_order='SELECT diamonds_orders.id AS order_id, diamondid, diamond_price, jewellery_price, jewellery_id, buyer_id, downpayment, in_stock_confirmed, order_sent, sold_price_total, extra_info, ordered_at, users.id AS idofuser, users.name AS nameofuser, email, recommenedby FROM diamonds_orders, users WHERE diamonds_orders.buyer_id = users.id AND ( diamonds_orders.buyer_id = '.$accountID.' OR (users.recommenedby = '.$accountID.' AND users.agent = "NO")) AND transaction_status = "SOLD" ORDER BY ordered_at ASC';

$stmt_h_d=$conn->query($sql_history_diamond_order);
$h_d_found=$stmt_h_d->rowCount();
if($h_d_found){
	require_once('_includes/functions/detail_fordiamond_byshape.php');
	foreach($stmt_h_d as $r_h_d){
		$diamondid=$r_h_d['diamondid'];
		$jewellery_id=$r_h_d['jewellery_id'];
		$sold_price_total=$r_h_d['sold_price_total'];
		$ordered_at=$r_h_d['ordered_at'];
		
		$idofuser=$r_h_d['idofuser'];
		$nameofuser=$r_h_d['nameofuser'];
		$crr_email=$r_h_d['email'];
		
		if(empty($nameofuser)){
			$nameofuser=$crr_email;
		}
		
		$buyernametxt='';
		$buyerclass='';
		
		if($idofuser!=$accountID){
			//$buyerclass=' orderedbyclientuser';
			$buyernametxt='<p class="orderredbytxt"><span class="thetxtofusername"><span class="glyphicon glyphicon-user"></span> 预定人: '.$nameofuser.'</span></p>';
		}
		
		$total_price+=$sold_price_total;
		$total_price_diamond+=$sold_price_total;
		$listclass='';
		if(strtotime($ordered_at)>=$last_year){
			$total_price_one_year+=$sold_price_total;
			$listclass=' oneyear';
			$total_piece_one_year+=1;
		}
		
		$crr_dia_shape='不详';
		$crr_dia_weight='不详';
		$crr_dia_color='不详';
		$crr_dia_clarity='不详';
		
		$sql_diamond_crr='SELECT * FROM diamonds WHERE stock_ref = "'.$diamondid.'"';
		foreach($conn->query($sql_diamond_crr) as $r_dia){
			$crr_dia_shape=$r_dia['shape'];
			$crr_dia_weight=$r_dia['carat'];
			$crr_dia_color=$r_dia['color'];
			$crr_dia_clarity=$r_dia['clarity'];
		}
		
		if($jewellery_id!=NULL && $jewellery_id!=''){
			$sql_jewelry_crr='SELECT * FROM jewelry WHERE id = '.$jewellery_id;
			$stmt_j_c=$conn->query($sql_jewelry_crr);
			$crr_j_found=$stmt_j_c->rowCount();
			if($crr_j_found){
				foreach($stmt_j_c as $rj){
					$jewelry_pic=$rj['image1'];
					$jewelry_name=$rj['name'];
				}
			}
		}else{
			$jewelry_pic=NULL;
			$jewelry_name=NULL;
		}
		
		if(isset($jewelry_pic)){
			$crr_order_picture='/_images/jewelry/thumbs/'.$jewelry_pic;
			$crr_dia_pic_for_jew='/_images/constant/'.detail_forDiamond_byShape($crr_dia_shape, 'PICTURE');
		}else{			
			$crr_order_picture='/_images/constant/'.detail_forDiamond_byShape($crr_dia_shape, 'PICTURE');
			$crr_dia_pic_for_jew=NULL;
		}
?>
<li class="dia-order-history<?php echo $listclass; ?>" id="diaorder-<?php echo $diamondid; ?>">
<?php echo $buyernametxt; ?>

<!--
<span class="picfororder" style="background-image:url('<?php echo $crr_order_picture; ?>');">
<?php
if(!empty($crr_dia_pic_for_jew)){
?>
<img class="img-dia-for-jew" src="<?php echo $crr_dia_pic_for_jew; ?>" />
<?php
}
?>
</span>
-->

<span class="dia-description"><?php echo detail_forDiamond_byShape($crr_dia_shape, 'NAMECN'); ?>钻石: <?php echo $crr_dia_weight ?>克拉，<?php echo $crr_dia_color; ?>色， 净度 <?php echo $crr_dia_clarity; ?></span>
<span class="jewelry-description">配饰:
<?php
if(isset($jewelry_name)){
	echo $jewelry_name;
}else{
	echo '无';
}
?>
</span>
<span class="totalpricecrr">
<?php require_once('_includes/functions/currency_calculator.php'); ?>
总价格: <?php echo round($sold_price_total, 2); ?>欧元
</span>




<span class="orderedtime">
<?php 
if(strtotime($ordered_at)<$last_year){
	echo '超一年 不积分 ';
}else{
	echo '<span class="glyphicon glyphicon-ok"></span> 积分有效 ';
}
echo $ordered_at; 
?>
</span>
<br class="clear" />
</li>
<?php
	}
}else{
	?>
  <li class="no-order">没有可积分的钻石购买纪录</li>
  <?php
}










############################################################
# for records no - points
############################################################

//$sql_history_diamond='SELECT * FROM diamonds_orders WHERE buyer_id = '.$accountID.' AND transaction_status = "SOLD" ORDER BY ordered_at DESC';
$sql_history_diamond_order_noprofit='SELECT diamonds_orders_no_profit.id AS order_id, diamondid, diamond_price, jewellery_price, jewellery_id, buyer_id, downpayment, in_stock_confirmed, order_sent, sold_price_total, extra_info, ordered_at, users.id AS idofuser, users.name AS nameofuser, email, recommenedby FROM diamonds_orders_no_profit, users WHERE diamonds_orders_no_profit.buyer_id = users.id AND diamonds_orders_no_profit.buyer_id = '.$accountID.' AND transaction_status = "SOLD" ORDER BY ordered_at ASC';

$stmt_h_d=$conn->query($sql_history_diamond_order_noprofit);
$h_d_found=$stmt_h_d->rowCount();
if($h_d_found){
?>
<li id="no-profit-title-box">&darr; 不能积分的交易纪录 &darr;</li>
<?php
	require_once('_includes/functions/detail_fordiamond_byshape.php');
	foreach($stmt_h_d as $r_h_d){
		$diamondid=$r_h_d['diamondid'];
		$jewellery_id=$r_h_d['jewellery_id'];
		$sold_price_total=$r_h_d['sold_price_total'];
		$ordered_at=$r_h_d['ordered_at'];
		
		$idofuser=$r_h_d['idofuser'];
		$nameofuser=$r_h_d['nameofuser'];
		$crr_email=$r_h_d['email'];
		
		if(empty($nameofuser)){
			$nameofuser=$crr_email;
		}
		
		$buyernametxt='';
		$buyerclass='';
		
		if($idofuser!=$accountID){
			//$buyerclass=' orderedbyclientuser';
			$buyernametxt='<p class="orderredbytxt"><span class="thetxtofusername"><span class="glyphicon glyphicon-user"></span> 预定人: '.$nameofuser.'</span></p>';
		}
		
		
		$listclass='';
		if(strtotime($ordered_at)>=$last_year){
			$listclass=' oneyear';
		}
		
		$crr_dia_shape='不详';
		$crr_dia_weight='不详';
		$crr_dia_color='不详';
		$crr_dia_clarity='不详';
		
		$sql_diamond_crr='SELECT * FROM diamonds WHERE stock_ref = "'.$diamondid.'"';
		foreach($conn->query($sql_diamond_crr) as $r_dia){
			$crr_dia_shape=$r_dia['shape'];
			$crr_dia_weight=$r_dia['carat'];
			$crr_dia_color=$r_dia['color'];
			$crr_dia_clarity=$r_dia['clarity'];
		}
		
		if($jewellery_id!=NULL && $jewellery_id!=''){
			$sql_jewelry_crr='SELECT * FROM jewelry WHERE id = '.$jewellery_id;
			$stmt_j_c=$conn->query($sql_jewelry_crr);
			$crr_j_found=$stmt_j_c->rowCount();
			if($crr_j_found){
				foreach($stmt_j_c as $rj){
					$jewelry_pic=$rj['image1'];
					$jewelry_name=$rj['name'];
				}
			}
		}else{
			$jewelry_pic=NULL;
			$jewelry_name=NULL;
		}
		
		if(isset($jewelry_pic)){
			$crr_order_picture='/_images/jewelry/thumbs/'.$jewelry_pic;
			$crr_dia_pic_for_jew='/_images/constant/'.detail_forDiamond_byShape($crr_dia_shape, 'PICTURE');
		}else{			
			$crr_order_picture='/_images/constant/'.detail_forDiamond_byShape($crr_dia_shape, 'PICTURE');
			$crr_dia_pic_for_jew=NULL;
		}
?>
<li class="dia-order-history noprofit">
<?php echo $buyernametxt; ?>

<!--
<span class="picfororder" style="background-image:url('<?php echo $crr_order_picture; ?>');">
<?php
if(!empty($crr_dia_pic_for_jew)){
?>
<img class="img-dia-for-jew" src="<?php echo $crr_dia_pic_for_jew; ?>" />
<?php
}
?>
</span>
-->

<span class="dia-description"><?php echo detail_forDiamond_byShape($crr_dia_shape, 'NAMECN'); ?>钻石: <?php echo $crr_dia_weight ?>克拉，<?php echo $crr_dia_color; ?>色， 净度 <?php echo $crr_dia_clarity; ?></span>
<span class="jewelry-description">配饰:
<?php
if(isset($jewelry_name)){
	echo $jewelry_name;
}else{
	echo '无';
}
?>
</span>
<span class="totalpricecrr">
<?php require_once('_includes/functions/currency_calculator.php'); ?>
总价格: <?php echo round($sold_price_total, 2); ?>欧元
</span>




<span class="orderedtime">
该交易不可用于积分
<?php 
echo $ordered_at; 
?>
</span>
<br class="clear" />
</li>
<?php
	}
}
?>



</ul>
</div>

<div class="historybox history-diamond">
<h4>首饰购买历史记录</h4>
<ul>
<?php
$total_price_jewelry=0;
//$sql_history_jewelry='SELECT * FROM jewelry_orders WHERE buyer_id = '.$accountID.' AND transaction_status = "SOLD" ORDER BY ordered_at DESC';
$sql_history_jewelry_order='SELECT jewelry_orders.id AS order_id, jewellery_id, jewellery_price, buyer_id, downpayment, order_sent, sold_price_total, extra_info, ordered_at, users.id AS idofuser, users.name AS nameofuser, email, recommenedby FROM jewelry_orders, users WHERE jewelry_orders.buyer_id = users.id AND ( jewelry_orders.buyer_id = '.$accountID.' OR (users.recommenedby = '.$accountID.' AND users.agent = "NO")) AND transaction_status = "SOLD" ORDER BY ordered_at ASC';
$stmt_h_j=$conn->query($sql_history_jewelry_order);
$h_j_found=$stmt_h_j->rowCount();
if($h_j_found){
	foreach($stmt_h_j as $r_h_j){
		$jewellery_id=$r_h_j['jewellery_id'];
		$sold_price_total=$r_h_j['sold_price_total'];
		$ordered_at=$r_h_j['ordered_at'];
		
		$idofuser=$r_h_j['idofuser'];
		$nameofuser=$r_h_j['nameofuser'];
		$crr_email=$r_h_j['email'];
		
		if(empty($nameofuser)){
			$nameofuser=$crr_email;
		}
		
		$buyernametxt='';
		$buyerclass='';
		
		if($idofuser!=$accountID){
			//$buyerclass=' orderedbyclientuser';
			$buyernametxt='<p class="orderredbytxt"><span class="thetxtofusername"><span class="glyphicon glyphicon-user"></span> 预定人: '.$nameofuser.'</span></p>';
		}
		
		
		
		$sql_jewelry_crr='SELECT jewelry.category, need_diamond, jewelry.name AS JEWELRYNAME, image1, category_en, category_cn FROM jewelry, jewelry_category WHERE jewelry.id = '.$jewellery_id.' AND jewelry.category = jewelry_category.id';
		foreach($conn->query($sql_jewelry_crr) as $r_jew){
			$crr_jew_img=$r_jew['image1'];
			$crr_jew_category=$r_jew['category_cn'];
			$crr_jew_name=$r_jew['JEWELRYNAME'];
			$need_diamond=$r_jew['need_diamond'];
		}
		
		
		
		$total_price+=$sold_price_total;
		$total_price_jewelry+=$sold_price_total;
		$listclass='';
		if(strtotime($ordered_at)>=$last_year){
			$listclass=' oneyear';
			$total_price_one_year+=$sold_price_total;
			if($need_diamond=='NO'){
			  $total_piece_one_year+=1;
			}
		}
		
		$crr_order_picture='/_images/jewelry/thumbs/'.$crr_jew_img;
		
?>
<li class="dia-order-history<?php echo $listclass; ?>">
<?php echo $buyernametxt; ?>
<!--
<span class="picfororder" style="background-image:url('<?php echo $crr_order_picture; ?>');"></span>
-->
<span class="jewelry-name"><?php echo $crr_jew_category; ?></span>
<span class="jewelry-description"><?php echo $crr_jew_name; ?></span>
<span class="totalpricecrr">
<?php require_once('_includes/functions/currency_calculator.php'); ?>
总价格: <?php echo round($sold_price_total, 2); ?>欧元
</span>




<span class="orderedtime">
<?php
if(strtotime($ordered_at)<$last_year){
	echo '超一年 不积分 ';
}else{
	echo '<span class="glyphicon glyphicon-ok"></span> 积分有效 ';
	if($need_diamond=='YES'){
		echo '(空托不累计数量)';
	}
}
echo $ordered_at; ?>
</span>
<br class="clear" />
</li>
<?php
	}
}else{
	?>
  <li class="no-order">尚未购买任何首饰</li>
  <?php
}
?>
</ul>
</div>

<div class="historysummary">
<p>您一年内总的购买金额为 <span class="thepricevalue"><?php echo round($total_price_one_year, 2); ?></span>欧元</p>
<p>您一年内共购买<span class="thepricevalue"><?php echo $total_piece_one_year; ?></span>件合适商品。</p>

<?php
if($agent_level_locked=='NO'){

if($agent_level==1){
	if($total_price_one_year<40000 && $total_piece_one_year<12){
		$agent_level_words='您再继续购买'.(12-$total_piece_one_year).'件合适商品，并且这些商品的价值超过'.(40000-$total_price_one_year).'欧元，您即可升级为二级代理，且得到一年消费总额的5%作为返点。';
	}else	if($total_price_one_year<40000 && $total_piece_one_year>=12){
		$agent_level_words='您再继续消费'.(40000-$total_price_one_year).'欧元，您即可升级为二级代理，享受85折价格，且得到一年消费总额的5%作为返点。';
	}else if($total_price_one_year>=40000 && $total_piece_one_year<12){
		$agent_level_words='您再继续购买'.(12-$total_piece_one_year).'件合适商品，您即可升级为二级代理，享受85折价格，且得到一年消费总额的5%作为返点。';
	}else if($total_piece_one_year>=12 && $total_price_one_year>=40000 && ($total_piece_one_year<20 || $total_price_one_year<100000)){
		$agent_level_words='<span class="glyphicon glyphicon-thumbs-up"></span> 您已经可以升级为二级代理，享受85折价格，且得到一年消费总额的5%(须根据购买时的实际折扣核算)作为返点。您在几个小时内会收到确认消息。';
	}else if($total_piece_one_year>=20 && $total_price_one_year>=100000){
		$agent_level_words='<span class="glyphicon glyphicon-thumbs-up"></span> 您已经可以升级为三级代理，享受83折价格，且得到一年消费总额的7%(须根据购买时的实际折扣核算)作为返点。您在几个小时内会收到确认消息。';
	}
}else if($agent_level==2){
	if($total_price_one_year<40000 || $total_piece_one_year<12){
		$agent_level_words='<span class="glyphicon glyphicon-bell"></span> 您一年内的购买量不足二级代理的最低限度，您即将成为一级代理。';
	}else if($total_piece_one_year>=12 && $total_price_one_year>=40000 && ($total_piece_one_year<20 || $total_price_one_year<100000)){
		if($total_piece_one_year>=20){
			$agent_level_words='您再继续消费'.(100000-$total_price_one_year).'欧元，您即可升级为三级代理，享受83折价格，且得到一年消费总额的2%作为返点。';
		}else if($total_price_one_year>=100000){
			$agent_level_words='您再继续购买'.(20-$total_piece_one_year).'件合适商品，您即可升级为三级代理，享受83折价格，且得到一年消费总额的2%作为返点。';
		}else{
		  $agent_level_words='您再继续购买'.(20-$total_piece_one_year).'件合适商品，并且金额超过'.(100000-$total_price_one_year).'，您即可升级为三级代理，享受83折价格，且得到一年消费总额的2%作为返点。';
		}
	}else if($total_piece_one_year>=20 && $total_price_one_year>=100000){
		$agent_level_words='<span class="glyphicon glyphicon-thumbs-up"></span> 您已经可以升级为三级代理，享受83折价格，且得到一年消费总额的2%(须根据购买时的实际折扣核算)作为返点。我们会尽快联系您和您确认。';
	}
}else if($agent_level==3){
	if($total_price_one_year<40000 || $total_piece_one_year<12){
		$agent_level_words='<span class="glyphicon glyphicon-bell"></span> 您一年内的购买量不足三级代理和二级代理的最低限度，您即将成为一级代理。';
	}else if($total_piece_one_year>=12 && $total_price_one_year>=40000 && ($total_piece_one_year<20 || $total_price_one_year<100000)){
		$agent_level_words='<span class="glyphicon glyphicon-bell"></span> 您一年内的购买量不足三级代理的最低限度，您即将成为二级代理。';
	}else if($total_piece_one_year>=20 && $total_price_one_year>=100000){
		$agent_level_words='<span class="glyphicon glyphicon-thumbs-up"></span> 您目前是三级代理，享受83折价格（最高折扣），感谢您带来的傲人成绩。';
	}
}

}else{
	$agent_level_words='您的代理级别被设定为恒定'.$agent_level.'级。永久享有折扣'.priceforagent($agent_level, 1).'。';
}//if($agent_level_locked=='NO'){
?>

<p class="agentlevelwords"><?php echo $agent_level_words; ?></p>


<p>您的全部购买金额为 <span class="thepricevalue"><?php echo round($total_price, 2); ?></span>欧元</p>


</div>

</div>

<!-- history recommended users-->
<div class="history-recommenedusers" id="myclients">
<h3>我的客户</h3>
<div class="theuserslist">
<?php


//$sql_rec_users='SELECT id, name, email, recommenedby FROM users WHERE belong_to_agent='.$accountID;
$sql_rec_users='SELECT id, name, email, icon, recommenedby FROM users WHERE recommenedby='.$accountID;
$stmt_r_u=$conn->query($sql_rec_users);
$r_u_found=$stmt_r_u->rowCount();
if($r_u_found){
?>
<ul class="recommenduserslist">
<?php
foreach($stmt_r_u as $r_ru){
	$ru_id=$r_ru['id'];
	$ru_name=$r_ru['name'];
	$ru_email=$r_ru['email'];
	$ru_icon=$r_ru['icon'];
	
	if(empty($ru_icon)){
		$ru_icon='/_images/constant/noicon.png';
	}
	
	if($ru_name==NULL || $ru_name==''){
		$crr_ru_name=$ru_email;
	}else{
		$crr_ru_name=$ru_name;
	}
	
	$sql_crr_u_order_d='SELECT COUNT(*) AS order_d_num FROM diamonds_orders WHERE buyer_id = '.$ru_id.' AND transaction_status = "ORDERED"';
	foreach($conn->query($sql_crr_u_order_d) as $ruod){
		$order_d_num=$ruod['order_d_num'];
	}
	$sql_crr_u_order_j='SELECT COUNT(*) AS order_j_num FROM jewelry_orders WHERE buyer_id = '.$ru_id.' AND transaction_status = "ORDERED"';
	foreach($conn->query($sql_crr_u_order_j) as $ruoj){
		$order_j_num=$ruoj['order_j_num'];
	}
	$totalordernum=$order_d_num+$order_j_num;
	
	
	$sql_crr_ru_all_d='SELECT SUM(sold_price_total) AS totalprice_d_for_crr_user_all, COUNT(*) AS pieces_d FROM diamonds_orders WHERE buyer_id = '.$ru_id.' AND transaction_status = "SOLD"';
	foreach($conn->query($sql_crr_ru_all_d) as $crr_u_a){
		$crr_user_total_all_d=$crr_u_a['totalprice_d_for_crr_user_all'];
		$crr_user_total_pieces_d=$crr_u_a['pieces_d'];
	}
	
	
	$sql_crr_ru_all_j='SELECT SUM(sold_price_total) AS totalprice_j_for_crr_user_all, COUNT(*) AS pieces_j FROM jewelry_orders WHERE buyer_id = '.$ru_id.' AND transaction_status = "SOLD"';
	foreach($conn->query($sql_crr_ru_all_j) as $crr_u_a){
		$crr_user_total_all_j=$crr_u_a['totalprice_j_for_crr_user_all'];
		$crr_user_total_pieces_j=$crr_u_a['pieces_j'];
	}
	
	$crr_user_total_all=$crr_user_total_all_d+$crr_user_total_all_j;
	$crr_user_total_pieces=$crr_user_total_pieces_d+$crr_user_total_pieces_j;
?>
<li>



<span class="clientpic" style="background-image:url(<?php echo $ru_icon; ?>)"></span>

<span class="ru_name"><?php echo $crr_ru_name; ?></span>




<?php
if(!empty($totalordernum)){
?>
<span class="neworder"> ( <span class="glyphicon glyphicon-gift"></span> 有新预订 ) </span>
<?php
}
?>

<span class="usertotalamout">
一年内消费总额: <?php echo $crr_user_total_all; ?>欧元 (共购买<?php echo $crr_user_total_pieces; ?>件商品)

<a target="_blank" class="userhistorybtn" href="userhistory.php?id=<?php echo $ru_id; ?>">
<span class="userdetailbtn">详情 &raquo;</span>
</a>
</span>

<a class="chatlinker" href="/message-agent.php?c=<?php echo $ru_id; ?>"><span class="glyphicon glyphicon-comment"></span> 消息
<?php
$sql_message_num='SELECT COUNT(*) AS mn FROM message WHERE user_id = '.$ru_id.' AND message_read = "NO"';
$stmt_mn=$conn->query($sql_message_num);
foreach($stmt_mn as $rowmn){
	$usermn=$rowmn['mn'];
}
if($usermn){
	$mn_class_add=' newmessage';	
	$mntxt=$usermn.'则新消息';
}else{
	$mn_class_add='';
	$mntxt='无新消息';
}
?>
 
<span class="messagemn<?php echo $mn_class_add; ?>">
[ <?php echo $mntxt; ?> ]
</span>

</a>


<br class="clear" />
</li>
<?php	
}
?>
</ul>
<?php
}else{
?>
<p>您尚未推荐任何用户</p>
<?php
}
?>
</div>
</div>




<!-- history extended users-->
<div class="history-recommenedusers extendedclients">
<h3>客户发展的客人</h3>
<div class="theuserslist">
<?php
$sql_rec_e_users='SELECT id, name, email, recommenedby FROM users WHERE belong_to_agent='.$accountID.' AND relation_level > 1  AND relation_level <= 6';
$stmt_r_u_e=$conn->query($sql_rec_e_users);
$r_u_e_found=$stmt_r_u_e->rowCount();
if($r_u_e_found){
?>
<ul class="recommenduserslist">
<?php
foreach($stmt_r_u_e as $r_rue){
	$rue_id=$r_rue['id'];
	$rue_name=$r_rue['name'];
	$rue_email=$r_rue['email'];
	
	if($rue_name==NULL || $rue_name==''){
		$crr_rue_name=$rue_email;
	}else{
		$crr_rue_name=$rue_name;
	}
	
	$crr_usere_total_all_d=0;
	$crr_usere_total_all_j=0;
	
	$sql_crr_rue_all_d='SELECT SUM(sold_price_total) AS totalprice_d_for_crr_user_all, COUNT(*) AS pieces_d FROM diamonds_orders WHERE buyer_id = '.$rue_id.' AND transaction_status = "SOLD"';
	foreach($conn->query($sql_crr_rue_all_d) as $crr_ue_a){
		$crr_usere_total_all_d=$crr_ue_a['totalprice_d_for_crr_user_all'];
	}
	
	
	$sql_crr_rue_all_j='SELECT SUM(sold_price_total) AS totalprice_j_for_crr_user_all, COUNT(*) AS pieces_j FROM jewelry_orders WHERE buyer_id = '.$rue_id.' AND transaction_status = "SOLD"';
	foreach($conn->query($sql_crr_rue_all_j) as $crr_ue_a){
		$crr_usere_total_all_j=$crr_ue_a['totalprice_j_for_crr_user_all'];
	}
	
	$crr_usere_total_all=$crr_usere_total_all_d+$crr_usere_total_all_j;
?>
<li>



<span class="ru_name"><?php echo $crr_rue_name; ?></span>

<?php
if(!empty($crr_usere_total_all)){
?>
<span class="ru_price newpoints">(总消费: <?php echo $crr_usere_total_all; ?>欧元)</span>
<?php
}
?>
<a target="_blank" class="userhistorybtn" href="userhistory.php?id=<?php echo $rue_id; ?>">
详情 &raquo;</a>
<br class="clear" />
</li>
<?php	
}
?>
</ul>
<?php
}else{
?>
<p>您的客人尚未推荐任何用户</p>
<?php
}
?>
</div>
</div>





<!-- discount ticket -->
<div class="generalinfo heshibi-box discountticketbox" id="coupon">
<h3>我的打折券</h3>
<?php
$sql_dt='SELECT * FROM discount_tickets WHERE status = "ACTIVE" AND user = '.$accountID.' ORDER BY id DESC';
$stmt_dt=$conn->query($sql_dt);
$dtfound=$stmt_dt->rowCount();
if($dtfound){
foreach($stmt_dt as $rdt){
	$crr_dt_value=$rdt['discount_value'];
	$crr_dt_created=$rdt['created'];
?>
<span class="discountticketcontainer">
<img src="/_images/constant/<?php echo 'discount'.($crr_dt_value*100); ?>.jpg" class="dt-pic" />
  <span class="dt-time">
      <span class="glyphicon glyphicon-time"></span> 有效期: 
      <?php 
      echo date('Y年m月d日', strtotime(date('Y-m-d', strtotime($crr_dt_created)+8*60*60).' +1 Month')); 
      ?>
  </span>
</span>
<?php
}
}else{
	echo '没有可以使用的打折券';
}
?>
<p class="discountticketexplwords">
您购买商品时，系统会自动选择最早的一张打折卡添加到您的订单中。<br />
打折券的折扣和您自己折扣不能累加，以两者中折扣较高的为准。
</p>

<h4 class="heshibitranferbtnbox dt"><a href="senddiscountticket.php"><span class="glyphicon glyphicon-hand-right"></span> 发打折券给我的客人 &raquo; </a></h4>
<p class="heshibitransfer-expl"><span class="glyphicon glyphicon-info-sign"></span> 您可以给您的客人发打折券作为他们购买时的折扣凭证</p>
</div>



<!-- heshibi -->
<div class="generalinfo heshibi-box">
<h3>我的合适币（代金劵）</h3>
<div class="my-heshibi">
<div class="my-heshibi-amount-box">
余额: <span id="heshibi-balance-value"><?php echo round($balance, 2); ?></span>€
</div>
</div>

<p>合适币余额: <?php echo round($balance, 2); ?>欧元 (<?php echo round(EuroToYuan($balance), 2); ?>元人民币, <?php echo round(EuroToDollar($balance), 2); ?>美元)</p>

<h4 class="heshibitranferbtnbox"><a href="heshibitransfer.php"><span class="glyphicon glyphicon-hand-right"></span> 发合适币给我的客人 &raquo; </a></h4>
<p class="heshibitransfer-expl"><span class="glyphicon glyphicon-info-sign"></span> 如果客人一星期内不使用，金额自动返回</p>

</div>

  </div>
</template>
