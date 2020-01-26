package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/fx"
)

var MuxModule = fx.Provide(runtime.NewServeMux, NewGatewayMarshallerOption)

type GatewayMarshallerParams struct {
	fx.In

	EnumsAsInts  bool   `name:"gateway_marshaller_enums_as_ints" optional:"true"`
	EmitDefaults bool   `name:"gateway_marshaller_emit_defaults" optional:"true"`
	Indent       string `name:"gateway_marshaller_indent" optional:"true"`
	OrigName     bool   `name:"gateway_marshaller_orig_name" optional:"true"`
}

func NewGatewayMarshallerOption(p GatewayMarshallerParams) runtime.ServeMuxOption {
	return runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true})
}
