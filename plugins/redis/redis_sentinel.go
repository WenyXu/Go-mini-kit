package redis

import (
	r "github.com/go-redis/redis"
	"strings"
)
// redisSentinelConfig redisConfig sentinel config
type redisSentinelConfig struct {
	Enabled bool   `json:"enabled"`
	Master  string `json:"master"`
	XNodes  string `json:"nodes"`
	nodes   []string
}
//GetNodes get nodes in Redis Sentinel mode
func (s *redisSentinelConfig) GetNodes() []string {
	if len(s.XNodes) != 0 {
		for _, v := range strings.Split(s.XNodes, ",") {
			v = strings.TrimSpace(v)
			s.nodes = append(s.nodes, v)
		}
	}
	return s.nodes
}

func initSentinel(c *redisConfig) {
	_client = r.NewFailoverClient(&r.FailoverOptions{
		MasterName:    c.Sentinel.Master,
		SentinelAddrs: c.Sentinel.GetNodes(),
		DB:            c.DBNum,
		Password:      c.Password,
	})

}