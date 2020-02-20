package server

import (
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/fx"
)

var GatewayMuxModule = fx.Options(
	fx.Provide(runtime.NewServeMux, NewServeMuxMarshallerOption),
	fx.Invoke(RegisterProtoMux),
)

type ServeMuxMarshallerParams struct {
	fx.In

	EnumsAsInts  bool               `name:"gateway_marshaller_enums_as_ints" optional:"true"`
	EmitDefaults bool               `name:"gateway_marshaller_emit_defaults" optional:"true"`
	Indent       string             `name:"gateway_marshaller_indent" optional:"true"`
	OrigName     bool               `name:"gateway_marshaller_orig_name" optional:"true"`
	AnyResolver  jsonpb.AnyResolver `name:"gateway_marshaller_any_resolver" optional:"true"`
}

func NewServeMuxMarshallerOption(p ServeMuxMarshallerParams) runtime.ServeMuxOption {
	marshaller := &runtime.JSONPb{
		EnumsAsInts:  p.EnumsAsInts,
		EmitDefaults: p.EmitDefaults,
		Indent:       p.Indent,
		OrigName:     p.OrigName,
		AnyResolver:  p.AnyResolver,
	}

	return runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaller)
}

func RegisterProtoMux(mux *http.ServeMux, gatewayMux *runtime.ServeMux) {
	mux.Handle("/", gatewayMux)
}
