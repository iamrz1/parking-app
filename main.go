package main

import (
	"supertal-tha-parking-app/api"
	"supertal-tha-parking-app/config"
	"supertal-tha-parking-app/conn"
	rDB "supertal-tha-parking-app/data/rdbms"
	"supertal-tha-parking-app/logger"
	"supertal-tha-parking-app/model"
	"sync"

	"gorm.io/gorm"
)

// @title Parking App
// @version v1.0
// @description This is a parking app
// @contact.name Rezoan Tamal
// @contact.email rezoan(.)tamal@gmail.com
// @host localhost:8080
// @BasePath /

func main() {
	appCnf, err := config.AppCnf()
	if err != nil {
		logger.GetLogger().Fatal(err.Error())
		return
	}

	logger.InitLogger(appCnf.Environment, appCnf.LogLevel)

	dbConfig, err := config.RDBMSCnf()
	if err != nil {
		logger.GetLogger().Fatal(err.Error())
	}

	db, dbCloser, err := conn.Connect(&dbConfig)
	if err != nil {
		logger.GetLogger().Fatal(err.Error())
	}

	defer dbCloser()

	model.AutoMigrate(db)
	ensureAdminUser(db)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go runServer(&wg, api.NewServer("api", appCnf.ServerPort, appCnf.Timeout, api.NewAPIRouter(&appCnf, db)))
	go runServer(&wg, api.NewServer("system", appCnf.SystemPort, appCnf.Timeout, api.NewSystemRouter()))
	wg.Wait()
}

func runServer(wg *sync.WaitGroup, server *api.Server) {
	defer wg.Done()
	server.Run()
}

func ensureAdminUser(db *gorm.DB) {
	uStore := rDB.NewUserStore(db)
	isParkingManger := true
	_, _ = uStore.Register(&model.UserCreateReq{
		Name:             "Admin",
		Username:         "admin",
		Password:         "admin",
		IsParkingManager: &isParkingManger,
	})
}
