package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `gorm:"type:varchar(64);uniqueIndex" json:"username"`
	Email       string `gorm:"type:varchar(128);uniqueIndex" json:"email"`
	Bio         string `gorm:"type:text" json:"bio"`
	Gender      string `gorm:"type:varchar(16)" json:"gender"`
	Nationality string `gorm:"type:varchar(64)" json:"nationality"`
}
