package boot

import (
	"Go-mini-kit/boot/config"
)

var (
	pluginFuncs []func()
)

type Options struct {
	EnableDB    bool
	EnableRedis bool
	cfgOps      []config.Option
}

type Option func(o *Options)

func Init(opts ...config.Option) {
	// Initializing config
	config.Init(opts...)

	// Initializing plugin's init
	for _, f := range pluginFuncs {
		f()
	}
}

// Register func
func Register(f func()) {
	pluginFuncs = append(pluginFuncs, f)
}
