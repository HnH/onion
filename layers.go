package onion

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/HnH/di"
)

// Layer is an Application layer abstraction
type Layer interface {
	Initialize(context.Context) error
}

// Shutdowner is a layer that has his state shutdown logic
type Shutdowner interface {
	Shutdown() <-chan error
}

type layers []Layer

func (lrs layers) initialize(ctx context.Context) (err error) {
	slog.Info("starting to initialize layers", slog.Int("layers", len(lrs)))

	for _, l := range lrs {
		if err = l.Initialize(ctx); err != nil {
			return
		}

		if t, ok := l.(di.Provider); ok {
			if err = t.Provide(di.Ctx(ctx).Container()); err != nil {
				return
			}
		}

		slog.Info("layer initialized", slog.String("layer", reflect.TypeOf(l).String()))
	}

	slog.Info("all layers successfully initialized", slog.Int("layers", len(lrs)))

	return
}

func (lrs layers) shutdown(ctx context.Context) (err error) {
	slog.Info("starting to gracefully shut down layers", slog.Int("layers", len(lrs)))

	// Loop and close layers in reversed order
	for n := len(lrs) - 1; n >= 0; n-- {
		if t, ok := lrs[n].(Shutdowner); ok {
			for err = range t.Shutdown() {
				if err != nil {
					slog.Error(err.Error())
				}
			}
		}

		slog.Info("layer gracefully shut down", slog.String("layer", reflect.TypeOf(lrs[n]).String()))
		lrs = lrs[:n] // pop
	}

	slog.Info("all layers gracefully shut down")

	return
}
