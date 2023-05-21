package dao

import (
	"github.com/FlowRays/WhisperTrail/model"
)

func CreateLandmark(landmark *model.Landmark, db *model.Database) error {
	err := db.DB.Create(&landmark).Error
	return err
}

func GetLandmarksByUserStatus(isLoggedIn bool, userID uint, landmarks *[]model.Landmark, db *model.Database) error {
	var err error
	if isLoggedIn {
		err = db.DB.Where("is_logged_in = ? AND user_id = ?", isLoggedIn, userID).Find(&landmarks).Error
	} else {
		err = db.DB.Where("is_logged_in = ?", isLoggedIn).Find(&landmarks).Error
	}
	return err
}

func GetLandmarkByID(isLoggedIn bool, userID uint, landmark *model.Landmark, db *model.Database) error {
	var err error
	if isLoggedIn {
		err = db.DB.Where("id = ? AND is_logged_in = ? AND user_id = ?", landmark.ID, isLoggedIn, userID).First(&landmark).Error
	} else {
		err = db.DB.Where("id = ? AND is_logged_in = ?", landmark.ID, isLoggedIn).First(&landmark).Error
	}
	return err
}
