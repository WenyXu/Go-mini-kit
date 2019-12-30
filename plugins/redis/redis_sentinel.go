package redis

import (
	r "github.com/go-redis/redis"
	"strings"
)
// redisSentinel redis sentinel config
type redisSentinel struct {
	Enabled bool   `json:"enabled"`
	Master  string `json:"master"`
	XNodes  string `json:"nodes"`
	nodes   []string
}


//GetNodes get nodes in Redis Sentinel mode
func (s *redisSentinel) GetNodes() []string {
	if len(s.XNodes) != 0 {
		for _, v := range strings.Split(s.XNodes, ",") {
			v = strings.TrimSpace(v)
			s.nodes = append(s.nodes, v)
		}
	}
	return s.nodes
}

func initSentinel(redisConfig *redis) {
	_client = r.NewFailoverClient(&r.FailoverOptions{
		MasterName:    redisConfig.Sentinel.Master,
		SentinelAddrs: redisConfig.Sentinel.GetNodes(),
		DB:            redisConfig.DBNum,
		Password:      redisConfig.Password,
	})

}