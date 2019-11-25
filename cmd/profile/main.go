package main

import (
	"context"
	"log"
	"net"
	"net/http"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	internal "github.com/reviz0r/http-api/internal/profile"
	pkg "github.com/reviz0r/http-api/pkg/profile"
)

var logger = logrus.NewEntry(logrus.New())

func main() {
	// Start GRPC server
	go func () {
		err := startGRPC(":50051")
		if err != nil {
			log.Fatalf("startGRPC: %v", err)
		}
	}()

	// Start GRPC-Gateway server
	go func () {
		err := startHTTP(":8081")
		if err != nil {
			log.Fatalf("startHTTP: %v", err)
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
			grpcLogrus.PayloadStreamServerInterceptor(logger,
				func(ctx context.Context, fullMethodName string, servingObject interface{}) bool { return true }),
			grpcRecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcLogrus.UnaryServerInterceptor(logger),
			grpcLogrus.PayloadUnaryServerInterceptor(logger,
				func(ctx context.Context, fullMethodName string, servingObject interface{}) bool { return true }),
			grpcRecovery.UnaryServerInterceptor(),
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
