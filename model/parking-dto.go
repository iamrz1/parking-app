package model

import "time"

type ParkingLotCreateReq struct {
	Name          string
	NumberOfSlots int
}

type ParkingLotUpdateReq struct {
	Name string
}

type ParkingLotRes struct {
	ID    uint
	Name  string
	Slots []ParkingSlotRes `json:"Slots,omitempty"`
}

type ParkingSlotCreateReq struct {
	Name string
}

type ParkingSlotUpdateReq struct {
	Name string
}

type ParkingSlotRes struct {
	ID               uint
	UnderMaintenance *bool
	Booked           *bool
	BookedSince      time.Time `json:"BookedAt,omitempty"`
}

type BookingReq struct {
	LotID uint
}

type BookingRes struct {
	LotID    uint
	SlotID   uint
	ParkedAt time.Time
}

type UnBookRes struct {
	Charge string
}

type MaintenanceStatusReq struct {
	SlotID          uint
	MaintenanceMode *bool
}
