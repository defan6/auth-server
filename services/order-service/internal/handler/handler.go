package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/defan6/market/services/order-service/internal/dto"
	"github.com/go-chi/chi/v5"
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
	CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetOrder(ctx context.Context, id int64) (*dto.OrderResponse, error)
}

func (h *defaultOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateOrder"

	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(op, slog.String("error", err.Error()))
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateOrder(r.Context(), &req)
	if err != nil {
		h.log.Error(op, slog.String("error", err.Error()))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(op, slog.String("error", err.Error()))
	}
}

func (h *defaultOrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetOrder"

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Error(op, slog.String("error", "invalid order id"))
		http.Error(w, `{"error":"invalid order id"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		h.log.Error(op, slog.String("error", err.Error()))
		http.Error(w, `{"error":"order not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(op, slog.String("error", err.Error()))
	}
}
