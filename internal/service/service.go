package service

import "frappuccino/internal/repo"

type Service struct {
	CustomerService  CustomerServiceInf
	InventoryService InventoryServiceInf
	MenuService      MenuServiceInf
	OrderService     OrderServiseInf
}

func New(repo *repo.Repository) *Service {
	var service Service
	service.CustomerService = NewCustomerService(repo.CustomerRepo)
	service.InventoryService = NewInventoryService(repo.InventoryRepo)
	service.MenuService = NewMenuService(repo.MenuRepo)
	service.OrderService = NewOrderService(repo.OrderRepo)
	return &service
}
