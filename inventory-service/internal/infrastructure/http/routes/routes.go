package routes

import (
	"inventory-service/internal/infrastructure/http/handlers"

	"github.com/gorilla/mux"
)

// NewRouter configura e retorna o roteador principal da aplicação
func NewRouter(productHandler *handlers.ProductHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/products", productHandler.Create).Methods("POST")
	router.HandleFunc("/products/{id}/reserve-stock", productHandler.ReserveStock).Methods("POST")
	router.HandleFunc("/products/{id}/confirm-stock", productHandler.ConfirmStock).Methods("POST")
	router.HandleFunc("/products/{id}/cancel-reserve", productHandler.CancelReservation).Methods("POST")
	return router
}
