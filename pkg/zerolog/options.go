package zerolog

// Option represents single option type
type Option func(Options)

// Options represents a target for applying an Option
type Options interface {
	SetReplaceGlobalLogger()
	SetLevel(string)
}

// WithLevel creates a level Option
func WithLevel(level string) Option {
	return func(o Options) {
		o.SetLevel(level)
	}
}

// WithReplaceGlobalLogger creates a replace global logger Option
func WithReplaceGlobalLogger() Option {
	return func(o Options) {
		o.SetReplaceGlobalLogger()
	}
}
