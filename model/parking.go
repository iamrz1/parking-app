package model

import (
	"time"

	"gorm.io/gorm"
)

type ParkingLot struct {
	gorm.Model
	Name  string        `gorm:"unique"`
	Slots []ParkingSlot `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ParkingSlot struct {
	gorm.Model
	UnderMaintenance *bool
	Booked           *bool
	ParkingLotID     uint `gorm:"index"`
	BookedSince      time.Time
	UserID           uint // a parking lot is either available or belong to a user
	User             User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
