package zerolog

import (
	"log/slog"
	"strings"

	"github.com/HnH/di"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

var New = zerolog.New
var NewConsoleWriter = zerolog.NewConsoleWriter

func NewProvider(zlog zerolog.Logger, options ...Option) di.Provider {
	var p = &Provider{
		log: zlog,
		lvl: "info",
	}

	for _, opt := range options {
		opt(p)
	}

	return p
}

type Provider struct {
	log zerolog.Logger
	lvl string
	def bool
}

func (p *Provider) getSlogLevel() slog.Level {
	switch strings.ToLower(p.lvl) {
	case "debug":
		return slog.LevelDebug
	case "trace":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

func (p *Provider) getZerologLevel() zerolog.Level {
	switch strings.ToLower(p.lvl) {
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

func (p *Provider) SetReplaceGlobalLogger() {
	p.def = true
}

func (p *Provider) SetLevel(lvl string) {
	p.lvl = lvl
}

func (p *Provider) Provide(container di.Container) (err error) {
	var (
		zlogger = p.log.Level(p.getZerologLevel())
		slogger = slog.New(slogzerolog.Option{
			Level:  p.getSlogLevel(),
			Logger: &zlogger,
		}.NewZerologHandler())
	)

	if err = container.Implementation(slogger); err != nil {
		return
	}

	if err = container.Implementation(zlogger); err != nil {
		return
	}

	if p.def {
		slog.SetDefault(slogger)
		zerolog.DefaultContextLogger = &zlogger
	}

	return nil
}
