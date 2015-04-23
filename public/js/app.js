function freshCaptcha() {
	var $captchaImg = $('#captchaImg');
	var src = $captchaImg.attr("src");
	var freshStr = "?reload=" + (new Date()).getTime();
	$captchaImg.attr("src", src + freshStr);
}

function redirectTo(link, delay) {
	setTimeout(function() {
		location.href = link
	}, delay)
}

function SendOrder() {
	var name = "刘丽萍"
	var phone = "18615774413"
	$.ajax({
	    type: "POST",
	    async: false,
	    cache: false,
	    dataType: 'json',
	    contentType: 'application/json; charset=utf-8',
	    url: "naiping.pospal.cn/order/place",
	    data: JSON.stringify({
	        customerName: name,
	        customerPhoneNum: phone,
	        customerAddress: "test",
	        orderComment: "",
	        paymentMethod: "CustomerBalance"
	    }),
	    error: function (XmlHttpRequest, textStatus, errorThrown) { alert("下单出现错误！"); },
	    success: function (data) {
	        if (data.status) {
	            if (data.isHtml) {
	                window.open(data.redirectUrl);
	                $("body").append(data.html);
	            }
	            else {
	                window.location.replace(data.redirectUrl);
	            }
	        }
	        else {
	            alert(data.message);
	        }
		}
	});
}