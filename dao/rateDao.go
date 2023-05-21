package dao

import (
	"github.com/FlowRays/WhisperTrail/model"
)

func GetAvgRating(id uint, db *model.Database) (float64, error) {
	var ratings []model.Rate
	result := db.DB.Find(&ratings, "landmark_id = ?", id)
	if result.Error != nil {
		return 0, result.Error
	}

	var totalRating uint
	for _, rating := range ratings {
		totalRating += rating.Rating
	}
	var avgRating float64
	if len(ratings) == 0 {
		avgRating = 0
	} else {
		avgRating = float64(totalRating) / float64(len(ratings))
	}

	return avgRating, nil
}

func CreateRate(rate *model.Rate, db *model.Database) error {
	err := db.DB.Create(&rate).Error
	return err
}
