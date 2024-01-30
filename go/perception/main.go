package perception

import (
	"bytes"
	"io"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type perception struct {
	serviceName string
}

func Init(serviceName string) (*perception, error) {

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())

	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	return &perception{serviceName: serviceName}, nil
}

func (p *perception) TraceHTTPHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Start a span for the HTTP request
		ctx, span := otel.Tracer("perception/main").Start(r.Context(), "HTTP "+r.Method)
		defer span.End()

		if r.ContentLength > 0 {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				span.RecordError(err)
			} else {
				span.SetAttributes(attribute.String("http.request.body", string(bodyBytes)))
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.scheme", r.URL.Scheme),
			attribute.String("http.host", r.URL.Host),
			attribute.String("http.path", r.URL.Path),
		)

		for name, values := range r.Header {
			span.SetAttributes(attribute.String("http.header."+name, values[0]))
		}

		// Call the original handler with the enriched context
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}
