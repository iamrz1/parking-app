package data

import "supertal-tha-parking-app/model"

type ParkingStore interface {
	CreateParkingLot(*model.ParkingLotCreateReq) (*model.ParkingLotRes, error)
	GetParkingLots() ([]*model.ParkingLotRes, error)
	GetParkingLot(lotID uint) (*model.ParkingLotRes, error)
	SetParkingSlotStatus(*model.MaintenanceStatusReq) error
	CreateBookingForUser(lotID, userID uint) (*model.BookingRes, error)
	RemoveBookingForUser(userID uint) (*model.UnBookRes, error)
}
