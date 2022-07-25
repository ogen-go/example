package app

import (
	"context"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

const EnvLogLevel = "LOG_LEVEL"

const (
	exitCodeOk       = 0
	exitCodeWatchdog = 1
)

const (
	shutdownTimeout = time.Second * 5
	watchdogTimeout = shutdownTimeout + time.Second*5
)

// Run f until interrupt.
func Run(f func(ctx context.Context, log *zap.Logger) error) {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := zap.NewProductionConfig()

	if s := os.Getenv(EnvLogLevel); s != "" {
		var lvl zapcore.Level
		if err := lvl.UnmarshalText([]byte(s)); err != nil {
			panic(err)
		}
		cfg.Level.SetLevel(lvl)
	}

	lg, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		if err := f(ctx, lg); err != nil {
			return err
		}
		return nil
	})
	go func() {
		// Guaranteed way to kill application.
		<-ctx.Done()

		// Context is canceled, giving application time to shut down gracefully.
		lg.Info("Waiting for application shutdown")
		time.Sleep(watchdogTimeout)

		// Probably deadlock, forcing shutdown.
		lg.Warn("Graceful shutdown watchdog triggered: forcing shutdown")
		os.Exit(exitCodeWatchdog)
	}()

	if err := wg.Wait(); err != nil {
		if err == context.Canceled {
			lg.Info("Graceful shutdown")
			return
		}
		lg.Fatal("Failed",
			zap.Error(err),
		)
	}

	os.Exit(exitCodeOk)
}
