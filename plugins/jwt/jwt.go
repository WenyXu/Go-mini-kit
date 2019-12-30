package jwt

import "go-mini-kit.com/boot"

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
