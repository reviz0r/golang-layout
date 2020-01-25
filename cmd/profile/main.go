package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/reviz0r/golang-layout/pkg/config"
	"github.com/reviz0r/golang-layout/pkg/db"
	"github.com/reviz0r/golang-layout/pkg/grace"
	"github.com/reviz0r/golang-layout/pkg/log"
	"github.com/reviz0r/golang-layout/pkg/server"

	profile_internal "github.com/reviz0r/golang-layout/internal/profile"
	profile_pkg "github.com/reviz0r/golang-layout/pkg/profile"
)

// TODO: make config
const (
	grpcPort = ":50051"
	httpPort = ":8081"
)

func main() {
	app := fx.New(
		config.Module,
		log.Module,
		db.Module,
		server.Module,
		server.InterceptorsModule,
		server.GrpcLoggingPayloadModule,

		// logic modules
		profile_internal.Module,
	)

	go app.Run()

	mainCtx, mainWg := grace.Grace()

	// Start GRPC-Gateway server
	mainWg.Add(1)
	go startHTTP(mainCtx, mainWg, httpPort, grpcPort, grpc.WithInsecure())

	mainWg.Wait()
}

func startHTTP(ctx context.Context, wg *sync.WaitGroup, port, grpcPort string, opts ...grpc.DialOption) {
	defer wg.Done()

	protoMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true}))

	err := profile_pkg.RegisterUserServiceHandlerFromEndpoint(
		ctx, protoMux, "localhost"+grpcPort, opts)
	if err != nil {
		panic(fmt.Errorf("cannot register user service: %v", err))
	}

	mux := http.NewServeMux()
	mux.Handle("/", protoMux)

	mux.HandleFunc("/docs/profile/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "api/proto/profile/profile.api.swagger.json")
		})

	s := http.Server{Addr: port, Handler: mux}

	go func() {
		<-ctx.Done()
		s.Shutdown(ctx)
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(fmt.Errorf("cannot listen http port %s %v", port, err))
	}

	err = s.Serve(lis)
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("cannot serve http: %v", err))
	}
}
