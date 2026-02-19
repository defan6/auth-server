package storage

import (
	"fmt"

	"github.com/defan6/market/services/order-service/internal/config"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(cfg *config.DBConfig) *Database {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return &Database{db: db}
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
