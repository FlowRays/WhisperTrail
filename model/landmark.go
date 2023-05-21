package model

type Landmark struct {
	ID         uint   `gorm:"primaryKey"`
	Path       string `gorm:"size:255"`
	Latitude   string `gorm:"not null"`
	Longitude  string `gorm:"not null"`
	Text       string
	IsLoggedIn bool `gorm:"not null"`
	UserID     uint
}
