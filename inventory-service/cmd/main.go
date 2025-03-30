// cmd/main.go
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"inventory-service/internal/application/product"
	"inventory-service/internal/infrastructure/http/handlers"
	"inventory-service/internal/infrastructure/persistence"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	productRepo := persistence.NewProductRepository(db)
	productService := product.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	router := mux.NewRouter()
	router.HandleFunc("/products", productHandler.Create).Methods("POST")
	router.HandleFunc("/products/{id}/reserve-stock", productHandler.ReserveStock).Methods("POST")
	router.HandleFunc("/products/{id}/confirm-stock", productHandler.ConfirmStock).Methods("POST")
	router.HandleFunc("/products/{id}/cancel-reserve", productHandler.CancelReservation).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func setupDatabase() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost user=postgres password=postgres dbname=inventory sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
