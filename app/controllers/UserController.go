package controllers

import "github.com/revel/revel"

type User struct {
	BaseController
}

func (c *User) Account() revel.Result {
	return c.Render()
}
