package mockdb

import (
	"context"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/fx"
)

// Module register mock database connection in DI container
var Module = fx.Provide(NewDatabase)

// NewDatabase gives new mocked database connection
func NewDatabase(lc fx.Lifecycle) (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db, mock, nil
}
