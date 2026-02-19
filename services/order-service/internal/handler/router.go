package handler

import (
	"log/slog"
	"net/http"
	"time"

	apphttp "github.com/defan6/market/services/order-service/internal/app/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *defaultOrderHandler, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(apphttp.Logging(log))
	r.Route("/api/v1/orders", func(r chi.Router) {
		r.Post("/", handler.CreateOrder)
		r.Get("/", handler.GetOrder)
	})

	return r
}
