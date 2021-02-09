package models

import (
	"github.com/garbein/lottery-golang/apps"
)

type User struct {
	Id     int
	Mobile string
}

func (user *User) GetByMobile(mobile string) error {
	err := apps.App.DB.First(&user, "mobile = ?", mobile).Error
	return err
}

func (user *User) Create() error {
	err := apps.App.DB.Create(&user).Error
	return err
}
