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

type SessionUser struct {
	Email   string
	VipCode string
	Name    string
	Id      int64
}

func (s SessionUser) DisplayName() string {
	if len(s.Name) == 0 {
		return s.Email
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
		Email:   user.Email,
		VipCode: user.CardId,
		Name:    user.Name,
		Id:      user.Id,
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
	DoUserLogin(user *entity.User) error
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
	return user, ok && err != nil
}

func (this *defaultUserService) DoUserLogin(user *entity.User) error {
	user.LastSignAt = time.Now()
	_, err := this.session.Id(user.Id).Update(user)
	return err
}
