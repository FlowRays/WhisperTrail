package service

import (
	"github.com/FlowRays/WhisperTrail/dao"
	"github.com/FlowRays/WhisperTrail/model"
)

func GetRate(id uint, db *model.Database) (float64, error) {
	rating, err := dao.GetAvgRating(id, db)
	if err != nil {
		return 0, err
	}

	return rating, err
}

func CreateRate(rate *model.Rate, db *model.Database) error {
	err := dao.CreateRate(rate, db)

	return err
}
