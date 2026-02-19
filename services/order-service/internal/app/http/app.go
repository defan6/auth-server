package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/defan6/market/services/order-service/internal/config"
)

type App struct {
	log        *slog.Logger
	httpServer http.Server
	port       int
}

func New(log *slog.Logger, cfg config.HttpConfig, handler http.Handler) *App {
	httpServer := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return &App{
		log:        log,
		httpServer: httpServer,
		port:       cfg.Port,
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("HTTP server is running", slog.String("addr", a.httpServer.Addr))
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	const op = "httpapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping HTTP server", slog.Int("port", a.port))

	return a.httpServer.Shutdown(ctx)
}
