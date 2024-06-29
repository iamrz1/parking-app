package api

import (
	"encoding/json"
	"net/http"
	"supertal-tha-parking-app/data"
	_ "supertal-tha-parking-app/docs"
	cerror "supertal-tha-parking-app/error"
	"supertal-tha-parking-app/logger"
	"supertal-tha-parking-app/middleware"
	"supertal-tha-parking-app/model"

	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"github.com/iamrz1/rutils"
)

type public struct {
	chi.Router
	UserStore data.UserStore
}

// NewPublicHandler contains all the apis under parking manager
func NewPublicHandler(uStore data.UserStore) http.Handler {
	h := &public{
		chi.NewRouter(),
		uStore,
	}

	h.registerMiddleware()
	h.registerEndpoints()

	return h
}

func (api *public) registerMiddleware() {
	api.Use(logger.GenReqID)
	api.Use(chiware.Logger)
	api.Use(middleware.RequestLogger())
	api.Use(middleware.ResponseLogger())
}

func (api *public) registerEndpoints() {
	api.Group(func(r chi.Router) {
		r.Post("/register", api.Register)
		r.Post("/login", api.Login)
	})
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags Public
// @Accept  json
// @Produce  json
// @Param Body body model.UserCreateReq true "All fields are mandatory"
// @Success 200 {object} model.UserCreateRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /api/v1/register [post]
func (api *public) Register(w http.ResponseWriter, r *http.Request) {
	var input model.UserCreateReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.GetLogger().Errorf("error decoding input: %v", err)
		rutils.HandleObjectError(w, rutils.NewValidationError("invalid JSON", nil))
		return
	}

	if err = input.Validate(); err != nil {
		rutils.HandleObjectError(w, rutils.NewValidationError("", nil))
		return
	}

	res, err := api.UserStore.Register(&input)
	if err != nil {
		logger.GetLogger().Errorf("error registering user: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusCreated, "Registered User", res, nil, true)
}

// Login godoc
// @Summary Login as user
// @Description Login as user
// @Tags Public
// @Accept  json
// @Produce  json
// @Param Body body model.LoginReq true "All fields are mandatory"
// @Success 200 {object} model.LoginRes
// @Failure 400 {object} utils.GenericErrorResponse
// @Failure 422 {object} utils.GenericErrorResponse
// @Failure 500 {object} utils.GenericErrorResponse
// @Router /api/v1/login [post]
func (api *public) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.GetLogger().Errorf("error decoding input: %v", err)
		rutils.HandleObjectError(w, rutils.NewValidationError("invalid JSON", nil))
		return
	}

	if err = input.Validate(); err != nil {
		rutils.HandleObjectError(w, rutils.NewValidationError("", nil))
		return
	}

	res, err := api.UserStore.Login(&input)
	if err != nil {
		logger.GetLogger().Errorf("error login: %v", err)
		rutils.HandleObjectError(w, cerror.GetValidationErr(err))
		return
	}

	rutils.ServeJSONObject(w, http.StatusOK, "Logged In", res, nil, true)
}
