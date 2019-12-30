package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/micro/go-micro/util/log"
	"go-mini-kit.com/boot"
)

var (
	_initialized bool
	_mysqlDB     *sql.DB
	_config		 =&db{}
	m            sync.RWMutex
)

func init() {
	boot.Register(initDB)
}

// initDB initialize DB
func initDB() {
	m.Lock()
	defer m.Unlock()

	var err error

	if _initialized {
		err = fmt.Errorf("[initDB] db initialized")
		log.Logf(err.Error())
		return
	}

	initMysql()

	_initialized = true
}

// GetDB get DB
func GetDB() *sql.DB {
	return _mysqlDB
}
