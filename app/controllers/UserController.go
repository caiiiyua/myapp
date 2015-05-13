package controllers

import (
	"encoding/json"
	"log"
	"myapp/app/models"

	"github.com/revel/revel"
)

type User struct {
	BaseController
}

type UserItems struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Qty  int64  `json:"qty"`
}

type Act struct {
	Card     string      `json:"card"`
	Action   int         `json:"action"` //"0" for show and "1" for add and "2" for reduce
	ActItems []UserItems `json:"items"`
}

func (c User) parseAct() (Act, error) {
	var act Act
	err := json.NewDecoder(c.Request.Body).Decode(&act)
	return act, err
}

func (c User) Account(id string) revel.Result {
	if ret := c.checkLogined(id); ret != nil {
		log.Println("id:", id, "was not logined!")
		return ret
	}

	user, ok := c.userService().CheckUserById(id)
	if !ok {
		return c.Redirect(Auth.Login)
	}

	name := models.ToSessionUser(user).DisplayName()
	cardNo := models.ToSessionUser(user).GetVipNo()
	userId := id

	account, _ := c.userService().GetUserItems(cardNo)
	log.Println("account:", account)
	return c.Render(account, name, cardNo, userId)
}

func (c User) Accounts() revel.Result {

	id, ok := c.Session["id"]
	log.Println("id:", id, " ok:", ok)
	if ok && id != "" && id != "0" {
		return c.Redirect("/users/%s", id)
	}
	return c.Redirect(Auth.Login)
}

func (c User) Join(id string) revel.Result {
	return c.Render(id)
}

func (c User) DoJoin(id, vipno, phoneno string) revel.Result {
	log.Println("dojoin with id:", id, " vipNo:", vipno, " phoneNo:", phoneno)
	c.Validation.Required(vipno).Message("请输入会员卡号").Key("vipno")
	c.Validation.Required(phoneno).Message("请输入手机号码").Key("phoneno")
	if ret := c.checkAuth(); ret != nil {
		return ret
	}

	user, ok := c.userService().CheckUserById(id)
	c.Validation.Required(ok).Message("未注册会员或会员信息错误，请重新登录。").Key("id")

	log.Println("set register session:", c.Session, "with user:", user)

	err := c.userService().JoinAccount(&user, vipno, phoneno)
	if err != nil {
		c.Validation.Error("Phone or Vip No. was not correct!").Key("join")
	}
	return c.Redirect("/users/%s", id)
}

func (c User) UpdateItems() revel.Result {
	act, _ := c.parseAct()
	log.Println("updateitems with account:", act)
	for _, item := range act.ActItems {
		if act.Action == 1 {
			c.userService().AddItems(act.Card, item.Id, item.Qty)
		} else if act.Action == 2 {
			c.userService().ReduceItems(act.Card, item.Id, item.Qty)
		}
	}

	account, _ := c.userService().GetUserItems(act.Card)
	log.Println("latest account:", account)
	retAct := Act{}
	retAct.Card = act.Card
	retAct.Action = 0
	for _, item := range account {
		i := UserItems{}
		i.Id = item.ItemId
		i.Name = item.Name
		i.Qty = item.Qty
		retAct.ActItems = append(retAct.ActItems, i)
	}

	return c.Render(retAct)
}

func (c User) Items() revel.Result {
	return c.Render()
}
