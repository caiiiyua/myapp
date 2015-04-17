package entity

type Item struct {
	Id          int64
	Code        string `xorm:"unique not null index"`
	Name        string
	Description string `xorm:"varchar(64)"`
	Price       float64
}
