package controllers

import (
	"log"
	"myapp/app/models"

	"github.com/revel/revel"
)

type User struct {
	BaseController
}

func (c User) Account(id string) revel.Result {
	if ret := c.checkLogined(id); ret != nil {
		log.Println("id:", id, "was not logined!")
		return ret
	}

	user, ok := c.userService().CheckUserById(id)
	if !ok {
		return c.Redirect(Auth.Login)
	}

	name := models.ToSessionUser(user).DisplayName()
	cardNo := models.ToSessionUser(user).GetVipNo()

	account, _ := c.userService().GetUserItems(cardNo)
	log.Println("account:", account)
	return c.Render(account, name, cardNo)
}

func (c User) Accounts() revel.Result {

	id, ok := c.Session["id"]
	log.Println("id:", id, " ok:", ok)
	if ok && id != "" && id != "0" {
		return c.Redirect("/users/%s", id)
	}
	return c.Redirect(Auth.Login)
}
