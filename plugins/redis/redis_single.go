package redis

import (
	r "github.com/go-redis/redis"
)

// redisConfig redisConfig config
type redisConfig struct {
	Enabled  bool                 `json:"enabled"`
	Conn     string               `json:"conn"`
	Password string               `json:"password"`
	DBNum    int                  `json:"dbNum"`
	Timeout  int                  `json:"timeout"`
	Sentinel *redisSentinelConfig `json:"sentinel"`
}

func initSingle(c *redisConfig) {
	_client = r.NewClient(&r.Options{
		Addr:     c.Conn,
		Password: c.Password, // no password set
		DB:       c.DBNum,    // use default DB
	})
}

// Redis get redisConfig
func Redis() *r.Client {
	return _client
}
