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

	searchMenu(ctx context.Context, q string, minPrice, maxPrice int) ([]models.SearchMenu, error)
	searchOrder(ctx context.Context, q string, minPrice, maxPrice int) ([]models.SearchOrder, error)
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
func (r *ApplicationRepository) Search(ctx context.Context, q string, filters []string, minPrice, maxPrice int) (models.Search, error) {
	var err error
	var searchMenu []models.SearchMenu
	var searchOrder []models.SearchOrder
	for _, filter := range filters {
		if filter == "menu" {
			searchMenu, err = r.searchMenu(ctx, q, minPrice, maxPrice)
			if err != nil {
				return models.Search{}, err
			}
		} else if filter == "order" {
			searchOrder, err = r.searchOrder(ctx, q, minPrice, maxPrice)
			if err != nil {
				return models.Search{}, err
			}
		} else if filter == "0" || filter == "all" {
			searchMenu, err = r.searchMenu(ctx, q, minPrice, maxPrice)
			if err != nil {
				return models.Search{}, err
			}
			searchOrder, err = r.searchOrder(ctx, q, minPrice, maxPrice)
			if err != nil {
				return models.Search{}, err
			}
		}

	}
	return models.Search{
		MenuItems:  searchMenu,
		OrderItems: searchOrder,
	}, nil
}
func (r *ApplicationRepository) searchOrder(ctx context.Context, q string, minPrice, maxPrice int) ([]models.SearchOrder, error) {
	query := `
SELECT o.order_id,c.full_name,oi.item_name,oi.price
FROM orders o 
JOIN order_item AS oi USING(order_id)
JOIN customers AS c USING(customer_id)
WHERE (c.full_name ILIKE  '%' || $1 || '%' OR  oi.item_name ILIKE '%' || $1 || '%')
`
	var err error
	var rows *sql.Rows
	if minPrice != -1 && maxPrice != -1 {
		query += `AND oi.price BETWEEN $2 AND $3`
		rows, err = r.db.QueryContext(ctx, query, q, minPrice, maxPrice)
	} else if minPrice != -1 {
		query += `AND oi.price>=$2`
		rows, err = r.db.QueryContext(ctx, query, q, minPrice)
	} else if maxPrice != -1 {
		query += `AND oi.price<= $2`
		rows, err = r.db.QueryContext(ctx, query, q, maxPrice)
	} else {
		rows, err = r.db.QueryContext(ctx, query, q)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allOrders []models.SearchOrder

	for rows.Next() {
		var order models.SearchOrder
		err := rows.Scan(&order.Id, &order.CustomerName, &order.ItemName, order.ItemPrice)
		if err != nil {
			return nil, err
		}
		allOrders = append(allOrders, order)
	}
	return allOrders, nil
}
func (r *ApplicationRepository) searchMenu(ctx context.Context, q string, minPrice, maxPrice int) ([]models.SearchMenu, error) {
	query := `SELECT menu_item_id,item_name,item_description,price
			  FROM menu_items
			  WHERE (item_name ILIKE '%' || $1 || '%' OR item_description ILIKE '%' || $1 || '%' )
`
	var err error
	var rows *sql.Rows
	if minPrice != -1 && maxPrice != -1 {
		query += `AND price BETWEEN $2 AND $3`
		rows, err = r.db.QueryContext(ctx, query, q, minPrice, maxPrice)
	} else if minPrice != -1 {
		query += `AND price>=$2`
		rows, err = r.db.QueryContext(ctx, query, q, minPrice)
	} else if maxPrice != -1 {
		query += `AND price<= $2`
		rows, err = r.db.QueryContext(ctx, query, q, maxPrice)
	} else {
		rows, err = r.db.QueryContext(ctx, query, q)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var allMenu []models.SearchMenu
	for rows.Next() {
		var menu models.SearchMenu
		err := rows.Scan(&menu.Id, &menu.Name, &menu.Description, &menu.Price)
		if err != nil {
			return nil, err
		}
		allMenu = append(allMenu, menu)
	}
	return allMenu, nil
}
