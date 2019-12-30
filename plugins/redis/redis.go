package redis

import (
	"sync"

	r "github.com/go-redis/redis"
	"github.com/micro/go-micro/util/log"
	"go-mini-kit.com/boot"
	"go-mini-kit.com/boot/config"
)

var (
	_client      *r.Client
	m            sync.RWMutex
	_config		 =&redisConfig{}
	_initialized bool
)
// init Initialize Redis
func init() {
	boot.Register(initRedis)
}

func initRedis() {
	m.Lock()
	defer m.Unlock()

	if _initialized {
		log.Log("[initRedis] Redis initialized...")
		return
	}

	log.Log("[initRedis] Initializing Redis...")

	c := config.GetInstance()
	err := c.Scan("redisConfig", _config)
	if err != nil {
		log.Logf("[initRedis] %s", err)
	}

	if !_config.Enabled {
		log.Logf("[initRedis] Redis disabled")
		return
	}

	// Sentinel Mode
	if _config.Sentinel != nil && _config.Sentinel.Enabled {
		log.Log("[initRedis] Initializing Redis，Sentinel Mode...")
		initSentinel(_config)
	} else { // Single Mode
		log.Log("[initRedis] Initializing Redis，Single Mode...")
		initSingle(_config)
	}

	log.Log("[initRedis] Initializing Redis，Testing connection...")

	pong, err := _client.Ping().Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Logf("[initRedis] Initializing Redis，Ping ... %s", pong)
}


