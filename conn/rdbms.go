package conn

import (
	"fmt"
	"supertal-tha-parking-app/config"
	"supertal-tha-parking-app/logger"
	"time"

	"gorm.io/driver/mysql" // enable mysql dialect
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect connects to a preconfigured database ...
func Connect(cnf *config.RDBMS) (*gorm.DB, func(), error) {
	logger.GetLogger().Info("connecting to database...")

	dialector, err := getDialector(cnf.Dialect, cnf.DSN)
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(*dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.GetLogger().Errorf("unable to get sql database: %v", err)
	}

	if cnf.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(cnf.MaxOpenConns)
	}

	if cnf.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(cnf.MaxIdleConns)
	}

	if cnf.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(cnf.ConnMaxLifetime))
	}

	logger.GetLogger().Info("database connected")
	return db, func() {
		logger.GetLogger().Info("disconnecting database...")
		if err := sqlDB.Close(); err != nil {
			logger.GetLogger().Errorf("unable to disconnect database: %v", err)
		}
		logger.GetLogger().Info("database disconnected")
	}, nil
}

func getDialector(dialect, dsn string) (*gorm.Dialector, error) {
	var dialector gorm.Dialector
	switch dialect {
	case "postgres":
		dialector = postgres.Open(dsn)
	case "mysql":
		dialector = mysql.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database dialect %s", dialect)
	}

	return &dialector, nil
}
