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

type VIP struct {
	id    string `json:"id"`
	name  string `json:"name"`
	phone string `json:"phone_num"`
	addr  string `json:"address"`
	wId   string `json:"wechat_id"`
}

func (c User) parseVIP() (VIP, error) {
	var v VIP
	err := json.NewDecoder(c.Request.Body).Decode(&v)
	return v, err
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

func (c User) Items(id string) revel.Result {
	act, _ := c.userService().CheckUserById(id)
	account, _ := c.userService().GetUserItems(act.CardId)
	log.Println("get account:", account)
	retAct := Act{}
	retAct.Card = act.CardId
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

func (c User) AddVip() revel.Result {
	vip, err := c.parseVIP()
	if err != nil {
		return c.RenderJson("")
	}
	ok := c.userService().AddVip(vip.id, vip.name, vip.addr, vip.phone, vip.wId)
	if ok {
		return c.Render(vip)
	} else {
		return c.RenderJson("")
	}
}

func (c User) UpdateVip(id string) revel.Result {
	vip, err := c.parseVIP()
	if err != nil {
		return c.RenderJson("")
	}
	ok := c.userService().UpdateVip(vip.id, vip.name, vip.addr, vip.phone, vip.wId)
	if ok {
		return c.Render(vip)
	} else {
		return c.RenderJson("")
	}
}

func (c User) FeedBack() revel.Result {
	var code, state, token, openId, nickName, sex, city string
	code = c.Params.Get("code")
	state = c.Params.Get("state")
	log.Println(code, state)
	id, ok := c.Session["id"]
	log.Println("id:", id, " ok:", ok)
	if ok && id != "" && id != "0" {
		return c.Redirect("/users/%s", id)
	} else {
		// code, state = c.WeChatGetCode("feedback")
		token, openId = c.WeChatGetAccessToken(code, state)
	}
	var result = make(map[string]string)
	result["code"] = code
	result["state"] = state
	result["token"] = token
	result["openId"] = openId

	nickName, sex, city = c.WeChatGetUserInfo(token, openId)
	result["nickName"] = nickName
	result["sex"] = sex
	result["city"] = city
	return c.RenderJson(result)
}

func (c User) Bind() revel.Result {
	return c.RenderJson("")
}
