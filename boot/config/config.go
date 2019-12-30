package config

import (
	"fmt"
	"sync"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/util/log"
)

var (
	m            sync.RWMutex
	_initialized bool
	c = &configurator{}
)

// IConfigurator interface
type IConfigurator interface {
	Scan(name string, config interface{}) (err error)
}

// configurator configurator
type configurator struct {
	conf config.Config
}

func (c *configurator) init(ops Options) (err error) {
	m.Lock()
	defer m.Unlock()

	if _initialized {
		log.Logf("[init] initialized")
		return
	}

	c.conf = config.NewConfig()
	err = c.conf.Load(ops.Sources...)
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

// Instance get Instance
func Instance() IConfigurator {
	return c
}

// Init initialize
func Init(opts ...Option) {

	ops := Options{}
	for _, o := range opts {
		o(&ops)
	}

	c = &configurator{}

	_ = c.init(ops)
}
