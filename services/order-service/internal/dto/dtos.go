package dto

import "time"

type CreateOrderRequest struct {
	PaymentMethod string `json:"payment_method"`
	UserID        int64  `json:"user_id"`
	Items         []*CreateOrderItemRequest
}

type CreateOrderItemRequest struct {
	Quantity  int64 `json:"quantity"`
	ProductID int64 `json:"product_id"`
}

type CreateOrderResponse struct {
	ID            int64      `json:"id"`
	PaymentMethod string     `json:"payment_method"`
	TaxPrice      float64    `json:"tax_price"`
	ShippingPrice float64    `json:"shipping_price"`
	TotalPrice    float64    `json:"total_price"`
	UserID        int64      `json:"user_id"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Items         []*CreateOrderItemResponse
}

type CreateOrderItemResponse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Quantity  int64   `json:"quantity"`
	Image     string  `json:"image"`
	Price     float64 `json:"price"`
	ProductID int64   `json:"product_id"`
	OrderID   int64   `json:"order_id"`
}

type OrderItemResponse struct {
	ID        int64   `db:"id"`
	Name      string  `db:"name"`
	Quantity  int64   `db:"quantity"`
	Image     string  `db:"image"`
	Price     float64 `db:"price"`
	ProductID int64   `db:"product_id"`
	OrderID   int64   `db:"order_id"`
}

type OrderResponse struct {
	ID            int64      `db:"id"`
	PaymentMethod string     `db:"payment_method"`
	TaxPrice      float64    `db:"tax_price"`
	ShippingPrice float64    `db:"shipping_price"`
	TotalPrice    float64    `db:"total_price"`
	UserID        int64      `db:"user_id"`
	Status        string     `db:"status"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	Items         []*OrderItemResponse
}
