package http

import (
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
		Addr:         fmt.Sprintf(":%d", cfg.Port),
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
