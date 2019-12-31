module Go-mini-kit.com/user-srv

go 1.13

require (
	github.com/go-redis/redis v6.15.6+incompatible // indirect
	github.com/golang/protobuf v1.3.2
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1

)

replace (
	go-mini-kit.com/boot => ../boot
	go-mini-kit.com/plugins => ../plugins
	go-mini-kit.com/user-srv => ../user-srv
)
