package storage

import (
	"context"
	"database/sql"

	"github.com/defan6/market/services/order-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type defaultOrderStorage struct {
	db *sqlx.DB
}

func NewDefaultOrderStorage(db *sqlx.DB) *defaultOrderStorage {
	return &defaultOrderStorage{
		db: db,
	}
}

type execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type querier interface {
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
}

func executor(ctx context.Context, db *sqlx.DB) execer {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return db
}

func querierExec(ctx context.Context, db *sqlx.DB) querier {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return db
}

func (r *defaultOrderStorage) CreateOrder(ctx context.Context, order *domain.Order) error {
	err := querierExec(ctx, r.db).QueryRowxContext(ctx,
		`INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, user_id, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		order.PaymentMethod, order.TaxPrice, order.ShippingPrice,
		order.TotalPrice, order.UserID, order.Status, order.CreatedAt,
	).Scan(&order.ID)
	return err
}

func (r *defaultOrderStorage) CreateOrderItems(ctx context.Context, items []*domain.OrderItem) error {
	for i, item := range items {
		err := querierExec(ctx, r.db).QueryRowxContext(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, price, name, image) 
			 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
			item.OrderID, item.ProductID, item.Quantity, item.Price, item.Name, item.Image,
		).Scan(&items[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *defaultOrderStorage) GetOrder(ctx context.Context, id int64) (*domain.Order, error) {
	order := &domain.Order{}
	err := r.db.GetContext(ctx, order,
		`SELECT id, payment_method, tax_price, shipping_price, total_price, user_id, status, created_at, updated_at
		 FROM orders WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	return order, nil
}
