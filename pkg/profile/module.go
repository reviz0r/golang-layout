package profile

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var GatewayModule = fx.Invoke(UserServiceGateway)

func UserServiceGateway(config *viper.Viper, mux *runtime.ServeMux) error {
	return RegisterUserServiceHandlerFromEndpoint(
		context.TODO(), mux, config.GetString("gateway.profile_service_endpoint"), []grpc.DialOption{grpc.WithInsecure()})
}

var SwaggerModule = fx.Invoke(RegisterProfileSwagger)

func RegisterProfileSwagger(mux *http.ServeMux) {
	mux.HandleFunc("/docs/profile/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "api/proto/profile/profile.api.swagger.json")
		})
}
