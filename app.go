package onion

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/HnH/di"
)

// Application interface
type Application interface {
	Listen() error
}

// New creates new application
func New(ctx context.Context, options ...Option) Application {
	var app = &application{
		context:    ctx,
		chanSignal: make(chan os.Signal, 1),
	}

	for _, opt := range options {
		opt(app)
	}

	return app
}

type application struct {
	context    context.Context
	container  di.Container
	providers  []di.Provider
	layers     layers
	chanSignal chan os.Signal
}

// SetContainer sets container for Application
func (self *application) SetContainer(container di.Container) {
	self.container = container
}

// SetProviders sets providers for di.Container
func (self *application) SetProviders(providers ...di.Provider) {
	self.providers = providers
}

// SetLayers sets layers for Application
func (self *application) SetLayers(layers ...Layer) {
	self.layers = layers
}

func (self *application) init() (err error) {
	if self.container == nil {
		self.container = di.NewContainer()
	}

	self.context = di.Ctx(self.context).SetContainer(self.container).Raw()

	for _, p := range self.providers {
		if err = p.Provide(self.container); err != nil {
			return
		}
	}

	//if err = self.container.Implementation(self.logger); err != nil {
	//	return
	//}

	return self.layers.initialize(self.context)
}

// Listen for system signals
func (self *application) Listen() (err error) {
	if err = self.init(); err != nil {
		return
	}

	slog.Info("application instance started")

	signal.Notify(self.chanSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var sig = <-self.chanSignal

	slog.Info("process termination signal received", slog.String("SIGNAL", sig.String()))

	if err = self.layers.shutdown(self.context); err != nil {
		slog.Error(err.Error())
	}

	signal.Stop(self.chanSignal)

	slog.Info("application instance terminated")

	return
}
