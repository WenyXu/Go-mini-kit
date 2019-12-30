package jwt

import "Go-mini-kit/boot"

// jwt config
type jwt struct {
	SecretKey string `json:"secretKey"`
}

// init Initialize
func init() {
	boot.Register(initJwt)
}

func initJwt() {

}
