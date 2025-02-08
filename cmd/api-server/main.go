package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/app"
	"github.com/go-faster/sdk/zctx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"example/internal/api"
	"example/internal/httpmiddleware"
	"example/internal/oas"
)

const shutdownTimeout = 15 * time.Second

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger, m *app.Telemetry) error {
		var arg struct {
			Addr string
		}
		flag.StringVar(&arg.Addr, "addr", "0.0.0.0:8080", "listen address")
		flag.Parse()

		lg.Info("Initializing",
			zap.String("http.addr", arg.Addr),
		)
		oasServer, err := oas.NewServer(api.Handler{},
			oas.WithTracerProvider(m.TracerProvider()),
			oas.WithMeterProvider(m.MeterProvider()),
		)
		if err != nil {
			return errors.Wrap(err, "server init")
		}

		// Using OpenTelemetry instrumentation for HTTP server.
		routeFinder := httpmiddleware.MakeRouteFinder(oasServer)
		httpServer := http.Server{
			ReadHeaderTimeout: time.Second,
			Addr:              arg.Addr,
			Handler: httpmiddleware.Wrap(oasServer,
				httpmiddleware.InjectLogger(zctx.From(ctx)),
				httpmiddleware.Instrument("api", routeFinder, m),
				httpmiddleware.LogRequests(routeFinder),
				httpmiddleware.Labeler(routeFinder),
			),
		}
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			// Wait until g ctx canceled, then try to shut down server.
			<-ctx.Done()

			lg.Info("Shutting down", zap.Duration("timeout", shutdownTimeout))

			shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()
			return httpServer.Shutdown(shutdownCtx)
		})
		g.Go(func() error {
			defer lg.Info("Server stopped")
			if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return errors.Wrap(err, "http")
			}
			return nil
		})

		return g.Wait()
	})
}
