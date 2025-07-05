package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type AggregationHandler struct {
	aggregationService service.AggregationServiceInf
}

func NewAggregationHandler(service service.AggregationServiceInf) *AggregationHandler {
	return &AggregationHandler{aggregationService: service}
}

func (h *AggregationHandler) PopularItems(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.aggregationService.PopularItems()
	if err != nil {
		http.Error(w, "Failed to get popular items", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(popularItems); err != nil {
		http.Error(w, "Failed to encode popular items", http.StatusInternalServerError)
		return
	}
}

func (h *AggregationHandler) TotalPrice(w http.ResponseWriter, r *http.Request) {
	totalPrice, err := h.aggregationService.TotalPrice()
	if err != nil {
		http.Error(w, "Failed to get total price", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]float64{"total_price": totalPrice}); err != nil {
		http.Error(w, "Failed to encode total price", http.StatusInternalServerError)
		return
	}
}

func (h *AggregationHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := query.Get("q")
	filters := strings.Split(query.Get("filters"), ",")
	minPrice, _ := strconv.ParseFloat(query.Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(query.Get("max_price"), 64)

	searchResults, err := h.aggregationService.Search(r.Context(), q, filters, minPrice, maxPrice)
	if err != nil {
		http.Error(w, "Failed to perform search", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(searchResults); err != nil {
		http.Error(w, "Failed to encode search results", http.StatusInternalServerError)
		return
	}
}
func (h *AggregationHandler) OrderedItemByPeriod(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	period := query.Get("period")
	month := query.Get("month")
	year := query.Get("year")

	orderedItems, err := h.aggregationService.OrderedItemByPeriod(period, month, year)
	if err != nil {
		http.Error(w, "Failed to get ordered items by period", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orderedItems); err != nil {
		http.Error(w, "Failed to encode ordered items", http.StatusInternalServerError)
		return
	}
}
