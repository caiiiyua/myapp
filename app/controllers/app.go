package controllers

import "github.com/revel/revel"

import (
	"myapp/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
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
