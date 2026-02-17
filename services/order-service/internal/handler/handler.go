package handler

import (
	"context"
	"log/slog"
	"net/http"
)

type defaultOrderHandler struct {
	log     *slog.Logger
	service OrderService
}

func NewDefaultOrderHandler(log *slog.Logger, service OrderService) *defaultOrderHandler {
	return &defaultOrderHandler{
		log:     log,
		service: service,
	}
}

type OrderService interface {
	CreateOrder(ctx context.Context, request CreateOrderRequest) (CreateOrderResponse, error)
	GetOrder(ctx context.Context) (OrderResponse, error)
}

func createOrder(r *http.Request, w http.ResponseWriter) {

}
