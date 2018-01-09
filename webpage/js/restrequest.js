 //Get inContact Token
 var accessToken = '123123';
 var baseURI = 'http://localhost:8080/';

 function userregister()
 {
		var auth_token = 'AppName@VendorName:BusinessUnit';
		
		var url_base = baseURI + 'api/users';
		
		var requestPayload = {
			'cellphone': $('input#tel').val(),
			'username': $('input#realname').val(),
			'password':  $('input#newaccountpassword').val(),
			'email': $('input#newaccountemail').val(),
			'user_type': 'admin'
		};

	$.ajax({
		'url': url_base,
		'type': 'POST',
		'content-Type': 'multipart/form-data',
		'headers': {
			// Use access_token previously retrieved from inContact token 
			// service.
			'Authorization': 'basic ' + auth_token
		},
		'data': requestPayload,
		'success': function (result) {
			//Process success actions
		//	accessToken = result.access_token;
			baseURI = result.resource_server_base_uri;
			document.getElementById('userloginform').innerHTML = "注册成功！";
			setTimeout(document.location.href = 'login.html',"5000")
			return result;
		},
		'error': function (XMLHttpRequest, responseText) {
			//Process error actions
			document.getElementById('userloginform').innerHTML ='Error: ' + XMLHttpRequest.responseText ;
			sleep(2000);
			document.location.href = 'register.html';
			return false;
		}
 	});
}


 function userlogin()
 {
		var auth_token = 'AppName@VendorName:BusinessUnit';
		
		var url_base = baseURI + 'api/login';
		
		var requestPayload = {
			'username': $('input#loginkey').val(),
			'password':  $('input#accountpassword').val(),
		};

	$.ajax({
		'url': url_base,
		'type': 'POST',
		'content-Type': 'multipart/form-data',
		'headers': {
			// Use access_token previously retrieved from inContact token 
			// service.
			'Authorization': 'basic ' + auth_token
		},
		'data': requestPayload,
		'success': function (result) {
			//Process success actions
		//	accessToken = result.access_token;
			baseURI = result.resource_server_base_uri;
			//document.getElementById('userloginform').innerHTML = "注册成功！";
			setTimeout(document.location.href = 'home.html',"5000")
			return result;
		},
		'error': function (XMLHttpRequest, responseText) {
			//Process error actions
			document.getElementById('userloginform').innerHTML ='Error: ' + XMLHttpRequest.responseText ;
			sleep(2000);
			document.location.href = 'login.html';
			return false;
		}
 	});
}

 function getToken() {
	 var url_base = 
		'';

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
 
 
 
 

