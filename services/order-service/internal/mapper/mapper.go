package mapper

import (
	"github.com/defan6/market/services/order-service/internal/domain"
	"github.com/defan6/market/services/order-service/internal/dto"
)

func MapToOrderFromCreateOrderRequest(orderReq *dto.CreateOrderRequest) *domain.Order {
	var orderItems []*domain.OrderItem

	for _, item := range orderReq.Items {
		orderItem := mapToOrderItemFromOrderItemRequest(item)
		orderItems = append(orderItems, orderItem)
	}

	return &domain.Order{
		PaymentMethod: orderReq.PaymentMethod,
		UserID:        orderReq.UserID,
		Items:         orderItems,
	}
}

func mapToOrderItemFromOrderItemRequest(orderItemReq *dto.CreateOrderItemRequest) *domain.OrderItem {
	return &domain.OrderItem{
		ProductID: orderItemReq.ProductID,
		Quantity:  orderItemReq.Quantity,
	}
}

func mapToOrderItemResponseFromOrderItem(orderItem *domain.OrderItem) *dto.OrderItemResponse {
	return &dto.OrderItemResponse{
		ID:        orderItem.ID,
		Name:      orderItem.Name,
		Quantity:  orderItem.Quantity,
		Image:     orderItem.Image,
		Price:     orderItem.Price,
		ProductID: orderItem.ProductID,
		OrderID:   orderItem.OrderID,
	}
}

func MapToOrderResponseFromOrder(order *domain.Order) *dto.OrderResponse {
	var orderItemsRes []*dto.OrderItemResponse

	for _, item := range order.Items {
		orderItemRes := mapToOrderItemResponseFromOrderItem(item)
		orderItemsRes = append(orderItemsRes, orderItemRes)
	}

	return &dto.OrderResponse{
		ID:            order.ID,
		PaymentMethod: order.PaymentMethod,
		TaxPrice:      order.TaxPrice,
		ShippingPrice: order.ShippingPrice,
		TotalPrice:    order.TotalPrice,
		Items:         orderItemsRes,
		UserID:        order.UserID,
		Status:        order.Status,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}
}
