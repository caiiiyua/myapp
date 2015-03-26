package controllers

import (
	"database/sql"
	"log"
	"myapp/app"
	"myapp/app/models/entity"
	"myapp/app/utils"
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

func (x *XOrmTnController) Begin() revel.Result {
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
	app.Engine.DropTables("t_user")

	err = app.Engine.Sync2(new(entity.User), new(entity.UserRole), new(entity.UserLevel),
		new(entity.Location), new(entity.UserDetail), new(entity.UserItem))
	if err != nil {
		log.Fatal("Sync2 with error:", err)
	}

	// do init
	tryInitData()
}

func tryInitData() {
	total, err := app.Engine.Count(&entity.User{})
	if total > 0 && err == nil {
		revel.INFO.Println("total users:", total)
		return
	}
	var acts []app.Account
	acts = app.ImportAccounts()
	log.Println("accounts:", len(acts))
	users := []entity.User{}
	for _, act := range acts {
		user := entity.User{}
		user.CardId = act.CardId
		user.Name = act.Name
		user.Mobile = act.Mobile
		users = append(users, user)
	}
	for _, user := range users {
		if len(user.CardId) <= 0 || len(user.CardId) > 7 {
			continue
		}
		// log.Println(user)
		_, err := app.Engine.Insert(&user)
		if err != nil {
			log.Println(err)
		}
	}
}

func driverInfoFromConfig() (driver, spec string) {
	driver = utils.ForceGetConfig("db.driver")
	log.Println("db driver:", driver)
	spec = utils.ForceGetConfig("db.spec")
	log.Println("db sepc:", hidePassword(spec))
	return
}

func hidePassword(spec string) string {
	re := regexp.MustCompile(":.*@")
	return re.ReplaceAllString(spec, ":*******@")
}
