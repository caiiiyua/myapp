package models

import (
	"log"
	"myapp/app/models/entity"

	"github.com/go-xorm/xorm"
)

type UserService interface {
	Total() int64
	ListUsers() []entity.User
}

func DefaultUserService(session *xorm.Session) UserService {
	return &defaultUserService{session}
}

type defaultUserService struct {
	session *xorm.Session
}

func (this *defaultUserService) Total() int64 {
	ret, err := this.session.Count(&entity.User{})
	if err != nil {
		log.Println("get count failed:", err)
	}
	return ret
}

func (this *defaultUserService) ListUsers() (users []entity.User) {
	this.session.Find(&users)
	return
}
