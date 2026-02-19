package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/defan6/market/services/order-service/internal/dto"
)

var ErrProductNotFound = errors.New("product not found")

// ProductClient интерфейс для взаимодействия с product-service
type ProductClient interface {
	GetProductsByIDs(ctx context.Context, ids []int64) ([]*dto.ExternalProduct, error)
}

// StubProductClient заглушка для разработки
// TODO: заменить на HTTP/gRPC клиент после создания product-service
type StubProductClient struct{}

func NewStubProductClient() *StubProductClient {
	return &StubProductClient{}
}

// GetProductsByIDs временная реализация с заглушками
// TODO: заменить на HTTP запрос к product-service: GET /api/v1/products?ids=1,2,3
func (c *StubProductClient) GetProductsByIDs(ctx context.Context, ids []int64) ([]*dto.ExternalProduct, error) {
	products := make([]*dto.ExternalProduct, 0, len(ids))

	for _, id := range ids {
		// Заглушка: возвращаем фейковые данные
		product := &dto.ExternalProduct{
			ID:           id,
			Name:         fmt.Sprintf("Product-%d", id),
			Price:        100.0 * float64(id), // Фейковая цена для тестирования
			CountInStock: 100,                 // Фейковое наличие
			Image:        "/images/product.png",
			Description:  "Stub product description",
		}
		products = append(products, product)
	}

	return products, nil
}
