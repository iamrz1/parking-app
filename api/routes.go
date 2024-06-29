package api

import (
	"net/http"
	"supertal-tha-parking-app/config"
	rDB "supertal-tha-parking-app/data/rdbms"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

// NewAPIRouter ...
func NewAPIRouter(appConfig *config.App, db *gorm.DB) http.Handler {
	apiRouter := chi.NewRouter()
	apiRouter.Use(cors.AllowAll().Handler)

	pStore := rDB.NewParkingStore(db)
	uStore := rDB.NewUserStore(db)

	apiRouter.Get("/doc/*", httpSwagger.Handler())
	apiRouter.Route("/api/v1", func(r chi.Router) {
		r.Mount("/manager", NewParkingManagerHandler(pStore))
		r.Mount("/user", NewUserHandler(uStore, pStore))
		r.Mount("/public", NewPublicHandler(uStore))
	})

	return apiRouter
}

// NewSystemRouter ...
func NewSystemRouter() http.Handler {
	systemRouter := chi.NewRouter()
	systemRouter.Mount("/system", NewSystemHandler())

	systemRouter.Mount("/debug", middleware.Profiler())

	return systemRouter
}
