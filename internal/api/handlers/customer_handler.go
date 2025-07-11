package handlers

import (
	"encoding/json"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log"
	"net/http"
)

type CustomerHandler struct {
	customerService service.CustomerServiceInf
}

func NewCustomerHandler(service service.CustomerServiceInf) *CustomerHandler {
	return &CustomerHandler{customerService: service}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var input models.Customer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Printf("input %v", input)

	err := h.customerService.Create(r.Context(), &input)
	if err != nil {
		log.Printf("failed to create Customer: %v", err) // <- вот здесь логируем ошибку
		http.Error(w, "failed to create Customer", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	customer, err := h.customerService.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to get Customers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // FIX
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	customer, err := h.customerService.GetCustomerByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var input models.Customer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.customerService.UpdateCustomerByID(r.Context(), &input)
	if err != nil {
		http.Error(w, "failed to update customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Customer updated successfully"}`))
}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id") // FIX
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	err := h.customerService.DeleteCustomerByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, "failed to delete customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Customer deleted successfully"}`))
}
