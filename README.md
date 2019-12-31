# 使用 Go-micro 构建微服务-2 gRPC 配置中心服务

# 前言
[使用 Go-micro 构建微服务-1 Quick Start](https://www.yuque.com/dd6hnm/mhce0q/vas58o)
## 概览
本章将演示实现一个 gRPC 配置中心服务端及客户端调用
## 包含内容
gRPC 配置中心示例
## 不包含内容
负载均衡

# 准备工作
# Redis (optional)
docker-compose.yaml
```yaml
version: '3'
services:
    redis:
      image: redis
      container_name: redis
      command: redis-server --requirepass 123456
      ports:
        - "6379:6379"
      volumes:
        - ./data:/data
```

# Go 1.13
> The only change you have to make is add a dot to the first element of the module path. If you don't intend to make the module available on the internet, you can use a reserved TLD like `.localhost`that signifies that. For example, the following module path would work:
> ```
my-api-server.localhost/my-utils/uuid
```

ref:[https://github.com/golang/go/issues/35020](https://github.com/golang/go/issues/35020)

在 go1.13 中， go module 名称规范要求路径的第一部分必须满足域名规范，否则会报`malformed module path "XXXXXX": missing dot in first path element` 

## 我们的解决方案
STEP 1: 重命名目录
`mv Go-mini-kit go-mini-kit.com`
STEP 2: go.mod 中使用 replace
```go
replace (
	go-mini-kit.com/boot => ../boot
	go-mini-kit.com/plugins => ../plugins
	go-mini-kit.com/user-srv => ../user-srv
)
```
 
# 源码
[https://github.com/WenyXu/Go-mini-kit/tree/Part2](https://github.com/WenyXu/Go-mini-kit/tree/Part2)
业务代码非常简单，在此不再赘述。
# config-grpc-srv
我们使用 micro/go-micro 生态内 micro/go-plugins 下 config/source/grpc 来实现一个 gRPC 服务
micro/go-plugins/config/source/grpc
目录结构
```shell
.
├── README.md
├── grpc.go
├── options.go
├── proto
│   ├── grpc.pb.go
│   └── grpc.proto
├── util.go
└── watcher.go
```

## grpc.proto
[github.com/micro/go-plugins/config/source/grpc/proto/grpc.proto](https://github.com/micro/go-plugins/blob/master/config/source/grpc/proto/grpc.proto)
```protobuf
syntax = "proto3";

service Source {
	rpc Read(ReadRequest) returns (ReadResponse) {};
	rpc Watch(WatchRequest) returns (stream WatchResponse) {};
}

message ChangeSet {
	bytes data = 1;
	string checksum = 2;
	string format = 3;
	string source = 4;
	int64 timestamp = 5;
}

message ReadRequest {
	string path = 1;
}

message ReadResponse {
	ChangeSet change_set = 1;	
}

message WatchRequest {
	string path = 1;
}

message WatchResponse {
	ChangeSet change_set = 1;
}
```

## 服务端 main.go 实现
回到我们的目录实现这个 gRPC 服务的 Read & Watch func
config-grpc-srv 目录结构
```shell
.
├── conf
│   └── micro.yml
└── main.go
```
代码实现
[config-grpc-srv/main.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/config-grpc-srv/main.go)
```go
package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/util/log"
	proto "github.com/micro/go-plugins/config/source/grpc/proto"
	"google.golang.org/grpc"
)

var (
	mux        sync.RWMutex
	configMaps = make(map[string]*proto.ChangeSet)
	// config file list
	apps = []string{"micro"}
)
// Service struct
type Service struct{}

func main() {
	// recover
	defer func() {
		if r := recover(); r != nil {
			log.Logf("[main] Recovered in f %v", r)
		}
	}()

	// loading and watch config file
	err := loadAndWatchConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	// create a new gRPC server
	service := grpc.NewServer()
	proto.RegisterSourceServer(service, new(Service))
	ts, err := net.Listen("tcp", ":9600")
	if err != nil {
		log.Fatal(err)
	}
	log.Logf("configServer started")

	// start server
	err = service.Serve(ts)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Service) Read(ctx context.Context, req *proto.ReadRequest) (rsp *proto.ReadResponse, err error) {
	appName := parsePath(req.Path)

	rsp = &proto.ReadResponse{
		ChangeSet: getConfig(appName),
	}
	return
}
// Watch func
func (s Service) Watch(req *proto.WatchRequest, server proto.Source_WatchServer) (err error) {
	appName := parsePath(req.Path)
	rsp := &proto.WatchResponse{
		ChangeSet: getConfig(appName),
	}
	if err = server.Send(rsp); err != nil {
		log.Logf("[Watch] watching failed，%s", err)
		return err
	}

	return
}

func loadAndWatchConfigFile() (err error) {
	// loading files in dir
	for _, app := range apps {
		if err := config.Load(file.NewSource(
			file.WithPath("./conf/" + app + ".yml"),
		)); err != nil {
			log.Fatalf("[loadAndWatchConfigFile] loading files in dir failed，%s", err)
			return err
		}
	}

	// watching the modification of files
	watcher, err := config.Watch()
	if err != nil {
		log.Fatalf("[loadAndWatchConfigFile] watching the modification of files is failed，%s", err)
		return err
	}

	go func() {
		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatalf("[loadAndWatchConfigFile] watching the modification of files is failed， %s", err)
				return
			}

			log.Logf("[loadAndWatchConfigFile] files modified，%s", string(v.Bytes()))
		}
	}()

	return
}

func getConfig(appName string) *proto.ChangeSet {
	bytes := config.Get(appName).Bytes()

	log.Logf("[getConfig] appName：%s", appName)
	return &proto.ChangeSet{
		Data:      bytes,
		Checksum:  fmt.Sprintf("%x", md5.Sum(bytes)),
		Format:    "yml",
		Source:    "file",
		Timestamp: time.Now().Unix(),
	}
}

func parsePath(path string) (appName string) {
	paths := strings.Split(path, "/")

	if paths[0] == "" && len(paths) > 1 {
		return paths[1]
	}

	return paths[0]
}

```
# user-srv
在 user-srv 调用 gRPC Client
```shell
.
├── Dockerfile
├── Makefile
├── README.md
├── generate.go
├── go.mod
├── go.sum
├── handler
│   └── user.go
├── main.go
├── model
│   ├── model.go
│   └── user
│       └── user.go
├── plugin.go
└── proto
    └── user
        ├── user.pb.go
        ├── user.pb.micro.go
        └── user.proto
```
## user-srv/main.go
[user-srv/main.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/user-srv/main.go)
```go
import (
	"github.com/micro/go-plugins/config/source/grpc"
)

func main(){
    ...
    initConfig()
    ...
}
func initConfig(){

    source := grpc.NewSource(
		grpc.WithAddress("127.0.0.1:9600"),
		grpc.WithPath("micro"),
	)
    boot.Init(config.WithSource(source))
	err := config.GetInstance().Scan(_appName, _config)
    ...
}

```
## grpc.go
[github.com/micro/go-plugins/config/source/grpc/grpc.go](https://github.com/micro/go-plugins/blob/master/config/source/grpc/grpc.go)
```go
package grpc

import (
	"context"
	"crypto/tls"

	"github.com/micro/go-micro/config/source"
	proto "github.com/micro/go-plugins/config/source/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type grpcSource struct {
	addr      string
	path      string
	opts      source.Options
	...
    client    *grpc.ClientConn
}

var (
	DefaultPath    = "/micro/config"
	DefaultAddress = "localhost:8080"
)

func (g *grpcSource) Read() (set *source.ChangeSet, err error) {

	var opts []grpc.DialOption

	...

	g.client, err = grpc.Dial(g.addr, opts...)
	if err != nil {
		return nil, err
	}
	cl := proto.NewSourceClient(g.client)
	rsp, err := cl.Read(context.Background(), &proto.ReadRequest{
		Path: g.path,
	})
	if err != nil {
		return nil, err
	}
	return toChangeSet(rsp.ChangeSet), nil
}

func (g *grpcSource) Watch() (source.Watcher, error) {
	cl := proto.NewSourceClient(g.client)
	rsp, err := cl.Watch(context.Background(), &proto.WatchRequest{
		Path: g.path,
	})
	if err != nil {
		return nil, err
	}
	return newWatcher(rsp)
}

func (g *grpcSource) String() string {
	return "grpc"
}

func NewSource(opts ...source.Option) source.Source {
	var options source.Options
	for _, o := range opts {
		o(&options)
	}

	addr := DefaultAddress
	path := DefaultPath

	if options.Context != nil {
		a, ok := options.Context.Value(addressKey{}).(string)
		if ok {
			addr = a
		}
		p, ok := options.Context.Value(pathKey{}).(string)
		if ok {
			path = p
		}
	}

	return &grpcSource{
		addr: addr,
		path: path,
		opts: options,
	}
}

```


user-srv/main.go 中调用 NewSource 返回一个 grpcSource struct (实现了 source.Source interface )

[grpc/grpc.go](https://github.com/micro/go-plugins/blob/master/config/source/grpc/grpc.go) 中 grpcSource Read & Watch 调用 gRPC client 发起请求

[micro/go-micro/config/source](https://github.com/micro/go-micro/blob/master/config/source/source.go)
```go
package source

...

// Source is the source from which config is loaded
type Source interface {
	Read() (*ChangeSet, error)
	Write(*ChangeSet) error
	Watch() (Watcher, error)
	String() string
}
...
```

再来关注到 [user-srv/main.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/user-srv/main.go) （文中）Line 16 boot.Init(config.WithSource(source)) 

[boot/config/options.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/boot/config/options.go)
```go
package config

import "github.com/micro/go-micro/config/source"

type Options struct {
	Apps    map[string]interface{}
	Sources []source.Source
}

type Option func(ops *Options)

//WithSource func
func WithSource(src source.Source) Option {
	return func(ops *Options) {
		ops.Sources = append(ops.Sources, src)
	}
}
```
[boot/boot.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/boot/boot.go) Init func
```go
func Init(options ...config.Option) {
	// Initializing config
	config.Init(options...)

	// Initializing plugin's init
	for _, f := range pluginFuncs {
		f()
	}
}
```
[boot/config.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/boot/config/config.go) Init funct
```go
func Init(optionList ...Option) {

	options := Options{}
	for _, option := range optionList {
		option(&options)
	}

	_configurator = &configurator{}

	_ = _configurator.init(options)
}
```
[boot/config.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/boot/config/config.go) func (c *configurator) init(options Options) (err error)
```go
func (c *configurator) init(options Options) (err error) {
	m.Lock()
	defer m.Unlock()

	if _initialized {
		log.Logf("[init] initialized")
		return
	}

	c.conf = config.NewConfig()
	err = c.conf.Load(options.Sources...)
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		log.Logf("[init] start to watching modification of files ...")

		// start to watching
		watcher, err := c.conf.Watch()
		if err != nil {
			log.Fatal(err)
		}

		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatal(err)
			}

			log.Logf("[init] modification of files : %v", string(v.Bytes()))
		}
	}()

	_initialized = true
	return
}
```

这里我们再来看看到 [boot/config.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/boot/config/config.go) func (c *configurator) init(options Options) (err error) （上文）Line 10 config.NewConfig()
NewConfig() 返回了 Config interface

[go-micro/config/config/config.go](https://github.com/micro/go-micro/tree/master/config/config.go) func NewConfig(opts ...Option) Config
```go
package config

// NewConfig returns new config
func NewConfig(opts ...Option) Config {
	return newConfig(opts...)
}
```
[go-micro/config/config/config.go](https://github.com/micro/go-micro/tree/master/config/config.go) type Config interface
```go
package config

...

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// provide the reader.Values interface
	reader.Values
	// Stop the config loader/watcher
	Close() error
	// Load config sourc￿es
	Load(source ...source.Source) error
	// Force a source changeset sync
	Sync() error
	// Watch a value for changes
	Watch(path ...string) (Watcher, error)
}
```
[go-micro/config/config/default.go](https://github.com/micro/go-micro/blob/master/config/default.go) func newConfig(opts ...Option) Config
```go
package config

func newConfig(opts ...Option) Config {
	options := Options{
		Loader: memory.NewLoader(),
		Reader: json.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	options.Loader.Load(options.Source...)
	snap, _ := options.Loader.Snapshot()
	vals, _ := options.Reader.Values(snap.ChangeSet)

	c := &config{
		exit: make(chan bool),
		opts: options,
		snap: snap,
		vals: vals,
	}

	go c.run()

	return c
}
```
[go-micro/config/config/default.go](https://github.com/micro/go-micro/blob/master/config/default.go) type config struct
```go
package config

type config struct {
	exit chan bool
	opts Options

	sync.RWMutex
	// the current snapshot
	snap *loader.Snapshot
	// the current values
	vals reader.Values
}
```
[go-micro/config/config/default.go](https://github.com/micro/go-micro/blob/master/config/default.go) func (c *config) Load(sources ...source.Source) error 
```go
package config

func (c *config) Load(sources ...source.Source) error {
	if err := c.opts.Loader.Load(sources...); err != nil {
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	vals, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.vals = vals

	return nil
}

```

这里 c.opts.Loader 即 memory.NewLoader() 返回 loader.Loader interface

[go-micro/config/loader/loader.go](https://github.com/micro/go-micro/blob/master/config/loader/loader.go) type Loader interface 
```go
package load

type Loader interface {
	// Stop the loader
	Close() error
	// Load the sources
	Load(...source.Source) error
	// A Snapshot of loaded config
	Snapshot() (*Snapshot, error)
	// Force sync of sources
	Sync() error
	// Watch for changes
	Watch(...string) (Watcher, error)
	// Name of loader
	String() string
}
```

[go-micro/config/loader/memory ](https://github.com/micro/go-micro/blob/master/config/loader/memory/memory.go)func (m *memory) Load(sources ...source.Source) error
```go
package memory

func (m *memory) Load(sources ...source.Source) error {
	var gerrors []string

	for _, source := range sources {
		set, err := source.Read()
		if err != nil {
			gerrors = append(gerrors,
				fmt.Sprintf("error loading source %s: %v",
					source,
					err))
			// continue processing
			continue
		}
		m.Lock()
		m.sources = append(m.sources, source)
		m.sets = append(m.sets, set)
		idx := len(m.sets) - 1
		m.Unlock()
		go m.watch(idx, source)
	}

	if err := m.reload(); err != nil {
		gerrors = append(gerrors, err.Error())
	}

	// Return errors
	if len(gerrors) != 0 {
		return errors.New(strings.Join(gerrors, "\n"))
	}
	return nil
}

```

Line 7 `source.Read()`还记得我们最开始 grpcSource struct (实现了 source.Source interface)吗？我们在 user-srv/main.go 创建 grpcSource (实现了 source.Source interface) 被传入到 config.Load(source ...source.Source) 中，然后在 Load() 中调用实现 source.Source  interface 对象的 Read() ，这个 Read() 便是正是 [grpc/grpc.go](https://github.com/micro/go-plugins/blob/master/config/source/grpc/grpc.go) 中实现的。
# boot
## boot/config.go
至此所有流程都梳理完了
[boot/config.go](https://github.com/WenyXu/Go-mini-kit/blob/Part2/boot/config/config.go) 
```go
package config

import (
	"fmt"
	"sync"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/util/log"
)

var (
	m             sync.RWMutex
	_initialized  bool
	_configurator = &configurator{}
)

// IConfigurator interface
type IConfigurator interface {
	Scan(name string, config interface{}) (err error)
}

// configurator configurator
type configurator struct {
	conf config.Config
}

func (c *configurator) init(options Options) (err error) {
	m.Lock()
	defer m.Unlock()

	if _initialized {
		log.Logf("[init] initialized")
		return
	}

	c.conf = config.NewConfig()
	err = c.conf.Load(options.Sources...)
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		log.Logf("[init] start to watching modification of files ...")

		// start to watching
		watcher, err := c.conf.Watch()
		if err != nil {
			log.Fatal(err)
		}

		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatal(err)
			}

			log.Logf("[init] modification of files : %v", string(v.Bytes()))
		}
	}()

	_initialized = true
	return
}
// Scan func
// the config interface will get value
func (c *configurator) Scan(name string, config interface{}) (err error) {

	v := c.conf.Get(name)
	if v != nil {
		err = v.Scan(config)
	} else {
		err = fmt.Errorf("[Scan] config is not exist ，err：%s", name)
	}
	return
}

// GetInstance get GetInstance
func GetInstance() IConfigurator {
	return _configurator
}

// Init initialize
func Init(optionList ...Option) {

	options := Options{}
	for _, option := range optionList {
		option(&options)
	}

	_configurator = &configurator{}

	_ = _configurator.init(options)
}

```

# 测试
## Run
```shell
cd $GOPATH/src/go-mini-kit.com
cd config-grpc-srv
go run main.go
cd ..
cd user-srv
go run main.go plugin.go
```
## gRPC
```shell
weny@weny-server:~$ micro --registry=etcd call im.terminal.go.srv.user User.QueryUserByName '{"userName":"micro"}'
{
	"success": true,
	"user": {
		"id": 1,
		"name": "micro"
	}
}
```

