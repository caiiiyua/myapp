package models

import (
	"errors"
	"fmt"
	"log"
	"myapp/app/models/entity"
	"time"

	"myapp/app/utils"

	"github.com/go-xorm/xorm"
)

type UserService interface {
	Total() int64
	ListUsers() []entity.User
	// RegisterUser(username, password string) (entity.User, error)
	RegisterUserByEmail(email, password string) (entity.User, error)
	// RegisterUserByPhone(mobile, password string) (entity.User, error)
	ExistsUserByEmail(email string) bool
	// ExistsUserByPhone(mobile string) bool
	// ExistsUser(username string) bool
	Activate(email, code string) (entity.User, error)
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

func (this *defaultUserService) RegisterUserByEmail(email, password string) (user entity.User, err error) {
	user.Email = email
	user.CryptedPassword = password
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
