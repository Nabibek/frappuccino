package service

import (
	"context"
	"frappuccino/internal/repo"
	"frappuccino/models"
	"log"
)

type AggregationServiceInf interface {
	TotalPrice() (float64, error)
	PopularItems() (models.PopularItems, error)
	Search(ctx context.Context, q string, filters []string, minPrice, maxPrice float64) (models.Search, error)
	OrderedItemByPeriod(period string, month string, year string) (models.ListOrderedItemByPeriods, error)
}

type AggregationService struct {
	aggregationRepo repo.AggregationRepo
}

func NewAggregationService(aggregationRepo repo.AggregationRepo) *AggregationService {
	return &AggregationService{aggregationRepo: aggregationRepo}
}

func (s *AggregationService) TotalPrice() (float64, error) {
	log.Println("Count TotalPrice")
	res, err := s.aggregationRepo.TotalPrice()
	if err != nil {
		log.Printf("Failed to get TotalPrice: %v", err)
		return 0, err
	}
	log.Printf("TotalPrice: %v", res)
	return res, nil
}
func (s *AggregationService) PopularItems() (models.PopularItems, error) {
	log.Println("Get PopulatItems")
	res, err := s.aggregationRepo.PopularItems()
	if err != nil {
		log.Printf("Failed to get PoplarItem: %v", err)
		return models.PopularItems{}, err
	}
	log.Println("Succes to get popular item")
	return res, nil
}

func (s *AggregationService) Search(ctx context.Context, q string, filters []string, minPrice, maxPrice float64) (models.Search, error) {
	log.Printf("Search by filer: %v", filters)
	res, err := s.aggregationRepo.Search(ctx, q, filters, minPrice, maxPrice)
	if err != nil {
		log.Printf("Failed Search: %v", err)
		return models.Search{}, err
	}
	log.Panicln("Success to search")
	return res, nil
}
func (s *AggregationService) OrderedItemByPeriod(period string, month string, year string) (models.ListOrderedItemByPeriods, error) {
	log.Println("Get OrderedItem by period")
	res, err := s.aggregationRepo.OrderedItemByPeriod(period, month, year)
	if err != nil {
		log.Printf("Failed get ordered itema by period")
		return models.ListOrderedItemByPeriods{}, err
	}
	log.Println("success to get ordered item by period")
	return res, nil
}
