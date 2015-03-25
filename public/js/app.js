function freshCaptcha() {
	var $captchaImg = $('#captchaImg');
	var src = $captchaImg.attr("src");
	$captchaImg.attr("src", src);
}