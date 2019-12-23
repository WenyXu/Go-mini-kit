package config

// MysqlConfig Interface
type IMysqlConfig interface {
	GetURL() string
	GetEnabled() bool
	GetMaxIdleConnection() int
	GetMaxOpenConnection() int
}

// mysqlConfig Struct
type mysqlConfig struct {
	URL               string `json:"url"`
	Enable            bool   `json:"enabled"`
	MaxIdleConnection int    `json:"maxIdleConnection"`
	MaxOpenConnection int    `json:"maxOpenConnection"`
}

// Get connect url
func (config mysqlConfig) GetURL() string {
	return config.URL
}


func (config mysqlConfig) GetEnabled() bool {
	return config.Enable
}

// Get MaxIdleConnection
func (config mysqlConfig) GetMaxIdleConnection() int {
	return config.MaxIdleConnection
}

// Get MaxOpenConnection
func (config mysqlConfig) GetMaxOpenConnection() int {
	return config.MaxOpenConnection
}
