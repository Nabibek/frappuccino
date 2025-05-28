package repo

import (
	"context"
	"database/sql"
	"fmt"
	"frappuccino/models"
)

type ApplicationRepo interface {
	TotalPrice() (float64, error)
	Search(ctx context.Context)
}

type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) TotalPrice() (float64, error) {
	var totalPrice float64

	err := r.db.QueryRow(`
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
func (r *ApplicationRepository) Search(ctx context.Context, q string, filters []string, minPrice, maxPrice *float64) (SearchResult, error) {
	for _, filter := range filters {
		if filter == "menu" {
			searchMenu, err := SearchMenu(ctx, q, filter, minPrice, maxPrice)
			if err != nil {
				return nil, err
			}
		} else if filter == "order" {
			searchOrder, err := SearchOrder(ctx, q, filter, minPrice, maxPrice)
			if err != nil {
				return nil, err
			}
		} else if filter == "0" || filter == "all" {
			searchMenu, err := SearchMenu(ctx, q, filter, minPrice, maxPrice)
			if err != nil {
				return nil, err
			}
			searchOrder, err := SearchOrder(ctx, q, filter, minPrice, maxPrice)
			if err != nil {
				return nil, err
			}
		}

	}
	return models.Search{
		MenuItems:  []searchMenu,
		OrderItems: []searchOrder,
	}
}
func (r *ApplicationRepository) searchOrder(ctx context.Context, q string, filters []string, minPrice, maxPrice *float64) ([]models.SearchOrder, error) {
	var allOrders []models.SearchOrder
	query := `
SELECT o.order_id,c.customer_name,oi.item_name, 		
`
}
