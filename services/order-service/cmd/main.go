package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/defan6/market/services/order-service/internal/config"
	"github.com/defan6/market/services/shared/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// load config
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// setup logger
	log := setupLogger(cfg.Env)

	// setup app

	// run app

	// graceful shutdown
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOptions: &slog.HandlerOptions{Level: slog.LevelDebug}
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
