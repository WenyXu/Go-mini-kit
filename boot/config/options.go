package config

import "github.com/micro/go-micro/config/source"

type Options struct {
	Apps    map[string]interface{}
	Sources []source.Source
}

type Option func(ops *Options)

//WithSource func
func WithSource(src source.Source) Option {
	return func(ops *Options) {
		ops.Sources = append(ops.Sources, src)
	}
}
