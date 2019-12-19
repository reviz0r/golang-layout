package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"

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
	// Start GRPC server
	go func() {
		err := startGRPC(grpcPort)
		if err != nil {
			logger.Fatalf("startGRPC: %v", err)
		}
	}()

	// Start GRPC-Gateway server
	go func() {
		err := startHTTP(context.TODO(), httpPort, grpcPort, grpc.WithInsecure())
		if err != nil {
			logger.Fatalf("startHTTP: %v", err)
		}
	}()

	select {}
}

func startGRPC(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

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

	logger.Printf("Start grpc server on port %s", port)
	return grpcServer.Serve(lis)
}

func startHTTP(ctx context.Context, port, grpcPort string, opts ...grpc.DialOption) error {
	protoMux := runtime.NewServeMux()

	mux := http.NewServeMux()
	mux.Handle("/", protoMux)

	err := pkg.RegisterUserServiceHandlerFromEndpoint(
		ctx, protoMux, "localhost"+grpcPort, opts)
	if err != nil {
		return err
	}

	mux.HandleFunc("/docs/profile/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "api/proto/profile/profile.api.swagger.json")
		})

	logger.Printf("Start http server on port %s", port)
	return http.ListenAndServe(port, mux)
}
