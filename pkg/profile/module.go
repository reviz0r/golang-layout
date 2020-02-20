package profile

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var GatewayModule = fx.Invoke(UserServiceGateway)

type GatewayParams struct {
	fx.In

	Mux      *runtime.ServeMux
	Endpoint string            `name:"gateway_user_service_endpoint"`
	Opts     []grpc.DialOption `group:"gateway_user_service_dial_options"`
}

func UserServiceGateway(p GatewayParams) error {
	return RegisterUserServiceHandlerFromEndpoint(context.TODO(), p.Mux, p.Endpoint, p.Opts)
}

var GatewayInsecureDialModule = fx.Provide(DialInsecure)

type DialInsecureResult struct {
	fx.Out

	Op grpc.DialOption `group:"gateway_user_service_dial_options"`
}

func DialInsecure() DialInsecureResult {
	return DialInsecureResult{Op: grpc.WithInsecure()}
}

var SwaggerModule = fx.Invoke(RegisterProfileSwagger)

func RegisterProfileSwagger(mux *http.ServeMux) {
	mux.HandleFunc("/docs/profile/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "api/proto/profile/profile.api.swagger.json")
		})
}
