package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log"
	"net/http"
)

type OrderHandler struct {
	orderServise service.OrderServise
}

func NewOrderHandler(service service.OrderServise) *OrderHandler {
	return &OrderHandler{orderServise: service}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input models.Order
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Println("input %v", input)
	order, err := h.orderServise.Create(r.Context(), input)
	if err != nil {
		log.Printf("failed to create ingredient: %v", err) // <- вот здесь логируем ошибку
		http.Error(w, "failed to create ingredient", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
