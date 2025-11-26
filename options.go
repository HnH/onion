package onion

import (
	"github.com/HnH/di"
)

// Option represents single option type
type Option func(Options)

// Options represents a target for applying an Option
type Options interface {
	SetContainer(di.Container)
	SetProviders(...di.Provider)
	SetLayers(...Layer)
}

// WithContainer creates a container Option
func WithContainer(container di.Container) Option {
	return func(o Options) {
		o.SetContainer(container)
	}
}

// WithProviders creates a providers Option
func WithProviders(providers ...di.Provider) Option {
	return func(o Options) {
		o.SetProviders(providers...)
	}
}

// WithLayers creates a layers Option
func WithLayers(layers ...Layer) Option {
	return func(o Options) {
		o.SetLayers(layers...)
	}
}
