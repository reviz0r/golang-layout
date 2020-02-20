package main

import (
	"go.uber.org/fx"

	"github.com/reviz0r/golang-layout/pkg/config"
	"github.com/reviz0r/golang-layout/pkg/db"
	"github.com/reviz0r/golang-layout/pkg/log"
	"github.com/reviz0r/golang-layout/pkg/server"
	"github.com/reviz0r/golang-layout/pkg/server/gateway"

	profileInternal "github.com/reviz0r/golang-layout/internal/profile"
	profilePkg "github.com/reviz0r/golang-layout/pkg/profile"
)

func main() {
	app := fx.New(
		fx.NopLogger,

		config.Module,
		log.Module,
		db.Module,
		server.Module,
		server.InterceptorsModule,
		server.GrpcLoggingPayloadModule,

		profilePkg.GatewayModule,
		profilePkg.GatewayInsecureDialModule,
		gateway.MuxModule,
		server.HTTPModule,

		// logic modules
		profileInternal.Module,
		profilePkg.SwaggerModule,
	)

	app.Run()
}
