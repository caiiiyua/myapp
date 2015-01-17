package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"myapp/app/models"
)

type ItemCtrl struct {
	GorpController
}

func (i *ItemCtrl) parseItem() (models.Item, error) {
	item := models.Item{}
	err := json.NewDecoder(i.Request.Body).Decode(&item)
	return item, err
}

func (i *ItemCtrl) Add() revel.Result {
	if item, err := i.parseItem(); err != nil {
		return i.RenderText("Unable to parse item from JSON")
	} else {
		// validate from model
		item.Validate(i.Validation)
		if i.Validation.HasErrors() {
			return i.RenderText("Something wrong in your Item")
		} else {
			if err := i.Txn.Insert(&item); err != nil {
				return i.RenderText("Error insert Item into Datebase")
			} else {
				return i.RenderJson(item)
			}
		}
	}
}

func (i *ItemCtrl) Get(id int64) revel.Result {
	item := new(models.Item)
	err := i.Txn.SelectOne(item,
		`select * from Item where id = ?`, id)
	if err != nil {
		return i.RenderText("Item [%v] doesn't exist.", id)
	}
	return i.RenderJson(item)
}

func (i *ItemCtrl) List() revel.Result {
	lastId := parseIntOrDefault(i.Params.Get("lid"), -1)
	limit := parseUintOrDefault(i.Params.Get("limit"), uint64(25))
	items, err := i.Txn.Select(models.Item{},
		`select * from Item where id > ? limit ?`, lastId, limit)
	if err != nil {
		return i.RenderText("Error when trying to get Items from Database")
	}
	return i.RenderJson(items)
}

func (i *ItemCtrl) Update(id int64) revel.Result {
	item, err := i.parseItem()
	if err != nil {
		return i.RenderText("Unable to parse item from JSON")
	}
	// Ensure the id is set
	item.Id = id
	success, err := i.Txn.Update(&item)
	if err != nil || success == 0 {
		return i.RenderText("Unable to update Item [%v]", id)
	}
	return i.RenderText("Updated Item [%v]", id)
}

func (i *ItemCtrl) Delete(id int64) revel.Result {
	success, err := i.Txn.Delete(&models.Item{Id: id})
	if err != nil || success == 0 {
		return i.RenderText("Failed to delete Item [%v]", id)
	}
	retrn i.RenderText("Deleted Item [%v]", id)
}
