package config

// EtcdConfig Interface
type IEtcdConfig interface {
	GetEnabled() bool
	GetPort() int
	GetHost() string
}

// etcdConfig struct
type etcdConfig struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}

func (config etcdConfig) GetPort() int {
	return config.Port
}

func (config etcdConfig) GetEnabled() bool {
	return config.Enabled
}

func (config etcdConfig) GetHost() string {
	return config.Host
}
