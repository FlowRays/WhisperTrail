package model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null" form:"username" json:"username"`
	Password string `gorm:"not null" form:"password" json:"password"`
}
