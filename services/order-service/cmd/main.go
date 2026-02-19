package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	apphttp "github.com/defan6/market/services/order-service/internal/app/http"
	"github.com/defan6/market/services/order-service/internal/config"
	"github.com/defan6/market/services/order-service/internal/handler"
	"github.com/defan6/market/services/order-service/internal/service"
	"github.com/defan6/market/services/order-service/internal/storage"
	database "github.com/defan6/market/services/order-service/storage"
	"github.com/defan6/market/services/shared/logger/handlers/slogpretty"
	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// setup logger
	log := setupLogger(cfg.Env)

	log.Info("starting app", slog.String("env", cfg.Env))

	// init storage
	db := database.NewDatabase(&cfg.DB)

	// init layers
	orderStorage := storage.NewDefaultOrderStorage(db.GetDB())
	txManager := storage.NewTxManager(db.GetDB())
	productClient := service.NewStubProductClient() // заглушка до создания product-service
	orderService := service.NewDefaultOrderService(log, orderStorage, txManager, productClient)
	orderHandler := handler.NewDefaultOrderHandler(log, orderService)

	// init router
	router := handler.NewRouter(orderHandler, log)

	// init app
	app := apphttp.New(log, cfg.Server, router)

	// run app
	go app.MustRun()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	log.Info("stopping application...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Error("failed to stop app", slog.String("err", err.Error()))
	}

	log.Info("app stopped")
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
		SlogOptions: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
