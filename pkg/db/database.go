package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// import postgresql driver
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

// Module register database connection in DI container
var Module = fx.Provide(NewDatabase)

// DatabaseParams .
type DatabaseParams struct {
	fx.In

	DSN             string        `name:"database_dsn"`
	ConnMaxLifetime time.Duration `name:"database_conn_max_lifetime" optional:"true"`
	MaxIdleConns    int           `name:"database_max_idle_conns" optional:"true"`
	MaxOpenConns    int           `name:"database_max_open_conns" optional:"true"`

	PingOnStart bool `name:"database_ping_on_start" optional:"true"`
}

// NewDatabase gives new predefined database connection
func NewDatabase(lc fx.Lifecycle, p DatabaseParams) (*sql.DB, error) {
	dbconn, err := sql.Open("postgres", p.DSN)
	if err != nil {
		return nil, fmt.Errorf("cannot open connection to database: %v", err)
	}

	if p.ConnMaxLifetime != 0 {
		dbconn.SetConnMaxLifetime(p.ConnMaxLifetime)
	}

	if p.MaxIdleConns != 0 {
		dbconn.SetMaxIdleConns(p.MaxIdleConns)
	}

	if p.MaxOpenConns != 0 {
		dbconn.SetMaxOpenConns(p.MaxOpenConns)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if p.PingOnStart {
				err = dbconn.PingContext(ctx)
				if err != nil {
					return fmt.Errorf("cannot ping database connection: %v", err)
				}
			}
			return nil
		},

		OnStop: func(ctx context.Context) error {
			err := dbconn.Close()
			if err != nil {
				return fmt.Errorf("cannot close database connection: %v", err)
			}
			return nil
		},
	})

	return dbconn, nil
}
