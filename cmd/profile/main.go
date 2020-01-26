package main

import (
	"go.uber.org/fx"

	"github.com/reviz0r/golang-layout/pkg/config"
	"github.com/reviz0r/golang-layout/pkg/db"
	"github.com/reviz0r/golang-layout/pkg/log"
	"github.com/reviz0r/golang-layout/pkg/server"
	"github.com/reviz0r/golang-layout/pkg/server/gateway"

	profile_internal "github.com/reviz0r/golang-layout/internal/profile"
	profile_pkg "github.com/reviz0r/golang-layout/pkg/profile"
)

func main() {
	app := fx.New(
		config.Module,
		log.Module,
		db.Module,
		server.Module,
		server.InterceptorsModule,
		server.GrpcLoggingPayloadModule,

		profile_pkg.GatewayModule,
		profile_pkg.GatewayInsecureDialModule,
		gateway.MuxModule,
		server.HTTPModule,

		// logic modules
		profile_internal.Module,
	)

	app.Run()
}
