package repo

import (
	"context"
	"database/sql"
	"fmt"
	"frappuccino/models"
	"frappuccino/utils"

	"github.com/lib/pq"
)

type OrderRepo interface {
	Create(ctx context.Context, order models.Order) (models.Order, error)
	Orders(ctx context.Context) ([]models.Order, error)
	GetItemByID(ctx context.Context, orderId string) (models.Order, error)
	UpdateItemByID(ctx context.Context, order models.Order) error
	DeleteItemByID(ctx context.Context, orderId string) error
	checkIngregients(tx *sql.Tx, orderItems []models.OrderItems) error
	minusInventory(tx *sql.Tx, orderItems []models.OrderItems) error
}

type orderRepo struct {
	*Repository
}

func NewOrderRepository(db *sql.DB) OrderRepo {
	return &orderRepo{NewRepository(db)}
}

func (r *orderRepo) Create(ctx context.Context, order models.Order) (models.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	// check inventory
	err = r.checkIngregients(tx, order.OrderItems)
	if err != nil {
		return models.Order{}, err
	}
	err = r.db.QueryRowContext(ctx,
		`INSERT INTO orders (customer_id,special_instructions,order_payment_method)
		VALUES ($1,$2,$3,)
		RETURNING order_id,created_at,updated_at;`, order.CustomerId, order.SpecialInstructions, order.PaymentMethod).Scan(order.OrderId, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return models.Order{}, err
	}
	var totalPrice utils.DEC

	for _, items := range order.OrderItems {
		totalPrice += items.Quantity * items.UnitPrice // add unit_price from menu Items
	}
	order.TotalPrice = totalPrice
	err = r.minusInventory(tx, order.OrderItems)

	order.OrderStatus = "pending"

	return order, tx.Commit()
}

func (r *orderRepo) checkIngregients(tx *sql.Tx, orderItems []models.OrderItems) error {
	query1 := `
	SELECT i.quantity>= mii.quantity * $1
	FROM inventory i
	JOIN menu_item_ingredients mii USING(ingredient_id)
	WHERE mii.menu_item_id = $2`

	for _, item := range orderItems {
		var have bool
		err := tx.QueryRow(query1, item.Quantity, item.MenuItemId).Scan(&have)

		if err != nil || have == false {
			return fmt.Errorf("Doesn't have inventory for menu item %d: %w",
				item.MenuItemId, err)
		}
	}

	return nil
}

func (r *orderRepo) minusInventory(tx *sql.Tx, orderItems []models.OrderItems) error {
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

func (r *orderRepo) Orders(ctx context.Context) ([]models.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM orders`)
	if err != nil {
		return []models.Order{}, err
	}
	defer rows.Close()
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.OrderId, &order.CustomerId,pq.Array(&order.OrderItems), &order.SpecialInstructions, &order.TotalPrice, &order.OrderStatus, &order.PaymentMethod, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}


func (r *orderRepo)GetItemByID(ctx context.Context, orderId string) (models.Order, error){
	r.db.QueryContext(ctx,`SELECT * FROM orders WHERE order_id = $1`,orderId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.MenuItems{}, fmt.Errorf("Item not found: %w", err)
		}
		return models.MenuItems{}, fmt.Errorf("failed to get Item: %w", err)
	}
	return item, nil
}
func (r *orderRepo)UpdateItemByID(ctx context.Context, order models.Order) error{
	tx, err :+ r.db.BeginTx(ctx,nil)
	if err != nil{
		return err
	}
	defer tx.Rollback()
	res, err  := tx.ExecContext(ctx, `
	UPDATE order_item
	SET order_item_id,
	order_id,
	menu_item_id,
	customizations,
	quantity,
	`)
}
DeleteItemByID(ctx context.Context, orderId string) error