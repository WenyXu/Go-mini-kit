package db

import (
	"database/sql"
	"fmt"
	"sync"

	"Go-mini-kit/user-srv/basic/config"
	"github.com/micro/go-micro/util/log"
)

var (
	inited  bool
	mysqlDB *sql.DB
	m       sync.RWMutex
)

// Init DB
func Init() {
	m.Lock()
	defer m.Unlock()

	var err error

	if inited {
		err = fmt.Errorf("[Init] db inited")
		log.Logf(err.Error())
		return
	}
	
	// if the config enabled is true
	if config.GetMysqlConfig().GetEnabled() {
		initMysql()
	}

	inited = true
}

// GetDB 
func GetDB() *sql.DB {
	return mysqlDB
}
