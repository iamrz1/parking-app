package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username          string `gorm:"unique;index"`
	EncryptedPassword string
	Name              string
	IsParkingManager  *bool
}
