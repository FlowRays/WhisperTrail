package model

import "gorm.io/gorm"

type Database struct {
	DB *gorm.DB
}