package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	auth "Go-mini-kit/auth/proto/auth"
)

var (
	accessService access.Service
)

// Init 初始化handler
func Init() {
	var err error
	accessService, err = access.GetService()
	if err != nil {
		log.Fatal("[Init] 初始化Handler错误，%s", err)
		return
	}
}
// Service struct 
type Service struct{}

// CreateUserAccessToken 生成token
func (s *Service) CreateUserAccessToken(ctx context.Context, req *auth.Request, rsp *auth.Response) error {
	log.Log("[MakeAccessToken] 收到创建token请求")

	token, err := accessService.MakeAccessToken(&access.Subject{
		ID:   strconv.FormatUint(req.UserId, 10),
		Name: req.UserName,
	})
	if err != nil {
		rsp.Error = &auth.Error{
			Detail: err.Error(),
		}

		log.Logf("[MakeAccessToken] token生成失败，err：%s", err)
		return err
	}

	rsp.Token = token
	return nil
}

// DeleteUserAccessToken 清除用户token
func (s *Service) DeleteUserAccessToken(ctx context.Context, req *auth.Request, rsp *auth.Response) error {
	log.Log("[DelUserAccessToken] 清除用户token")
	err := accessService.DelUserAccessToken(req.Token)
	if err != nil {
		rsp.Error = &auth.Error{
			Detail: err.Error(),
		}

		log.Logf("[DelUserAccessToken] 清除用户token失败，err：%s", err)
		return err
	}

	return nil
}
