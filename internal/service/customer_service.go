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
	GetCustomerByID(ctx context.Context, CustomerId string) (models.Customer, error)
	UpdateCustomerByID(ctx context.Context, customer *models.Customer) error
	DeleteCustomerByID(ctx context.Context, CustomerId string) error
}

type CustomerService struct {
	customerRepo repo.CustomerRepo
}

func NewCustomerService(customerRepo repo.CustomerRepo) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) Create(ctx context.Context, customer *models.Customer) error {
	log.Println("Creating new Customer :", customer.FullName)
	err := s.customerRepo.Create(ctx, customer)
	if err != nil {
		log.Printf("Failed to create menu customer '%s': %v", customer.FullName, err)
		return fmt.Errorf("could not create menu customer: %w", err)
	}
	log.Println("Customer  created successfully:", customer.CustomerId)
	return nil
}

func (s *CustomerService) GetAll(ctx context.Context) ([]models.Customer, error) {
	log.Println("Fetching all Customers")
	customers, err := s.customerRepo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to fetch Customer: %v", err)
		return nil, fmt.Errorf("could not retrieve menu: %w", err)
	}
	log.Printf("Retrieved %d customer", len(customers))
	return customers, nil
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, CustomerId string) (models.Customer, error) {
	log.Printf("Fetching customer by ID: %s", CustomerId)
	customer, err := s.customerRepo.GetCustomerByID(ctx, CustomerId)
	if err != nil {
		log.Printf("Failed to fetch customer [%s]: %v", CustomerId, err)
		return models.Customer{}, fmt.Errorf("could not get customer: %w", err)
	}
	log.Printf("Retrieved customer [%s]: %s", customer.CustomerId, customer.FullName)
	return customer, nil
}

func (s *CustomerService) UpdateCustomerByID(ctx context.Context, customer *models.Customer) error {
	log.Printf("Updating  Customer [%s]", customer.CustomerId)
	err := s.customerRepo.UpdateCustomerByID(ctx, customer)
	if err != nil {
		log.Printf("Failed to update  Customer [%s]: %v", customer.CustomerId, err)
		return fmt.Errorf("could not update  Customer: %w", err)
	}
	log.Printf("Customer [%s] updated successfully", customer.CustomerId)
	return nil
}

func (s *CustomerService) DeleteCustomerByID(ctx context.Context, CustomerId string) error {
	log.Printf("Deleting  Customer [%s]", CustomerId)
	err := s.customerRepo.DeleteCustomerByID(ctx, CustomerId)
	if err != nil {
		log.Printf("Failed to  Customer [%s]: %v", CustomerId, err)
		return fmt.Errorf("could not delete  Customer: %w", err)
	}
	log.Printf("Customer [%s] deleted successfully", CustomerId)
	return nil
}
