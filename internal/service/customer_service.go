package service

import (
	"context"
	"fmt"
	"frappuccino/internal/repo"
	"frappuccino/models"
	"log"
)

type CustomerServiceInf interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetAll(ctx context.Context) ([]models.Customer, error)
	GetItemByID(ctx context.Context, CustomerId string) (models.Customer, error)
	UpdateItemByID(ctx context.Context, customer *models.Customer) error
	DeleteItemByID(ctx context.Context, CustomerId string) error
}

type CustomerService struct {
	customerRepo repo.CustomerRepo
}

func NewCustomerService(customerRepo repo.CustomerRepo) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) Create(ctx context.Context, customer *models.Customer) error {
	log.Println("Creating new Customer item:", customer.FullName)
	err := s.customerRepo.Create(ctx, customer)
	if err != nil {
		log.Printf("Failed to create menu customer '%s': %v", customer.FullName, err)
		return fmt.Errorf("could not create menu customer: %w", err)
	}
	log.Println("Customer item created successfully:", customer.CustomerId)
	return nil
}

func (s *CustomerService) GetAll(ctx context.Context) ([]models.Customer, error) {
	log.Println("Fetching all menu items")
	menu, err := s.customerRepo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to fetch menu items: %v", err)
		return nil, fmt.Errorf("could not retrieve menu: %w", err)
	}
	log.Printf("Retrieved %d menu items", len(menu))
	return menu, nil
}

func (s *CustomerService) GetItemByID(ctx context.Context, CustomerId string) (models.Customer, error) {
	log.Printf("Fetching menu item by ID: %s", CustomerId)
	customer, err := s.customerRepo.GetItemByID(ctx, CustomerId)
	if err != nil {
		log.Printf("Failed to fetch menu item [%s]: %v", CustomerId, err)
		return models.Customer{}, fmt.Errorf("could not get menu item: %w", err)
	}
	log.Printf("Retrieved menu item [%s]: %s", customer.CustomerId, customer.FullName)
	return customer, nil
}

func (s *CustomerService) UpdateItemByID(ctx context.Context, customer *models.Customer) error {
	log.Printf("Updating menu item [%s]", customer.CustomerId)
	err := s.customerRepo.UpdateItemByID(ctx, customer)
	if err != nil {
		log.Printf("Failed to update menu item [%s]: %v", customer.CustomerId, err)
		return fmt.Errorf("could not update menu item: %w", err)
	}
	log.Printf("Menu item [%s] updated successfully", customer.CustomerId)
	return nil
}

func (s *CustomerService) DeleteItemByID(ctx context.Context, CustomerId string) error {
	log.Printf("Deleting menu item [%s]", CustomerId)
	err := s.customerRepo.DeleteItemByID(ctx, CustomerId)
	if err != nil {
		log.Printf("Failed to delete menu item [%s]: %v", CustomerId, err)
		return fmt.Errorf("could not delete menu item: %w", err)
	}
	log.Printf("Menu item [%s] deleted successfully", CustomerId)
	return nil
}
