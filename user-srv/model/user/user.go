package user

import (
	"fmt"
	"sync"

	"go-mini-kit.com/plugins/db"
	userProto "go-mini-kit.com/user-srv/proto/user"
	"github.com/micro/go-micro/util/log"
)

var (
	srv *service
	m sync.RWMutex
)

type service struct{
}

type IService interface {
	QueryUserByName(userName string)(res *userProto.User,err error)
}

func Init()  {
	m.Lock()
	defer m.Unlock()

	if srv !=nil{
		return
	}
	srv = &service{}
}

func GetService()(IService,error){
	if srv ==nil{
		return nil,fmt.Errorf(("[GetService] GetService srv was not inited"))
	}
	return srv,nil
}

//QueryUserByName TODO: move these funcs into a new file
func (s *service) QueryUserByName(userName string) (res *userProto.User, err error) {
	queryString := `SELECT user_id, user_name FROM user WHERE user_name = ?`

	// connect DB
	o := db.GetDB()

	res = &userProto.User{}

	// Query
	err = o.QueryRow(queryString, userName).Scan(&res.Id, &res.Name)
	if err != nil {
		log.Logf("[QueryUserByName] Query failed，err：%s", err)
		return
	}
	return
}
