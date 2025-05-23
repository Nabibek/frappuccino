package service

import (
	"context"
	"frappuccino/internal/repo"
	"frappuccino/models"
	"log"
)

type OrderServise interface {
	Create(ctx context.Context, order models.Order) (models.Order, error)
	Orders(ctx context.Context) ([]models.Order, error)
	GetOrderByID(ctx context.Context, orderId string) (models.Order, error)
	UpdateOrdeItemrByID(ctx context.Context, orderItems models.OrderItems) error
	DeleteOrderByID(ctx context.Context, orderId string) error
}

type orderServise struct {
	orderRepo repo.OrderRepo
}

func NewOrderService(orderRepo repo.OrderRepo) OrderServise {
	return &orderServise{orderRepo: orderRepo}
}

func (s *orderServise) Create(ctx context.Context, order models.Order) (models.Order, error) {
	log.Println("Create new order", order.OrderId)
	created, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		log.Println("Failed to create order")
		return models.Order{}, err
	}
	log.Println("Order created successfully", created.OrderId)
	return created, nil
}

func (s *orderServise) Orders(ctx context.Context) ([]models.Order, error) {
	log.Println("Get orders ")
	orders, err := s.Orders(ctx)
	if err != nil {
		log.Println("Failed to get orders")
		return nil, err
	}
	log.Println("Get orders successfully")
	return orders, nil
}

func (s *orderServise) GetOrderByID(ctx context.Context, orderId string) (models.Order, error) {
	log.Println("Get order BY id")
	order, err := s.GetOrderByID(ctx, orderId)
	if err != nil {
		log.Println("Failed to get order")
		return models.Order{}, err
	}
	log.Println("Get order successfully")
	return order, nil
}

func (s *orderServise) UpdateOrdeItemrByID(ctx context.Context, orderItems models.OrderItems) error {
	log.Println("updateing order items")
	err := s.UpdateOrdeItemrByID(ctx, orderItems)
	if err != nil {
		log.Println("Failed update order item ")
		return err
	}
	log.Println("Updating order item successfully")
	return nil
}

func (s *orderServise) DeleteOrderByID(ctx context.Context, orderId string) error {
	log.Println("Deleting order")
	err := s.DeleteOrderByID(ctx, orderId)
	if err != nil {
		log.Println("Failed to delete order")
		return err
	}
	return nil
}
