package controllers

import "github.com/revel/revel"
import _ "github.com/go-sql-driver/mysql"

func init() {
	revel.OnAppStart(InitDB)

	revel.InterceptMethod((*XOrmController).Begin, revel.BEFORE)
	revel.InterceptMethod((*XOrmTnController).Begin, revel.BEFORE)
	revel.InterceptMethod((*XOrmTnController).Commit, revel.AFTER)
	revel.InterceptMethod((*XOrmTnController).Rollback, revel.PANIC)
}
