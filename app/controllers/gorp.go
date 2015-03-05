package controllers

import (
	"database/sql"
	"fmt"
	"myapp/app/models"
	"strings"

	r "github.com/revel/revel"
	"gopkg.in/gorp.v1"
)

func getParamString(param string, defaultValue string) string {
	p, found := r.Config.String(param)
	if !found {
		if defaultValue == "" {
			r.ERROR.Fatal("Cound not find parameter: " + param)
		} else {
			return defaultValue
		}
	}
	return p
}

func getConnectionString() string {
	host := getParamString("db.host", "localhost")
	port := getParamString("db.port", "3306")
	user := getParamString("db.user", "")
	pass := getParamString("db.password", "")
	dbname := getParamString("db.name", "test")
	protocol := getParamString("db.protocol", "tcp")
	dbargs := getParamString("dbargs", " ")

	if strings.Trim(dbargs, " ") != "" {
		dbargs = "?" + dbargs
	} else {
		dbargs = ""
	}
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, protocol, host, port, dbname, dbargs)
}

func defineItemTable(dbm *gorp.DbMap) {
	t := dbm.AddTable(models.Item{}).SetKeys(true, "id")
	// t.ColMap("name").SetMaxSize(25)
	fmt.Println("table name is ", t.TableName)
}

var InitDB func() = func() {
	connectionString := getConnectionString()
	if db, err := sql.Open("mysql", connectionString); err != nil {
		r.ERROR.Fatal(err)
	} else {
		Dbm = &gorp.DbMap{
			Db:      db,
			Dialect: gorp.MySQLDialect{"InnoDB", "utf8"}}
		defineItemTable(Dbm)

		if err := Dbm.CreateTablesIfNotExists(); err != nil {
			r.ERROR.Fatal(err)
		}
	}

}

var (
	Dbm *gorp.DbMap
)

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
