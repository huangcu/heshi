

<?php
/*===================session========================*/
session_start();

require_once('admin_authorizing_part.php');

require_once('../_includes/functions/connection.php');
$conn=dbConnect('write','pdo');
$conn->query("SET NAMES 'utf8'");
/*

if(isset($_POST['dia_shape'])){
	exit();
}
*/

if(isset($_POST['name']) && isset($_POST['username']) && isset($_POST['password'])){
	
	
	$name=$_POST['name'];
	$username=$_POST['username'];
	$password=$_POST['password'];
	$icon=$_POST['icon'];
	
	
	
	
	if(empty($name) || empty($username) || empty($password)){
		$message_db='请添写所有信息';
	}
	
	
	$sql_admin_check='SELECT * FROM admins WHERE username = "'.$username.'"';
	$stmt_check=$conn->query($sql_admin_check);
	$found=$stmt_check->rowCount();
	if($found){
		$message_db='该用户名已经使用了，请使用其它的用户名';
	}
	
	
	if(empty($icon)){
		$icon=NULL;
	}
	
	
	if(!isset($message_db)){
		
	$sql_insert='INSERT INTO admins (username, password, real_name, icon) 
	VALUES(:username, :password, :real_name, :icon)';
	
	
	$stmt=$conn->prepare($sql_insert);	  
	$stmt->bindParam(':username', $username, PDO::PARAM_STR);
	$stmt->bindParam(':password', $password, PDO::PARAM_STR);
	$stmt->bindParam(':real_name', $name, PDO::PARAM_STR);
	$stmt->bindParam(':icon', $icon, PDO::PARAM_STR);
	
	$stmt->execute();
	$OK=$stmt->rowCount();
	
	
	if($OK){
		$message_db="添加成功";
		
		
		//////////############################################################
		//////////############################################################
		//////////############################################################ 现在加到微信客服里 !!!!!!!!!!!!!!!!!!
		//////////############################################################
		//////////############################################################
		
		
		//1添加客服
		require_once('../_weixin/get_access_token_function.php');
		//$theaccesstoken
		$theaccesstoken=getAccessToken(false);
		
		$data='{
				 "kf_account" : "'.$username.'@admin",
				 "nickname" : "'.$name.'",
				 "password" : "beyoudiamond-antwerp-1234",
		}';
		
		
		$kfURL="https://api.weixin.qq.com/customservice/kfaccount/add?access_token=".$theaccesstoken;
		
		$ch = curl_init(); 
		
		curl_setopt($ch, CURLOPT_URL, $kfURL); 
		curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "POST");
		curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, FALSE); 
		curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, FALSE);
		curl_setopt($ch, CURLOPT_FOLLOWLOCATION, 1);
		curl_setopt($ch, CURLOPT_AUTOREFERER, 1); 
		curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true); 
		
		$result = curl_exec($ch);
		
		if (curl_errno($ch)) {
				echo 'Errno'.curl_error($ch);
		}
		
		curl_close($ch);
		
		$obj_reply = json_decode($result);
		$reply_feedback=$obj_reply->{'errmsg'};
		$reply_errorcode=$obj_reply->{'errcode'};
		
		if(trim($reply_feedback)=='ok'){
			$message_db="添加成功,并且已经添加微信公众号客服";
			
			$sql_update_admin='UPDATE admins SET wechat_kefu = "YES" WHERE username = "'.$username.'"';
			$stmt_admin_update=$conn->query($sql_update_admin);
		}
		
		
		//2 添加头像				
		if(!empty($icon)){
			$file_name_with_full_path = realpath('../_images/admins/'.$icon);
			$data=array('extra_info' => 'kefu-icon','file_contents'=>'@'.$file_name_with_full_path);
			
			
			$iconURL="https://api.weixin.qq.com/customservice/kfaccount/uploadheadimg?access_token=".$theaccesstoken.'&kf_account='.$username.'@admin';

			
			$ch = curl_init(); 
			
			curl_setopt($ch, CURLOPT_URL, $iconURL); 
			curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "POST");
			curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, FALSE); 
			curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, FALSE);
			curl_setopt($ch, CURLOPT_FOLLOWLOCATION, 1);
			curl_setopt($ch, CURLOPT_AUTOREFERER, 1); 
			curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
			curl_setopt($ch, CURLOPT_RETURNTRANSFER, true); 
			
			$info = curl_exec($ch);
			
			if (curl_errno($ch)) {
					echo 'Errno'.curl_error($ch);
			}
			
			curl_close($ch);
		}
		
		//echo "db ok";
	}else{
		//echo "db no ok";
		$error=$stmt->errorInfo();
			if(isset($error[2])){
				$error=$error[2];
				//echo $error;
			}
	}
	}
}

?>


<template>
  <div>
    <vue-title title="添加管理人员"></vue-title>
    <h1>添加普通管理人员</h1>
    <p style="text-align:center;">普通管理人员权限：仅网站库存</p>
    <p v-if="failure" class='alert'>{{ failure}}</p>
    <p v-if="error" class='alert'>{{ error}}</p>
    <p v-if="message_db" class='alert'>{{ message_db}}</p>

  </div>
</template>
<script src="./NewAdmin.js"></script>
<style src="./NewAdmin.css" scoped></style>





</head>

<body>
<!--
<div class="mnavic">
<a class="mnavi" href="index.php">SUBMISSION</a><a class="mnavi" href="list.php">MANAGEMENT</a>
<a class="mnavi" href="banner.php">BANNER</a>
</div>
-->
<?php
include('navi.php');
?>
<hr />

<h1>添加普通管理人员</h1>
<p style="text-align:center;">普通管理人员权限：仅网站库存</p>
<?php 
if(isset($failure)){
    echo "<p class='alert'>"."$failure"."</p>";
}
if(isset($error)){
	//print_r($error);
    echo "<p class='alert'>"."$error"."</p>";
}
if(isset($message_db)){
	echo "<p class='alert'>"."$message_db"."</p>";
}


?>





<div id="contentwrapper">


<form action="" method="post" id="addnewadmin">
  
  

  
<p>  
<label>姓名: </label><span>(必填)</span>

  <input name="name" id="name" type="text" class="formbox" size="200"/>
  </p>
  
  <p>
  <label>用户名: </label><span>(英文字母，必填)</span>
  <input name="username" id="username" type="text" class="formbox" size="200"/>
  </p>
  
  <p>
  <label>密码: </label><span>(英文字母或数字，必填)</span>
  <input name="password" id="password" type="text" class="formbox" size="200"/>
  </p>
  
  <input type="hidden" name="icon" id="icon" value="" />
  
  
</form> 



<div class="picture" style="position:relative;">
<div class="picture-preview-box"></div>
    <form action="uploadingimage.php" method="post" enctype="multipart/form-data" id="uploadingadminicon" class="uploadImage">
    <input type="hidden" name="picturewhere" value="adminicon" /> 
    <label>添加头像:</label>
    <input type="file" name="image" id="image" class="image_selecting">                
    </form>
    <br style="clear:both;" />
    <div class="progress" style="width:300px; position:relative; top:-320px; left:488px;">
    <p>上传中</p>
        <div class="bar"></div >
        <div class="percent">0%</div >
    </div>    
</div>




<p style="padding-top:30px; text-align:center; padding-left:0; background:none;">
<button onclick="submittheform()" type="button" id="submit_article" style="font-size:18px; font-weight:bold; padding:12px 50px; background-color:#393; color:#FFF; border-width:1px;"><span class="glyphicon glyphicon-send" style="color:#FFF; left:0;"></span> &nbsp;  提交 </button>
</p>
  
  
</div>
  

<div id="status" style="display:none;"></div>
<div id="iframecontainer" style="display:none;"></div>
<br /><br /><br /><br /><br /><br />

<?php
include_once('footer.php');
?>
</body>
</html>