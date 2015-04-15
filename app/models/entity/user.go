package entity

import "time"

// UserLevel level for users to distingush VIPs
type UserLevel struct {
	Id          int64
	Sort        int
	Name        string
	Code        string `xorm:"unique"`
	ScoreStart  int    `xorm:"int default 0"`
	ScoreEnd    int    `xorm:"int default 0"`
	Description string
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated`
}

// UserRole role of users (admin, customer, stuff...)
type UserRole struct {
	Id        int64
	Sort      int
	Name      string
	Code      string    `xorm:"unique"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated`
}

type User struct {
	Id                      int64  `db:"id" json:"id"`
	Username                string `db:"username" json:"username" xorm:"index"`                              // login name
	CryptedPassword         string `json:"-" xorm:"varchar(64)"`                                             //	encrypted password
	Name                    string `db:"Name" json:"Name" xorm:"varchar(100)"`                               // real name of user
	Email                   string `db:"email" json:"email" xorm:varchar(100)`                               // email address of user
	Mobile                  string `db:"mobilephone" json:"mobile_phone" xorm:"varchar(100) not null index"` // mobile phone of user
	CardId                  string `db:"CardId" json:"card_id" xorm:"varchar(64) index"`                     // vip card no of user
	Scores                  int    `json:"scores" xorm:"int default 0"`                                      // total scores of user currently
	Level                   string `json:"level" xorm:"varchar(20)"`                                         // vip level of user
	Gender                  string `json:"gender" xorm:"varchar(100)"`                                       // gender of user
	Age                     int64  `json:"age" xorm:"index"`
	Birthday                time.Time
	Address                 string    `db:"address" json:"addr"`
	Status                  string    `db:"status" json:"status"`
	Verified                bool      `json:"verified xorm:"verified"` // user vailation
	NonExpired              bool      // user not expired
	CredentialsNonExpired   bool      // user's credentials not expired
	NonLocked               bool      // user not locked
	ActivationCode          string    `json:"activation_code" xorm:"activation_code"` // activation code for user
	ActivationCodeCreatedAt time.Time // time of activation code created
	PasswordResetCode       string    // reset code for password finding
	CreateTime              time.Time `db:"createtime" json:"createtime" xorm:"created"`
	UpdateTime              time.Time `db:"updatetime" json:"updatetime" xorm:"updated"`
	LastSignAt              time.Time // last time logged in
	DataVersion             int       `json:"-" xorm:"version '_version'"`
}

// Location of user
type Location struct {
	Id       int64
	Province string
	City     string
	Area     string
	Address  string
}

// User info
type UserDetail struct {
	Id         int64
	UserId     int64  // related user id
	Role       string // role of user
	LocationId string // location id of user
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
	UserId    int64  // related user id
	CardId    string // related user's card no
	ItemId    int64  // related item id
	Qty       int64
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
