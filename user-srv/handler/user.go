package handler

import (
	"context"
	"github.com/micro/go-micro/util/log"

	userModel "Go-mini-kit/user-srv/model/user"
	userProto "Go-mini-kit/user-srv/proto/user"
)

var (
	userService userModel.IService
)
// Init Handler
func Init() {

    var err error
    userService, err = userModel.GetService()
    if err != nil {
        log.Fatal("[Init] Init Handler failed")
        return
    }
}

type Service struct{}

// QueryUserByName 
func (srv *Service) QueryUserByName(ctx context.Context, req *userProto.Request, rsp *userProto.Response) error {

    user, err := userService.QueryUserByName(req.UserName)

    if err != nil {
        rsp.Success = false
        rsp.Error = &userProto.Error{
            Code:   500,
            Detail: err.Error(),
        }
        return nil
    }
    rsp.User = user
    rsp.Success = true

    return nil
}