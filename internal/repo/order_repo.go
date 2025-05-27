package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"frappuccino/models"
	"frappuccino/utils"

	"github.com/lib/pq"
)

type OrderRepo interface {
	Create(ctx context.Context, order *models.Order) error
	Orders(ctx context.Context) ([]models.Order, error)
	GetOrderByID(ctx context.Context, orderId string) (models.Order, error)
	UpdateOrdeItemrByID(ctx context.Context, orderItems *models.OrderItems) error
	DeleteOrderByID(ctx context.Context, orderId string) error
	checkIngregients(tx *sql.Tx, orderItems []models.OrderItems) error
	minusInventory(tx *sql.Tx, orderItems []models.OrderItems) error
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	// check inventory
	err = r.checkIngregients(tx, order.OrderItems)
	if err != nil {
		return err
	}
	err = r.db.QueryRowContext(ctx,
		`INSERT INTO orders (customer_id,special_instructions,order_payment_method)
		VALUES ($1,$2,$3)
		RETURNING order_id,created_at,updated_at;`, order.CustomerId, order.SpecialInstructions, order.PaymentMethod).Scan(&order.OrderId, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	var totalPrice utils.DEC

	for _, items := range order.OrderItems {
		err = tx.QueryRow(`SELECT mi.price FROM menu_items mi WHERE menu_item_id = $1`, items.MenuItemId).Scan(&items.UnitPrice)
		if err != nil {
			return err
		}

		totalPrice += items.Quantity * items.UnitPrice // add unit_price from menu Items
	}
	order.TotalPrice = totalPrice
	err = r.minusInventory(tx, order.OrderItems)

	order.OrderStatus = "PENDING"

	return tx.Commit()
}

func (r *OrderRepository) checkIngregients(tx *sql.Tx, orderItems []models.OrderItems) error {
	query1 := `
	SELECT i.quantity>= mii.quantity * $1
	FROM inventory i
	JOIN menu_item_ingredients mii USING(ingredient_id)
	WHERE mii.menu_item_id = $2`

	for _, item := range orderItems {
		var have bool
		err := tx.QueryRow(query1, item.Quantity, item.MenuItemId).Scan(&have)

		if err != nil || have == false {
			return fmt.Errorf("Doesn't have inventory for menu item %d: %w", item.MenuItemId, err)
		}
	}

	return nil
}

func (r *OrderRepository) minusInventory(tx *sql.Tx, orderItems []models.OrderItems) error {
	for _, item := range orderItems {
		_, err := tx.Exec(`
            WITH ingredients AS (
                SELECT ingredient_id, quantity 
                FROM menu_item_ingredients 
                WHERE menu_item_id = $1
            )
            UPDATE inventory i
            SET quantity = i.quantity - (ing.quantity * $2)
            FROM ingredients ing
            WHERE i.ingredient_id = ing.ingredient_id`,
			item.MenuItemId, item.Quantity,
		)
		if err != nil {
			return fmt.Errorf("failed to deduct ingredient from inventory: %w", err)
		}
	}
	return nil
}

func (r *OrderRepository) Orders(ctx context.Context) ([]models.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM orders`)
	if err != nil {
		return []models.Order{}, err
	}
	defer rows.Close()
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.OrderId, &order.CustomerId, pq.Array(&order.OrderItems), &order.SpecialInstructions, &order.TotalPrice, &order.OrderStatus, &order.PaymentMethod, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderId string) (models.Order, error) {
	var order models.Order
	err := r.db.QueryRowContext(ctx, `SELECT * FROM orders WHERE order_id = $1`, orderId).Scan(&order.OrderId, &order.CustomerId, pq.Array(&order.OrderItems), &order.SpecialInstructions, &order.TotalPrice, &order.OrderStatus, &order.PaymentMethod, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, fmt.Errorf("Item not found: %w", err)
		}
		return models.Order{}, fmt.Errorf("failed to get Item: %w", err)
	}
	return order, nil
}

func (r *OrderRepository) UpdateOrdeItemrByID(ctx context.Context, orderItems *models.OrderItems) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `
	UPDATE order_item
	SET 
	menu_item_id = $3,
	customizations = $4,
	quantity = $5,
	WHERE order_id = $1 AND order_item_id= $2
	`, orderItems.OrderId, orderItems.OrderItemId, orderItems.MenuItemId, orderItems.Customizations, orderItems.Quantity)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *OrderRepository) DeleteOrderByID(ctx context.Context, orderId string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `DELETE FROM orders WHERE id = $1`, orderId)
	if err != nil {
		return fmt.Errorf("failed to delete Orders: %w", err)
	}

	// Verify exactly one row was deleted
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Commit transaction if everything succeeded
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
