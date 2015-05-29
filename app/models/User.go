package models

import (
	"errors"
	"fmt"
	"log"
	"myapp/app/models/entity"
	"strconv"
	"time"

	"myapp/app/utils"

	"github.com/go-xorm/xorm"
)

type ItemAccount struct {
	entity.UserItem `xorm:"extends"`
	Name            string
}

type SessionUser struct {
	Email    string
	VipCode  string
	Name     string
	Nickname string
	Id       int64
}

func (s SessionUser) DisplayName() string {
	if len(s.Name) == 0 {
		if len(s.Email) == 0 {
			return s.Nickname
		} else {
			return s.Email
		}
	}
	return s.Name
}

func (s SessionUser) GetId() string {
	return strconv.FormatInt(s.Id, 10)
}

func (s SessionUser) GetVipNo() string {
	return s.VipCode
}

func ToSessionUser(user entity.User) SessionUser {
	return SessionUser{
		Email:    user.Email,
		VipCode:  user.CardId,
		Name:     user.Name,
		Id:       user.Id,
		Nickname: user.Nickname,
	}
}

type UserService interface {
	Total() int64
	ListUsers() []entity.User
	// RegisterUser(username, password string) (entity.User, error)
	RegistUserByEmail(email, password string) (entity.User, error)
	// RegisterUserByPhone(mobile, password string) (entity.User, error)
	ExistsUserByEmail(email string) bool
	// ExistsUserByPhone(mobile string) bool
	// ExistsUser(username string) bool
	Activate(email, code string) (entity.User, error)
	CheckUser(email, password string) (entity.User, bool)
	CheckUserByEmail(email string) (entity.User, bool)
	CheckUserById(id string) (entity.User, bool)
	DoUserLogin(user *entity.User) error
	CheckUserLogin(userId int64, userName string) (entity.User, bool)
	JoinAccount(user *entity.User, vipNo, phoneNo string) error
	GetUserItems(vipNo string) ([]ItemAccount, bool)
	GetUserItem(vipNo string, itemId int64) (ItemAccount, bool)
	AddItems(vipNo string, itemId int64, qty int64) (ItemAccount, bool)
	ReduceItems(vipNo string, itemId int64, qty int64) (ItemAccount, bool)
	AddVip(id, name, addr, phone, wId string) bool
	UpdateVip(id, name, addr, phone, wId string) bool

	// User by Wechat
	AddWeChatMember(openId, unionId, nickname, sex, city, province, headimgUrl string) (entity.User, error)
	UpdateWeChatMember(id int64, openId, unionId, nickname, sex, city, province, headimgUrl string) (user entity.User, err error)
	CheckWeChatMemberByOpenId(openId string) (user entity.User, ok bool)
	CheckWeChatMemberByUnionId(unionId string) (user entity.User, ok bool)
}

func DefaultUserService(session *xorm.Session) UserService {
	return &defaultUserService{session}
}

type defaultUserService struct {
	session *xorm.Session
}

func (this *defaultUserService) Total() int64 {
	s := this.session
	if s != nil {
		fmt.Println("session is valid")
	}
	fmt.Println("session is", s)
	ret, err := s.Count(&entity.User{})
	if err != nil {
		log.Println("get count failed:", err)
	}
	return ret
}

func (this *defaultUserService) ListUsers() (users []entity.User) {
	this.session.Find(&users)
	return
}

func (this *defaultUserService) AddWeChatMember(openId, unionId, nickname, sex,
	city, province, headimgUrl string) (user entity.User, err error) {
	user.OpenId = openId
	user.UnionId = unionId
	user.Gender = sex
	user.Nickname = nickname
	user.City = city
	user.Province = province
	user.HeadimgUrl = headimgUrl

	_, err = this.session.Insert(&user)
	return
}

func (this *defaultUserService) UpdateWeChatMember(id int64, openId, unionId, nickname, sex,
	city, province, headimgUrl string) (user entity.User, err error) {
	user.Gender = sex
	user.Nickname = nickname
	user.City = city
	user.Province = province
	user.HeadimgUrl = headimgUrl

	_, err = this.session.Id(id).Update(&user)
	return
}

func (this *defaultUserService) CheckWeChatMemberByOpenId(openId string) (user entity.User, ok bool) {
	ok, err := this.session.Where("openid=?", openId).Get(&user)
	return user, ok && err == nil
}

func (this *defaultUserService) CheckWeChatMemberByUnionId(unionId string) (user entity.User, ok bool) {
	ok, err := this.session.Where("unionid=?", unionId).Get(&user)
	return user, ok && err == nil
}

// ========================================================================================

func (this *defaultUserService) RegistUserByEmail(email, password string) (user entity.User, err error) {
	user.Email = email
	user.CryptedPassword = utils.Sha1(password)
	user.ActivationCode = utils.Uuid()
	user.ActivationCodeCreatedAt = time.Now()

	_, err = this.session.Insert(&user)
	return
}

func (this *defaultUserService) ExistsUserByEmail(email string) bool {
	total, _ := this.session.Where("email=?", email).Count(&entity.User{})
	return total > 0
}

func (this *defaultUserService) Activate(email, code string) (user entity.User, err error) {
	if len(code) == 0 {
		err = errors.New("Activation code invalid")
		return
	}
	var users []entity.User
	err = this.session.Where("email=? and activation_code=?", email, code).Find(&users)
	if err != nil {
		return
	}
	if len(users) > 0 {
		user = users[0]
		user.Verified = true
		user.ActivationCode = ""
		this.session.Id(user.Id).Cols("verified", "activation_code").Update(&user)
		return
	} else {
		err = errors.New("no user found	!")
		return
	}
}

func (this *defaultUserService) CheckUser(email, password string) (user entity.User, ok bool) {
	log.Println("email:", email, "pwd:", utils.Sha1(password))
	ok, err := this.session.Where("email=? and crypted_password=?", email, utils.Sha1(password)).Get(&user)
	log.Println("user:", user, "ok:", ok, "err:", err)
	return user, ok && err == nil
}

func (this *defaultUserService) CheckUserByEmail(email string) (user entity.User, ok bool) {
	ok, err := this.session.Where("email=?", email).Get(&user)
	return user, ok && err == nil
}

func (this *defaultUserService) CheckUserById(id string) (user entity.User, ok bool) {
	ok, err := this.session.Where("id=?", id).Get(&user)
	log.Println("user:", user, "ok:", ok, "err:", err)
	return user, ok && err == nil
}

func (this *defaultUserService) DoUserLogin(user *entity.User) error {
	user.LastSignAt = time.Now()
	_, err := this.session.Id(user.Id).Update(user)
	return err
}

func (this *defaultUserService) CheckUserLogin(userId int64, userName string) (user entity.User, ok bool) {
	ok, err := this.session.Where("id=? and email=?", userId, userName).Get(&user)
	return user, ok && err == nil
}

func (this *defaultUserService) JoinAccount(user *entity.User, vipNo, phoneNo string) error {
	var vipInfo entity.User
	ok, err := this.session.Where("mobile=? and card_id=?", phoneNo, vipNo).Get(&vipInfo)
	if ok && err == nil {
		user.CreateTime = vipInfo.CreateTime
		user.CardId = vipInfo.CardId
		user.Mobile = vipInfo.Mobile
		user.Address = vipInfo.Address
		user.Name = vipInfo.Name
		log.Println("New user:", user)
		this.session.Id(vipInfo.Id).Delete(&vipInfo)
		_, updateErr := this.session.Id(user.Id).Update(user)
		if updateErr != nil {
			log.Println("Update failed:", updateErr)
			return updateErr
		}
		itemSql := "update t_user_item set user_id=? where card_id=?"
		this.session.Exec(itemSql, user.Id, user.CardId)
	}
	return err
}

func (this *defaultUserService) GetUserItems(vipNo string) (items []ItemAccount, ok bool) {
	this.session.Sql("select t_user_item.*, t_item.Name from t_user_item, t_item where t_user_item.card_id=? and t_user_item.item_id = t_item.code", vipNo).Find(&items)
	return items, true
}

func (this *defaultUserService) GetUserItem(vipNo string, itemId int64) (item ItemAccount, ok bool) {
	this.session.Sql("select t_user_item.*, t_item.Name from t_user_item, t_item where t_user_item.card_id=? and t_user_item.item_id=?", vipNo, itemId).Find(&item)
	return item, true
}

func (this *defaultUserService) AddItems(vipNo string, itemId int64, qty int64) (userItem ItemAccount, ok bool) {
	var item entity.UserItem
	ok, err := this.session.Where("card_id=? and item_id=?", vipNo, itemId).Get(&item)
	if !ok || err != nil {
		return userItem, false
	}
	item.Qty += qty
	_, updateErr := this.session.Cols("qty").Where("card_id=? and item_id=?", vipNo, itemId).Update(item)
	if updateErr != nil {
		log.Println("Update failed:", updateErr)
		return userItem, false
	}
	return this.GetUserItem(vipNo, itemId)
}

func (this *defaultUserService) ReduceItems(vipNo string, itemId int64, qty int64) (userItem ItemAccount, ok bool) {
	var item entity.UserItem
	ok, err := this.session.Where("card_id=? and item_id=?", vipNo, itemId).Get(&item)
	if !ok || err != nil {
		return userItem, false
	}
	item.Qty -= qty
	if item.Qty < 0 {
		return userItem, false
	}
	_, updateErr := this.session.Cols("qty").Where("card_id=? and item_id=?", vipNo, itemId).Update(item)
	if updateErr != nil {
		log.Println("Update failed:", updateErr)
		return userItem, false
	}
	return this.GetUserItem(vipNo, itemId)
}

func (this *defaultUserService) AddVip(id, name, addr, phone, wId string) bool {
	var user entity.User
	var idx int64
	this.session.Sql("select max(card_id)+1 from t_user").Find(&idx)
	id = fmt.Sprintf("%07d", idx)
	log.Println("new vip card_id is:", id)
	user.Name = name
	user.CardId = id
	user.Address = addr
	user.Mobile = phone
	// user.wId = wId
	_, err := this.session.Insert(&user)
	if err != nil {
		log.Println("Add Vip failed:", err)
		return false
	}
	return true
}

func (this *defaultUserService) UpdateVip(id, name, addr, phone, wId string) bool {
	var user entity.User
	log.Println("updated vip card_id is:", id)
	ok, err := this.session.Where("card_id=?", id).Get(&user)
	if !ok || err != nil {
		return false
	}
	user.Name = name
	user.Address = addr
	user.Mobile = phone
	// user.wId = wId
	_, updateErr := this.session.Update(&user)
	if updateErr != nil {
		log.Println("Add Vip failed:", updateErr)
		return false
	}
	return true
}
