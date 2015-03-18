package controllers

import (
	"database/sql"
	"log"
	"myapp/app"
	"myapp/app/models/entity"
	"regexp"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

var Db *sql.DB

type XOrmController struct {
	*revel.Controller
	Engine *xorm.Engine
}

type XOrmTnController struct {
	*revel.Controller
	XOrmSession *xorm.Session
}

func (x *XOrmController) Begin() revel.Result {
	log.Println("Begin .....")
	if app.Engine == nil {
		log.Println("App.Engine can not be nil")
	}
	// assertion instead ?

	x.Engine = app.Engine
	return nil
}

func (x *XOrmTnController) Beign() revel.Result {
	log.Println("Tn Begin .....")
	if app.Engine == nil {
		log.Println("App.Engine can not be nil")
	}
	x.XOrmSession = app.Engine.NewSession()
	x.XOrmSession.Begin()
	return nil
}

func (x *XOrmTnController) Commit() revel.Result {
	log.Println("Tn Commit .....")
	if x.XOrmSession == nil {
		log.Println("XOrmSession can not be nil")
	}

	x.XOrmSession.Commit()
	x.XOrmSession.Close()

	return nil
}

func (x *XOrmTnController) Rollback() revel.Result {
	log.Println("Tn Rollback .....")
	if x.XOrmSession == nil {
		log.Println("XOrmSession can not be nil")
	}

	x.XOrmSession.Rollback()
	x.XOrmSession.Close()

	return nil
}

var InitDB func() = func() {
	driver, spec := driverInfoFromConfig()
	var err error
	Db, err = sql.Open(driver, spec)
	if err != nil {
		log.Fatal("SQL Open failed:", err)
	}

	app.Engine, err = xorm.NewEngine(driver, spec)
	if err != nil {
		log.Fatal("NewEngine failed:", err)
	}
	app.Engine.ShowSQL = true // revel.Config.BoolDefault("db.show_sql", true)
	app.Engine.ShowDebug = true

	app.Engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, "t_"))

	err = app.Engine.Sync2(new(entity.User))
	if err != nil {
		log.Fatal("Sync2 with error:", err)

	}

	// do init
	tryInitData()
}

func tryInitData() {

}

func driverInfoFromConfig() (driver, spec string) {
	var exist bool
	driver, exist = revel.Config.String("db.driver")
	if !exist {
		log.Println("driver does not exist")
	}
	log.Println("db driver:", driver)

	spec, exist = revel.Config.String("db.spec")
	if !exist {
		log.Println("spec does not exist")
	}
	log.Println("db sepc:", hidePassword(spec))
	return
}

func hidePassword(spec string) string {
	re := regexp.MustCompile(":.*@")
	return re.ReplaceAllString(spec, ":*******@")
}
