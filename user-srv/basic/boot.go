package basicBoot

import (
	"Go-mini-kit/user-srv/basic/config"
	"Go-mini-kit/user-srv/basic/db"
)

func Init() {
	config.Init()
	db.Init()
}
