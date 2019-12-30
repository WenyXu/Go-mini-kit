package redis

import (
	r "github.com/go-redis/redis"
)


// redis redis config
type redis struct {
	Enabled  bool           `json:"enabled"`
	Conn     string         `json:"conn"`
	Password string         `json:"password"`
	DBNum    int            `json:"dbNum"`
	Timeout  int            `json:"timeout"`
	Sentinel *redisSentinel `json:"sentinel"`
}

func initSingle(redisConfig *redis) {
	_client = r.NewClient(&r.Options{
		Addr:     redisConfig.Conn,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.DBNum,    // use default DB
	})
}

// Redis get redis
func Redis() *r.Client {
	return _client
}
