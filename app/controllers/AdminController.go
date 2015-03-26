package controllers

import "github.com/revel/revel"

type Admin struct {
	BaseController
}

func (c *Admin) Users() revel.Result {
	c.RenderArgs["users_count"] = c.userService().Total()
	c.RenderArgs["users"] = c.userService().ListUsers()

	return c.Render()
}
