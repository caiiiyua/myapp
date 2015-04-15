package controllers

import "github.com/revel/revel"

type User struct {
	BaseController
}

func (c *User) Account(id string) revel.Result {
	if ret := c.checkAuth(); ret != nil {
		return ret
	}
	return c.Render()
}
