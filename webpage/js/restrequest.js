 var accessToken = 'Jbm6XfXQj/KqmMTqz6c4GQWl9U6JMLQ/T4LzPWIEi2W2Q23GDkuIfxvbUC/rar8ZJIWWSVo68fZ/hv6n0oAeXaQKEfhKmGUZ8m8JHm5TteBZwqZuqXAbOeowTJVBn8aaUhfSfZbmgNnXwDEnhjZ1DZ8jG2Khy9uzoHu5ogwbVHQ=';
 var baseURI = 'http://localhost:8080/';

 function userregister()
 {
		var auth_token = accessToken;
		
		var url_base = baseURI + 'api/users';
		
		var requestPayload = {
			'cellphone': $('input#tel').val(),
			'username': $('input#realname').val(),
			'password':  $('input#newaccountpassword').val(),
			'email': $('input#newaccountemail').val(),
			'user_type': 'CUSTOMER'
		};

	$.ajax({
		'url': url_base,
		'type': 'POST',
		'content-Type': 'multipart/form-data',
		'headers': {
			// Use access_token 
			'X-Auth-Token': auth_token
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
		var auth_token = accessToken;
		
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
			'X-Auth-Token': auth_token
		},
		'data': requestPayload,
		'success': function (result) {
			//Process success actions
		//	accessToken = result.access_token;
			baseURI = result.resource_server_base_uri;
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
		   document.getElementById('pageDiv').innerHTML = result.access_token;
		   return result;
		 },
		 'error': function (XMLHttpRequest, textStatus, errorThrown) {
		   //Process error actions
		   console.log(XMLHttpRequest.status + ' ' + 
			   XMLHttpRequest.statusText);
		   return false;
		 }
	 });
   }

 // List products
 function getproductsList(Jtype) {
	var auth_token = accessToken;
	var url_base = baseURI + 'api/products/'+ Jtype;
	 $.ajax({
		 'url': url_base,
		 'type': 'GET',
		 'content-Type': 'x-www-form-urlencoded',
		 'dataType': 'json',
		 'headers': {
		   // Use access_token 
		   'X-Auth-Token': auth_token
		 },
		 'success': function (result) {
		   //Process success actions
		   var returnResult = JSON.stringify(result);
		   document.getElementById('productsList').innerHTML = returnResult;
		   return result;
		 },
		 'error': function (XMLHttpRequest, textStatus, errorThrown) {
		   //Process error actions
			 document.getElementById('productsList').innerHTML = XMLHttpRequest.status + ' ' + XMLHttpRequest.statusText;
		   return false;
		 }
	 });
 }

 function getproductbyID(Jtype, ID) {
	var auth_token = accessToken;
	var url_base = baseURI + 'api/products/search/'+ Jtype;
	var requestPayload = { 
		// search prosuct by Stok_ID
		'ref': $('input[name=searchref]').val(),
	};
	 $.ajax({
		 'url': url_base,
		 'type': 'POST',
		 'content-Type': 'x-www-form-urlencoded',
		 'dataType': 'json',
		 'data': requestPayload,
		 'headers': {
		   // Use access_token 
		   'X-Auth-Token': auth_token
		 },
		 
		 'success': function (result) {
		   //Process success actions
		   var returnResult = JSON.stringify(result);
		   document.getElementById('productsList').innerHTML = returnResult;
		   return result;
		 },
		 'error': function (XMLHttpRequest, textStatus, errorThrown) {
		   //Process error actions
			 document.getElementById('productsList').innerHTML = XMLHttpRequest.status + ' ' + XMLHttpRequest.statusText;
		   return false;
		 }
	 });
 }

 function getproductsbyCategory(Jtype, category) {
	var auth_token = accessToken;
	var url_base = baseURI + 'api/products/filter/'+ Jtype;
	var requestPayload = {
		// search prosuct by Stok_ID
		'category': category,
	};
	 $.ajax({
		 'url': url_base,
		 'type': 'POST',
		 'content-Type': 'x-www-form-urlencoded',
		 'dataType': 'json',
		 'data': requestPayload,
		 'headers': {
		   // Use access_token 
		   'X-Auth-Token': auth_token
		 },
		 
		 'success': function (result) {
		   //Process success actions
		   var returnResult = JSON.stringify(result);
		   document.getElementById('productsList').innerHTML = returnResult;
		   return result;
		 },
		 'error': function (XMLHttpRequest, textStatus, errorThrown) {
		   //Process error actions
			 document.getElementById('productsList').innerHTML = XMLHttpRequest.status + ' ' + XMLHttpRequest.statusText;
		   return false;
		 }
	 });
 }

 function getUserinfro()
 {
	var auth_token = accessToken;
		
	var url_base = baseURI + '/api/admin/users';
	
	var requestPayload = {	};

$.ajax({
	'url': url_base,
	'type': 'POST',
	'content-Type': 'multipart/form-data',
	'headers': {
		// Use access_token previously retrieved from inContact token 
		// service.
		'X-Auth-Token': auth_token
	},
	'data': requestPayload,
	'success': function (result) {
		//Process success actions
	//	accessToken = result.access_token;
		baseURI = result.resource_server_base_uri;
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


  

 
 
 
 

