package model

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&ParkingLot{})
	db.AutoMigrate(&ParkingSlot{})
}
