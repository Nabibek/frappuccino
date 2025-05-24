package handlers

import "frappuccino/internal/service"

type Handler struct {
	CustomerHandler  *CustomerHandler
	InventoryHandler *InventoryHandler
	MenuHandler      *MenuHandler
	OrderHandler     *OrderHandler
}

func New(service *service.Service) *Handler {
	return &Handler{
		CustomerHandler:  NewCustomerHandler(service.CustomerService),
		InventoryHandler: NewInventoryHandler(service.InventoryService),
		MenuHandler:      NewMenuHandler(service.MenuService),
		OrderHandler:     NewOrderHandler(service.OrderService),
	}
}
