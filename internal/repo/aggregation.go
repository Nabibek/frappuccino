package repo

import (
	"context"
	"database/sql"
	"fmt"
	"frappuccino/models"
)

type AggregationRepo interface {
	TotalPrice() (float64, error)
	PopularItems() (models.PopularItems, error)
	Search(ctx context.Context, q string, filters []string, minPrice, maxPrice float64) (models.Search, error)
	OrderedItemByPeriod(period string, month string, year string) (models.ListOrderedItemByPeriods, error)
	searchMenu(ctx context.Context, q string, minPrice, maxPrice float64) ([]models.SearchMenu, error)
	searchOrder(ctx context.Context, q string, minPrice, maxPrice float64) ([]models.SearchOrder, error)
}

type AggregationRepository struct {
	db *sql.DB
}

func NewAggregationRepository(db *sql.DB) *AggregationRepository {
	return &AggregationRepository{db: db}
}

func (r *AggregationRepository) TotalPrice() (float64, error) {
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
func (r *AggregationRepository) PopularItems() (models.PopularItems, error) {
	query := `SELECT item_name, count(item_name)
			 FROM order_items
			 GROUP BY item_name
			 ORDER BY count(item_name) DESC
			 LIMIT 10`
	var popularItems models.PopularItems

	rows, err := r.db.Query(query)
	if err != nil {
		return models.PopularItems{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var popularItem models.PopularItem

		err = rows.Scan(&popularItem.ItemName, &popularItem.OrderedTimes)
		if err != nil {
			return models.PopularItems{}, err
		}
		popularItems.Items = append(popularItems.Items, popularItem)
	}
	return popularItems, nil
}

func (r *AggregationRepository) Search(ctx context.Context, q string, filters []string, minPrice, maxPrice float64) (models.Search, error) {
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

func (r *AggregationRepository) OrderedItemByPeriod(period string, month string, year string) (models.ListOrderedItemByPeriods, error) {
	var list models.ListOrderedItemByPeriods
	if period == "day" {
		query := `WITH base_date AS (SELECT TO_DATE ($1 || ' ' || $2, 'Month YYYY')AS first_day), days AS (SELECT generate series( 1,date_part('days', (DATE_TRUNC('month',(SELECT first_day FROM base_date))+interval '1 month -1 day '))::int )AS n)
		SELECT 
		d.n,
		count(o.order_id)
		FROM days d
		LEFT JOIN orders o ON DATE(o.created_at)=(SELECT first_day FROM base_date) + (d.n -1) * INTERVAL '1 day'
		GROUP BY d.n
		ORDER BY d.n;`
		rows, err := r.db.Query(query, month, year)
		if err != nil {
			return models.ListOrderedItemByPeriods{}, err
		}
		defer rows.Close()
		for rows.Next() {
			var item models.OrderedItemByPeriod
			err = rows.Scan(&item.Date, &item.Count)
			if err != nil {
				return models.ListOrderedItemByPeriods{}, err
			}
			list.Items = append(list.Items, item)
		}
	} else if period == "month" {
		query := `WITH all_months AS (
		SELECT generate_series(1,12) AS month_num
	),
		monthly_orders AS (
		SELECT 
			EXTRACT (MONTH FROM o.created_at)::int AS month_num,
			COUNT (o.order_id) AS order_count
		FROM orders o
		WHERE EXTRACT(YEAR FROM o.created_at) = $1
		GROUP BY EXTRACT (MONTH FROM o.created_at)
	)
		SELECT 
			TO_CHAR(TO_DATE(m.month_num::text, 'MM'),'Month')AS month_name,
			COALESCE (mo.order_count,0) AS total_orders
		FROM all_months m
		LEFT JOIN monthly_orders mo ON m.month_num = mo.month_num
		ORDER BY m.month_num;`
		rows, err := r.db.Query(query, year)
		if err != nil {
			return models.ListOrderedItemByPeriods{}, err
		}
		defer rows.Close()
		for rows.Next() {
			var item models.OrderedItemByPeriod
			err = rows.Scan(&item.Date, &item.Count)
			if err != nil {
				return models.ListOrderedItemByPeriods{}, err
			}
			list.Items = append(list.Items, item)
		}
	}
	return list, nil
}

func (r *AggregationRepository) searchOrder(ctx context.Context, q string, minPrice, maxPrice float64) ([]models.SearchOrder, error) {
	query := `
SELECT o.order_id,c.full_name,oi.item_name,oi.price
FROM orders o 
JOIN order_items AS oi USING(order_id)
JOIN customers AS c USING(customer_id)
WHERE (c.full_name ILIKE  '%' || $1 || '%' OR  oi.item_name ILIKE '%' || $1 || '%';)
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
		err := rows.Scan(&order.Id, &order.CustomerName, &order.ItemName, &order.ItemPrice)
		if err != nil {
			return nil, err
		}
		allOrders = append(allOrders, order)
	}
	return allOrders, nil
}
func (r *AggregationRepository) searchMenu(ctx context.Context, q string, minPrice, maxPrice float64) ([]models.SearchMenu, error) {
	query := `SELECT menu_item_id,item_name,item_description,price
			  FROM menu_items
			  WHERE (item_name ILIKE '%' || $1 || '%' OR item_description ILIKE '%' || $1 || '%'; )
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
