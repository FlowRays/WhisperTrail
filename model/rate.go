package model

type Rate struct {
	ID         uint   `gorm:"primaryKey"`
	LandmarkID string `gorm:"not null" form:"landmark_id" json:"landmark_id"`
	Rating     uint   `gorm:"not null" form:"rating" json:"rating"`
}
