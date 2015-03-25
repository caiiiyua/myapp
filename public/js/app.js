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