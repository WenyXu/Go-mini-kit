package db

import (
	"database/sql"
	"go-mini-kit.com/boot/config"

	"github.com/micro/go-micro/util/log"
)

type db struct {
	Mysql mysql `json:"mysql"`
}


// mysql mySQL config
type mysql struct {
	URL               string `json:"url"`
	Enable            bool   `json:"enabled"`
	MaxIdleConnection int    `json:"maxIdleConnection"`
	MaxOpenConnection int    `json:"maxOpenConnection"`
}

func initMysql() {
	log.Logf("[initMysql] initializing Mysql")

	c := config.GetInstance()
	//_config := &db{}

	err := c.Scan("db", _config)
	if err != nil {
		log.Logf("[initMysql] %s", err)
	}

	if !_config.Mysql.Enable {
		log.Logf("[initMysql] Mysql disabled")
		return
	}

	// create connection
	_mysqlDB, err = sql.Open("mysql", _config.Mysql.URL)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// set max connection
	_mysqlDB.SetMaxOpenConns(_config.Mysql.MaxOpenConnection)

	// set max idle connection
	_mysqlDB.SetMaxIdleConns(_config.Mysql.MaxIdleConnection)

	// ping
	if err = _mysqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Logf("[initMysql] mysql connected")
}
