package services

import (
	"errors"

	"github.com/garbein/lottery-golang/models"
	"github.com/garbein/lottery-golang/utils"
	"gorm.io/gorm"
)

func Login(mobile string) (string, error) {
	var user models.User
	err := user.GetByMobile(mobile)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		user = models.User{Mobile: mobile}
		user.Create()
	} else if err != nil {
		return "", err
	}
	token, err := utils.GenerateAccessToken(user.Id)
	if err == nil {
		return token, nil
	}
	return "", err
}
