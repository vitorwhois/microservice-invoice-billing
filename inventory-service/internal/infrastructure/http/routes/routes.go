package routes

import (
	"github.com/vitorwhois/microservice-invoice-billing/inventory-service/internal/infrastructure/http/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(productHandler *handlers.ProductHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/products", productHandler.Create).Methods("POST")
	router.HandleFunc("/products/{id}/reserve-stock", productHandler.ReserveStock).Methods("POST")
	router.HandleFunc("/products/{id}/confirm-stock", productHandler.ConfirmStock).Methods("POST")
	router.HandleFunc("/products/{id}/cancel-reserve", productHandler.CancelReservation).Methods("POST")
	router.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	return router
}
