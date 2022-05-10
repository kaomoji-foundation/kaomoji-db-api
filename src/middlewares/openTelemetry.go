package middlewares

import (
	"github.com/gofiber/fiber/v2"
	fiberOtel "github.com/psmarcin/fiber-opentelemetry/pkg/fiber-otel"
	"go.opentelemetry.io/otel/trace"
)

func OpenTelemery() func(*fiber.Ctx) error {
	// create new middleware using custom config
	otelMiddleware := fiberOtel.New(fiberOtel.Config{
		// name for root span in trace on request
		SpanName: "http/request",
		// array of span options for root span
		TracerStartAttributes: []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindConsumer),
		},
		// key name for context store in fiber.Ctx
		LocalKeyName: "OpenTelemetryCtx",
	})
	return otelMiddleware
}
