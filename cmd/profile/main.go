package main

import (
	"context"
	"log"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	internal "github.com/reviz0r/http-api/internal/profile"
	pkg "github.com/reviz0r/http-api/pkg/profile"
)

var logger = logrus.NewEntry(logrus.New())

func main() {
	// Start GRPC server
	go startGRPC(":50051")

	// Start GRPC-Gateway server
	go startHTTP(":8081")

	select {}
}

func startGRPC(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(logger),
			grpc_logrus.PayloadStreamServerInterceptor(logger,
				func(ctx context.Context, fullMethodName string, servingObject interface{}) bool { return true }),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logger),
			grpc_logrus.PayloadUnaryServerInterceptor(logger,
				func(ctx context.Context, fullMethodName string, servingObject interface{}) bool { return true }),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	pkg.RegisterUserServiceServer(grpcServer, new(internal.UserService))

	log.Printf("Start grpc server on port %s", port)
	return grpcServer.Serve(lis)
}

func startHTTP(port string) error {
	protoMux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pkg.RegisterUserServiceHandlerFromEndpoint(
		context.TODO(), protoMux, "localhost:50051", opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", protoMux)

	mux.HandleFunc("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/proto/profile/profile.api.swagger.json")
	})

	log.Printf("Start http server on port %s", port)
	return http.ListenAndServe(port, mux)
}
