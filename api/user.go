package api

import (
	"encoding/json"
	"log"
	"net/http"
	"supertal-tha-parking-app/data"
	_ "supertal-tha-parking-app/docs"
	cerror "supertal-tha-parking-app/error"
	"supertal-tha-parking-app/logger"
	"supertal-tha-parking-app/middleware"
	"supertal-tha-parking-app/model"
	"supertal-tha-parking-app/utils"

	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"github.com/iamrz1/rutils"
)

type user struct {
	chi.Router
	UserStore    data.UserStore
	ParkingStore data.ParkingStore
}

// NewUserHandler contains all the apis under parking manager
func NewUserHandler(uStore data.UserStore, pStore data.ParkingStore) http.Handler {
	h := &user{
		chi.NewRouter(),
		uStore,
		pStore,
	}

	h.registerMiddleware()
	h.registerEndpoints()

	return h
}

func (api *user) registerMiddleware() {
	api.Use(utils.UserOnly)
	api.Use(logger.GenReqID)
	api.Use(chiware.Logger)
	api.Use(middleware.RequestLogger())
	api.Use(middleware.ResponseLogger())
}

func (api *user) registerEndpoints() {
	api.Group(func(r chi.Router) {
		r.Get("/parking-lots", api.getParkingLots)
		r.Post("/park", api.bookParkingSlot)
		r.Post("/unpark", api.unbookParkingSlot)
	})
}

// getParkingLots godoc
// @Summary Find parking lots
// @Description Find parking lots
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ParkingLotRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /api/v1/user/parking-lots [get]
func (api *user) getParkingLots(w http.ResponseWriter, r *http.Request) {
	res, err := api.ParkingStore.GetParkingLots()
	if err != nil {
		logger.GetLogger().Errorf("error getting parking lots: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusOK, "Fetched Parking Lots", res, nil, true)
}

// bookParkingSlot godoc
// @Summary Book a parking slot in a parking lot
// @Description Book a parking slot by lot ID
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} model.BookingRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /api/v1/user/park [post]
func (api *user) bookParkingSlot(w http.ResponseWriter, r *http.Request) {
	var input model.BookingReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.GetLogger().Errorf("error decoding input: %v", err)
		rutils.HandleObjectError(w, rutils.NewValidationError("invalid JSON", nil))
		return
	}

	if err = input.Validate(); err != nil {
		rutils.HandleObjectError(w, rutils.NewValidationError("", err))
		return
	}

	currentUser, err := api.UserStore.GetUser(r.Header.Get(utils.UsernameKey))
	if err != nil {
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	res, err := api.ParkingStore.CreateBookingForUser(input.LotID, currentUser.ID)
	if err != nil {
		logger.GetLogger().Errorf("error booking parking slot: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusOK, "Booked Parking Slot", res, nil, true)
}

// unbookParkingSlot godoc
// @Summary Book a parking slot in a parking lot
// @Description Book a parking slot by lot ID
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} model.BookingRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /api/v1/user/unpark [post]
func (api *user) unbookParkingSlot(w http.ResponseWriter, r *http.Request) {
	currentUser, err := api.UserStore.GetUser(r.Header.Get(utils.UsernameKey))
	if err != nil {
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	log.Println("got current user")

	res, err := api.ParkingStore.RemoveBookingForUser(currentUser.ID)
	if err != nil {
		logger.GetLogger().Errorf("error unbooking parking slot: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusOK, "Unbooked Parking Slot", res, nil, true)
}
