package models

import "time"

type User struct {
	Id         int64     `db:"id" json:"id"`
	Username   string    `db:"name" json:"name"`
	Name       string    `db:"Name" json:"Name"`
	Email      string    `db:"email json:email`
	Mobile     string    `db:"mobilephone" json:"mobilephone"`
	Card       string    `db:"cardno" json:"cardno"`
	Verified   bool      `db:"verified" json:"verified"`
	CreateTime time.Time `db:"createtime" json:"createtime"`
}

type UserAccount struct {
	Id            int64
	UserId        int64
	Card          string
	SaveAmount    int64
	ConsumeAmount int64
	ConsumeCount  int64

	LastConsumeTime time.Time
}

type UserAddress struct {
	Id       int64
	UserId   int64
	Address  string
	Street   string
	City     string
	PostCode string
}

type UserItem struct {
	Id     int64
	UserId int64
	Card   string
	ItemId int64
	Qty    int64
}
