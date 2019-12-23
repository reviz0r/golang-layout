package profile

import "github.com/volatiletech/sqlboiler/boil"

//go:generate sqlboiler --config ../../sqlboiler.toml psql

func init() {
	boil.DebugMode = false
}
