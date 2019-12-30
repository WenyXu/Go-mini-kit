package db

import (
	"go-mini-kit.com/boot/config"
	"database/sql"

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

	c := config.Instance()
	cfg := &db{}

	err := c.Scan("db", cfg)
	if err != nil {
		log.Logf("[initMysql] %s", err)
	}

	if !cfg.Mysql.Enable {
		log.Logf("[initMysql] Mysql disabled")
		return
	}

	// create connection
	mysqlDB, err = sql.Open("mysql", cfg.Mysql.URL)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// set max connection
	mysqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConnection)

	// set max idle connection
	mysqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConnection)

	// ping
	if err = mysqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Logf("[initMysql] mysql connected")
}
