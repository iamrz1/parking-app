package data

import (
	"fmt"
	"supertal-tha-parking-app/data"
	"supertal-tha-parking-app/model"
	"supertal-tha-parking-app/utils"
	"time"

	"github.com/iamrz1/rutils"
	"gorm.io/gorm"
)

type parkingStore struct {
	db *gorm.DB
}

func NewParkingStore(db *gorm.DB) data.ParkingStore {
	return &parkingStore{db: db}
}

func (p *parkingStore) CreateParkingLot(req *model.ParkingLotCreateReq) (*model.ParkingLotRes, error) {
	parkingLot := model.ParkingLot{Name: req.Name}

	err := p.db.Create(&parkingLot).Error
	if err != nil {
		return nil, err
	}

	parkingSlots := make([]model.ParkingSlot, req.NumberOfSlots)

	for i := range parkingSlots {
		parkingSlots[i] = model.ParkingSlot{
			UnderMaintenance: utils.BoolP(false),
			Booked:           utils.BoolP(false),
			ParkingLotID:     parkingLot.ID,
		}
	}

	err = p.db.Create(&parkingSlots).Error
	if err != nil {
		return nil, err
	}

	parkingSlotRes := make([]model.ParkingSlotRes, len(parkingSlots))
	for i := range parkingSlots {
		parkingSlotRes[i] = model.ParkingSlotRes{
			ID:               parkingSlots[i].ID,
			UnderMaintenance: parkingSlots[i].UnderMaintenance,
			Booked:           parkingSlots[i].Booked,
			BookedSince:      parkingSlots[i].BookedSince,
		}
	}

	response := model.ParkingLotRes{
		ID:    parkingLot.ID,
		Name:  parkingLot.Name,
		Slots: parkingSlotRes,
	}

	return &response, nil
}

func (p *parkingStore) GetParkingLots() ([]*model.ParkingLotRes, error) {
	var parkingLots []model.ParkingLot

	db := p.db

	err := db.Find(&parkingLots).Error
	if err != nil {
		return nil, err
	}

	res := make([]*model.ParkingLotRes, len(parkingLots))

	for i, _ := range parkingLots {
		res[i] = &model.ParkingLotRes{
			ID:   parkingLots[i].ID,
			Name: parkingLots[i].Name,
		}
	}

	return res, err
}

func (p *parkingStore) GetParkingLot(id uint) (*model.ParkingLotRes, error) {
	var parkingLot model.ParkingLot

	err := p.db.Preload("Slots").First(&parkingLot, id).Error
	if err != nil {
		return nil, err
	}

	parkingSlotRes := make([]model.ParkingSlotRes, len(parkingLot.Slots))
	for i := range parkingLot.Slots {
		parkingSlotRes[i] = model.ParkingSlotRes{
			ID:               parkingLot.Slots[i].ID,
			UnderMaintenance: parkingLot.Slots[i].UnderMaintenance,
			Booked:           parkingLot.Slots[i].Booked,
			BookedSince:      parkingLot.Slots[i].BookedSince,
		}
	}

	res := model.ParkingLotRes{
		ID:    parkingLot.ID,
		Name:  parkingLot.Name,
		Slots: parkingSlotRes,
	}

	return &res, err
}

func (p *parkingStore) SetParkingSlotStatus(req *model.MaintenanceStatusReq) error {
	return p.db.Model(&model.ParkingSlot{}).Where("id = ?", req.SlotID).
		Update("under_maintenance", *req.MaintenanceMode).Error
}

func (p *parkingStore) CreateBookingForUser(lotID, userID uint) (*model.BookingRes, error) {
	var parkingSlot model.ParkingSlot

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// check if already booked
	var count int64
	err := tx.Model(&model.ParkingSlot{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if count > 0 {
		tx.Rollback()
		return nil, rutils.NewValidationError("Already Booked", nil)
	}

	// find available spot
	target := model.ParkingSlot{
		ParkingLotID:     lotID,
		UnderMaintenance: utils.BoolP(false),
		Booked:           utils.BoolP(false),
	}
	err = tx.Where(&target).Order("id").First(&parkingSlot).Error
	if err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, rutils.NewValidationError("No slot available", nil)
		}

		return nil, err
	}

	// save parking spot
	parkingSlot.UserID = userID
	parkingSlot.Booked = utils.BoolP(true)
	parkingSlot.BookedSince = time.Now().UTC()
	err = tx.Save(&parkingSlot).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res := &model.BookingRes{
		LotID:    parkingSlot.ParkingLotID,
		SlotID:   parkingSlot.ID,
		ParkedAt: parkingSlot.BookedSince,
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return res, err
}

func (p *parkingStore) RemoveBookingForUser(userID uint) (*model.UnBookRes, error) {
	tx := p.db.Begin()

	var parkingSlot model.ParkingSlot
	err := tx.Model(&model.ParkingSlot{}).Where("user_id = ?", userID).First(&parkingSlot).Error
	if err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, rutils.NewValidationError("Not parked", nil)
		}

		return nil, err
	}

	timeElapsed := time.Now().UTC().Sub(parkingSlot.BookedSince)
	charge := utils.ChargePerHour * utils.DurationToHours(timeElapsed)

	// update parking spot
	parkingSlot.UserID = 0
	parkingSlot.Booked = utils.BoolP(false)
	parkingSlot.BookedSince = time.Time{}
	err = tx.Save(&parkingSlot).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res := &model.UnBookRes{
		Charge: fmt.Sprintf("Rs. %d", charge),
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return res, err
}
