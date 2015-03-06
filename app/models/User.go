package models

import "time"

type User struct {
	Id          int64     `db:"id" json:"id"`
	Username    string    `db:"name" json:"name"`
	UsernameRaw string    `db:"nameraw" json:"nameraw"`
	Email       string    `db:"email json:email`
	Verified    bool      `db:"verified" json:"verified"`
	CreateTime  time.Time `db:"createtime" json:"createtime"`
}
