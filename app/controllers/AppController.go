package controllers

import (
	"fmt"

	"github.com/revel/revel"
)

import (
	"myapp/app"
	"myapp/app/models"
)

type App struct {
	BaseController
}

func (c *App) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}

func (c App) Index() revel.Result {
	locale := c.SetLocale()
	fmt.Println(locale)
	c.RenderArgs["users_count"] = c.userService().Total()
	// c.RenderArgs["users"] = c.userService().ListUsers()
	c.RenderArgs["version"] = app.Version
	//	return c.RenderTemplate("home/index_" + locale + ".html")
	// return c.RenderTemplate("home/index.html")
	return c.Render()
}

func (c App) Inventory() revel.Result {
	var item models.Item
	itemName := "cup"
	itemStock := "10"
	return c.Render(item, itemName, itemStock)
}

func (c App) Item() revel.Result {
	return c.Render()
}
