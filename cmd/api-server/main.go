package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"example/internal/api"
	"example/internal/app"
	"example/internal/oas"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		var arg struct {
			Addr        string
			MetricsAddr string
		}
		flag.StringVar(&arg.Addr, "addr", "127.0.0.1:8080", "listen address")
		flag.StringVar(&arg.MetricsAddr, "metrics.addr", "127.0.0.1:9090", "metrics listen address")
		flag.Parse()

		lg.Info("Initializing",
			zap.String("http.addr", arg.Addr),
			zap.String("metrics.addr", arg.MetricsAddr),
		)

		m, err := app.NewMetrics(lg, app.Config{
			Addr: arg.MetricsAddr,
			Name: "api",
		})
		if err != nil {
			return errors.Wrap(err, "metrics")
		}

		oasServer, err := oas.NewServer(api.Handler{},
			oas.WithTracerProvider(m.TracerProvider()),
			oas.WithMeterProvider(m.MeterProvider()),
		)
		if err != nil {
			return errors.Wrap(err, "server init")
		}
		httpServer := http.Server{
			Addr:    arg.Addr,
			Handler: oasServer,
		}

		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			return m.Run(ctx)
		})
		g.Go(func() error {
			<-ctx.Done()
			return httpServer.Shutdown(ctx)
		})
		g.Go(func() error {
			defer lg.Info("Server stopped")
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return errors.Wrap(err, "http")
			}
			return nil
		})

		return g.Wait()
	})
}
