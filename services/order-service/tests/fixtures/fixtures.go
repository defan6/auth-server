package fixtures

import (
	"os"
	"path/filepath"
	"testing"
)

// GetFixturePath возвращает абсолютный путь к файлу фикстуры
func GetFixturePath(fileName string) string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "fixtures", fileName)
}

// LoadOrderFixture загружает JSON фикстуру для заказа
func LoadOrderFixture(t *testing.T, fileName string) []byte {
	t.Helper()

	path := GetFixturePath(fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load fixture %s: %v", fileName, err)
	}

	return data
}

// FixturePaths содержит пути ко всем фикстурам
var FixturePaths = struct {
	// Valid запросы
	CreateOrderValid string

	// Invalid запросы
	CreateOrderEmptyPayment     string
	CreateOrderEmptyItems       string
	CreateOrderInvalidQuantity  string
	CreateOrderInvalidProductID string
	CreateOrderMultipleItems    string

	// Ответы
	OrderResponse string

	// Ошибки
	ErrorEmptyPayment  string
	ErrorEmptyItems    string
	ErrorOrderNotFound string
}{
	CreateOrderValid:            "create_order_valid.json",
	CreateOrderEmptyPayment:     "create_order_empty_payment.json",
	CreateOrderEmptyItems:       "create_order_empty_items.json",
	CreateOrderInvalidQuantity:  "create_order_invalid_quantity.json",
	CreateOrderInvalidProductID: "create_order_invalid_product_id.json",
	CreateOrderMultipleItems:    "create_order_multiple_items.json",
	OrderResponse:               "order_response.json",
	ErrorEmptyPayment:           "error_empty_payment.json",
	ErrorEmptyItems:             "error_empty_items.json",
	ErrorOrderNotFound:          "error_order_not_found.json",
}
