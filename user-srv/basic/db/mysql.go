package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"Go-mini-kit/user-srv/basic/config"
	"github.com/micro/go-micro/util/log"
)

func initMysql() {
	var err error

	mysqlDB, err = sql.Open("mysql", config.GetMysqlConfig().GetURL())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	mysqlDB.SetMaxOpenConns(config.GetMysqlConfig().GetMaxOpenConnection())

	mysqlDB.SetMaxIdleConns(config.GetMysqlConfig().GetMaxIdleConnection())

	if err = mysqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
}
