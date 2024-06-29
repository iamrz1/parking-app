package config

import (
	"fmt"
	cerror "supertal-tha-parking-app/error"
	"sync"

	"github.com/spf13/viper"
)

// RDBMS ...
type RDBMS struct {
	DSN             string
	Dialect         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // in minutes
}

func (cnf *RDBMS) validate() error {
	err := cerror.ValidationError{}
	if cnf.DSN == "" {
		err.Add("dsn", "required")
	}
	if cnf.Dialect == "" {
		err.Add("dialect", "required")
	}

	if len(err) > 0 {
		return fmt.Errorf("rdbms configuration error: %v", err)
	}

	return nil
}

var rdbmsCnf RDBMS
var rdbmsErr error
var rdbmsOnce = sync.Once{}

func loadRDBMS() {
	rdbmsCnf = RDBMS{
		DSN:             viper.GetString("rdbms.dsn"),
		Dialect:         viper.GetString("rdbms.dialect"),
		MaxOpenConns:    viper.GetInt("rdbms.max_open_conns"),
		MaxIdleConns:    viper.GetInt("rdbms.max_idle_conns"),
		ConnMaxLifetime: viper.GetInt("rdbms.conn_max_lifetime"),
	}
}

// RDBMSCnf ...
func RDBMSCnf() (RDBMS, error) {
	rdbmsOnce.Do(func() {
		read()
		loadRDBMS()
		rdbmsErr = rdbmsCnf.validate()
	})
	return rdbmsCnf, rdbmsErr
}
