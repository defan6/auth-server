package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/defan6/market/services/order-service/internal/domain"
	"github.com/defan6/market/services/order-service/internal/dto"
	"github.com/defan6/market/services/order-service/internal/mapper"
	"github.com/defan6/market/services/order-service/internal/storage"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidQuantity    = errors.New("invalid quantity")
	ErrInvalidProductID   = errors.New("invalid product id")
	ErrEmptyOrderItems    = errors.New("order items is empty")
	ErrPaymentMethodEmpty = errors.New("payment method is required")
)

type defaultOrderService struct {
	log           *slog.Logger
	storage       OrderStorage
	txManager     *storage.TxManager
	productClient ProductClient
}

type OrderStorage interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	CreateOrderItems(ctx context.Context, items []*domain.OrderItem) error
	GetOrder(ctx context.Context, id int64) (*domain.Order, error)
}

func NewDefaultOrderService(
	log *slog.Logger,
	storage OrderStorage,
	txManager *storage.TxManager,
	productClient ProductClient,
) *defaultOrderService {
	return &defaultOrderService{
		log:           log,
		storage:       storage,
		txManager:     txManager,
		productClient: productClient,
	}
}

func (s *defaultOrderService) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	const op = "service.CreateOrder"

	// 1. Валидация запроса
	if err := validateCreateOrderReq(request); err != nil {
		return nil, fmt.Errorf("%s: invalid request: %w", op, err)
	}

	// 2. Получаем продукты из БД (цены не доверяем клиенту!)
	domainItems, itemsPrice, err := s.processOrderRequest(ctx, request, op)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to process order request: %w", op, err)
	}

	// 3. Рассчитываем цены
	taxPrice, shippingPrice, totalPrice := calculateOrderPrices(itemsPrice)

	// 4. Создаём заказ
	order := &domain.Order{
		PaymentMethod: request.PaymentMethod,
		TaxPrice:      taxPrice,
		ShippingPrice: shippingPrice,
		TotalPrice:    totalPrice,
		UserID:        request.UserID,
		Status:        domain.OrderStatusPending,
		CreatedAt:     time.Now(),
		Items:         domainItems,
	}

	var createdOrder *domain.Order

	err = s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.storage.CreateOrder(ctx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Проставляем OrderID для элементов
		for _, item := range order.Items {
			item.OrderID = order.ID
		}

		if err := s.storage.CreateOrderItems(ctx, order.Items); err != nil {
			return fmt.Errorf("failed to create order items: %w", err)
		}

		createdOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	response := mapper.MapToOrderResponseFromOrder(createdOrder)
	return response, nil
}

func (s *defaultOrderService) GetOrder(ctx context.Context, id int64) (*dto.OrderResponse, error) {
	order, err := s.storage.GetOrder(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	response := mapper.MapToOrderResponseFromOrder(order)
	return response, nil
}

// processOrderRequest получает продукты из product-service, проверяет наличие и создаёт элементы заказа
func (s *defaultOrderService) processOrderRequest(ctx context.Context, req *dto.CreateOrderRequest, op string) ([]*domain.OrderItem, float64, error) {
	// Собираем уникальные ID продуктов
	productIDs := make([]int64, 0, len(req.Items))
	uniqueProductIDs := make(map[int64]struct{})
	for _, item := range req.Items {
		if _, exists := uniqueProductIDs[item.ProductID]; !exists {
			uniqueProductIDs[item.ProductID] = struct{}{}
			productIDs = append(productIDs, item.ProductID)
		}
	}

	// Получаем продукты из product-service (с реальными ценами)
	products, err := s.productClient.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products from product-service: %w", err)
	}

	// Проверяем, что все продукты найдены
	if len(products) != len(uniqueProductIDs) {
		return nil, 0, fmt.Errorf("%s: not all products found", op)
	}

	// Создаём мапу для быстрого доступа
	productMap := make(map[int64]*dto.ExternalProduct, len(products))
	for _, product := range products {
		productMap[product.ID] = product
	}

	// Создаём элементы заказа с реальными ценами
	domainItems := make([]*domain.OrderItem, 0, len(req.Items))
	var itemsPrice float64 = 0

	for _, item := range req.Items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return nil, 0, fmt.Errorf("product %d not found", item.ProductID)
		}

		// Проверяем наличие на складе
		if product.CountInStock < item.Quantity {
			return nil, 0, fmt.Errorf("not enough stock for product %d: requested %d, available %d",
				item.ProductID, item.Quantity, product.CountInStock)
		}

		orderItem := &domain.OrderItem{
			Name:      product.Name,
			Quantity:  item.Quantity,
			Image:     product.Image,
			Price:     product.Price, // ← Цена из product-service, не от клиента!
			ProductID: product.ID,
		}
		domainItems = append(domainItems, orderItem)

		itemsPrice += product.Price * float64(item.Quantity)
	}

	return domainItems, itemsPrice, nil
}

// calculateOrderPrices рассчитывает налог, доставку и итоговую сумму
func calculateOrderPrices(itemsPrice float64) (float64, float64, float64) {
	const taxRate = 0.10      // 10% налог
	const shippingPrice = 150 // Фиксированная доставка

	taxPrice := itemsPrice * taxRate
	totalPrice := itemsPrice + taxPrice + shippingPrice

	return taxPrice, shippingPrice, totalPrice
}

func validateCreateOrderReq(req *dto.CreateOrderRequest) error {
	if req.PaymentMethod == "" {
		return ErrPaymentMethodEmpty
	}

	if req.Items == nil || len(req.Items) == 0 {
		return ErrEmptyOrderItems
	}

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}
		if item.ProductID <= 0 {
			return ErrInvalidProductID
		}
	}

	return nil
}
