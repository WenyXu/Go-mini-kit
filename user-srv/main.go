package main

import (
	"fmt"
	"github.com/micro/go-plugins/config/source/grpc"
	"go-mini-kit.com/user-srv/model"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/util/log"

	"go-mini-kit.com/boot"
	"go-mini-kit.com/boot/common"
	"go-mini-kit.com/boot/config"
	"go-mini-kit.com/user-srv/handler"
	userProto "go-mini-kit.com/user-srv/proto/user"
)

var (
	_appName = "user_srv"
	_config  = &userConfig{}
)
type userConfig struct {
	common.AppConfig
}

func main() {

	// Basic Config Init
	initConfig()
	//boot.Init()
	
	reg := etcd.NewRegistry(registryOptions)

	// New Service
	service := micro.NewService(
		micro.Name("im.terminal.go.srv.user"),
		micro.Registry(reg),
		micro.Version("latest"),
	)

	// Initialize service
	service.Init(
		micro.Action(func(c *cli.Context) {
			// Initialize model
			model.Init()
			// Initialize handler
			handler.Init()
		}),
	)

	// Register Handler
	if err:=userProto.RegisterUserHandler(service.Server(), new(handler.Service)); err!=nil{
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	etcdCfg := &common.EtcdConfig{}
	err := config.GetInstance().Scan("etcd", etcdCfg)
	if err != nil {
		panic(err)
	}
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.Host, etcdCfg.Port)}
}

func initConfig(){
	source := grpc.NewSource(
		grpc.WithAddress("127.0.0.1:9600"),
		grpc.WithPath("micro"),
	)

	boot.Init(config.WithSource(source))

	err := config.GetInstance().Scan(_appName, _config)
	if err != nil {
		panic(err)
	}

	log.Logf("[initCfg] config，_config：%v", _config)

	return
}
