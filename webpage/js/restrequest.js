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
			document.getElementById('productsList').innerHTML = renderJresult(result, "没有对应的产品")
		
			 return true
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
		
		 document.getElementById('productsList').innerHTML = renderJresult(result, "没有对应ID的产品")
			return true
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
			 document.getElementById('productsList').innerHTML = renderJresult(result, "没有对应的产品")
			 return true
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

 function getproductsJbyfilter(Jtype) {
	var auth_token = accessToken;
	var url_base = baseURI + 'api/products/filter/'+ Jtype;
	var requestPayload = {
		// search prosuct by Stok_ID
		'material':  document.getElementById('materialselection').value,
		'price': document.getElementById('priceselection').value,
		'mounting': document.getElementById('mountingtypeselection').value,
		'diashape': document.getElementById('diashapeselection').value,
		'smalldiaschoice': document.getElementById('smalldiaschoiceselection').value,
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
			 document.getElementById('productsList').innerHTML = renderJresult(result, "没有对应的产品")
			 return true
		 },
		 'error': function (XMLHttpRequest, textStatus, errorThrown) {
		   //Process error actions
			 document.getElementById('productsList').innerHTML = XMLHttpRequest.status + ' ' + XMLHttpRequest.statusText;
		   return false;
		 }
	 });
 }


 function renderJresult(result, message)
 {
	if (result == null)
	{
		Prosuctbox = message
	}
	else 
	{
	 var Prosuctbox = '';

	 for(var i=0;i<result.length;i++)
	 {
		 var prosuctDataline = '<div class="jewelrybox complete"><a class="seedetailbtn-big demo-box" href="jewelrydetail.html?id='+ result[i]["id"] + '">' +
		 '<span class="imageholder" style="background-image:url("/pic/jewelry/thumbs/' + result[i]["images"] + '")"></span>'+ 
		 '<span class="jewelryname">'+ result[i]["name"]+'</span>'+
		 '<span class="jewelryprice">'+ result[i]["price"]+' EUR </span>'+
		 '<span class="stocknum">'+ result[i]["stock_quantity"]+'</span>' +
		 '</a><p class="actionbox"><a class="seedetailbtn" href="jewelrydetail.php?id=1299"><span class="glyphicon glyphicon-eye-open"></span> 详情</a><a class="choosebtn" href=""><span class="glyphicon glyphicon-gift"></span> 购买</a></p><span class="indi-icons-container"><span class="glyphicon glyphicon-film"></span></span></div>';
		 Prosuctbox = Prosuctbox+prosuctDataline;
	 }   
	
	}
	return Prosuctbox
 }
