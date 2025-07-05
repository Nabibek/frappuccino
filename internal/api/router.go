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

	mux.HandleFunc("POST /menu", handlers.MenuHandler.CreateMenuItem)
	mux.HandleFunc("GET /menu", handlers.MenuHandler.GetAllMenu)
	mux.HandleFunc("PUT /menu/{id}", handlers.MenuHandler.UpdateMenuItem)
	mux.HandleFunc("DELETE /menu/{id}", handlers.MenuHandler.DeleteMenuItem)
	mux.HandleFunc("GET /menu/{id}", handlers.MenuHandler.GetIngredientByID)
	mux.HandleFunc("POST /order", handlers.OrderHandler.CreateOrder)

	mux.HandleFunc("GET /order", handlers.OrderHandler.Orders)
	mux.HandleFunc("GET /order/{id}", handlers.OrderHandler.GetOrderByID)
	mux.HandleFunc("PUT /order/{id}", handlers.OrderHandler.UpdateOrderItem)
	mux.HandleFunc("DELETE /order/{id}", handlers.OrderHandler.DeleteOrderByID)
	mux.HandleFunc("PUT /order/status/{id}", handlers.OrderHandler.UpdateStatusOrder)

	mux.HandleFunc("POST /customer", handlers.CustomerHandler.CreateCustomer)
	mux.HandleFunc("GET /customer", handlers.CustomerHandler.GetAllCustomers)
	mux.HandleFunc("GET /customer/{id}", handlers.CustomerHandler.GetCustomerByID)
	mux.HandleFunc("PUT /customer/{id}", handlers.CustomerHandler.UpdateCustomer)
	mux.HandleFunc("DELETE /customer/{id}", handlers.CustomerHandler.DeleteCustomer)

	mux.HandleFunc("GET /totalprice", handlers.AggregationHandler.TotalPrice)
	mux.HandleFunc("GET /popularitems", handlers.AggregationHandler.PopularItems)
	mux.HandleFunc("GET /search", handlers.AggregationHandler.Search)
	mux.HandleFunc("GET /orderedItems", handlers.AggregationHandler.OrderedItemByPeriod)
	return mux
}
