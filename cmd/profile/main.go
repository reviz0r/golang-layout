package main

import (
	"go.uber.org/fx"

	"github.com/reviz0r/golang-layout/pkg/config"
	"github.com/reviz0r/golang-layout/pkg/db"
	"github.com/reviz0r/golang-layout/pkg/log"
	"github.com/reviz0r/golang-layout/pkg/server"

	profileInternal "github.com/reviz0r/golang-layout/internal/profile"
	profilePkg "github.com/reviz0r/golang-layout/pkg/profile"
)

func main() {
	app := fx.New(
		fx.NopLogger,

		config.Module,
		log.Module,
		db.Module,

		// grpc modules
		server.Module,
		server.InterceptorsModule,
		server.GrpcLoggingPayloadModule,

		// gateway modules
		server.GatewayMuxModule,
		server.HTTPModule,

		// logic modules
		profileInternal.Module,
		profilePkg.GatewayModule,
		profilePkg.SwaggerModule,
	)

	app.Run()
}
