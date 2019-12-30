package boot

import (
	"go-mini-kit.com/boot/config"
)

var (
	pluginFuncs []func()
)

func Init(options ...config.Option) {
	// Initializing config
	config.Init(options...)

	// Initializing plugin's init
	for _, f := range pluginFuncs {
		f()
	}
}

// Register func
func Register(f func()) {
	pluginFuncs = append(pluginFuncs, f)
}
