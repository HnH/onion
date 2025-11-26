package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type logger struct {
}

func NewLogger() tracelog.Logger {
	return &logger{}
}

func (pl *logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	var zlevel zerolog.Level
	switch level {
	case tracelog.LogLevelNone:
		zlevel = zerolog.NoLevel
	case tracelog.LogLevelError:
		zlevel = zerolog.ErrorLevel
	case tracelog.LogLevelWarn:
		zlevel = zerolog.WarnLevel
	case tracelog.LogLevelInfo:
		zlevel = zerolog.InfoLevel
	case tracelog.LogLevelDebug:
		zlevel = zerolog.DebugLevel
	default:
		zlevel = zerolog.DebugLevel
	}

	var pgxlog = *zerolog.Ctx(ctx)
	if v, ok := data["time"]; ok {
		// log as string, else zerolog just writes it as int
		if duration, okk := v.(time.Duration); okk {
			data["exec_time"] = duration.String()
		} else {
			data["exec_time"] = v
		}

		delete(data, "time")
	}

	if v, ok := data["args"]; ok {
		if slice, okk := v.([]any); okk && len(slice) > 0 {
			// databaseSQLResultFormats github.com/jackc/pgx/v5@v5.5.3/stdlib/sql.go:90
			if _, okk = slice[0].(pgx.QueryResultFormatsByOID); okk {
				slice = slice[1:]
			}

			for i, arg := range slice {
				slice[i] = fmt.Sprint(arg)
			}
		}
	}

	if len(data) > 0 {
		pgxlog = pgxlog.With().Fields(data).Logger()
	}

	pgxlog.WithLevel(zlevel).Msg(msg)
}
