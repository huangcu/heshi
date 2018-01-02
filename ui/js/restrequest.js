 //Get inContact Token
 var accessToken = '';
 var baseURI = '';

 function getToken() {
	 var url_base = 
		'https://api.incontact.com/InContactAuthorizationServer/Token';

	 // The auth_token is the base64 encoded string for the API 
	 // application.
	 var auth_token = 'AppName@VendorName:BusinessUnit';
	 auth_token = window.btoa(auth_token);
	 var requestPayload = {
		 // Enter your inContact credentials for the 'username' and 
		 // 'password' fields.
		 'grant_type': 'password',
		 'username': 'YourUsernameHere',
		 'password': 'YourPasswordHere',
		 'scope': ''
	 }
	 $.ajax({
		 'url': url_base,
		 'type': 'POST',
		 'content-Type': 'x-www-form-urlencoded',
		 'dataType': 'json',
		 'headers': {
		   // Use access_token previously retrieved from inContact token 
		   // service.
		   'Authorization': 'basic ' + auth_token
		 },
		 'data': requestPayload,
		 'success': function (result) {
		   //Process success actions
		   accessToken = result.access_token;
		   baseURI = result.resource_server_base_uri;
		   alert('Success!\r\nAccess Token:\r' + accessToken + 
			   '\r\nBase URI:\r' + baseURI)
		   document.getElementById('pageDiv').innerHTML = result.access_token;
		   return result;
		 },
		 'error': function (XMLHttpRequest, textStatus, errorThrown) {
		   //Process error actions
		   alert('Error: ' + errorThrown);
		   console.log(XMLHttpRequest.status + ' ' + 
			   XMLHttpRequest.statusText);
		   return false;
		 }
	 });
   }

 // PUT CALL BELOW HERE!!!

 // BU Agents List
 function getAgentList() {
	 // The baseURI variable is created by the result.base_server_base_uri 
	 // which is returned when getting a token and should be used to 
	 // create the url_base.
	 var url_base = baseURI;
	 $.ajax({
		 'url': url_base + '/services/{version}/agents',
		 'type': 'GET',
		 'content-Type': 'x-www-form-urlencoded',
		 'dataType': 'json',
		 'headers': {
		   // Use access_token previously retrieved from inContact token 
		   // service.
		   'Authorization': 'bearer ' + accessToken
		 },
		 'success': function (result) {
		   //Process success actions
		   var returnResult = JSON.stringify(result);
		   alert('Success!\r\n' + returnResult);
		   document.getElementById('callResults').innerHTML = returnResult;
		   return result;
		 },
		 'error': function (XMLHttpRequest, textStatus, errorThrown) {
		   //Process error actions
		   alert('Error: ' + errorThrown);
		   console.log(XMLHttpRequest.status + ' ' + 
			   XMLHttpRequest.statusText);
		   return false;
		 }
	 });
 }
  
 //END CALL ABOVE HERE




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