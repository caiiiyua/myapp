package controllers

import (
	"fmt"
	"log"
	"myapp/app"
	"myapp/app/controllers/api"

	"github.com/dchest/captcha"
	"github.com/revel/revel"
)

type Auth struct {
	BaseController
}

func (a *Auth) checkAuth() bool {
	if a.Validation.HasErrors() {
		a.Validation.Keep()
		a.FlashParams()
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

	a.Validation.Required(!a.userService().ExistsUserByEmail(email)).Message(
		"userHasBeenRegistered", email)
	if ret := a.checkAuth(); !ret {
		ok.Ok = false
		ok.Msg = "exists"
		log.Println(email, "has already been registed")
		return a.RenderJson(ok)
	}

	user, err := a.userService().RegistUserByEmail(email, pwd)
	if err != nil {
		a.Flash.Error(err.Error())
		ok.Ok = false
		ok.Msg = "error"
		log.Println(email, "registed with error:", err)
		return a.RenderJson(ok)
	}

	// send activation mail to user
	SendMail(a.Message("activationMail"), fmt.Sprintf(
		`<a href="http://localhost:9000/activate?activationCode=%s&email=%s">%s</a>`,
		user.ActivationCode, email, a.Message("activation")), email)

	ret := api.RigisterResponse{}
	ret.Username = email
	ret.EmailProvider = EmailProvider(email)
	ret.Email = email
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

	if !a.checkAuth() {
		ok.Msg = "login"
		ok.Ok = false
	} else {
		ok.Next = app.NextJson{"href", "/index"}
		a.Session["user"] = email
		if email == "879939101@qq.com" {
			ok.Msg = "admin"
		}
	}

	log.Println("set register session:", a.Session, "with resp:", ok)
	return a.RenderJson(ok)
}

func (a *Auth) Activate(activationCode, email string) revel.Result {
	revel.INFO.Println("activationCode:", activationCode)
	user, err := a.userService().Activate(email, activationCode)
	revel.INFO.Println("Activate user:", user)

	if err != nil {
		a.Flash.Error(err.Error())
	} else {
		a.Flash.Success(a.Message("activationSuccess"))
		a.RenderArgs["activated"] = true
		a.RenderArgs["loginName"] = user.Name
		a.RenderArgs["email"] = user.Email
	}
	return a.RenderTemplate("home/activate.html")
}
