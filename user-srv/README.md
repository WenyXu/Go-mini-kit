# 构建 Go-micro 构建微服务-1 Quick Start

# 前言
今天我们很高兴，能和越来越多的地方政府、企业合作伙伴落地更多颇具规模的项目。同时今年我们做了一个大胆的决定，逐步使用 Go 构建微服务并替换我们后端业务。为此就有了这个系列的文章，意在帮助我们同学快速上手同时提供一些示例。

# 前提
了解 Golang 语法

# 概览
本章将演示实现一个 Tiny 用户服务
## 包含内容
使用 Go-micro/micro ，使用 ProtoBuf 构建 gRPC 服务，读取本地配置文件示例。
## 不包含内容
配置中心，日志持久化，扩容，熔断，降级，容错，健康检查，链路追踪，应用容器化。

# 知识地图
Golang

[Document](https://golang.org/doc/)
[Document/zh-cn](https://go-zh.org/doc/) 
gRPC
[Document ](https://www.grpc.io/docs/)
[Quick Start with Go](https://www.grpc.io/docs/quickstart/go/)
Etcd
[Github](https://github.com/etcd-io/etcd)
Mysql 5.7
[Document](https://dev.mysql.com/doc/refman/5.7/en/tutorial.html)
Docker
[Document](https://docs.docker.com/)
[Docker 入门教程-阮一峰](https://www.ruanyifeng.com/blog/2018/02/docker-tutorial.html)
Go-micro/Micro
[Github](https://github.com/micro/go-micro)
[Document](https://micro.mu/docs/framework.html)
[Document/zh-cn](https://micro.mu/docs/cn/go-micro.html)
[Micro Tutorials/zh-cn](https://github.com/micro-in-cn/tutorials)


# 约定
## Interface 
### 命名
Prefix: I
Name Style: {Prefix}+{Name}
e.g.

```go
type IService interface {
	...
}
```

## Model 
### Import 别名
Suffix: Model
Name Style {name}+{Suffix}
e.g.

```go
import (
	userModel "user-srv/model/user"
)
```

## Proto 
### Import 别名
Suffix: Proto
Name Style {name}+{Suffix}

```go
import (
	userProto "user-srv/proto/user"
)
```

# 本文开发环境
```shell
Distributor ID: Ubuntu
Description:    Ubuntu 18.04.2 LTS
Release:        18.04
Codename:       bionic
```

```shell
go version go1.13.5 linux/amd64
```

```shell
micro version 1.18.0
```

# 准备工作
## Golang
### 安装
[Document/Install](https://golang.org/doc/install)
[Document/zh-cn/Install](https://go-zh.org/doc/install)
## gRPC
### 安装
[Document/Install](https://www.grpc.io/docs/quickstart/go/)
## Docker
### 安装
[Doucment](https://docs.docker.com/get-docker/)
## Micro
### 安装
**请使用 Go modules ，环境变量 GO111MODULE=on**
```shell
go get github.com/micro/go-micro
go get github.com/micro/micro
```
## Mysql
### 安装
这里我们给出一个简单的 docker-compose.yaml
```yaml
version: '3'
services:
  mysql:
    network_mode: "host"
    environment:
    	MYSQL_ROOT_PASSWORD: "root"
    image: "docker.io/mysql:5.7" 
    restart: always
    volumes:
    	- "./db:/var/lib/mysql"
    	- "./conf/my.cnf:/etc/my.cnf"

```
同一目录下创建 conf ,db 文件夹，以及/conf/my.cnf 文件
```shell
├── conf
│   └── my.cnf
├── db
└── docker-compose.yaml
```
/conf/my.cnf demo
```shell
[mysqld]
user=mysql
default-storage-engine=INNODB
character-set-server=utf8mb4
[client]
default-character-set=utf8mb4
[mysql]
default-character-set=utf8mb4
```

### 运行
```shell
docker-compose up -d
```

进入 mysql e.g.
```shell
docker exec -t -i {CONTAINER ID or NAMES} /bin/bash
```

```shell
weny@weny-server:~$ docker ps
CONTAINER ID        IMAGE                 COMMAND                  CREATED             STATUS              PORTS
     NAMES
c27980918233        mysql:5.7             "docker-entrypoint.s…"   11 hours ago        Up 11 hours
     mysql_mysql_1
95ee68e8b6a1        netdata/netdata       "/usr/sbin/run.sh"       21 hours ago        Up 21 hours         0.0.0.0:19999->19999/tcp
     netdata_netdata_1
658dffa0a871        quay.io/coreos/etcd   "/usr/local/bin/etcd…"   21 hours ago        Up 21 hours         0.0.0.0:2379->2379/tcp, 2380/tcp   etcd_proxy_1
da6ad7972a85        quay.io/coreos/etcd   "/usr/local/bin/etcd…"   21 hours ago        Up 21 hours         2379-2380/tcp
     etcd_node3_1
95fb0c950bdb        quay.io/coreos/etcd   "/usr/local/bin/etcd…"   21 hours ago        Up 21 hours         2379-2380/tcp
     etcd_node1_1
3f36ba028069        quay.io/coreos/etcd   "/usr/local/bin/etcd…"   21 hours ago        Up 21 hours         2379-2380/tcp
     etcd_node2_1
weny@weny-server:~$ docker exec -t -i c27980918233 /bin/bash
root@weny-server:/# mysql -u root -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 4
Server version: 5.7.28 MySQL Community Server (GPL)

Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> exit
Bye
```

## Etcd
### 安装
这里我们给出一个简单的 docker-compose.yaml
```yaml
version: "3"
services:

  proxy:
    image: quay.io/coreos/etcd
    networks:
      etcd_cluster_net:
    environment:
      - ETCDCTL_API=3
    ports:
      - "2379:2379"
    command:
      - /usr/local/bin/etcd
      - grpc-proxy
      - start
      - --listen-addr
      - 0.0.0.0:2379
      - --endpoints
      - 172.16.238.100:2379,172.16.238.101:2379,172.16.238.102:2379

  node1:
    image: quay.io/coreos/etcd
    volumes:
      - node1-data:/etcd-data
    expose:
      - 2379
      - 2380   
    networks:
      etcd_cluster_net:
        ipv4_address: 172.16.238.100
    environment:
      - ETCDCTL_API=3
    command:
      - /usr/local/bin/etcd
      - --data-dir=/etcd-data
      - --name
      - node1
      - --initial-advertise-peer-urls
      - http://172.16.238.100:2380
      - --listen-peer-urls
      - http://0.0.0.0:2380
      - --advertise-client-urls
      - http://172.16.238.100:2379
      - --listen-client-urls
      - http://0.0.0.0:2379
      - --initial-cluster
      - node1=http://172.16.238.100:2380,node2=http://172.16.238.101:2380,node3=http://172.16.238.102:2380
      - --initial-cluster-state
      - new
      - --initial-cluster-token
      - docker-etcd

  node2:
    image: quay.io/coreos/etcd
    volumes:
      - node2-data:/etcd-data
    networks:
      etcd_cluster_net:
        ipv4_address: 172.16.238.101
    environment:
      - ETCDCTL_API=3
    expose:
      - 2379
      - 2380
    command:
      - /usr/local/bin/etcd
      - --data-dir=/etcd-data
      - --name
      - node2
      - --initial-advertise-peer-urls
      - http://172.16.238.101:2380
      - --listen-peer-urls
      - http://0.0.0.0:2380
      - --advertise-client-urls
      - http://172.16.238.101:2379
      - --listen-client-urls
      - http://0.0.0.0:2379
      - --initial-cluster
      - node1=http://172.16.238.100:2380,node2=http://172.16.238.101:2380,node3=http://172.16.238.102:2380
      - --initial-cluster-state
      - new
      - --initial-cluster-token
      - docker-etcd

  node3:
    image: quay.io/coreos/etcd
    volumes:
      - node3-data:/etcd-data
    networks:
      etcd_cluster_net:
        ipv4_address: 172.16.238.102
    environment:
      - ETCDCTL_API=3
    expose:
      - 2379
      - 2380
    command:
      - /usr/local/bin/etcd
      - --data-dir=/etcd-data
      - --name
      - node3
      - --initial-advertise-peer-urls
      - http://172.16.238.102:2380
      - --listen-peer-urls
      - http://0.0.0.0:2380
      - --advertise-client-urls
      - http://172.16.238.102:2379
      - --listen-client-urls
      - http://0.0.0.0:2379
      - --initial-cluster
      - node1=http://172.16.238.100:2380,node2=http://172.16.238.101:2380,node3=http://172.16.238.102:2380
      - --initial-cluster-state
      - new
      - --initial-cluster-token
      - docker-etcd

volumes:
  node1-data:
  node2-data:
  node3-data:

networks:
  etcd_cluster_net:
    driver: bridge
    ipam:
      driver: default
      config:
      -
        subnet: 172.16.238.0/24
```
### 运行
```shell
docker-compose up -d
```
### 验证
```shell
docker exec -t {CONTAINER ID or NAMES} etcdctl member list
daf3fd52e3583ff, started, node3, http://172.16.238.102:2380, http://172.16.238.102:2379
422a74f03b622fef, started, node1, http://172.16.238.100:2380, http://172.16.238.100:2379
ed635d2a2dbef43d, started, node2, http://172.16.238.101:2380, http://172.16.238.101:2379
```

# 开始
## 预期
服务层 user-srv 实现 QueryUserByName 
## 相关阅读
micro new command 
[Document/zh-cn](https://micro.mu/docs/cn/new.html) （推荐，较于英文文档示例较多)
[Document](https://micro.mu/docs/new.html)
## 使用脚手架生成
```shell
micro new --namespace=im.terminal.go --type=srv --alias=user Go-mini-kit/user-srv
```

```shell
Creating service im.terminal.go.srv.user in /home/weny/Projects/Golang/src/Go-mini-kit/user-srv

.
├── main.go
├── generate.go
├── plugin.go
├── handler
│   └── user.go
├── subscriber
│   └── user.go
├── proto/user
│   └── user.proto
├── Dockerfile
├── Makefile
├── README.md
└── go.mod
```
+新增 basic 目录
+新增 conf 目录
+新增 model 目录
-删除 subscribe 目录
```shell
.
├── basic* 新增 basic 目录												
│   ├── boot.go
│   ├── config
│   │   ├── config.go
│   │   ├── etcd.go
│   │   ├── mysql.go
│   │   └── profiles.go
│   └── db
│       ├── db.go
│       └── mysql.go
├── conf* 新增 conf 目录
│   ├── application-db.yml
│   ├── application-etcd.yml
│   └── application.yml
├── Dockerfile
├── generate.go
├── go.mod
├── go.sum
├── handler
│   └── user.go
├── main.go
├── Makefile
├── model* 新增 model 目录
│   ├── boot.go
│   └── user
│       └── user.go
├── plugin.go
├── proto
│   └── user
│       ├── user.pb.go
│       ├── user.pb.micro.go
│       └── user.proto
└── README.md
```

## User Proto
### 相关阅读
[proto buffers overview](https://developers.google.com/protocol-buffers/docs/overview)
[gRPC getting started with go](https://www.grpc.io/docs/quickstart/go/)
### Code
/proto/user/user.proto
```go
syntax = "proto3";

package im.terminal.go.srv.user;

service User {
    rpc QueryUserByName (Request) returns (Response) {
    }
}

message user {
    int64 id = 1;
    string name = 2;
	uint64 createdTime = 3;
    uint64 updatedTime = 4;
}

message Error {
    int32 code = 1;
    string detail = 2;
}

message Request {
    string userID = 1;
    string userName = 2;
}

message Response {
    bool success = 1;
    Error error = 2;
    user user = 3;
}
```
### Generate gRPC code
#### 相关阅读
[protoc-gen-micro](https://github.com/micro/protoc-gen-micro)
```shell
weny@weny-server:~/Projects/Golang/src/Go-mini-kit/user-srv$ protoc --proto_path=. --go_out=. --micro_out=. proto/user/user.proto
```
### DB table
```shell
CREATE TABLE `user`
(
    `id`           int(10) unsigned NOT NULL AUTO_INCREMENT ,
    `user_id`      int(10) unsigned DEFAULT NULL ,
    `user_name`    varchar(20) CHARACTER ,
    `created_time` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_time` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_user_name_uindex` (`user_name`),
    UNIQUE KEY `user_user_id_uindex` (`user_id`)
) ENGINE = InnoDB;
```
### Insert sth
By command line
```shell
INSERT INTO user (user_id, user_name) VALUE (1, 'micro');
```
Or GUI 

## Basic
读取配置文件
```shell
├── basic				
│   ├── boot.go
│   ├── config
│   │   ├── config.go
│   │   ├── etcd.go
│   │   ├── mysql.go
│   │   └── profiles.go
│   └── db
│       ├── db.go
│       └── mysql.go
```
### 相关阅读
[document/config](https://micro.mu/docs/go-config.html)
[document/zh-cn/config](https://micro.mu/docs/cn/go-config.html)
### Code
[https://github.com/WenyXu/Go-mini-kit/tree/master/user-srv/basic](https://github.com/WenyXu/Go-mini-kit/tree/master/user-srv/basic)


## Conf
配置文件
```shell
├── conf
│   ├── application-db.yml
│   ├── application-etcd.yml
│   └── application.yml
```
### Files
[https://github.com/WenyXu/Go-mini-kit/tree/master/user-srv/conf](https://github.com/WenyXu/Go-mini-kit/tree/master/user-srv/conf)


## User Model
```shell
├── model* 新增 model 目录
│   ├── boot.go
│   └── user
│       └── user.go
```

/model/boot.go
```go
package modelBoot

import "Go-mini-kit/user-srv/model/user"

func Init(){
	user.Init()
}
```

/model/user/user.go
```go
package user

import(
	"fmt"
	"sync"

	userProto "Go-mini-kit/user-srv/proto/user"
	"github.com/micro/go-micro/util/log"
	"Go-mini-kit/user-srv/basic/db"
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

//TODO: move these funcs into a new file
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
```

## Handler
handler/user.go
```go
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
```

## Main.go
/main.go
```go
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
	
    //Use etcd as Registry
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
```

## Run
```shell
go run main.go plugin.go
```

```shell
weny@weny-server:~/Projects/Golang/src/Go-mini-kit/user-srv$ go run main.go plugin.go
2019-12-23 11:15:03.769178 I | [Init] Loading config file, path: /home/weny/Projects/Golang/src/Go-mini-kit/user-srv/conf/application.yml, {Include:etcd, db}
2019-12-23 11:15:03.769199 I | [Init] Loading config file, path: /home/weny/Projects/Golang/src/Go-mini-kit/user-srv/conf/application-etcd.yml
2019-12-23 11:15:03.769209 I | [Init] Loading config file, path: /home/weny/Projects/Golang/src/Go-mini-kit/user-srv/conf/application-db.yml
2019-12-23 11:15:03.770365 I | Transport [http] Listening on [::]:39417
2019-12-23 11:15:03.770406 I | Broker [http] Connected to [::]:37017
2019-12-23 11:15:03.771701 I | Registry [etcd] Registering node: im.terminal.go.srv.user-fd1b7ed3-488d-4bfe-8c68-e86b93b80bc5
```

## Test gRPC
### 相关阅读
[document/micro-cli](https://micro.mu/docs/cli.html)
[document/zh-cn/micro-cli](https://micro.mu/docs/cn/cli.html)
**
**in this case {Server Name} : im.terminal.go.srv.user**
```shell
$ micro --registry=etcd call {Server Name} User.QueryUserByName '{"userName":"micro"}'
{
   "user": {
       "id": 1,
       "name": "micro",
   }
}
```

```shell
$ micro --registry=etcd call {Server Name} User.QueryUserByName '{"userName":"micro1"}'
{
        "error": {
                "code": 500,
                "detail": "sql: no rows in result set"
        }
}
```
# Source Code
[https://github.com/WenyXu/Go-mini-kit](https://github.com/WenyXu/Go-mini-kit)
