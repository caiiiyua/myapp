package controllers

import (
	"database/sql"
	"log"
	"myapp/app"
	"myapp/app/models/entity"
	"myapp/app/utils"
	"regexp"
	"strconv"

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
	app.Engine.DropTables("t_user_item")
	app.Engine.DropTables("t_item")
	app.Engine.DropTables("item")

	err = app.Engine.Sync2(new(entity.User), new(entity.UserRole), new(entity.UserLevel),
		new(entity.Location), new(entity.UserDetail), new(entity.UserItem), new(entity.Item))
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
	items := []entity.UserItem{}
	for _, act := range acts {
		user := entity.User{}
		user.CardId = act.CardId
		user.Name = act.Name
		user.Mobile = act.Mobile
		users = append(users, user)
		for _, item := range act.Items {
			i := entity.UserItem{}
			i.CardId = item.ActId
			i.ItemId, _ = strconv.ParseInt(item.ItemId, 10, 64)
			i.Qty = item.Quantity
			items = append(items, i)
		}
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

	app.Engine.Insert(&items)

	initItems()
}

func initItems() {
	var items []entity.Item
	milk250 := entity.Item{}
	milk250.Code = "1001"
	milk250.Price = 7.0
	milk250.Name = "250ml 鲜牛奶"
	milk250.Description = "250ml 巴氏鲜奶"
	items = append(items, milk250)

	milk500 := entity.Item{}
	milk500.Code = "1002"
	milk500.Price = 12.0
	milk500.Name = "500ml 鲜牛奶"
	milk500.Description = "500ml 巴氏鲜奶"
	items = append(items, milk500)

	yoghourt := entity.Item{}
	yoghourt.Code = "2001"
	yoghourt.Price = 8.0
	yoghourt.Name = "原味酸奶"
	yoghourt.Description = "原味酸奶"
	items = append(items, yoghourt)

	pudding := entity.Item{}
	pudding.Code = "1008"
	pudding.Price = 10.0
	pudding.Name = "布丁"
	pudding.Description = "焦糖布丁"
	items = append(items, pudding)

	shuangpi := entity.Item{}
	shuangpi.Code = "3001"
	shuangpi.Price = 8.0
	shuangpi.Name = "双皮奶"
	shuangpi.Description = "双皮奶"
	items = append(items, shuangpi)

	app.Engine.Insert(&items)
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
