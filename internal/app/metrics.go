package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/go-faster/errors"
	"github.com/povilasv/prommod"
	promClient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Metrics wraps application metrics and providers.
type Metrics struct {
	meterProvider *sdkmetric.MeterProvider
	prometheus    http.Handler

	tracerProvider *sdktrace.TracerProvider
	jaeger         *jaeger.Exporter

	resource *resource.Resource
	mux      *http.ServeMux
	srv      *http.Server
}

// Config for metrics.
type Config struct {
	Addr string // address for metrics server, optional
	Name string // default name of the service, optional
}

func (m *Metrics) registerProfiler() {
	// Routes for pprof.
	m.mux.HandleFunc("/debug/pprof/", pprof.Index)
	m.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Manually add support for paths linked to by index page at /debug/pprof/.
	m.mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	m.mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	m.mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	m.mux.Handle("/debug/pprof/block", pprof.Handler("block"))
}

func (m *Metrics) registerPrometheus() {
	// Route for prometheus metrics from registry.
	m.mux.Handle("/metrics", m.prometheus)
}

func (m *Metrics) MeterProvider() metric.MeterProvider {
	return m.meterProvider
}

func (m *Metrics) TracerProvider() trace.TracerProvider {
	return m.tracerProvider
}

func (m *Metrics) Shutdown(ctx context.Context) error {
	return m.srv.Shutdown(ctx)
}

func (m *Metrics) registerRoot() {
	m.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Briefly describe exported endpoints for admin or devops that has
		// only curl and hope for miracle.
		var b strings.Builder
		b.WriteString("Service is up and running.\n\n")
		b.WriteString("Resource:\n")
		for _, a := range m.resource.Attributes() {
			b.WriteString(fmt.Sprintf("  %-32s %s\n", a.Key, a.Value.AsString()))
		}
		b.WriteString("\nAvailable debug endpoints:\n")
		for _, s := range []struct {
			Name        string
			Description string
		}{
			{"/metrics", "prometheus metrics"},
			{"/debug/pprof", "exported pprof"},
		} {
			b.WriteString(fmt.Sprintf("%-20s - %s\n", s.Name, s.Description))
		}
		_, _ = fmt.Fprintln(w, b.String())
	})
}

func (m *Metrics) Run(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		if err := m.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	wg.Go(func() error {
		// Wait until g ctx canceled, then try to shut down server.
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return m.Shutdown(ctx)
	})

	return wg.Wait()
}

// NewMetrics returns new Metrics.
func NewMetrics(log *zap.Logger, cfg Config) (*Metrics, error) {
	if cfg.Addr == "" {
		cfg.Addr = os.Getenv("METRICS_ADDR")
	}
	if cfg.Addr == "" {
		cfg.Addr = "127.0.0.1:9090"
	}

	// The Resource uses environmental variables for default resource attributes:
	//
	// - OTEL_RESOURCE_ATTRIBUTES
	// - OTEL_SERVICE_NAME
	if _, ok := os.LookupEnv("OTEL_SERVICE_NAME"); !ok && cfg.Name != "" {
		// Default service name to provided name.
		_ = os.Setenv("OTEL_SERVICE_NAME", cfg.Name)
	}
	res := resource.Default()

	registry := promClient.NewPedanticRegistry()
	// Register legacy prometheus-only runtime metrics.
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
		collectors.NewBuildInfoCollector(),
		prommod.NewCollector("server"),
	)

	promExporter, err := prometheus.New(
		prometheus.WithRegisterer(registry),
	)
	if err != nil {
		return nil, errors.Wrap(err, "prometheus")
	}
	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(promExporter),
	)

	// The jaeger.WithAgentEndpoint uses environment variables to configure endpoints:
	//
	// - OTEL_EXPORTER_JAEGER_AGENT_HOST is used for the agent address host
	// - OTEL_EXPORTER_JAEGER_AGENT_PORT is used for the agent address port
	//
	// See jaeger.WithAgentEndpoint docs for more info.
	jaegerExporter, err := jaeger.New(jaeger.WithAgentEndpoint())
	if err != nil {
		return nil, errors.Wrap(err, "jaeger")
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(jaegerExporter),
	)

	mux := http.NewServeMux()
	m := &Metrics{
		meterProvider: metricProvider,
		prometheus:    promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),

		tracerProvider: tracerProvider,
		jaeger:         jaegerExporter,

		resource: res,
		mux:      mux,
		srv: &http.Server{
			Handler: mux,
			Addr:    cfg.Addr,
		},
	}

	// Register global OTEL providers.
	otel.SetMeterProvider(m.MeterProvider())
	otel.SetTracerProvider(m.tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)

	m.registerRoot()
	m.registerProfiler()
	m.registerPrometheus()

	log.Info("Metrics initialized",
		zap.Stringer("otel.resource", res),
		zap.String("http.addr", cfg.Addr),
	)

	return m, nil
}
