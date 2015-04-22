package controllers

import (
	"fmt"
	"log"
	"myapp/app"
	"myapp/app/controllers/api"
	"myapp/app/models"
	"myapp/app/utils"

	"github.com/dchest/captcha"
	"github.com/revel/revel"
)

type Auth struct {
	BaseController
}

// func (a Auth) checkReg() revel.Result {
// 	return a.check(Auth.Register)
// }

// func (a Auth) checkAuth() revel.Result {
// 	return a.check(Auth.Login)
// }

// func (a Auth) check(f interface{}) revel.Result {
// 	if a.Validation.HasErrors() {
// 		a.Validation.Keep()
// 		a.FlashParams()
// 		return a.Redirect(f)
// 	}
// 	return nil
// }

func (a Auth) Register2() revel.Result {
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

func (a Auth) Register(userType string) revel.Result {
	revel.INFO.Println("userType: %v", userType)

	Captcha := struct {
		CaptchaId string
	}{
		captcha.New(),
	}

	return a.Render(Captcha)
}

func (a Auth) DoRegister(email, pwd, pwd2, validateCode, captchaId string) revel.Result {
	revel.INFO.Println("DoRegister:", email, pwd, validateCode, captchaId)
	a.Validation.Required(captcha.VerifyString(captchaId,
		validateCode)).Message(a.Message("inputCaptcha")).Key("validateCode")
	a.Validation.Required(email).Message("inputEmail")
	a.Validation.Email(email).Message("wrongEmail")
	a.Validation.MaxSize(email, 100).Message("wrongEmail")
	a.Validation.Required(pwd).Message("inputPassword")
	a.Validation.MinSize(pwd, 6).Message("notGoodPassword")
	a.Validation.MaxSize(pwd, 50).Message("notGoodPassword")
	a.Validation.Required(pwd2).Message("inputPassword2")
	a.Validation.Required(pwd == pwd2).Message("confirmPassword").Key("confirmPassword")

	if ret := a.checkReg(); ret != nil {
		return ret
	}

	a.Validation.Required(!a.userService().ExistsUserByEmail(email)).Message("userHasBeenRegistered").Key("email")
	if ret := a.checkReg(); ret != nil {
		return ret
	}

	user, err := a.userService().RegistUserByEmail(email, pwd)
	if err != nil {
		a.Flash.Error(err.Error())
		a.Redirect(Auth.Register)
	}

	// data := struct {
	// 	ActivationCode string
	// 	Email          string
	// }{
	// 	user.ActivationCode,
	// 	user.Email,
	// }
	data := make(map[string]interface{})
	data["ActivationCode"] = user.ActivationCode
	data["Email"] = user.Email

	// TODO: Do with timeout
	SendMail(a.Message("activationMail"), utils.RenderTemplateToString("home/user_activate_tmpl.html",
		data),
		email)
	a.RenderArgs["emailProvider"] = EmailProvider(email)
	a.RenderArgs["email"] = email

	go a.userService().DoUserLogin(&user)

	a.Session["user"] = models.ToSessionUser(user).DisplayName()
	a.Session["id"] = models.ToSessionUser(user).GetId()
	return a.Render()
}

func (a Auth) DoRegister2(email, pwd, captcha, captchaId string) revel.Result {
	ok := app.NewOk()

	a.Validation.Required(!a.userService().ExistsUserByEmail(email)).Message(
		"userHasBeenRegistered", email)
	if ret := a.checkAuth(); ret != nil {
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
	if email == "879939101@qq.com" {
		ok.Msg = "admin"
	}
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
	return a.Render()
}

func (a Auth) Login2() revel.Result {
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
	revel.INFO.Println("email:", email, "validationCode:", validationCode, "captchaId:", captchaId)
	a.Validation.Required(captcha.VerifyString(captchaId,
		validationCode)).Message(a.Message("inputCaptcha")).Key("validateCode")
	a.Validation.Required(email).Message(a.Message("inputEmail"))
	a.Validation.Required(pwd).Message(a.Message("inputPassword"))
	a.Validation.MinSize(pwd, 6).Message(a.Message("passwordTips"))

	if ret := a.checkAuth(); ret != nil {
		return ret
	}

	user, ok := a.userService().CheckUser(email, pwd)
	a.Validation.Required(ok).Message(a.Message("wrongUsernameOrPassword")).Key("email")
	if ret := a.checkAuth(); ret != nil {
		log.Println("Need login first!")
		return ret
	}

	go a.userService().DoUserLogin(&user)

	a.Session["user"] = models.ToSessionUser(user).DisplayName()
	a.Session["id"] = models.ToSessionUser(user).GetId()
	log.Println("set register session:", a.Session, "with user:", user)
	return a.Redirect(App.Index)
}

func (a Auth) DoLogin2(email, pwd, validationCode, captchaId string) revel.Result {
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

	if ret := a.checkAuth(); ret != nil {
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

func (a Auth) Activate(activationCode, email string) revel.Result {
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

func (a Auth) Join() revel.Result {
	a.RenderArgs["needCaptcha"] = "true"
	a.RenderArgs["openRegister"] = "true"
	Captcha := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	a.RenderArgs["Captcha"] = Captcha
	log.Println("captchaId:", Captcha.CaptchaId)
	return a.Render()
}

func (a *Auth) DoJoin(email, pwd, validationCode, captchaId, vipno, phoneno string) revel.Result {
	revel.INFO.Println("email:", email, "validationCode:", validationCode, "captchaId:", captchaId,
		"vipno:", vipno, "phoneno:", phoneno)
	a.Validation.Required(captcha.VerifyString(captchaId,
		validationCode)).Message(a.Message("inputCaptcha")).Key("validateCode")
	a.Validation.Required(email).Message(a.Message("inputEmail"))
	a.Validation.Required(pwd).Message(a.Message("inputPassword"))
	a.Validation.MinSize(pwd, 6).Message(a.Message("passwordTips"))

	if ret := a.checkAuth(); ret != nil {
		return ret
	}

	user, ok := a.userService().CheckUser(email, pwd)
	a.Validation.Required(ok).Message(a.Message("wrongUsernameOrPassword")).Key("email")
	if ret := a.checkAuth(); ret != nil {
		log.Println("Need login first!")
		return ret
	}
	go a.userService().DoUserLogin(&user)

	a.Session["user"] = models.ToSessionUser(user).DisplayName()
	a.Session["id"] = models.ToSessionUser(user).GetId()
	log.Println("set register session:", a.Session, "with user:", user)

	err := a.userService().JoinAccount(&user, vipno, phoneno)
	if err != nil {
		a.Validation.Error("Phone or Vip No. was not correct!").Key("join")
		return a.Redirect(Auth.Join)
	}
	id := models.ToSessionUser(user).GetId()
	return a.Redirect("/users/%s", id)
}
