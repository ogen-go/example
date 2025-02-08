package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/app"
	"github.com/go-faster/sdk/zctx"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"

	"example/internal/api"
	"example/internal/oas"
)

func run(ctx context.Context, lg *zap.Logger, m *app.Telemetry) error {
	var arg struct {
		BaseURL string
		ID      int64
	}
	flag.StringVar(&arg.BaseURL, "url", "http://server:8080", "target server url")
	flag.Int64Var(&arg.ID, "id", 1337, "pet id to request")
	flag.Parse()

	// For route finding.
	oasServer, err := oas.NewServer(api.Handler{})
	if err != nil {
		return errors.Wrap(err, "server init")
	}

	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport,
			otelhttp.WithTracerProvider(m.TracerProvider()),
			otelhttp.WithMeterProvider(m.MeterProvider()),
			otelhttp.WithPropagators(m.TextMapPropagator()),
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				route, ok := oasServer.FindPath(r.Method, r.URL)
				if !ok {
					return operation
				}
				return route.OperationID()
			}),
		),
	}
	client, err := oas.NewClient(arg.BaseURL,
		oas.WithClient(httpClient),
		oas.WithMeterProvider(m.MeterProvider()),
		oas.WithTracerProvider(m.TracerProvider()),
	)
	if err != nil {
		return errors.Wrap(err, "client")
	}

	tracer := m.TracerProvider().Tracer("example")
	fetchPet := func(ctx context.Context) error {
		ctx, span := tracer.Start(ctx, "tick")
		defer span.End()
		res, err := client.GetPetById(ctx, oas.GetPetByIdParams{
			PetId: arg.ID,
		})
		if err != nil {
			return errors.Wrap(err, "get pet")
		}
		zctx.From(ctx).Info("Got pet", zap.Any("pet", res))
		return nil
	}
	tick := func() {
		if err := fetchPet(ctx); err != nil {
			zctx.From(ctx).Error("Failed to fetch pet", zap.Error(err))
		}
	}
	tick()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			tick()
		}
	}
}

func main() {
	app.Run(run)
}
