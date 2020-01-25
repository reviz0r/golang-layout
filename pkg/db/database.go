package db

import (
	"context"
	"database/sql"
	"fmt"

	// import postgresql driver
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

// Module register database connection in DI container
var Module = fx.Provide(NewDatabase)

// DatabaseParams .
type DatabaseParams struct {
	fx.In

	DatabaseDSN         string `name:"database_dsn"`
	DatabasePingOnStart bool   `name:"database_ping_on_start"`
}

// NewDatabase gives new predefined database connection
func NewDatabase(lc fx.Lifecycle, p DatabaseParams) (*sql.DB, error) {
	dbconn, err := sql.Open("postgres", p.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("cannot open connection to database: %v", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if p.DatabasePingOnStart {
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
