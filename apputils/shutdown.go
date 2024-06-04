package apputils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofreego/goutils/logger"
)

type Application interface {
	Name() string
	Shutdown(ctx context.Context)
}

func GracefulShutdown(ctx context.Context, apps ...Application) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-sigCh
	logger.Info(ctx, "Shutting down... please wait ....")

	for _, app := range apps {
		logger.Info(ctx, "Shutting down %s", app.Name())
		app.Shutdown(ctx)
		logger.Info(ctx, "%s is down", app.Name())
	}
}
