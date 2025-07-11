package repo

import "database/sql"

type Repository struct {
	CustomerRepo    CustomerRepo
	InventoryRepo   InventoryRepo
	MenuRepo        MenuRepo
	OrderRepo       OrderRepo
	AggregationRepo AggregationRepo
}

func New(db *sql.DB) *Repository {
	return &Repository{
		CustomerRepo:    NewCustomerRepository(db),
		InventoryRepo:   NewInventoryRepository(db),
		MenuRepo:        NewMenuRepository(db),
		OrderRepo:       NewOrderRepository(db),
		AggregationRepo: NewAggregationRepository(db),
	}
}
