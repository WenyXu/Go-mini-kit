package access

import (
	"fmt"
	"sync"

	redis "github.com/go-redis/redis"
)

var (
	s  *service
	ca *redis.Client
	m  sync.RWMutex
)

// service struct
type service struct {
}

// IService Interface
type IService interface {
	// CreateAccessToken 
	CreateUserAccessToken(subject *Subject) (ret string, err error)

	// GetCachedAccessToken 
	GetCachedAccessToken(subject *Subject) (ret string, err error)

	// DelUserAccessToken 
	DeleteUserAccessToken(token string) (err error)
}

// GetService func
func GetService() (IService, error) {
	if s == nil {
		return nil, fmt.Errorf("[GetService] GetService 未初始化")
	}
	return s, nil
}

// Init func
func Init() {
	m.Lock()
	defer m.Unlock()

	if s != nil {
		return
	}

	ca = redis.GetRedis()

	s = &service{}
}
