package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log"
	"net/http"
)

type OrderHandler struct {
	orderServise service.OrderServiseInf
}

func NewOrderHandler(service service.OrderServiseInf) *OrderHandler {
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
	err := h.orderServise.Create(r.Context(), &input)
	if err != nil {
		log.Printf("failed to create ingredient: %v", err) // <- вот здесь логируем ошибку
		http.Error(w, "failed to create ingredient", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) Orders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderServise.Orders(r.Context())
	if err != nil {
		http.Error(w, "failed to get orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) UpdateOrdeItem(w http.ResponseWriter, r *http.Request) {
	var input models.OrderItems
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err := h.orderServise.UpdateOrdeItemrByID(r.Context(), &input)
	if err != nil {
		http.Error(w, "failed to update order Item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Order item updated successfully"}`))
}

func (h *OrderHandler) DeleteOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // FIX
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	err := h.orderServise.DeleteOrderByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, "failed to delete order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Order deleted successfully"}`))
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // FIX
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	order, err := h.orderServise.GetOrderByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
func (h *OrderHandler) UpdateStatusOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // FIX
	status := r.URL.Query().Get("status")
	if status != "CANCELLED" || status != "PENDING" || status != "COMPLETED" {
		http.Error(w, "incorrect status", http.StatusBadRequest)
		return
	}
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	err := h.orderServise.UpdateStatusOrder(r.Context(), idStr, status)
	if err != nil {
		http.Error(w, "Can not update status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Order status updated successfully"}`))
}
