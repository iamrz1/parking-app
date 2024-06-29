package api

import (
	"encoding/json"
	"net/http"
	"strconv"
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

type parkingManager struct {
	chi.Router
	ParkingStore data.ParkingStore
}

// NewParkingManagerHandler contains all the apis under parking manager
func NewParkingManagerHandler(pStore data.ParkingStore) http.Handler {
	h := &parkingManager{
		chi.NewRouter(),
		pStore,
	}

	h.registerMiddleware()
	h.registerEndpoints()

	return h
}

func (api *parkingManager) registerMiddleware() {
	api.Use(utils.ManagerOnly)
	api.Use(logger.GenReqID)
	api.Use(chiware.Logger)
	api.Use(middleware.RequestLogger())
	api.Use(middleware.ResponseLogger())
}

func (api *parkingManager) registerEndpoints() {
	api.Group(func(r chi.Router) {
		r.Mount("/parking-lots", api.getParkingLotEndpoints())
		r.Post("/parking-slot-status", api.switchParkingSlotStatus)
	})
}

func (api *parkingManager) getParkingLotEndpoints() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", api.createParkingLot)
	r.Get("/", api.getParkingLots)
	r.Get("/{id}", api.getParkingLot)

	return r
}

// createParkingLot godoc
// @Summary Create a parking lot
// @Description Takes a name and if the name is unique, it creates a Parking Lot with that name
// @Tags Manager
// @Accept  json
// @Produce  json
// @Param Body body model.ParkingLotCreateReq true "All fields are mandatory"
// @Success 200 {object} model.ParkingLotRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /api/v1/manager/parking-lots [post]
func (api *parkingManager) createParkingLot(w http.ResponseWriter, r *http.Request) {
	var input model.ParkingLotCreateReq
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

	res, err := api.ParkingStore.CreateParkingLot(&input)
	if err != nil {
		logger.GetLogger().Errorf("error creating parking lot: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusCreated, "Created Parking Lot", res, nil, true)
}

// getParkingLots godoc
// @Summary Find parking lots
// @Description Find parking lots
// @Tags Manager
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ParkingLotRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /parking-lots [get]
func (api *parkingManager) getParkingLots(w http.ResponseWriter, r *http.Request) {
	res, err := api.ParkingStore.GetParkingLots()
	if err != nil {
		logger.GetLogger().Errorf("error getting parking lots: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusOK, "Fetched Parking Lots", res, nil, true)
}

// getParkingLot godoc
// @Summary Find parking a parking lot
// @Description Find parking a parking lot  by ID
// @Tags Manager
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ParkingLotRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /parking-lots/{id} [get]
func (api *parkingManager) getParkingLot(w http.ResponseWriter, r *http.Request) {
	lotIDText := chi.URLParam(r, "id")
	lotID, _ := strconv.Atoi(lotIDText)
	if lotID == 0 {
		rutils.HandleObjectError(w, rutils.NewValidationError("Invalid Lot ID", nil))
		return
	}

	res, err := api.ParkingStore.GetParkingLot(uint(lotID))
	if err != nil {
		logger.GetLogger().Errorf("error getting parking lot: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusOK, "Fetched Parking Lot", res, nil, true)
}

// switchParkingSlotStatus godoc
// @Summary Switch maintenance mode of a parking slot
// @Description Switch maintenance mode of a parking slot
// @Tags Manager
// @Accept  json
// @Produce  json
// @Param Body body model.MaintenanceStatusReq true "All fields are mandatory"
// @Success 200 {object} model.ParkingLotRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /parking-slot-status [post]
func (api *parkingManager) switchParkingSlotStatus(w http.ResponseWriter, r *http.Request) {
	var input model.MaintenanceStatusReq
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

	err = api.ParkingStore.SetParkingSlotStatus(&input)
	if err != nil {
		logger.GetLogger().Errorf("error setting parking slot status: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusOK, "Maintenance Status Updated", nil, nil, true)
}
