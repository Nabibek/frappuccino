package api

import (
	"frappuccino/internal/api/handlers"
	"net/http"
)

func Router(handlers *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /inventory", handlers.InventoryHandler.CreateInventoryIngredient)
	mux.HandleFunc("GET /inventory", handlers.InventoryHandler.GetInventory)
	mux.HandleFunc("GET /inventory/{id}", handlers.InventoryHandler.GetIngredientByID)
	mux.HandleFunc("PUT /inventory/{id}", handlers.InventoryHandler.UpdateIngredient)
	mux.HandleFunc("DELETE /inventory/{id}", handlers.InventoryHandler.DeleteIngredient)
	return mux
}
