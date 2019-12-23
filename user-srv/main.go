package main

import (
	"fmt"

	"github.com/micro/go-micro"
	"github.com/micro/cli"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/util/log"

	userProto "Go-mini-kit/user-srv/proto/user"
	"Go-mini-kit/user-srv/basic"
	"Go-mini-kit/user-srv/basic/config"
	"Go-mini-kit/user-srv/handler"
	"Go-mini-kit/user-srv/model"
)

func main() {

	// Basic Config Init 
	basicBoot.Init()
	
	reqistry := etcd.NewRegistry(registryOptions)

	// New Service
	service := micro.NewService(
		micro.Name("im.terminal.go.srv.user"),
		micro.Registry(reqistry),
		micro.Version("latest"),
	)

	// Initialize service
	service.Init(
		micro.Action(func(c *cli.Context) {
			// Initialize model
			modelBoot.Init()
			// Initialize handler
			handler.Init()
		}),
	)

	// Register Handler
	userProto.RegisterUserHandler(service.Server(), new(handler.Service))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	etcdCfg := config.GetEtcdConfig()
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.GetHost(), etcdCfg.GetPort())}
}

