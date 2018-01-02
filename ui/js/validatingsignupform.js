// JavaScript Document
function processthesignupform(){
	if($.trim($('input#newaccountemail').val())==''){
		alert('请输入您的邮箱!');
		return;
	}
	if($.trim($('input#newaccountpassword').val())==''){
		alert('请设定您的密码!');
		return;
	}
	if($.trim($('input#newaccountpassword').val()).length<6){
		alert('密码太短，请输入至少6位!');
		return;
	}
	if($.trim($('input#newaccountpassword').val())!=$.trim($('input#newaccountpasswordagain').val())){
		alert('两次输入的密码不一致，请重新输入');
		return;
	}
	if($.trim($('input#realname').val())==''){
		alert('请留下您的姓名！');
		return;
	}
	if($.trim($('input#wechatid').val())=='' && $.trim($('input#tel').val())==''){
		alert('请留下您的微信ID或电话号码，以确保我们能够联系到您');
		return;
	}
	$('form#the_signup_form').submit();
}