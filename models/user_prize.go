package models

import (
	"github.com/garbein/lottery-golang/apps"
)

type UserPrize struct {
	Id      int
	UserId  int
	PrizeId int
	Status  int
}

func (userPrize *UserPrize) Create() (int, error) {
	err := apps.App.DB.Create(&userPrize).Error
	return userPrize.Id, err
}
