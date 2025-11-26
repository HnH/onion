package pgx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/HnH/di"
	"github.com/XSAM/otelsql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

func NewConnection(ctx context.Context, cfg Config) (database *sqlx.DB, err error) {
	var pgxCfg *pgx.ConnConfig
	if pgxCfg, err = pgx.ParseConfig(cfg.String()); err != nil {
		return
	}

	pgxCfg.DefaultQueryExecMode = pgx.QueryExecModeExec
	pgxCfg.Tracer = &tracelog.TraceLog{
		Logger:   NewLogger(),
		LogLevel: tracelog.LogLevelInfo,
	}

	// OTEL instrumentation
	var otelOpts = []otelsql.Option{
		otelsql.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
			semconv.DBNamespace(pgxCfg.Database),
			attribute.String("db.user", pgxCfg.User), // for some reason deprecated in latest otel and no replacement until now
		),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			Ping:                 false,
			RowsNext:             false,
			DisableErrSkip:       false,
			DisableQuery:         false,
			RecordError:          nil,
			OmitConnResetSession: false,
			OmitConnPrepare:      false,
			OmitConnQuery:        false,
			OmitRows:             true,
			OmitConnectorConnect: false,
			SpanFilter:           nil,
		}),
	}

	var meter metric.MeterProvider
	if err = di.Ctx(ctx).Resolver().Resolve(&meter); err == nil {
		otelOpts = append(otelOpts, otelsql.WithMeterProvider(meter))
	}

	var tracer trace.TracerProvider
	if err = di.Ctx(ctx).Resolver().Resolve(&tracer); err == nil {
		otelOpts = append(otelOpts, otelsql.WithTracerProvider(tracer))
	}

	// Open connection
	var conn *sql.DB
	if conn, err = otelsql.Open(cfg.Driver, stdlib.RegisterConnConfig(pgxCfg), otelOpts...); err != nil {
		return nil, err
	}

	database = sqlx.NewDb(conn, cfg.Driver)
	database.SetMaxOpenConns(cfg.ConnectionLimit)
	database.SetMaxIdleConns(cfg.ConnectionLimit / 2)
	database.SetConnMaxLifetime(cfg.ConnectionLifetime)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, cfg.ConnectionTimeout)

	defer cancel()

	if meter != nil {
		if err = otelsql.RegisterDBStatsMetrics(conn, otelOpts...); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	var health gosundheit.Health
	if err = di.Ctx(ctx).Resolver().Resolve(&health); err == nil {
		var check gosundheit.Check
		if check, err = checks.NewPingCheck(fmt.Sprintf("db.check.%s", cfg.Database), database); err != nil {
			return nil, err
		}

		err = health.RegisterCheck(check, gosundheit.ExecutionPeriod(time.Minute/2), gosundheit.ExecutionTimeout(time.Second))
	}

	return
}
