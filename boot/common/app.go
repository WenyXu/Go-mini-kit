package common

import "strconv"

// AppConfig common config
type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (c *AppConfig) Addr() string {
	return c.Address + ":" + strconv.Itoa(c.Port)
}
