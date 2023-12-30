package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string
	Email     string
	Password  string
	UserPhoto UserPhoto `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
