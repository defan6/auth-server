package dto

// ExternalProduct DTO для получения данных о продукте из product-service
type ExternalProduct struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	CountInStock int64   `json:"count_in_stock"`
	Image        string  `json:"image"`
	Description  string  `json:"description"`
}
