package main

import (
	"Go-mini-kit/auth/handler"
	"Go-mini-kit/auth/subscriber"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	auth "Go-mini-kit/auth/proto/auth"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("im.terminal.go.srv.auth.srv.auth"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	_ = auth.RegisterServiceHandler(service.Server(), new(handler.Service))

	// Register Struct as Subscriber
	_ = micro.RegisterSubscriber("im.terminal.go.srv.auth.srv.auth", service.Server(), new(subscriber.Auth))

	// Register Function as Subscriber
	_ = micro.RegisterSubscriber("im.terminal.go.srv.auth.srv.auth", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
