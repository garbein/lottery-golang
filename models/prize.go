package models

import (
	"errors"

	"github.com/garbein/lottery-golang/apps"
	"gorm.io/gorm"
)

type Prize struct {
	Id             int
	Name           string
	TotalStock     int
	UsedStock      int
	LotteryPercent int
	Rule           string
	Status         int
}

type FormatPrize struct {
	Prize
	Low  int
	High int
}

func (prize *Prize) GetPrizeList() ([]*Prize, error) {
	var prizeList []*Prize
	err := apps.App.DB.Find(&prizeList, "status = ?", 1).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return prizeList, nil
}

func (prize *Prize) UpdateUsedStock() (int64, error) {
	result := apps.App.DB.Model(&prize).Update("used_stock", gorm.Expr("used_stock + ?", 1))
	return result.RowsAffected, result.Error
}
