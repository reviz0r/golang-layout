package server

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var GatewayMuxModule = fx.Options(
	fx.Provide(runtime.NewServeMux, NewServeMuxMarshallerOption),
	fx.Invoke(RegisterProtoMux),
)

type ServeMuxMarshallerParams struct {
	fx.In

	AnyResolver jsonpb.AnyResolver `name:"gateway_marshaller_any_resolver" optional:"true"`
}

func NewServeMuxMarshallerOption(p ServeMuxMarshallerParams, config *viper.Viper) runtime.ServeMuxOption {
	marshaller := &runtime.JSONPb{
		EnumsAsInts:  config.GetBool("gateway.marshaler.enums_as_ints"),
		EmitDefaults: config.GetBool("gateway.marshaler.emit_defaults"),
		Indent:       config.GetString("gateway.marshaler.indent"),
		OrigName:     config.GetBool("gateway.marshaler.orig_name"),
		AnyResolver:  p.AnyResolver,
	}

	return runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaller)
}

func RegisterProtoMux(mux *http.ServeMux, gatewayMux *runtime.ServeMux) {
	mux.Handle("/", gatewayMux)
}
