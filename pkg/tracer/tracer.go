package tracer

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
)

// Module register tracer in DI container
var Module = fx.Provide(NewTracer)

// NewTracer .
func NewTracer() opentracing.Tracer {
	return opentracing.NoopTracer{}
}
