package controllers

import (
	"log"
	"myapp/app"
	"myapp/app/controllers/api"

	"github.com/dchest/captcha"
	"github.com/revel/revel"
)

type Auth struct {
	BaseController
}

func (a *Auth) checkAuth(login bool, ok *app.Api) bool {
	if a.Validation.HasErrors() {
		a.Validation.Keep()
		a.FlashParams()
		if login {
			ok.Next = app.NextJson{"href", "/login"}
		} else {
			ok.Next = app.NextJson{"href", "/register"}
		}
		return false
	}
	return true
}

func (a Auth) Register() revel.Result {
	if v, ok := a.Session["user"]; ok {
		log.Println("register has session:", v)
	}
	// Captcha := struct {
	// 	CaptchaId string
	// }{
	// 	captcha.New(),
	// }
	// a.RenderArgs["Captcha"] = Captcha
	return a.RenderTemplate("home/register.html")
}

func (a Auth) DoRegister(email, pwd, captcha, captchaId string) revel.Result {
	ok := app.NewOk()
	ret := api.RigisterResponse{}
	ret.Username = "hello"
	ret.EmailProvider = "http://mail.qq.com"
	ret.Email = "hello@qq.com"
	ok.Item = ret
	a.Session["user"] = ret.Username
	log.Println("set register session:", a.Session)
	return a.RenderJson(ok)
	// return a.Redirect("/items")
}

func (a Auth) Logout() revel.Result {
	delete(a.Session, "user")
	return a.RenderTemplate("home/index.html")
}

func (a Auth) Login() revel.Result {
	a.RenderArgs["needCaptcha"] = "true"
	a.RenderArgs["openRegister"] = "true"
	Captcha := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	a.RenderArgs["Captcha"] = Captcha
	log.Println("captchaId:", Captcha.CaptchaId)
	return a.RenderTemplate("home/login.html")
}

func (a Auth) DoLogin(email, pwd, validationCode, captchaId string) revel.Result {
	log.Println("email:", email, "validationCode:", validationCode, "captchaId:", captchaId)
	ok := app.NewOk()
	ok.Ok = a.Validation.Required(captcha.VerifyString(captchaId, validationCode)).Ok
	if !ok.Ok {
		ok.Msg = "captcha"
		return a.RenderJson(ok)
	}
	ok.Ok = a.Validation.Required(email).Ok
	if !ok.Ok {
		ok.Msg = "email"
		return a.RenderJson(ok)
	}
	ok.Ok = a.Validation.Email(email).Ok
	if !ok.Ok {
		ok.Msg = "email"
		return a.RenderJson(ok)
	}

	if !a.checkAuth(true, &ok) {
		ok.Msg = "login"
		ok.Ok = false
	} else {
		ok.Next = app.NextJson{"href", "/index"}
		a.Session["user"] = "hello"
	}

	log.Println("set register session:", a.Session, "with resp:", ok)
	return a.RenderJson(ok)
}
