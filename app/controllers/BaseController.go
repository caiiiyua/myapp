package controllers

import (
	"myapp/app/models"
	"myapp/app/models/entity"
	"strings"

	"github.com/revel/revel"
)

// Base of all other controllers
type BaseController struct {
	*revel.Controller
	XOrmTnController
}

func (c BaseController) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}

func (c BaseController) GetUserId() string {
	if userid, ok := c.Session["UserId"]; ok {
		return userid
	}
	return ""
}

func (c BaseController) GetUsername() string {
	if username, ok := c.Session["Username"]; ok {
		return username
	}
	return ""
}

func (c BaseController) GetEmail() string {
	if email, ok := c.Session["Email"]; ok {
		return email
	}
	return ""
}

func (c BaseController) GetUserVerified() bool {
	if verified, ok := c.Session["Verified"]; ok {
		if verified == "1" {
			return true
		}
	}
	return false
}

func (c BaseController) GetSession(k string) string {
	if v, ok := c.Session[k]; ok {
		return v
	}
	return ""
}

func (c BaseController) SetSession(u entity.User) {
	if u.Id >= 0 {
		c.Session["UserId"] = string(u.Id)
		c.Session["Username"] = u.Username
		c.Session["Email"] = u.Email
		if u.Verified {
			c.Session["Verified"] = "1"
		} else {
			c.Session["Verified"] = "0"
		}
	}
}

func (c BaseController) SetLocale() string {
	locale := string(c.Request.Locale)
	lang := locale
	if strings.Contains(locale, "-") {
		pos := strings.Index(locale, "-")
		lang = locale[:pos]
	}
	if lang != "zh" {
		lang = "en"
	}
	return lang
}
