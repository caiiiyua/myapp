package entity

import "time"

type UserCategory struct {
	Id       int64
	Category string
}

type User struct {
	Id              int64  `db:"id" json:"id"`
	Username        string `db:"name" json:"name"`
	CryptedPassword string `json:"-" xorm:"varchar(64)"`
	Name            string `db:"Name" json:"Name" xorm:"varchar(100)"`
	Email           string `db:"email" json:"email" xorm:varchar(100)`
	Mobile          string `db:"mobilephone" json:"mobile_phone" xorm:"varchar(100) not null index"`
	CardId          string `db:"CardId" json:"card_id" xorm:"varchar(64) unique index"`
	Address         string `db:"address" json:"addr"`
	Status          string `db:"status" json:"status"`
	Verified        bool
	Age             int64 `xorm:"index"`
	Birthday        time.Time
	CreateTime      time.Time `db:"createtime" json:"createtime" xorm:"created"`
	UpdateTime      time.Time `db:"updatetime" json:"updatetime" xorm:"updated"`
}

type UserAccount struct {
	Id            int64
	UserId        int64
	CardId        string
	SaveAmount    float64
	ConsumeAmount float64
	ConsumeCount  int64
	LeftAmount    float64
	CreatedAt     time.Time `xorm:"created"`
	UpdatedAt     time.Time `xorm:"updated"`
}

// type UserAddress struct {
// 	Id       int64
// 	UserId   int64
// 	Address  string
// 	Street   string
// 	City     string
// 	PostCode string
// }

type UserItem struct {
	Id        int64
	UserId    int64
	CardId    string
	ItemId    int64
	Qty       int64
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
