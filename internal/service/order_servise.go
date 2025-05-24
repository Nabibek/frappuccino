package service

import (
	"context"
	"frappuccino/internal/repo"
	"frappuccino/models"
	"log"
)

type OrderServiseInf interface {
	Create(ctx context.Context, order *models.Order) error
	Orders(ctx context.Context) ([]models.Order, error)
	GetOrderByID(ctx context.Context, orderId string) (models.Order, error)
	UpdateOrdeItemrByID(ctx context.Context, orderItems *models.OrderItems) error
	DeleteOrderByID(ctx context.Context, orderId string) error
}

type OrderServise struct {
	orderRepo repo.OrderRepo
}

func NewOrderService(orderRepo repo.OrderRepo) *OrderServise {
	return &OrderServise{orderRepo: orderRepo}
}

func (s *OrderServise) Create(ctx context.Context, order *models.Order) error {
	log.Println("Create new order", order.OrderId)
	err := s.orderRepo.Create(ctx, order)
	if err != nil {
		log.Println("Failed to create order")
		return err
	}
	log.Println("Order created successfully", order.OrderId)
	return nil
}

func (s *OrderServise) Orders(ctx context.Context) ([]models.Order, error) {
	log.Println("Get orders ")
	orders, err := s.Orders(ctx)
	if err != nil {
		log.Println("Failed to get orders")
		return nil, err
	}
	log.Println("Get orders successfully")
	return orders, nil
}

func (s *OrderServise) GetOrderByID(ctx context.Context, orderId string) (models.Order, error) {
	log.Println("Get order BY id")
	order, err := s.GetOrderByID(ctx, orderId)
	if err != nil {
		log.Println("Failed to get order")
		return models.Order{}, err
	}
	log.Println("Get order successfully")
	return order, nil
}

func (s *OrderServise) UpdateOrdeItemrByID(ctx context.Context, orderItems *models.OrderItems) error {
	log.Println("updateing order items")
	err := s.UpdateOrdeItemrByID(ctx, orderItems)
	if err != nil {
		log.Println("Failed update order item ")
		return err
	}
	log.Println("Updating order item successfully")
	return nil
}

func (s *OrderServise) DeleteOrderByID(ctx context.Context, orderId string) error {
	log.Println("Deleting order")
	err := s.DeleteOrderByID(ctx, orderId)
	if err != nil {
		log.Println("Failed to delete order")
		return err
	}
	return nil
}
