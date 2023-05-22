package model

type Landmark struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Path       string `gorm:"size:255"  json:"path"`
	Latitude   string `gorm:"not null"  json:"latitude"`
	Longitude  string `gorm:"not null" json:"longitude"`
	Text       string `json:"text"`
	IsLoggedIn bool   `gorm:"not null" json:"is_logged_in"`
	UserID     uint   `json:"user_id"`
}
