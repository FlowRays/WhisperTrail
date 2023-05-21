package dao

import (
	"github.com/FlowRays/WhisperTrail/model"
)

func CreateUser(user *model.User, db *model.Database) error {
	err := db.DB.Create(&user).Error
	return err
}
func GetUserByName(user *model.User, db *model.Database) error {
	err := db.DB.Where("username = ?", user.Username).First(&user).Error
	return err
}

func GetUserByID(user *model.User, db *model.Database) error {
	err := db.DB.Where("id = ?", user.ID).First(&user).Error
	return err
}
