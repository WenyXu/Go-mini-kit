package db

import (
	"database/sql"
	"fmt"
	"sync"

	"Go-mini-kit/boot"
	"github.com/micro/go-micro/util/log"
)

var (
	_initialized bool
	mysqlDB      *sql.DB
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
	return mysqlDB
}
