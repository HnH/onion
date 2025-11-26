package pgx

import (
	"fmt"
	"time"
)

type Config struct {
	Protocol           string        `config:"database-protocol,required" yaml:"protocol" validate:"required"`
	Driver             string        `config:"database-driver,required" yaml:"driver" validate:"required"`
	Addr               string        `config:"database-addr,required" yaml:"addr" validate:"required,hostname_port"`
	Username           string        `config:"database-user,required" yaml:"username" validate:"required"`
	Password           string        `config:"database-password,required" yaml:"password" validate:"required"`
	Database           string        `config:"database-name,required" yaml:"database" validate:"required"`
	ConnectionLimit    int           `config:"database-connection-limit" yaml:"connectionLimit" validate:"required,gt=0"`
	ConnectionTimeout  time.Duration `config:"database-connection-timeout" yaml:"connectionTimeout" validate:"required,gt=1s"`
	ConnectionLifetime time.Duration `config:"database-connection-lifetime" yaml:"connectionLifetime" validate:"required,gt=1s"`
}

func NewConfig(addr, u, p string) Config {
	return Config{
		Protocol:           "postgres",
		Driver:             "pgx",
		Addr:               addr,
		Username:           u,
		Password:           p,
		ConnectionLimit:    10,
		ConnectionTimeout:  3 * time.Second,
		ConnectionLifetime: 5 * time.Minute,
	}
}

func (cfg Config) String() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s/%s?sslmode=disable",
		cfg.Protocol, cfg.Username, cfg.Password, cfg.Addr, cfg.Database,
	)
}
