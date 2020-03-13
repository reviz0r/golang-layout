package db

import (
	"context"
	"database/sql"
	"fmt"

	// import postgresql driver
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module register database connection in DI container
var Module = fx.Provide(NewDatabase)

// NewDatabase gives new predefined database connection
func NewDatabase(lc fx.Lifecycle, config *viper.Viper) (*sql.DB, error) {
	dbconn, err := sql.Open("postgres", config.GetString("database.dsn"))
	if err != nil {
		return nil, fmt.Errorf("cannot open connection to database: %v", err)
	}

	if connMaxLifetime := config.GetDuration("database.conn_max_lifetime"); connMaxLifetime != 0 {
		dbconn.SetConnMaxLifetime(connMaxLifetime)
	}

	if maxIdleConns := config.GetInt("database.max_idle_conns"); maxIdleConns != 0 {
		dbconn.SetMaxIdleConns(maxIdleConns)
	}

	if maxOpenConns := config.GetInt("database.max_open_conns"); maxOpenConns != 0 {
		dbconn.SetMaxOpenConns(maxOpenConns)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if config.GetBool("database.ping_on_start") {
				err = dbconn.PingContext(ctx)
				if err != nil {
					return fmt.Errorf("database: cannot ping connection: %v", err)
				}
			}
			return nil
		},

		OnStop: func(ctx context.Context) error {
			err := dbconn.Close()
			if err != nil {
				return fmt.Errorf("database: cannot close connection: %v", err)
			}
			return nil
		},
	})

	return dbconn, nil
}
