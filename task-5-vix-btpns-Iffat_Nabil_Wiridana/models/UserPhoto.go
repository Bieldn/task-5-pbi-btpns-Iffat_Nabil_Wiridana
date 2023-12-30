package models

import "gorm.io/gorm"

type UserPhoto struct {
	gorm.Model
	UserID   int	
	Title    string
	Caption  string
	PhotoUrl string
}
