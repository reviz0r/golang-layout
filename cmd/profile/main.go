package main

import (
	"go.uber.org/fx"

	"github.com/reviz0r/golang-layout/pkg/config"
	"github.com/reviz0r/golang-layout/pkg/db"
	"github.com/reviz0r/golang-layout/pkg/logger"
	"github.com/reviz0r/golang-layout/pkg/server"
	"github.com/reviz0r/golang-layout/pkg/tracer"

	profileInternal "github.com/reviz0r/golang-layout/internal/profile"
	profilePkg "github.com/reviz0r/golang-layout/pkg/profile"
)

func main() {
	app := fx.New(
		fx.NopLogger,

		config.Module,
		config.DefaultValues,
		logger.Module,

		db.Module,

		// grpc modules
		server.Module,
		server.InterceptorsModule,
		server.GrpcLoggingPayloadModule,
		server.PrometheusMetrics,
		tracer.Module,

		// gateway modules
		server.GatewayMuxModule,
		server.HTTPModule,
		server.PrometheusMetricsHandler,

		// logic modules
		profileInternal.Module,
		profilePkg.GatewayModule,
		profilePkg.SwaggerModule,
	)

	app.Run()
}
