package main

import (
	"database/sql"
	"fmt"
)

const dbname = "golang-layout"

// NewDatabase gives new predefined database connection
func NewDatabase() *sql.DB {
	dbconnString := fmt.Sprintf("user=postgres password=postgres host=localhost port=5432 sslmode=disable database=%s", dbname)
	dbconn, err := sql.Open("postgres", dbconnString)
	if err != nil {
		panic(fmt.Errorf("cannot connect to db: %v", err))
	}

	return dbconn
}
