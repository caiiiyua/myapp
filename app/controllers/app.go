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
	item.Instance("cup", 10)
	itemName := item.GetName()
	itemStock := item.GetStock()
	return c.Render(item, itemName, itemStock)
}
