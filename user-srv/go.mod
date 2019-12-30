module Go-mini-kit.com/user-srv

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	go-mini-kit.com/plugins v0.0.0-00010101000000-000000000000 // indirect
	go-mini-kit.com/user-srv v0.0.0-00010101000000-000000000000 // indirect

)

replace (
	go-mini-kit.com/boot => ../boot
	go-mini-kit.com/plugins => ../plugins
	go-mini-kit.com/user-srv => ../user-srv

)
