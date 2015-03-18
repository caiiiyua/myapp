package controllers

import (
	"github.com/revel/revel"
)

type Auth struct {
	BaseController
}

func (a Auth) Register() revel.Result {

	return a.RenderTemplate("home/register.html")
}

func (a Auth) DoRegister() revel.Result {
	return a.Redirect("/items")
}

func (a Auth) Logout() revel.Result {
	return a.Render()
}

func (a Auth) Login() revel.Result {
	return a.RenderTemplate("home/login.html")
}

func (a Auth) DoLogin() revel.Result {
	return a.Render()
}
