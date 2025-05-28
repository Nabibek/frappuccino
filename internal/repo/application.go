package repo

import (
	"context"
	"database/sql"
	"fmt"
)

type ApplicationRepo interface {
	TotalPrice(ctx context.Context) (float64, error)
	Search(ctx context.Context)
}

type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepo {
	return &ApplicationRepo{db: db}
}

func (r *ApplicationRepository) TotalPrice(ctx context.Context) (float64, error) {
	var totalPrice float64

	err := r.db.QueryRowContext(ctx, `
		SELECT SUM(total_price)
		FROM orders
		WHERE order_status = 'COMPLETED'
	`).Scan(&totalPrice)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to fetch total price: %w", err)
	}

	return totalPrice, nil
}
