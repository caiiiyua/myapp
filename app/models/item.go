package models

import (
	"fmt"

	"github.com/revel/revel"
)

type Item struct {
	Id       int64   `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	Stock    float64 `db:"stock" json:"stock"`
	Category string  `db:"category" json:"category"`
	Cost     float32 `db:"cost" json:"cost"`
	Price    float32 `db:"price" json:"price"`
}

func (i *Item) Validate(v *revel.Validation) {
	v.Check(i.Name,
		revel.ValidRequired(),
		revel.ValidMaxSize(25))
	v.Check(i.Stock,
		revel.ValidRequired())
	v.Check(i.Cost,
		revel.ValidRequired())
}

func (i *Item) SetName(newName string) {
	i.Name = newName
	fmt.Println("set item name to ", newName)
}

func (i *Item) GetName() string {
	fmt.Println("item name is ", i.Name)
	return i.Name
}

func (i *Item) GetId() int64 {
	fmt.Println("item id is ", i.Id)
	return i.Id
}

// /// Implement for inventory
// func (i *Item) Instance(name string, num float64) {
// 	i.Name = name
// 	i.Stock = num
// }

// func (i *Item) Init(num float64) {
// 	i.Instance("item", num)
// }

// func (i *Item) Destory() {
// 	i.Stock = 0
// }

// func (i *Item) SetStock(num float64) {
// 	i.Stock = num
// }

// func (i *Item) GetStock() float64 {
// 	return i.Stock
// }

// func (i *Item) AddStock(num float64) float64 {
// 	i.Stock = i.Stock + num
// 	return i.Stock
// }

// func (i *Item) ReduceStock(num float64) float64 {
// 	i.Stock = i.Stock - num
// 	if i.Stock < 0 {
// 		i.Stock = 0
// 	}
// 	return i.Stock
// }
