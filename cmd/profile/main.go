package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"sync"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	internal "github.com/reviz0r/http-api/internal/profile"
	pkg "github.com/reviz0r/http-api/pkg/profile"
)

// TODO: make config
const (
	grpcPort = ":50051"
	httpPort = ":8081"
)

var (
	logger *logrus.Entry
	db     *sql.DB
)

func init() {
	logger = NewLogger()
	db = NewDatabase()
}

func main() {
	logger.Debug("app starting")

	mainCtx, mainWg := Grace()

	// Start GRPC server
	mainWg.Add(1)
	go startGRPC(mainCtx, mainWg, grpcPort)

	// Start GRPC-Gateway server
	mainWg.Add(1)
	go startHTTP(mainCtx, mainWg, httpPort, grpcPort, grpc.WithInsecure())

	logger.Debug("app started")

	mainWg.Wait()

	logger.Debug("app is shutdown")
}

func startGRPC(ctx context.Context, wg *sync.WaitGroup, port string) {
	defer wg.Done()

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcLogrus.StreamServerInterceptor(logger),
			grpcLogrus.PayloadStreamServerInterceptor(logger, LoggingPayloadDecider()),
			grpcRecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcLogrus.UnaryServerInterceptor(logger),
			grpcLogrus.PayloadUnaryServerInterceptor(logger, LoggingPayloadDecider()),
			grpcRecovery.UnaryServerInterceptor(),
		)),
	)
	pkg.RegisterUserServiceServer(grpcServer, &internal.UserService{DB: db})

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		logger.Debug("grpc server is shutdown")
	}()

	logger.Debugf("Start grpc server on port %s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(fmt.Errorf("cannot listen grpc port %s %v", port, err))
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(fmt.Errorf("cannot serve grpc: %v", err))
	}
}

func startHTTP(ctx context.Context, wg *sync.WaitGroup, port, grpcPort string, opts ...grpc.DialOption) {
	defer wg.Done()

	protoMux := runtime.NewServeMux()

	err := pkg.RegisterUserServiceHandlerFromEndpoint(
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
		logger.Debug("http server is shutdown")
	}()

	logger.Debugf("Start http server on port %s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(fmt.Errorf("cannot listen http port %s %v", port, err))
	}

	err = s.Serve(lis)
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("cannot serve http: %v", err))
	}
}
