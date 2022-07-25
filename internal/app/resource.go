package app

import (
	"context"

	"github.com/go-faster/errors"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

// Resource returns new resource for application.
func Resource(ctx context.Context, namespace, name string) (*resource.Resource, error) {
	r, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithProcessRuntimeDescription(),
		resource.WithProcessRuntimeVersion(),
		resource.WithProcessRuntimeName(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(name),
			semconv.ServiceNamespaceKey.String(namespace),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "new")
	}
	return resource.Merge(resource.Default(), r)
}
