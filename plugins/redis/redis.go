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

	c := config.Instance()
	cfg := &redis{}
	err := c.Scan("redis", cfg)
	if err != nil {
		log.Logf("[initRedis] %s", err)
	}

	if !cfg.Enabled {
		log.Logf("[initRedis] Redis disabled")
		return
	}

	// Sentinel Mode
	if cfg.Sentinel != nil && cfg.Sentinel.Enabled {
		log.Log("[initRedis] Initializing Redis，Sentinel Mode...")
		initSentinel(cfg)
	} else { // Single Mode
		log.Log("[initRedis] Initializing Redis，Single Mode...")
		initSingle(cfg)
	}

	log.Log("[initRedis] Initializing Redis，Testing connection...")

	pong, err := _client.Ping().Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Logf("[initRedis] Initializing Redis，Ping ... %s", pong)
}


